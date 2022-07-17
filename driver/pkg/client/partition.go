// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	multiraftv1 "github.com/atomix/multi-raft/api/atomix/multiraft/v1"
	"github.com/atomix/runtime/pkg/errors"
	"github.com/atomix/runtime/pkg/grpc/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"sync"
)

func newPartition(client *Client, id multiraftv1.PartitionID) *Partition {
	return &Partition{
		client: client,
		id:     id,
	}
}

type Partition struct {
	client    *Client
	id        multiraftv1.PartitionID
	state     *PartitionState
	watchers  map[int]chan<- PartitionState
	watcherID int
	conn      *grpc.ClientConn
	mu        sync.RWMutex
}

func (p *Partition) ID() multiraftv1.PartitionID {
	return p.id
}

func (p *Partition) connect(ctx context.Context, config *multiraftv1.PartitionConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	address := fmt.Sprintf("%s:///%s:%d", resolverName, config.Host, config.Port)
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"raft"}`),
		grpc.WithResolvers(newResolver(p)),
		grpc.WithUnaryInterceptor(retry.RetryingUnaryClientInterceptor(retry.WithRetryOn(codes.Unavailable))),
		grpc.WithStreamInterceptor(retry.RetryingStreamClientInterceptor(retry.WithRetryOn(codes.Unavailable))))
	if err != nil {
		return errors.FromProto(err)
	}
	p.conn = conn

	client := multiraftv1.NewPartitionClient(p.conn)
	request := &multiraftv1.WatchPartitionRequest{}
	stream, err := client.Watch(context.Background(), request)
	if err != nil {
		return errors.FromProto(err)
	}

	go func() {
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Error(err)
			} else {
				switch e := event.Event.(type) {
				case *multiraftv1.WatchPartitionResponse_LeaderUpdated:
					var leader string
					var followers []string
					p.client.mu.RLock()
					config := p.client.config
					p.client.mu.RUnlock()
					for _, replica := range config.Replicas {
						address := fmt.Sprintf("%s:%d", replica.Host, replica.Port)
						if replica.NodeID == e.LeaderUpdated.Leader {
							leader = address
						} else {
							followers = append(followers, address)
						}
					}
					p.notify(&PartitionState{
						Leader:    leader,
						Followers: followers,
					})
				}
			}
		}
	}()
	return nil
}

func (p *Partition) notify(state *PartitionState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.state = state
	go func() {
		p.mu.RLock()
		defer p.mu.RUnlock()
		for _, watcher := range p.watchers {
			watcher <- *state
		}
	}()
}

func (p *Partition) watch(ctx context.Context, ch chan<- PartitionState) error {
	p.mu.Lock()
	p.watcherID++
	watcherID := p.watcherID
	p.watchers[watcherID] = ch
	p.mu.Unlock()

	go func() {
		<-ctx.Done()
		p.mu.Lock()
		delete(p.watchers, watcherID)
		p.mu.Unlock()
	}()
	return nil
}

type PartitionState struct {
	Leader    string
	Followers []string
}