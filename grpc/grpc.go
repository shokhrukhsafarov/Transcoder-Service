package grpc

import (
	"gitlab.com/transcodeuz/transcode-rest/config"
	"gitlab.com/transcodeuz/transcode-rest/genproto/transcoder_service"
	"gitlab.com/transcodeuz/transcode-rest/grpc/service"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/pkg/rabbitmq"
	"gitlab.com/transcodeuz/transcode-rest/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.Logger, strg storage.StorageI, rabbit *rabbitmq.RabbitMQ) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	transcoder_service.RegisterPipelineServiceServer(grpcServer, service.NewPipelineService(cfg, log, strg, rabbit))

	reflection.Register(grpcServer)
	return
}
