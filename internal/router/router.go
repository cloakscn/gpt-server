package router

import (
	"net/http"

	"github.com/cloakscn/gpt-server/internal/api"
)

func Router() {
	http.HandleFunc("/", api.SafeHandler(api.ListHandler))
	http.HandleFunc("/view", api.SafeHandler(api.ViewHandler))
	http.HandleFunc("/hello", api.SafeHandler(api.HelloHandler))
	http.HandleFunc("/upload", api.SafeHandler(api.UploadHandler))
}