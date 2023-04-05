package html

import (
	"html/template"
	"net/http"

	"github.com/cloakscn/gpt-server/internal/pkg/file"
)

const LIST_DIR = 0x0001

var (
	Templates map[string]*template.Template
)


func StaticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		filePath := staticDir + r.URL.Path[len(prefix)-1:]
		if (flags & LIST_DIR) == 0 {
			if exists := file.IsExists(filePath); !exists {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, filePath)
	})
}

func RenderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) (err error) {
	return Templates[tmpl + ".html"].Execute(w, locals)
}