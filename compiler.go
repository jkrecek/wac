package wac

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
)

type AssetFileType struct {
	glob       string
	command    string
	parameters []string
}

type AssetFileTypes []AssetFileType

type AssetCompiler struct {
	assetType   string
	contentType string
	fileTypes   AssetFileTypes
}

/*
var compilers = []AssetCompiler{
	AssetCompiler{
		assetType:   "css",
		contentType: "text/css",
		fileTypes: AssetFileTypes{
			AssetFileType{
				"*.css",
				"",
				nil,
			},
			AssetFileType{
				"*.scss",
				config.Data.Commands.Sass,
				[]string{"%f"},
			},
		},
	},
	AssetCompiler{
		assetType:   "js",
		contentType: "application/javascript",
		fileTypes: AssetFileTypes{
			AssetFileType{
				"*.js",
				"",
				nil,
			},
			AssetFileType{
				"*.coffee",
				config.Data.Commands.Coffee,
				[]string{"-p", "%f"},
			},
		},
	},
}
*/

func (container *StaticContainer) AssetCompile(assetType string) ([]byte, string) {
	var assetCompiler *AssetCompiler
	for _, ac := range container.assetsCompilers {
		if ac.assetType == assetType {
			assetCompiler = &ac
			break
		}
	}

	if assetCompiler != nil {
		return container.compile(assetCompiler), assetCompiler.contentType
	}

	return []byte{}, ""
}

func (container *StaticContainer) compile(ac *AssetCompiler) []byte {
	var fileNames []string
	contents := make(map[string][]byte)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for _, assetFileType := range ac.fileTypes {
		typeFilePaths, err := filepath.Glob(container.pathToStaticDir + "/assets/" + ac.assetType + "/" + assetFileType.glob)
		if err != nil {
			log.Println("Error globbing:", err)
		}

		wg.Add(len(typeFilePaths))

		for _, tfp := range typeFilePaths {
			go func(ac *AssetCompiler, filePath string, aft AssetFileType) {
				fileContent := ac.loadFile(filePath, &aft)

				_, file := filepath.Split(filePath)

				mutex.Lock()
				fileNames = append(fileNames, file)
				contents[file] = fileContent
				mutex.Unlock()

				wg.Done()
			}(ac, tfp, assetFileType)
		}
	}

	wg.Wait()

	sort.Strings(fileNames)
	buffer := new(bytes.Buffer)
	for _, fileName := range fileNames {
		buffer.WriteString("/** Contents of " + fileName + " **/\n")
		buffer.Write(contents[fileName])
		buffer.WriteString("\n\n")
	}

	return buffer.Bytes()
}

func (ac *AssetCompiler) loadFile(filePath string, aft *AssetFileType) []byte {
	var outputBytes []byte
	var err error
	if aft.command == "" {
		outputBytes, err = ioutil.ReadFile(filePath)
	} else {
		params := make([]string, len(aft.parameters))
		for k, para := range aft.parameters {
			if para == "%f" {
				params[k] = filePath
			} else {
				params[k] = para
			}
		}

		cmd := exec.Command(aft.command, params...)
		pipeOut, _ := cmd.StdoutPipe()
		pipeErr, _ := cmd.StderrPipe()

		cmd.Start()

		outputBytes, _ = ioutil.ReadAll(pipeOut)
		errOut, _ := ioutil.ReadAll(pipeErr)

		cmd.Wait()

		if len(errOut) != 0 {
			err = errors.New(string(errOut))
		}
	}

	if err != nil {
		log.Println("Error handling file ", filePath, ":", err)
	}

	return outputBytes
}
