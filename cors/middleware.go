package cors

import (
	"fmt"
	"net/http"
	"time"
)

func MiddlewareFunc(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//CORS headers:
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//
		fmt.Println("before handler; middlerware start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middlerware finished; %s\n", time.Since(start))
	})
}
