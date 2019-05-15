package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "List server")
	})

	dir, _ := os.Getwd()
	fmt.Println("Current working directory: " + dir)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	bind := ":" + os.Getenv("PORT")
	fmt.Printf("Listening on " + bind)
	http.ListenAndServe(bind, nil)
}
