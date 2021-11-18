package app

import (
	pb "github.com/XiovV/dokkup-agent/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) authenticate(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	header, _ := metadata.FromIncomingContext(stream.Context())

	apiKey := header["authorization"][0]

	if !s.config.CompareHash(apiKey) {
		stream.SendMsg(&pb.Response{Message: "invalid api key"})
		return nil
	}

	return handler(srv, stream)
}