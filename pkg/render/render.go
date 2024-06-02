package render

import (
	"bytes"
	"github.com/Joshua-Russel/bookings/pkg/config"
	"github.com/Joshua-Russel/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

func NewTemplate(conf *config.AppConfig) {
	app = conf
}
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(res http.ResponseWriter, file string, tmpldata *models.TemplateData) {
	//renderedFile, _ := template.ParseFiles("./templates/"+file, "./templates/base.layout.gohtml")
	//err := renderedFile.Execute(res, nil)
	//if err != nil {
	//	fmt.Println("error", err)
	//}
	var tmplcache map[string]*template.Template
	if app.UseCache {
		tmplcache = app.TmplCache
	} else {
		tmplcache, _ = CreateTemplate()
	}

	tmpl, isThere := tmplcache[file]
	if !isThere {
		log.Fatal("cannot read template from map")
	}
	buff := new(bytes.Buffer)
	tmpldata = AddDefaultData(tmpldata)
	err := tmpl.Execute(buff, tmpldata)
	if err != nil {
		log.Println(err)
	}
	//_, err = buff.WriteTo(os.Stdout)
	_, err = buff.WriteTo(res)

	if err != nil {
		log.Println(err)

	}

}

func CreateTemplate() (map[string]*template.Template, error) {
	tmplcache := map[string]*template.Template{}
	filenames, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		log.Println("error", err)
		return tmplcache, err
	}
	for _, name := range filenames {
		filename := filepath.Base(name)

		tmpl, err := template.New(filename).ParseFiles("./" + name)

		if err != nil {
			log.Println("error", err)
			return tmplcache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return tmplcache, err
		}
		if len(matches) > 0 {
			tmpl, err = tmpl.ParseGlob("./templates/*.layout.gohtml")

			if err != nil {
				return tmplcache, err
			}
		}
		tmplcache[filename] = tmpl

	}
	return tmplcache, nil
}

//var tempcache = make(map[string]*template.Template)

//func RenderTemplate(res http.ResponseWriter, file string) {
//	var err error
//	_, inMap := tempcache[file]
//
//	if !inMap {
//		log.Println("creating new Template")
//		err = createTemplate(file)
//		if err != nil {
//			fmt.Println("error", err)
//		}
//
//	} else {
//		log.Println("using cached template")
//	}
//	tmpl := tempcache[file]
//	err = tmpl.Execute(res, nil)
//	if err != nil {
//		log.Println("error", err)
//	}
//
//}
//func createTemplate(file string) error {
//	templates := []string{
//		fmt.Sprintf("./templates/%s", file),
//		"./templates/base.layout.gohtml",
//	}
//	tmpl, err := template.ParseFiles(templates...)
//	if err != nil {
//		return err
//	}
//	tempcache[file] = tmpl
//
//	return nil
//}
