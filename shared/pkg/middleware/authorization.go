package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

type Authorizer interface {
	Allow(subjectID uuid.UUID, tenantID string, permission string) (bool, error)
}

func Authorization(authz Authorizer, permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(jwtutil.Claims)
		if !ok {
			httputil.Error(w, http.StatusUnauthorized, "2001", "Unauthorized", nil)
			return
		}
		okay, err := authz.Allow(claims.UserID, claims.TenantID, permission)
		if err != nil {
			httputil.Error(w, http.StatusInternalServerError, "1001", "Internal error", nil)
			return
		}
		if !okay {
			httputil.Error(w, http.StatusForbidden, "3001", "Forbidden", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
