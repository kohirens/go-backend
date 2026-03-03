package backend

import (
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"text/template"

	"github.com/kohirens/www/storage"
)

type Renderer struct {
	store     storage.Storage
	location  string
	suffix    string
	Vars      map[string]any
	functions template.FuncMap
}

type TemplateManager interface {
	// AddFunctions template functions.
	AddFunctions(functions template.FuncMap)
	// AddVar Add an item to the variable map.
	AddVar(k string, v any)
	// AppendVars Appends a map to the Renderer.Vars map.
	//
	//	NOTE: When a key matches an existing key it will overwrite its value.
	AppendVars(vars map[string]any)
	// Load A template, but it will not render it, instead, the template.Template
	// object is returned so that you can render it when you want.
	Load(name string) (*template.Template, error)
	// LoadFiles Parse multiple templates that produces the desired output.
	LoadFiles(names ...string) (*template.Template, error)
	// Render Write a templates' content to a writer. You can provide vars
	// as a type `map[string]string` of key-value pairs; which will be used to fill
	// in string placeholders. Nothing more complex is supported at this time.
	// Also, remember that maps are by default passed by reference, so there is
	// no need to pass vars as a pointer.
	Render(name string, w io.Writer, vars map[string]any) error
	// RenderFiles Parse multiple templates that produces the desired output.
	//
	//	This uses LoadFiles which in turn uses template.ParseFiles,
	//	see https://pkg.go.dev/text/template#ParseFiles.
	RenderFiles(w io.Writer, vars map[string]any, names ...string) (*template.Template, error)
}

type Variables map[string]any

const ps = string(os.PathSeparator)

func NewTemplateManager(store storage.Storage, location, suffix string) TemplateManager {
	return &Renderer{
		location:  location,
		store:     store,
		suffix:    suffix,
		Vars:      make(map[string]any),
		functions: template.FuncMap{},
	}
}

// AddFunctions template functions.
func (m *Renderer) AddFunctions(functions template.FuncMap) {
	for name, function := range functions {
		m.functions[name] = function
	}
}

// AddVar Add an item to the variable map.
func (m *Renderer) AddVar(k string, v any) {
	m.Vars[k] = v
}

// AppendVars Appends a map to the Renderer.Vars map.
//
//	NOTE: When a key matches an existing key it will overwrite its value.
func (m *Renderer) AppendVars(vars map[string]any) {
	maps.Copy(m.Vars, vars)
}

// Load Parse a template into memory, but it will not render it, instead, the
// template.Template object is returned so that you can render it when you want.
func (m *Renderer) Load(name string) (*template.Template, error) {
	filename := buildFilename(m, name)

	Log.Infof(stdout.LoadTemplate, filename)

	tmplContent, e1 := m.store.Load(filename)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.FileNotFound, filename, e1.Error())
	}

	t, e2 := template.New(filename).Funcs(m.functions).Parse(string(tmplContent))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.TemplateParse, e2.Error())
	}

	Log.Dbugf(stdout.TemplateLoad, filename)

	return t, nil
}

// LoadFiles parse multiple templates that produces the desired output.
//
//	This loads each file from the storage first and passes it to template.Parse;
//	using the basename of each file as the name for the template. Similar to how
//	template.ParseFiles works. template.Parse was preferred to ensure this
//	system works with storage mediums that are not file based.
//	This uses template.ParseFiles, see
//	https://pkg.go.dev/text/template#Parse.
func (m *Renderer) LoadFiles(names ...string) (*template.Template, error) {
	tmplName := filepath.Base(buildFilename(m, names[0]))
	Log.Dbugf(stdout.LoadTemplate, tmplName)

	t := template.New(tmplName).Funcs(m.functions)
	var err error
	var tmp *template.Template

	for _, name := range names {
		filename := buildFilename(m, name)
		data, e2 := m.store.Load(filename)
		if e2 != nil {
			return nil, e2
		}
		tmp, err = t.Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf(stderr.TemplateParse, err.Error())
		}
		t = tmp
		tmp = nil
	}

	return t, nil
}

// Render Write a templates' content to a writer. You can provide vars
// as a type `map[string]string` of key-value pairs; which will be used to fill
// in string placeholders. Nothing more complex is supported at this time.
// Also, remember that maps are by default passed by reference, so there is
// no need to pass vars as a pointer.
func (m *Renderer) Render(name string, w io.Writer, vars map[string]any) error {
	t, e1 := m.Load(name)
	if e1 != nil {
		return e1
	}

	// Combine vars with any previously added to the manager.
	if vars != nil {
		m.AppendVars(vars)
	}

	return t.Execute(w, m.Vars)
}

// RenderFiles Parse multiple templates that produces the desired output.
//
//	This uses LoadFiles which in turn uses template.ParseFiles,
//	see https://pkg.go.dev/text/template#ParseFiles.
func (m *Renderer) RenderFiles(w io.Writer, vars map[string]any, names ...string) (*template.Template, error) {
	t, e1 := m.LoadFiles(names...)
	if e1 != nil {
		return nil, e1
	}

	// Combine vars with any previously added to the manager.
	if vars != nil {
		m.AppendVars(vars)
	}

	if e := t.Execute(w, m.Vars); e != nil {
		return nil, fmt.Errorf(stderr.RenderFiles, e.Error())
	}
	return t, nil
}

func buildFilename(m *Renderer, name string) string {
	if len(m.location) > 0 && m.location[len(m.location)-1] != '/' {
		return m.location + ps + name + "." + m.suffix
	}
	return m.location + name + "." + m.suffix
}
