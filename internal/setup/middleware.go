package setup

import (
	"log"
	"net/http"
	"runtime/debug"
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
			}
		}()

		next.ServeHTTP(w, req)
	})
}

func methodCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, ok := routes[req.Method]; !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}
		if _, ok := routes[req.Method][req.URL.Path]; !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		next.ServeHTTP(w, req)
	})
}
