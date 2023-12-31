package v1

import (
	pb_users_model "github.com/almalii/grpc-contracts/gen/go/users_service/model/v1"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/users_service/models"
	"notes-rew/internal/users_service/usecase"
)

func NewGetUserResponse(resp models.UserOutput) *pb_users_model.GetUserResponse {
	return &pb_users_model.GetUserResponse{
		Id:        resp.ID.String(),
		Name:      resp.Username,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt.String(),
		UpdatedAt: resp.UpdatedAt.String(),
	}
}

func NewUpdateUserInput(req *pb_users_model.UpdateUserRequest) usecase.UpdateUserInput {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.Error("error parsing userID: ", err)
		return usecase.UpdateUserInput{}
	}
	return usecase.UpdateUserInput{
		InitiatorID: userID,
		Username:    req.Name,
		Email:       req.Email,
		Password:    req.Password,
	}
}

func NewUpdateUserResponse(resp usecase.UpdateUserInput) *pb_users_model.UpdateUserResponse {
	return &pb_users_model.UpdateUserResponse{
		Id:    resp.InitiatorID.String(),
		Name:  *resp.Username,
		Email: *resp.Email,
	}
}
