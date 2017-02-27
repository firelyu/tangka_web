package main

import (
	"net/http"
	"github.com/firelyu/tangka_web/src/tangka"
	"html/template"
	"io/ioutil"
	"github.com/astaxie/beego/logs"
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
func parseOneTangkaHandle(fn func(http.ResponseWriter, *http.Request, *tangka.Tangka), name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs.Info("Handle " + name)
		id := r.FormValue("id")
		t, err := tangka.GetTangkaById(id)
		if err != nil {
			logs.Error(err.Error())
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

// save new/edited tangka
func saveHandle(w http.ResponseWriter, r *http.Request)  {
	logs.Info("Handle /save/")
	id := r.FormValue("id")
	t, _ := tangka.GetTangkaById(id)
	if t == nil {
		// New save tangka
		logs.Info("Save new tangak", t)
		t = &tangka.Tangka{Id:id}
	}

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
	logs.Info("Handle /list/")
	tangkaList, err := tangka.ListAllTangka()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "list.html", tangkaList)
}

//  the homepage
func homeHandle(w http.ResponseWriter, r *http.Request) {
	logs.Info("Handle /")
	h := tangka.NewHomepage("I Love Thangka.")
	templates.ExecuteTemplate(w, "index.html", h)
}

func addHandle(w http.ResponseWriter, r *http.Request) {
	logs.Info("Handle /add/")
	t := tangka.NewTangka("", "")
	templates.ExecuteTemplate(w, "add.html", t)
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
		if file.IsDir() {
			continue
		}
		tmplPathList = append(tmplPathList, dir + "/" + file.Name())
	}
	templates = template.Must(template.ParseFiles(tmplPathList...))

	return nil
}

func main() {
	if err := cacheTemplates(TEMPLATE_DIR); err != nil {
		logs.Error(err.Error())
		return
	}

	http.HandleFunc("/detail/", parseOneTangkaHandle(detailHandle, "/detail/"))
	http.HandleFunc("/edit/", parseOneTangkaHandle(editHandle, "/edit/"))
	http.HandleFunc("/save/", saveHandle)
	http.HandleFunc("/delete/", parseOneTangkaHandle(deleteHandle, "/delete/"))
	http.HandleFunc("/add/", addHandle)
	http.HandleFunc("/list/", listHandle)
	http.HandleFunc("/", homeHandle)

	http.ListenAndServe(LISTEN_PORT, nil)

}
