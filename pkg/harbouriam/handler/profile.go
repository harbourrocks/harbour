package handler

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbouriam/helper"
	"net/http"
)

// RefreshProfile extracts the latest user information from an id token
func RefreshProfile(w http.ResponseWriter, r *http.Request) {
	idToken := auth.GetIdTokenReq(r)

	_, err := helper.RefreshProfile(r.Context(), idToken)
	if err != nil {
		// error logged in RefreshProfile
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
