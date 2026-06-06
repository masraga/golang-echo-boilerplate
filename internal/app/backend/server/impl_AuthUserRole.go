package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/service/auth"
)

func (s *Server) AssignAuthUserRole(ctx echo.Context, userId string) error {
	var reqBody api.AssignAuthUserRoleRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.AssignAuthUserRole(ctx.Request().Context(), auth.AssignAuthUserRoleInput{
		UserId:    userId,
		RoleId:    reqBody.RoleId,
		CreatedBy: reqBody.CreatedBy,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.AssignAuthUserRoleResponse{
		UserId:        output.UserId,
		RoleId:        output.RoleId,
		RoleName:      output.RoleName,
		GrantedCount:  output.GrantedCount,
		UpdatedAtUtc0: output.UpdatedAtUtc0,
		CreatedBy:     output.CreatedBy,
	})
}

func (s *Server) DeleteAuthUserRole(ctx echo.Context, userId string) error {
	var reqBody api.DeleteAuthUserRoleJSONBody
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.DeleteAuthUserRole(ctx.Request().Context(), auth.DeleteAuthUserRoleInput{
		UserId:    userId,
		CreatedBy: reqBody.CreatedBy,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.DeleteAuthAccessResponse{IsSuccess: output.IsSuccess})
}
