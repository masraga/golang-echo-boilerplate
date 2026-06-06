package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func (s *Server) CreateAuthRole(ctx echo.Context) error {
	var reqBody api.CreateAuthRoleRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.CreateAuthRole(ctx.Request().Context(), auth.CreateAuthRoleInput{
		RoleName:    reqBody.RoleName,
		Description: reqBody.Description,
		OwnerId:     reqBody.OwnerId,
		CreatedBy:   reqBody.CreatedBy,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnCreated(ctx, authRoleToResponse(output))
}

func (s *Server) GetAuthRole(ctx echo.Context, id string) error {
	output, err := s.AuthService.GetAuthRole(ctx.Request().Context(), auth.GetAuthRoleInput{
		Id: id,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, authRoleToResponse(output))
}

func (s *Server) ListAuthRoles(ctx echo.Context) error {
	output, err := s.AuthService.ListAuthRoles(ctx.Request().Context(), auth.ListAuthRolesInput{})
	if err != nil {
		return returnError(ctx, err)
	}

	data := make([]api.AuthRole, 0, len(output.Data))
	for _, item := range output.Data {
		data = append(data, authRoleToResponse(item))
	}
	return returnOk(ctx, api.ListAuthRolesResponse{Data: data})
}

func (s *Server) UpdateAuthRole(ctx echo.Context, id string) error {
	var reqBody api.UpdateAuthRoleRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.UpdateAuthRole(ctx.Request().Context(), auth.UpdateAuthRoleInput{
		Id:          id,
		RoleName:    reqBody.RoleName,
		Description: reqBody.Description,
		OwnerId:     reqBody.OwnerId,
		CreatedBy:   reqBody.CreatedBy,
		IsActive:    reqBody.IsActive,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, authRoleToResponse(output))
}

func (s *Server) DeleteAuthRole(ctx echo.Context, id string) error {
	output, err := s.AuthService.DeleteAuthRole(ctx.Request().Context(), auth.DeleteAuthRoleInput{
		Id: id,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.DeleteAuthAccessResponse{IsSuccess: output.IsSuccess})
}

func authRoleToResponse(input auth.AuthRole) api.AuthRole {
	return api.AuthRole{
		Id:            input.Id,
		RoleName:      input.RoleName,
		Description:   input.Description,
		OwnerId:       input.OwnerId,
		CreatedAtUtc0: input.CreatedAtUtc0,
		UpdatedAtUtc0: input.UpdatedAtUtc0,
		CreatedBy:     input.CreatedBy,
		IsActive:      input.IsActive,
	}
}
