package server

import (
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	"github.com/masraga/kerp-api/internal/service/auth"
)

func (s *Server) CreateAuthRoleContractApi(ctx echo.Context, roleId string) error {
	var reqBody api.CreateAuthRoleContractApiRequest
	if err := bindOrReturnBadRequest(ctx, &reqBody); err != nil {
		return err
	}

	output, err := s.AuthService.CreateAuthRoleContractApi(ctx.Request().Context(), auth.CreateAuthRoleContractApiInput{
		RoleId:            roleId,
		AuthApiContractId: reqBody.AuthApiContractId,
		CreatedBy:         reqBody.CreatedBy,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnCreated(ctx, authRoleContractApiToResponse(output))
}

func (s *Server) ListAuthRoleContractApis(ctx echo.Context, roleId string) error {
	output, err := s.AuthService.ListAuthRoleContractApis(ctx.Request().Context(), auth.ListAuthRoleContractApisInput{
		RoleId: roleId,
	})
	if err != nil {
		return returnError(ctx, err)
	}

	data := make([]api.AuthRoleContractApi, 0, len(output.Data))
	for _, item := range output.Data {
		data = append(data, authRoleContractApiToResponse(item))
	}
	return returnOk(ctx, api.ListAuthRoleContractApisResponse{Data: data})
}

func (s *Server) DeleteAuthRoleContractApi(ctx echo.Context, roleId string, id string) error {
	output, err := s.AuthService.DeleteAuthRoleContractApi(ctx.Request().Context(), auth.DeleteAuthRoleContractApiInput{
		Id:     id,
		RoleId: roleId,
	})
	if err != nil {
		return returnError(ctx, err)
	}
	return returnOk(ctx, api.DeleteAuthAccessResponse{IsSuccess: output.IsSuccess})
}

func authRoleContractApiToResponse(input auth.AuthRoleContractApi) api.AuthRoleContractApi {
	return api.AuthRoleContractApi{
		Id:                input.Id,
		RoleId:            input.RoleId,
		AuthApiContractId: input.AuthApiContractId,
		CreatedAtUtc0:     input.CreatedAtUtc0,
		UpdatedAtUtc0:     input.UpdatedAtUtc0,
		CreatedBy:         input.CreatedBy,
		IsActive:          input.IsActive,
	}
}
