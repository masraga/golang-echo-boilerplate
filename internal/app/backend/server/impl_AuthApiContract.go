package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/service/auth"
)

func (s *Server) CreateAuthApiContract(ctx echo.Context) error {
	var reqBody api.CreateAuthApiContractRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.CreateAuthApiContract(ctx.Request().Context(), auth.CreateAuthApiContractInput{
		Id:             reqBody.Id,
		EndpointPath:   reqBody.EndpointPath,
		EndpointMethod: reqBody.EndpointMethod,
		Description:    reqBody.Description,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnCreated(ctx, authApiContractToResponse(output))
}

func (s *Server) GetAuthApiContract(ctx echo.Context, id string) error {
	output, err := s.AuthService.GetAuthApiContract(ctx.Request().Context(), auth.GetAuthApiContractInput{
		Id: id,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, authApiContractToResponse(output))
}

func (s *Server) ListAuthApiContracts(ctx echo.Context) error {
	output, err := s.AuthService.ListAuthApiContracts(ctx.Request().Context(), auth.ListAuthApiContractsInput{})
	if err != nil {
		return returnError(ctx, err)
	}

	data := make([]api.AuthApiContract, 0, len(output.Data))
	for _, item := range output.Data {
		data = append(data, authApiContractToResponse(item))
	}
	return returnOk(ctx, api.ListAuthApiContractsResponse{Data: data})
}

func (s *Server) UpdateAuthApiContract(ctx echo.Context, id string) error {
	var reqBody api.UpdateAuthApiContractRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.UpdateAuthApiContract(ctx.Request().Context(), auth.UpdateAuthApiContractInput{
		Id:             id,
		EndpointPath:   reqBody.EndpointPath,
		EndpointMethod: reqBody.EndpointMethod,
		Description:    reqBody.Description,
		IsActive:       reqBody.IsActive,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, authApiContractToResponse(output))
}

func (s *Server) DeleteAuthApiContract(ctx echo.Context, id string) error {
	output, err := s.AuthService.DeleteAuthApiContract(ctx.Request().Context(), auth.DeleteAuthApiContractInput{
		Id: id,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.DeleteAuthAccessResponse{IsSuccess: output.IsSuccess})
}

func authApiContractToResponse(input auth.AuthApiContract) api.AuthApiContract {
	return api.AuthApiContract{
		Id:             input.Id,
		EndpointPath:   input.EndpointPath,
		EndpointMethod: input.EndpointMethod,
		Description:    input.Description,
		CreatedAtUtc0:  input.CreatedAtUtc0,
		UpdatedAtUtc0:  input.UpdatedAtUtc0,
		IsActive:       input.IsActive,
	}
}
