package utils

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

const (
	mainDir = "./"
	wwwDir  = "www"
)

func RenderWwwResources(data map[string]interface{}) {
	renderIndexHTML(data)
	renderScriptJs(data)
}

func renderIndexHTML(data map[string]interface{}) {
	templatePath := fmt.Sprintf("%s/%s/indexTemplate.html", mainDir, wwwDir)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	indexHTMLPath := fmt.Sprintf("%s/%s/index.html", mainDir, wwwDir)
	indexFile, err := os.Create(indexHTMLPath)
	if err != nil {
		log.Fatalf("Error creating index.html file: %v", err)
	}
	defer indexFile.Close()

	err = tmpl.Execute(indexFile, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func renderScriptJs(data map[string]interface{}) {
	templatePath := fmt.Sprintf("%s/%s/scriptTemplate.js", mainDir, wwwDir)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	scriptJsPath := fmt.Sprintf("%s/%s/script.js", mainDir, wwwDir)
	indexFile, err := os.Create(scriptJsPath)
	if err != nil {
		log.Fatalf("Error creating script.js file: %v", err)
	}
	defer indexFile.Close()

	err = tmpl.Execute(indexFile, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
