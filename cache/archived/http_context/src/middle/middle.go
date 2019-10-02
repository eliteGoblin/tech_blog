package main

import (
	"fmt"
	"net/http"
	"strings"
	"utils"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("enter middleware")
		ctx := utils.NewContextWithRequestID(req.Context(), req)
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

func middleHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sayHello")
	resp, err := http.Get("http://localhost:8089")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("middle got %+v\n", *resp)
	}
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))

}

func main() {

	http.Handle("/", middleware(http.HandlerFunc(middleHandle)))
	if err := http.ListenAndServe(":8087", nil); err != nil {
		panic(err)
	}
}
`1`