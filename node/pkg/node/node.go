// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"fmt"
	multiraftv1 "github.com/atomix/multi-raft-storage/api/atomix/multiraft/v1"
	"github.com/atomix/multi-raft-storage/node/pkg/node/manager"
	"github.com/atomix/multi-raft-storage/node/pkg/node/server"
	"github.com/atomix/multi-raft-storage/node/pkg/primitive"
	"github.com/atomix/runtime/sdk/pkg/logging"
	"github.com/atomix/runtime/sdk/pkg/runtime"
	"google.golang.org/grpc"
	"os"
)

var log = logging.GetLogger()

func New(network runtime.Network, opts ...Option) *MultiRaftNode {
	var options Options
	options.apply(opts...)
	registry := primitive.NewRegistry()
	for _, primitiveType := range options.PrimitiveTypes {
		registry.Register(primitiveType)
	}
	return &MultiRaftNode{
		Options: options,
		network: network,
		manager: manager.NewNodeManager(registry, options.Config),
		server:  grpc.NewServer(),
	}
}

type MultiRaftNode struct {
	Options
	network runtime.Network
	manager *manager.NodeManager
	server  *grpc.Server
}

func (s *MultiRaftNode) Start() error {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	lis, err := s.network.Listen(address)
	if err != nil {
		return err
	}

	multiraftv1.RegisterNodeServer(s.server, server.NewNodeServer(s.manager))
	multiraftv1.RegisterPartitionServer(s.server, server.NewPartitionServer(s.manager))
	multiraftv1.RegisterSessionServer(s.server, server.NewSessionServer(s.manager))

	for _, primitiveType := range s.PrimitiveTypes {
		primitiveType.RegisterServices(s.server, newProtocol(s.manager))
	}
	go func() {
		if err := s.server.Serve(lis); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	return nil
}

func (s *MultiRaftNode) Stop() error {
	s.server.Stop()
	return nil
}
