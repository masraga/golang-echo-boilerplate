package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func (s *Server) CreateAuthUserApiContract(ctx echo.Context) error {
	var reqBody api.CreateAuthUserApiContractRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.CreateAuthUserApiContract(ctx.Request().Context(), auth.CreateAuthUserApiContractInput{
		UserId:        reqBody.UserId,
		ApiContractId: reqBody.ApiContractId,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnCreated(ctx, authUserApiContractToResponse(output))
}

func (s *Server) GetAuthUserApiContract(ctx echo.Context, id string) error {
	output, err := s.AuthService.GetAuthUserApiContract(ctx.Request().Context(), auth.GetAuthUserApiContractInput{
		Id: id,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, authUserApiContractToResponse(output))
}

func (s *Server) ListAuthUserApiContracts(ctx echo.Context) error {
	output, err := s.AuthService.ListAuthUserApiContracts(ctx.Request().Context(), auth.ListAuthUserApiContractsInput{})
	if err != nil {
		return returnError(ctx, err)
	}

	data := make([]api.AuthUserApiContract, 0, len(output.Data))
	for _, item := range output.Data {
		data = append(data, authUserApiContractToResponse(item))
	}
	return returnOk(ctx, api.ListAuthUserApiContractsResponse{Data: data})
}

func (s *Server) UpdateAuthUserApiContract(ctx echo.Context, id string) error {
	var reqBody api.UpdateAuthUserApiContractRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.UpdateAuthUserApiContract(ctx.Request().Context(), auth.UpdateAuthUserApiContractInput{
		Id:            id,
		UserId:        reqBody.UserId,
		ApiContractId: reqBody.ApiContractId,
		IsActive:      reqBody.IsActive,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, authUserApiContractToResponse(output))
}

func (s *Server) DeleteAuthUserApiContract(ctx echo.Context, id string) error {
	output, err := s.AuthService.DeleteAuthUserApiContract(ctx.Request().Context(), auth.DeleteAuthUserApiContractInput{
		Id: id,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.DeleteAuthAccessResponse{IsSuccess: output.IsSuccess})
}

func authUserApiContractToResponse(input auth.AuthUserApiContract) api.AuthUserApiContract {
	return api.AuthUserApiContract{
		Id:            input.Id,
		UserId:        input.UserId,
		ApiContractId: input.ApiContractId,
		CreatedAtUtc0: input.CreatedAtUtc0,
		UpdatedAtUtc0: input.UpdatedAtUtc0,
		IsActive:      input.IsActive,
	}
}
