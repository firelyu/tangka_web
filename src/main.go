package main

import (
	"net/http"
	"github.com/firelyu/tangka_web/src/tangka"
	"html/template"
	"io/ioutil"
	"fmt"
)

const (
	LISTEN_PORT = ":8080"
	TEMPLATE_DIR = "tmpl"
	TEMPLATE_SUFFIX = ".html"
)

// Read the template from the templates cache
func renderTemplate(w http.ResponseWriter, tmpl string, t *tangka.Tangka) {
	err := templates.ExecuteTemplate(w, tmpl+TEMPLATE_SUFFIX, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// parse the common post data
func parseOneTangkaHandle(fn func(http.ResponseWriter, *http.Request, *tangka.Tangka)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		t, err := tangka.GetTangkaById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fn(w, r, t)
	}
}

// edit the exist tangka
func editHandle(w http.ResponseWriter, r *http.Request, t *tangka.Tangka) {
	renderTemplate(w, "edit", t)
}

// save new tangka
func saveHandle(w http.ResponseWriter, r *http.Request, t *tangka.Tangka)  {
	newName := r.FormValue("name")
	t.Name = newName
	err := t.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/detail/?id=" + t.Id, http.StatusFound)
}

// remove the tangka
func deleteHandle(w http.ResponseWriter, r *http.Request, t *tangka.Tangka)  {
	err := t.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/list", http.StatusFound)
}

// show the detail of tangka
func detailHandle(w http.ResponseWriter, r *http.Request, t *tangka.Tangka) {
	renderTemplate(w, "detail", t)
}

// list all the tangka
func listHandle(w http.ResponseWriter, r *http.Request)  {
	tangkaList, err := tangka.ListAllTangka()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "list.html", tangkaList)
}

var (
	templates *template.Template
)

// cache the templates
func cacheTemplates(dir string) error {
	tmplList, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var tmplPathList []string
	for _, file := range tmplList {
		tmplPathList = append(tmplPathList, dir + "/" + file.Name())
	}
	templates = template.Must(template.ParseFiles(tmplPathList...))

	return nil
}

func main() {
	if err := cacheTemplates(TEMPLATE_DIR); err != nil {
		fmt.Println(err.Error())
		return
	}

	http.HandleFunc("/detail/", parseOneTangkaHandle(detailHandle))
	http.HandleFunc("/edit/", parseOneTangkaHandle(editHandle))
	http.HandleFunc("/save/", parseOneTangkaHandle(saveHandle))
	http.HandleFunc("/delete/", parseOneTangkaHandle(deleteHandle))
	http.HandleFunc("/list", listHandle)

	http.ListenAndServe(LISTEN_PORT, nil)

}
