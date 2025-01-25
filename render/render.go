package render

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

type Renderer struct {
	Renderer   string
	RootPath   string
	Secure     bool
	ServerName string
	Port       string
}

type TemplateData struct {
	IsAuthenticated bool
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	ServerName      string
	Port            string
	Secure          bool
	Data            map[string]interface{}
	CSRFToken       string
}

func (q *Renderer) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(q.Renderer) {
	case "go":
		q.GoPage(w, r, view, data)
	case "jet":

	}
	return nil
}

func (q *Renderer) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", q.RootPath, view))
	if err != nil {
		return err
	}
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	err = tmpl.Execute(w, &td)

	if err != nil {
		return err
	}
	return nil
}
