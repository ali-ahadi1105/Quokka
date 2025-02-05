package render

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Render struct {
	ServerName string
	Secure     bool
	Port       string
	Renderer   string
	RootPath   string
}

type TemplateData struct {
	ServerName      string
	Secure          bool
	Port            string
	CSRFToken       string
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	IsAuthenticated bool
}

func (q *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(q.Renderer) {
	case "go":
		return q.GoPage(w, r, view, data)
	case "jet":
	}
	return nil
}

func (q *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s", q.RootPath, view))
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
