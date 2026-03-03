package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	pb "gitlab.com/transcodeuz/transcode-rest/genproto/transcoder_service"
	"gitlab.com/transcodeuz/transcode-rest/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *pipelineServide) CreateCompanyAndProject(ctx context.Context, req *pb.CreateCompanyAndProjectRequest) (resp *emptypb.Empty, err error) {
	users, err := s.strg.Postgres().UserFind(ctx, &models.UserFindReq{
		Limit:            1,
		OrderByCreatedAt: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(users.Users) == 0 {
		return nil, errors.New("user is not found")
	}

	user := users.Users[0]

	company, err := s.strg.Postgres().CompanyCreate(ctx, &models.CompanyCreateReq{
		ID:      uuid.NewString(),
		Title:   req.Title,
		OwnerID: user.ID,
		Status:  "active",
	})
	if err != nil {
		return nil, err
	}

	storages, err := s.strg.Postgres().StorageFind(ctx, &models.StorageFindReq{
		Limit:            1,
		OrderByCreatedAt: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(storages.Storages) == 0 {
		return nil, errors.New("storage is not found")
	}

	storage := storages.Storages[0]

	_, err = s.strg.Postgres().ProjectCreate(ctx, &models.ProjectCreateReq{
		ID:        req.ProjectId,
		Title:     req.Title,
		AccessKey: user.Username,
		SecretKey: user.Username,
		CompanyID: company.ID,
		OwnerID:   user.ID,
		Status:    "active",
		StorageID: storage.ID,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
