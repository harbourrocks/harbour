package auth

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"net/http"
)

const OidcTokenStrKey = "oidcTokenStr"
const OidcTokenKey = "oidcToken"
const IdTokenKey = "idToken"

func GetOidcTokenStrReq(r *http.Request) string {
	return GetOidcTokenStrCtx(r.Context())
}

func GetOidcTokenStrCtx(ctx context.Context) string {
	oidcTokenStr := ctx.Value(OidcTokenStrKey).(string)
	return oidcTokenStr
}

func GetOidcTokenReq(r *http.Request) *oidc.IDToken {
	return GetOidcTokenCtx(r.Context())
}

func GetOidcTokenCtx(ctx context.Context) *oidc.IDToken {
	if oidcToken := ctx.Value(OidcTokenKey); oidcToken != nil {
		token := oidcToken.(*oidc.IDToken)
		return token
	} else {
		return nil
	}
}

func GetIdTokenReq(r *http.Request) *IdToken {
	return GetIdTokenCtx(r.Context())
}

func GetIdTokenCtx(ctx context.Context) *IdToken {
	if idToken := ctx.Value(IdTokenKey); idToken != nil {
		token := idToken.(IdToken)
		return &token
	} else {
		return nil
	}
}

func UseOidcTokenStr(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, OidcTokenStrKey, ExtractToken(r))

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}

func UseOidcToken(next http.HandlerFunc, oidcConfig OIDCConfig) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if idTokenStr := GetOidcTokenStrCtx(ctx); idTokenStr != "" {
			idToken, err := JwtAuth(ctx, idTokenStr, oidcConfig) // error logged in JwtAuth

			if err == nil {
				ctx = context.WithValue(ctx, OidcTokenKey, idToken)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}

func UseIdToken(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := GetOidcTokenCtx(ctx)
		log := logconfig.GetLogCtx(ctx)

		if token != nil {
			var idToken IdToken
			if err := token.Claims(&idToken); err != nil {
				log.WithError(err).Warn("Failed to extract claims")
			} else {
				ctx = context.WithValue(ctx, IdTokenKey, idToken)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}

func UseAuth(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if idToken := GetIdTokenCtx(ctx); idToken == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
