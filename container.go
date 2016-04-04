package wac

import "html/template"

type StaticContainer struct {
	assetsCompilers []AssetCompiler
	pathToStaticDir string
	isDebugMode bool
	templates *template.Template
}

var (
	container *StaticContainer = nil
)


func CreateStaticContainer(assetsCompilers []AssetCompiler, pathToStaticDir string) *StaticContainer {
	container = &StaticContainer{
		assetsCompilers: assetsCompilers,
		pathToStaticDir: pathToStaticDir,
		isDebugMode: false,
	}

	container.loadCompiledTemplates()

	return container
}

func (container *StaticContainer) SetDebugMode(debugMode bool) {
	container.isDebugMode = debugMode
}

func checkContainer() {
	if container == nil {
		panic("Static container was not created!")
	}
}
