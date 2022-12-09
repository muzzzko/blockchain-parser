package setup

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"blockchain-parser/internal/infrastructure/handler"
)

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("content-type", "application/json")

		next.ServeHTTP(w, req)
	})
}

func panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("stacktrace from panic: \n" + string(debug.Stack()))

				resp := handler.ErrorResponse{
					Message: "internal server error",
				}

				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(resp)
			}
		}()

		next.ServeHTTP(w, req)
	})
}

func methodCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, ok := routes[req.URL.Path]; ok {
			if _, ok := routes[req.URL.Path][req.Method]; !ok {
				w.WriteHeader(http.StatusMethodNotAllowed)

				return
			}
		}

		next.ServeHTTP(w, req)
	})
}
