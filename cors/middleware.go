package cors

import "net/http"

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Access-Control-Allow-Origin", "*")
		res.Header().Add("Content-Type", "application/json")
		res.Header().Set("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers", "*")
		handler.ServeHTTP(res, req)

	})
}
