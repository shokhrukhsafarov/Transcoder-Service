package service

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/transcodeuz/transcode-rest/config"
	pb "gitlab.com/transcodeuz/transcode-rest/genproto/transcoder_service"
	"gitlab.com/transcodeuz/transcode-rest/models"
	"gitlab.com/transcodeuz/transcode-rest/pkg/etc"
	"gitlab.com/transcodeuz/transcode-rest/pkg/logger"
	"gitlab.com/transcodeuz/transcode-rest/pkg/rabbitmq"
	"gitlab.com/transcodeuz/transcode-rest/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

type pipelineServide struct {
	strg   storage.StorageI
	cfg    config.Config
	log    logger.Logger
	rabbit *rabbitmq.RabbitMQ
	pb.UnimplementedPipelineServiceServer
}

func NewPipelineService(cfg config.Config, log logger.Logger, strg storage.StorageI, rabbit *rabbitmq.RabbitMQ) *pipelineServide {
	return &pipelineServide{
		strg:   strg,
		cfg:    cfg,
		log:    log,
		rabbit: rabbit,
	}
}

func (s *pipelineServide) Create(ctx context.Context, req *pb.CreatePipelineRequest) (resp *emptypb.Empty, err error) {
	sizeKbVideo, err := etc.GetFileSizeFromUrl(req.InputUrl)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	project, err := s.strg.Postgres().ProjectGet(ctx, &models.ProjectGetReq{ID: req.ProjectId})
	if err != nil {
		return &emptypb.Empty{}, err
	}

	pipelineRes, err := s.strg.Postgres().PipelineCreate(ctx, &models.PipelineCreateReq{
		ID:                uuid.NewString(),
		ProjectID:         project.ID,
		Stage:             "initial",
		InputURL:          req.InputUrl,
		OutputKey:         req.OutputKey,
		OutputPath:        req.BucketName,
		MaxResolution:     "1080p",
		SizeKB:            sizeKbVideo,
		ResolutionsString: req.Resolutions,
		KeyID:             req.KeyId,
		FieldSlug:         req.FieldSlug,
		TableSlug:         req.TableSlug,
	})
	if err != nil {
		return &emptypb.Empty{}, err
	}

	storage, err := s.strg.Postgres().StorageGet(ctx, &models.StorageGetReq{ID: project.StorageID})
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = s.rabbit.PublishPipeline(&models.PipelineRabbitMq{
		Id:           pipelineRes.ID,
		InputURI:     pipelineRes.InputURL,
		OutputKey:    pipelineRes.OutputKey,
		CdnUrl:       storage.DomainName,
		CdnAccessKey: storage.AccessKey,
		CdnSecretKey: storage.SecretKey,
		CdnRegion:    storage.Region,
		CdnBucket:    pipelineRes.OutputPath,
		CdnType:      "minio",
		Resolutions:  req.Resolutions,
		Language:     req.Language,
		LanguageCode: req.LanguageCode,
		KeyID:        req.KeyId,
	})
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return nil, nil
}

func (s *pipelineServide) GetList(ctx context.Context, req *pb.GetListPipelineRequest) (resp *pb.GetListPipelineResponse, err error) {
	s.log.Info("!!!GetListPipeline called->", req)

	pipelines, err := s.strg.Postgres().PipelinesFind(ctx, req)
	if err != nil {
		return nil, err
	}

	resp, err = etc.ConvertStructToStruct[*models.PipelinesFindResponse, *pb.GetListPipelineResponse](pipelines)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
