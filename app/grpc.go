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
	"strings"
)

type UpdaterServer struct {
	pb.UnimplementedUpdaterServer
	pb.UnimplementedRollbackServer

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
	pb.RegisterRollbackServer(grpcServer, s)

	fmt.Println("server started")
	return grpcServer.Serve(lis)
}

func (s *UpdaterServer) UpdateContainer(request *pb.UpdateRequest, stream pb.Updater_UpdateContainerServer) error {
	container, ok := s.controller.FindContainerByName(request.ContainerName)
	if !ok {
		s.sendMessage(stream, fmt.Sprintf("container with the name %s does not exist", request.ContainerName))
		return nil
	}

	imageParts := strings.Split(container.Image, ":")
	fmt.Println(imageParts)
	container.Image = imageParts[0] + ":" + request.Image
	fmt.Println(container.Image, request.Image)

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

	s.sendMessage(stream, fmt.Sprintf("updating %s to %s", request.ContainerName, container.Image))
	err = s.controller.UpdateContainer(request.ContainerName, container.Image, false)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrImageFormatInvalid):
			s.sendMessage(stream, "image format is invalid")
		case errors.Is(err, controller.ErrContainerNotFound):
			s.sendMessage(stream, "the requested container could not be found")
		default:
			s.sendMessage(stream, err.Error())
		}
		return nil
	}

	s.sendMessage(stream, "container updated successfully")
	return nil
}

func (s *UpdaterServer) RollbackContainer(request *pb.RollbackRequest, stream pb.Rollback_RollbackContainerServer) error {
	fmt.Println("rollback container", request.GetContainer())

	if request.GetContainer() == "" {
		stream.Send(&pb.Response{Message: "container value must not be empty"})
		return nil
	}

	err := s.controller.RollbackContainer(request.GetContainer())
	if err != nil {
		var containerStartFailedErr controller.ErrContainerStartFailed

		switch {
		case errors.Is(err, controller.ErrContainerNotFound):
			stream.Send(&pb.Response{Message: fmt.Sprintf("container with the name %s does not exist", request.GetContainer())})
		case errors.Is(err, controller.ErrRollbackContainerNotFound):
			stream.Send(&pb.Response{Message: fmt.Sprintf("%s does not have a rollback container (%s-rollback)", request.GetContainer(), request.GetContainer())})
		case errors.Is(err, controller.ErrContainerNotRunning):
			stream.Send(&pb.Response{Message: "the container failed to start"})
		case errors.As(err, &containerStartFailedErr):
			stream.Send(&pb.Response{Message: containerStartFailedErr.Reason.Error()})
		default:
			stream.Send(&pb.Response{Message: err.Error()})
		}
		return nil
	}

	stream.Send(&pb.Response{Message: "successfully rolled back container"})
	return nil
}

func (s *UpdaterServer) sendMessage(stream pb.Updater_UpdateContainerServer, message string) {

	stream.Send(&pb.Response{Message: message})
}