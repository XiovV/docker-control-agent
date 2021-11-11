package app

import (
	"errors"
	"fmt"
	"github.com/XiovV/dokkup-agent/config"
	"github.com/XiovV/dokkup-agent/controller"
	pb "github.com/XiovV/dokkup-agent/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

type UpdaterServer struct {
	pb.UnimplementedUpdaterServer

	controller controller.ContainerController
	config     *config.Config
}

func NewUpdaterServer(controller controller.ContainerController, config *config.Config) *UpdaterServer {
	return &UpdaterServer{controller: controller, config: config}
}

func (s *UpdaterServer) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterUpdaterServer(grpcServer, s)

	fmt.Println("server started")
	return grpcServer.Serve(lis)
}

func (s *UpdaterServer) UpdateContainer(request *pb.UpdateRequest, stream pb.Updater_UpdateContainerServer) error {

	container, ok := s.controller.FindContainerByName(request.ContainerName)
	if !ok {
		s.sendMessage(stream, "this container does not exist")
		return nil
	}

	s.sendMessage(stream, "container found")

	s.sendMessage(stream, fmt.Sprintf("pulling image %s", container.Image))

	alreadyExists, err := s.controller.PullImage(container.Image)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrImageFormatInvalid):
			s.sendMessage(stream,"image format is invalid")
		default:
			s.sendMessage(stream,err.Error())
		}
		return err
	}

	if alreadyExists {
		s.sendMessage(stream, "image already exists...")
	}


	return nil
}

func (s *UpdaterServer) sendMessage(stream pb.Updater_UpdateContainerServer, message string) {
	stream.Send(&pb.UpdateResponse{Message: message})
}