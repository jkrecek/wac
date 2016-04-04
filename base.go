package wac

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func TemplateRenderer(w http.ResponseWriter, template string) {
	checkContainer()

	container.renderTemplate(w, template, nil)
}

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	checkContainer()

	vars := mux.Vars(r)
	assetType := vars["type"]

	output, contentType := container.AssetCompile(assetType)
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.Header().Set("Cache-Control", "private, max-age=600")

	if container.isDebugMode {
		w.Write(output)
	} else {
		bfr := bytes.NewReader(output)

		if err := minifier.Minify(contentType, w, bfr); err != nil {
			log.Fatal("Minify error: ", err, contentType)
		}
	}
}
