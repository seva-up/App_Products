package routesAuth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/auth/deliveryAuth/httpAuth"
)

func NewRouter(authService auth.UserService) *mux.Router {
	router := mux.NewRouter()

	authHandler := httpAuth.NewAuthDelivery(authService)

	public := router.PathPrefix("/api/v1").Subrouter()
	{
		public.HandleFunc("/register", authHandler.Register).Methods("POST")
		public.HandleFunc("/health", healthCheck).Methods("GET")
	}

	return router
}
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
