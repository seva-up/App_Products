package auth

import "net/http"

type UserDelivery interface {
	Register(w http.ResponseWriter, r *http.Request)
}
