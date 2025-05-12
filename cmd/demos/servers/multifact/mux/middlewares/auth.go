package middlewares

import (
	"net/http"

	"golang-server/src/domain/auth"
)

func BearerAuthMiddleware(jwtService auth.JwtService) MiddleWare {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			authv := request.Header["Authorization"]

			isAuth := false
			if len(authv) == 1 {
				if authv[0][0:len("Bearer ")] == "Bearer " && len(authv[0]) > len("Bearer ") {
					token := authv[0][len("Bearer "):]
					_, err := jwtService.ValidateAndGetClaims(token)
					if err == nil {
						isAuth = true
					}
				}
			}

			if isAuth {
				next(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
			}
		}
	}
}
