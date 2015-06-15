package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	Base = "base.tmpl"
)

var Templates = make(map[string]*template.Template)

func ParseTemplates() {
	f, err := ioutil.ReadDir("templates")
	if err != nil {
		log.Fatalf("Unable to load template directory %v", err.Error())
	}
	for _, v := range f {
		if !v.IsDir() && v.Name() != Base {
			if err := ParseTemplate(v.Name()); err != nil {
				log.Println("Unable to parse files:", err.Error())
			} else {
				log.Println("Loaded template", v.Name())
			}
		}
	}
}

func ParseTemplate(filename string) error {
	var err error
	Templates[filename], err = template.ParseFiles(filepath.Join("templates", filename), filepath.Join("templates", Base))
	return err
}

func Render(rw io.Writer, templateName string, g interface{}) {
	if v, ok := Templates[templateName]; ok {
		if rw == nil {
			log.Println("Somehow the response writer given to Render was nil")
			return
		}
		if g == nil {
			log.Println("The interface given to Render as g was nil")
			return
		}
		if v == nil {
			log.Println("Templates gave an ok but v was nil")
			return
		}
		if err := v.Execute(rw, g); err != nil {
			log.Println("Execution of the template", templateName, "failed because", err.Error())
		}
	} else {
		log.Println("Caller to render called template that doesn't exist.", templateName, "valid templates are")
		for k, _ := range Templates {
			log.Println(k)
		}
	}
}

/*
	Assets stuff
*/

func UnpackAssets() {
	for _, v := range AssetNames() {
		unpackAsset(v)
	}
}
func unpackAsset(assName string) {
	d, _ := filepath.Split(assName)
	if err := os.MkdirAll(d, 0777); err != nil {
		log.Fatalf("Unable to make directory for asset: %v [%v]", d, err.Error())
	}
	data, err := Asset(assName)
	if err != nil {
		log.Fatalf("Unable to unpack asset from binary: %v [%v]", assName, err.Error())
	}
	if err := ioutil.WriteFile(assName, data, 0777); err != nil {
		log.Fatalf("Unable to write unpacked asset to directory: %v [%v]", assName, err.Error())
	}
}
