package api

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/cloakscn/gpt-server/internal/pkg/file"
	. "github.com/cloakscn/gpt-server/internal/pkg/error"
	"github.com/cloakscn/gpt-server/internal/pkg/html"
)

const (
	UPLOAD_DIR   = "./uploads"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1>Hello, world!</h1>")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := html.RenderHtml(w, "upload", nil); err != nil {
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

func ViewHandler(w http.ResponseWriter, r *http.Request) {
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

func ListHandler(w http.ResponseWriter, r *http.Request) {
	fi, err := ioutil.ReadDir(UPLOAD_DIR)
	if err != nil {
		Check(err)
		return
	}

	locals := make(map[string]interface{})
	images := []string{}
	for _, fileInfo := range fi {
		images = append(images, fileInfo.Name())
	}

	locals["images"] = images
	if err = html.RenderHtml(w, "list", locals); err != nil {
		Check(err)
		return
	}
}

func SafeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的 50x 错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// renderHtml(w, "error", e)

				// logging
				log.Printf("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}