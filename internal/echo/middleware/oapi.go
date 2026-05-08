package middleware

import (
	"fmt"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/masraga/kerp-api/generated/api"
	echomiddleware "github.com/oapi-codegen/echo-middleware"
)

func OapiGetSwagger() *openapi3.T {
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	swagger.Servers = nil
	return swagger
}

func OapiValidatorOpt() *echomiddleware.Options {
	return &echomiddleware.Options{
		Options: openapi3filter.Options{
			SchemaValidationOptions: []openapi3.SchemaValidationOption{
				openapi3.WithStringFormatValidator("email", openapi3.NewRegexpFormatValidator(openapi3.FormatOfStringForEmail)),
			},
		},
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			return c.JSON(err.Code, map[string]string{
				"message": oapiSimplifyError(err.Error()),
			})
		},
	}
}

func oapiSimplifyError(msg string) string {

	requiredRe := regexp.MustCompile(`property "([^"]+)" is missing`)
	requiredMatch := requiredRe.FindStringSubmatch(msg)

	if len(requiredMatch) > 1 {
		return fmt.Sprintf("%s is required", requiredMatch[1])
	}

	emailRe := regexp.MustCompile(`Error at "/([^"]+)": .*format "email"`)
	emailMatch := emailRe.FindStringSubmatch(msg)

	if len(emailMatch) > 1 {
		return fmt.Sprintf("%s format is invalid", emailMatch[1])
	}

	return "invalid request"
}
