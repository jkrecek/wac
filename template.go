package wac

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

var (
	minifier = createMinifier()
)

func (container *StaticContainer) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	if container.isDebugMode {
		container.loadCompiledTemplates()
	}

	buffer := new(bytes.Buffer)
	err := container.templates.ExecuteTemplate(buffer, tmpl + ".html", data)
	if err != nil {
		log.Println("Cannot execute", err)
	}

	if container.isDebugMode {
		w.Write(buffer.Bytes())
	} else {
		if err := minifier.Minify("text/html", w, buffer); err != nil {
			log.Fatal("Minify error html: ", err)
		}
	}
}

func (container *StaticContainer) loadCompiledTemplates() {
	container.templates = template.Must(template.ParseGlob(container.pathToStaticDir + "/templates/*.html"))
}

func createMinifier() *minify.Minify {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("application/javascript", js.Minify)
	return m
}
