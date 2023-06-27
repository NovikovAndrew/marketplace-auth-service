package auth

import "net/http"

func (ah *AuthenticationHandler) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}
