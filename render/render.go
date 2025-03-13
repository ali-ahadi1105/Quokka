package render

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	ServerName string
	Port       string
}

// data info that passing in template
type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	Secure          bool
	CSRFToken       string
	ServerName      string
	Port            string
}

func (q *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(q.Renderer) {
	case "go":
		q.GoPage(w, r, view, data)
	case "jet":
	}
	return nil
}

func (q *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", q.RootPath, view))
	if err != nil {
		return err
	}
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}
	return nil
}
