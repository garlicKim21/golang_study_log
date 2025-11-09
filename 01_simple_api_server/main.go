package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	h1 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #1!\n")
	}
	h2 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #2!\n")
	}
	h3 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #3!\n")
	}
	h4 := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #4!\n")
	}

	http.HandleFunc("/", h1)           // 이 HandleFunc는 '/' 경로로 시작하는 모든 요청을 처리합니다.
	http.HandleFunc("/{$}", h2)        // 이 HandleFunc는 '/' 경로와 정확이 일치하는 요청을 처리합니다.
	http.HandleFunc("/endpoint1", h3)  // 이 HandleFunc는 '/endpoint1' 경로와 정확히 일치하는 요청을 처리합니다.
	http.HandleFunc("/endpoint2/", h4) // 이 HandleFunc는 '/endpoint2' 경로와 정확히 일치하는 요청을 304 리다이렉트 처리하여 /endpoint2/ 경로로 리다이렉트 처리합니다. 그리고 /endpoint2/{anypath} 경로로 시작하는 모든 요청을 처리합니다.
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
