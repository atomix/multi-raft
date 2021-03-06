// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"context"
	api "github.com/atomix/multi-raft-storage/api/atomix/multiraft/map/v1"
	multiraftv1 "github.com/atomix/multi-raft-storage/api/atomix/multiraft/v1"
	"github.com/atomix/multi-raft-storage/driver/pkg/client"
	"github.com/atomix/multi-raft-storage/driver/pkg/util/async"
	mapv1 "github.com/atomix/runtime/api/atomix/runtime/map/v1"
	runtimev1 "github.com/atomix/runtime/api/atomix/runtime/v1"
	"github.com/atomix/runtime/sdk/pkg/runtime"
	"google.golang.org/grpc"
	"io"
)

const Type = "Map"
const APIVersion = "v1"

func NewServer(protocol *client.Protocol) mapv1.MapServer {
	return &Server{
		Protocol: protocol,
	}
}

type Server struct {
	*client.Protocol
}

func (s *Server) Create(ctx context.Context, request *mapv1.CreateRequest) (*mapv1.CreateResponse, error) {
	partitions := s.Partitions()
	err := async.IterAsync(len(partitions), func(i int) error {
		partition := partitions[i]
		session, err := partition.GetSession(ctx)
		if err != nil {
			return err
		}
		spec := multiraftv1.PrimitiveSpec{
			Type: multiraftv1.PrimitiveType{
				Name:       Type,
				ApiVersion: APIVersion,
			},
			Namespace: runtime.GetNamespace(),
			Name:      request.ID.Name,
		}
		return session.CreatePrimitive(ctx, spec)
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.CreateResponse{}
	return response, nil
}

func (s *Server) Close(ctx context.Context, request *mapv1.CloseRequest) (*mapv1.CloseResponse, error) {
	partitions := s.Partitions()
	err := async.IterAsync(len(partitions), func(i int) error {
		partition := partitions[i]
		session, err := partition.GetSession(ctx)
		if err != nil {
			return err
		}
		return session.ClosePrimitive(ctx, request.ID.Name)
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.CloseResponse{}
	return response, nil
}

func (s *Server) Size(ctx context.Context, request *mapv1.SizeRequest) (*mapv1.SizeResponse, error) {
	partitions := s.Partitions()
	sizes, err := async.ExecuteAsync[int](len(partitions), func(i int) (int, error) {
		partition := partitions[i]
		session, err := partition.GetSession(ctx)
		if err != nil {
			return 0, err
		}
		primitive, err := session.GetPrimitive(request.ID.Name)
		if err != nil {
			return 0, err
		}
		query := client.Query[*api.SizeResponse](primitive, "Size")
		output, err := query.Run(func(conn *grpc.ClientConn, headers *multiraftv1.QueryRequestHeaders) (*api.SizeResponse, error) {
			return api.NewMapClient(conn).Size(ctx, &api.SizeRequest{
				Headers:   *headers,
				SizeInput: api.SizeInput{},
			})
		})
		if err != nil {
			return 0, err
		}
		return int(output.Size_), nil
	})
	if err != nil {
		return nil, err
	}
	var size int
	for _, s := range sizes {
		size += s
	}
	response := &mapv1.SizeResponse{
		Size_: uint32(size),
	}
	return response, nil
}

func (s *Server) Put(ctx context.Context, request *mapv1.PutRequest) (*mapv1.PutResponse, error) {
	partition := s.PartitionBy([]byte(request.Key))
	session, err := partition.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	primitive, err := session.GetPrimitive(request.ID.Name)
	if err != nil {
		return nil, err
	}
	command := client.Command[*api.PutResponse](primitive, "Put")
	output, err := command.Run(func(conn *grpc.ClientConn, headers *multiraftv1.CommandRequestHeaders) (*api.PutResponse, error) {
		return api.NewMapClient(conn).Put(ctx, &api.PutRequest{
			Headers: *headers,
			PutInput: api.PutInput{
				Key: request.Key,
				Value: &api.Value{
					Value: request.Value.Value,
					TTL:   request.Value.TTL,
				},
				Index: multiraftv1.Index(request.IfTimestamp.GetLogicalTimestamp().Time),
			},
		})
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.PutResponse{
		Entry: mapv1.Entry{
			Key: output.Entry.Key,
			Value: &mapv1.Value{
				Value: output.Entry.Value.Value,
				TTL:   output.Entry.Value.TTL,
			},
			Timestamp: &runtimev1.Timestamp{
				Timestamp: &runtimev1.Timestamp_LogicalTimestamp{
					LogicalTimestamp: &runtimev1.LogicalTimestamp{
						Time: runtimev1.LogicalTime(output.Entry.Index),
					},
				},
			},
		},
	}
	return response, nil
}

func (s *Server) Insert(ctx context.Context, request *mapv1.InsertRequest) (*mapv1.InsertResponse, error) {
	partition := s.PartitionBy([]byte(request.Key))
	session, err := partition.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	primitive, err := session.GetPrimitive(request.ID.Name)
	if err != nil {
		return nil, err
	}
	command := client.Command[*api.InsertResponse](primitive, "Insert")
	output, err := command.Run(func(conn *grpc.ClientConn, headers *multiraftv1.CommandRequestHeaders) (*api.InsertResponse, error) {
		return api.NewMapClient(conn).Insert(ctx, &api.InsertRequest{
			Headers: *headers,
			InsertInput: api.InsertInput{
				Key: request.Key,
				Value: &api.Value{
					Value: request.Value.Value,
					TTL:   request.Value.TTL,
				},
			},
		})
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.InsertResponse{
		Entry: mapv1.Entry{
			Key: output.Entry.Key,
			Value: &mapv1.Value{
				Value: output.Entry.Value.Value,
				TTL:   output.Entry.Value.TTL,
			},
			Timestamp: &runtimev1.Timestamp{
				Timestamp: &runtimev1.Timestamp_LogicalTimestamp{
					LogicalTimestamp: &runtimev1.LogicalTimestamp{
						Time: runtimev1.LogicalTime(output.Entry.Index),
					},
				},
			},
		},
	}
	return response, nil
}

func (s *Server) Update(ctx context.Context, request *mapv1.UpdateRequest) (*mapv1.UpdateResponse, error) {
	partition := s.PartitionBy([]byte(request.Key))
	session, err := partition.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	primitive, err := session.GetPrimitive(request.ID.Name)
	if err != nil {
		return nil, err
	}
	command := client.Command[*api.UpdateResponse](primitive, "Update")
	output, err := command.Run(func(conn *grpc.ClientConn, headers *multiraftv1.CommandRequestHeaders) (*api.UpdateResponse, error) {
		return api.NewMapClient(conn).Update(ctx, &api.UpdateRequest{
			Headers: *headers,
			UpdateInput: api.UpdateInput{
				Key: request.Key,
				Value: &api.Value{
					Value: request.Value.Value,
					TTL:   request.Value.TTL,
				},
				Index: multiraftv1.Index(request.IfTimestamp.GetLogicalTimestamp().Time),
			},
		})
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.UpdateResponse{
		Entry: mapv1.Entry{
			Key: output.Entry.Key,
			Value: &mapv1.Value{
				Value: output.Entry.Value.Value,
				TTL:   output.Entry.Value.TTL,
			},
			Timestamp: &runtimev1.Timestamp{
				Timestamp: &runtimev1.Timestamp_LogicalTimestamp{
					LogicalTimestamp: &runtimev1.LogicalTimestamp{
						Time: runtimev1.LogicalTime(output.Entry.Index),
					},
				},
			},
		},
	}
	return response, nil
}

func (s *Server) Get(ctx context.Context, request *mapv1.GetRequest) (*mapv1.GetResponse, error) {
	partition := s.PartitionBy([]byte(request.Key))
	session, err := partition.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	primitive, err := session.GetPrimitive(request.ID.Name)
	if err != nil {
		return nil, err
	}
	query := client.Query[*api.GetResponse](primitive, "Get")
	output, err := query.Run(func(conn *grpc.ClientConn, headers *multiraftv1.QueryRequestHeaders) (*api.GetResponse, error) {
		return api.NewMapClient(conn).Get(ctx, &api.GetRequest{
			Headers: *headers,
			GetInput: api.GetInput{
				Key: request.Key,
			},
		})
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.GetResponse{
		Entry: mapv1.Entry{
			Key: output.Entry.Key,
			Value: &mapv1.Value{
				Value: output.Entry.Value.Value,
				TTL:   output.Entry.Value.TTL,
			},
			Timestamp: &runtimev1.Timestamp{
				Timestamp: &runtimev1.Timestamp_LogicalTimestamp{
					LogicalTimestamp: &runtimev1.LogicalTimestamp{
						Time: runtimev1.LogicalTime(output.Entry.Index),
					},
				},
			},
		},
	}
	return response, nil
}

func (s *Server) Remove(ctx context.Context, request *mapv1.RemoveRequest) (*mapv1.RemoveResponse, error) {
	partition := s.PartitionBy([]byte(request.Key))
	session, err := partition.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	primitive, err := session.GetPrimitive(request.ID.Name)
	if err != nil {
		return nil, err
	}
	command := client.Command[*api.RemoveResponse](primitive, "Remove")
	output, err := command.Run(func(conn *grpc.ClientConn, headers *multiraftv1.CommandRequestHeaders) (*api.RemoveResponse, error) {
		return api.NewMapClient(conn).Remove(ctx, &api.RemoveRequest{
			Headers: *headers,
			RemoveInput: api.RemoveInput{
				Key:   request.Key,
				Index: multiraftv1.Index(request.IfTimestamp.GetLogicalTimestamp().Time),
			},
		})
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.RemoveResponse{
		Entry: mapv1.Entry{
			Key: output.Entry.Key,
			Value: &mapv1.Value{
				Value: output.Entry.Value.Value,
				TTL:   output.Entry.Value.TTL,
			},
			Timestamp: &runtimev1.Timestamp{
				Timestamp: &runtimev1.Timestamp_LogicalTimestamp{
					LogicalTimestamp: &runtimev1.LogicalTimestamp{
						Time: runtimev1.LogicalTime(output.Entry.Index),
					},
				},
			},
		},
	}
	return response, nil
}

func (s *Server) Clear(ctx context.Context, request *mapv1.ClearRequest) (*mapv1.ClearResponse, error) {
	partitions := s.Partitions()
	err := async.IterAsync(len(partitions), func(i int) error {
		partition := partitions[i]
		session, err := partition.GetSession(ctx)
		if err != nil {
			return err
		}
		primitive, err := session.GetPrimitive(request.ID.Name)
		if err != nil {
			return err
		}
		command := client.Command[*api.ClearResponse](primitive, "Clear")
		_, err = command.Run(func(conn *grpc.ClientConn, headers *multiraftv1.CommandRequestHeaders) (*api.ClearResponse, error) {
			return api.NewMapClient(conn).Clear(ctx, &api.ClearRequest{
				Headers:    *headers,
				ClearInput: api.ClearInput{},
			})
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	response := &mapv1.ClearResponse{}
	return response, nil
}

func (s *Server) Events(request *mapv1.EventsRequest, server mapv1.Map_EventsServer) error {
	partitions := s.Partitions()
	return async.IterAsync(len(partitions), func(i int) error {
		partition := partitions[i]
		session, err := partition.GetSession(server.Context())
		if err != nil {
			return err
		}
		primitive, err := session.GetPrimitive(request.ID.Name)
		if err != nil {
			return err
		}
		command := client.StreamCommand[api.Map_EventsClient, *api.EventsResponse](primitive, "Events")
		stream, err := command.Open(func(conn *grpc.ClientConn, headers *multiraftv1.CommandRequestHeaders) (api.Map_EventsClient, error) {
			return api.NewMapClient(conn).Events(server.Context(), &api.EventsRequest{
				Headers: *headers,
				EventsInput: api.EventsInput{
					Key:    request.Key,
					Replay: request.Replay,
				},
			})
		})
		if err != nil {
			return err
		}
		for {
			output, err := command.Recv(stream.Recv)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			response := &mapv1.EventsResponse{
				Event: mapv1.Event{
					Type: mapv1.Event_Type(output.Event.Type),
					Entry: mapv1.Entry{
						Key: output.Event.Entry.Key,
						Value: &mapv1.Value{
							Value: output.Event.Entry.Value.Value,
							TTL:   output.Event.Entry.Value.TTL,
						},
						Timestamp: &runtimev1.Timestamp{
							Timestamp: &runtimev1.Timestamp_LogicalTimestamp{
								LogicalTimestamp: &runtimev1.LogicalTimestamp{
									Time: runtimev1.LogicalTime(output.Event.Entry.Index),
								},
							},
						},
					},
				},
			}
			if err := server.Send(response); err != nil {
				return err
			}
		}
	})
}

func (s *Server) Entries(request *mapv1.EntriesRequest, server mapv1.Map_EntriesServer) error {
	partitions := s.Partitions()
	return async.IterAsync(len(partitions), func(i int) error {
		partition := partitions[i]
		session, err := partition.GetSession(server.Context())
		if err != nil {
			return err
		}
		primitive, err := session.GetPrimitive(request.ID.Name)
		if err != nil {
			return err
		}
		query := client.StreamQuery[api.Map_EntriesClient, *api.EntriesResponse](primitive, "Entries")
		stream, err := query.Open(func(conn *grpc.ClientConn, headers *multiraftv1.QueryRequestHeaders) (api.Map_EntriesClient, error) {
			return api.NewMapClient(conn).Entries(server.Context(), &api.EntriesRequest{
				Headers:      *headers,
				EntriesInput: api.EntriesInput{},
			})
		})
		if err != nil {
			return err
		}
		for {
			output, err := query.Recv(stream.Recv)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			response := &mapv1.EntriesResponse{
				Entry: mapv1.Entry{
					Key: output.Entry.Key,
					Value: &mapv1.Value{
						Value: output.Entry.Value.Value,
						TTL:   output.Entry.Value.TTL,
					},
					Timestamp: &runtimev1.Timestamp{
						Timestamp: &runtimev1.Timestamp_LogicalTimestamp{
							LogicalTimestamp: &runtimev1.LogicalTimestamp{
								Time: runtimev1.LogicalTime(output.Entry.Index),
							},
						},
					},
				},
			}
			if err := server.Send(response); err != nil {
				return err
			}
		}
	})
}

var _ mapv1.MapServer = (*Server)(nil)
