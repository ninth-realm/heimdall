package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/ninth-realm/heimdall/docs/api"
)

//go:embed index.html
var index embed.FS

func main() {
	http.Handle(
		"/raw/",
		http.StripPrefix("/raw", http.FileServer(http.FS(api.Spec))),
	)
	http.Handle("/", http.FileServer(http.FS(index)))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
