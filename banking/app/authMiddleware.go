package app

import (
	"banking/domain"
	"bankinglib/errs"
	"fmt"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	repo domain.AuthRepository
}

func (a AuthMiddleware) authorizationHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentRoute := r.URL.Path
			currentRouteVars := r.URL.Query()

			routeVars := make(map[string]string)
			for key, value := range currentRouteVars {
				routeVars[key] = value[0]
			}

			authHeader := r.Header.Get("Authorization")

			if authHeader != "" {
				token := getTokenFromHeader(authHeader)

				isAuthorized := a.repo.IsAuthorized(token, currentRoute, routeVars)

				if isAuthorized {
					next.ServeHTTP(w, r)
				} else {
					appError := errs.AppError{Code: http.StatusForbidden, Message: "Unauthorized"}
					w.WriteHeader(appError.Code)
					fmt.Fprint(w, appError.AsMessage())

				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("missing token"))
			}
		})
	}
}

func getTokenFromHeader(header string) string {

	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}
