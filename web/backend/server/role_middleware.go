package server

import (
	"net/http"
)

type ServerMode string

const (
	ServerModeAll = ServerMode("all")
	ServerModeGet = ServerMode("get")
)

func SetupRoleMiddleware(mode ServerMode) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch mode {
			case ServerModeAll:
				next.ServeHTTP(w, r)
			case ServerModeGet:
				if r.Method == http.MethodGet {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			default:
				http.Error(w, "Invalid mode", http.StatusBadRequest)
			}
		})
	}
}
