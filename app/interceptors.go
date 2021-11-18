package app

import (
	pb "github.com/XiovV/dokkup-agent/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *UpdaterServer) authenticate(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	header, _ := metadata.FromIncomingContext(ss.Context())

	apiKey := header["authorization"][0]

	if !s.config.CompareHash(apiKey) {
		ss.SendMsg(&pb.Response{Message: "invalid api key"})
		return nil
	}

	return handler(srv, ss)
}