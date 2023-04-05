package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	httpx "github.com/cloakscn/gpt-server/internal/pkg/http"
	"github.com/cloakscn/gpt-server/internal/router"
)

const (
	TEMPLATE_DIR = "./views"
)

func init() {
	httpx.Templates = make(map[string]*template.Template)
	fi, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		panic(err)
	}

	var templateName, templatePath string
	for _, fileInfo := range fi {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading teamplate:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		httpx.Templates[templateName] = t
	}
}

func main() {
	mux := http.NewServeMux()
	// 静态资源和动态请求分离
	httpx.StaticDirHandler(mux, "/assets/", "./public", 0)
	router.PhotosRouter()
	router.GptRouter()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
}
