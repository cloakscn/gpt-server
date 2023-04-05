package router

import (
	"io"
	"net/http"
	"os"

	. "github.com/cloakscn/gpt-server/internal/pkg/error"
	"github.com/cloakscn/gpt-server/internal/pkg/file"
	httpx "github.com/cloakscn/gpt-server/internal/pkg/http"
)

func PhotosRouter() {
	http.HandleFunc("/", httpx.SafeHandler(listHandler))
	http.HandleFunc("/view", httpx.SafeHandler(viewHandler))
	http.HandleFunc("/hello", httpx.SafeHandler(helloHandler))
	http.HandleFunc("/upload", httpx.SafeHandler(uploadHandler))
}

const (
	UPLOAD_DIR = "./uploads"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1>Hello, world!</h1>")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := httpx.RenderHtml(w, "upload", nil); err != nil {
			Check(err)
			return
		}
		return
	}

	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		if err != nil {
			Check(err)
			return
		}
		filename := h.Filename
		defer f.Close()

		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		if err != nil {
			Check(err)
			return
		}
		defer t.Close()

		if _, err := io.Copy(t, f); err != nil {
			Check(err)
			return
		}
		http.Redirect(w, r, "/view?id="+filename, http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageID := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageID
	if exists := file.IsExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	// w.Header().Set("Content-Type", "image")
	w.Header().Set("Content-Type", "video/mpeg4")
	http.ServeFile(w, r, imagePath)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fi, err := os.ReadDir(UPLOAD_DIR)
	if err != nil {
		Check(err)
		return
	}

	locals := make(map[string]interface{})
	var images []string
	for _, fileInfo := range fi {
		images = append(images, fileInfo.Name())
	}

	locals["images"] = images
	if err = httpx.RenderHtml(w, "list", locals); err != nil {
		Check(err)
		return
	}
}
