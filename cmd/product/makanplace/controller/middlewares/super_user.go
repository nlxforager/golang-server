package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"slices"

	"golang-server/cmd/product/makanplace/controller/response_types"
	"golang-server/cmd/product/makanplace/service/mk_user_session"
)

func SuperUserMiddleware(mkService *mk_user_session.Service) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId := GetSessionIdFromRequest(r)
			session := mkService.GetSession(sessionId, false)
			log.Printf("checking IsSuperUser: %#v\n", session)
			// no gmails supplied or any of the gmails is not SU.
			if len(session.Gmails) == 0 || slices.ContainsFunc(session.Gmails, func(gmail string) bool {
				return !mkService.IsSuperUser(gmail)
			}) {
				response_types.ErrorNoBody(w, http.StatusUnauthorized, fmt.Errorf("not permitted to add outlet"))
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
