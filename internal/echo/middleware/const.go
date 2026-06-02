package middleware

const (
	AuthValidationFilterMethodPost = "POST"
	AuthValidationFilterMethodGet  = "GET"
)

var (
	skipAuthValidationFilterMap = map[string]AuthValidationFilterMethod{
		// "/api/v1/ping": {
		// 	Method: []string{
		// 		AuthValidationFilterMethodPost,
		// 		AuthValidationFilterMethodGet,
		// 	},
		// },
		"/api/v1/auth/register/phone": {
			Method: []string{
				AuthValidationFilterMethodPost,
			},
		},
		"/api/v1/auth/otp/verify": {
			Method: []string{
				AuthValidationFilterMethodPost,
			},
		},
		"/api/v1/auth/validate/pin": {
			Method: []string{
				AuthValidationFilterMethodPost,
			},
		},
		"/api/v1/crypto/encrypt": {
			Method: []string{
				AuthValidationFilterMethodPost,
			},
		},
	}
)
