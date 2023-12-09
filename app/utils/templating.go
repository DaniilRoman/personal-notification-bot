package modules

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

func RenderIndexHTML(data map[string]interface{}) {
	templatePath := fmt.Sprintf("%s/%s/template.html", mainDir, wwwDir)
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
