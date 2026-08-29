package backend

import (
	"io"
	"text/template"

	"github.com/kohirens/storage"
)

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

func NewTemplateManager(store storage.Storage, location, suffix string) TemplateManager {
	return &Renderer{
		location:  location,
		store:     store,
		suffix:    suffix,
		Vars:      make(map[string]any),
		functions: template.FuncMap{},
	}
}
