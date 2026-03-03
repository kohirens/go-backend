package backend

import (
	"bytes"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/kohirens/www/storage"
)

func TestHandler_Render(t *testing.T) {
	testWd, _ := filepath.Abs(fixtureDir)
	fixtures, _ := storage.NewLocalStorage(testWd)

	cases := []struct {
		name     string
		store    storage.Storage
		filename string
		vars     Variables
		wantW    string
		wantErr  bool
	}{
		{
			"simple-render",
			fixtures,
			"test-render-01",
			Variables{"TestVar": "1234"},
			"1234",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewTemplateManager(c.store, "", "tmpl")
			w := &bytes.Buffer{}
			err := h.Render(c.filename, w, c.vars)

			if (err != nil) != c.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if gotW := w.String(); gotW != c.wantW {
				t.Errorf("Render() gotW = %v, want %v", gotW, c.wantW)
				return
			}
		})
	}
}

func TestRenderer_AddFunctions(t *testing.T) {
	testWd, _ := filepath.Abs(fixtureDir)
	fixtures, _ := storage.NewLocalStorage(testWd)

	cases := []struct {
		name     string
		store    storage.Storage
		filename string
		vars     Variables
		funcs    template.FuncMap
		wantW    string
		wantErr  bool
	}{
		{
			"simple-render",
			fixtures,
			"test-function-render-01",
			Variables{"A": 1, "B": 2},
			template.FuncMap{"add": func(a, b int) int { return a + b }},
			"3",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewTemplateManager(c.store, "", "tmpl")
			h.AddFunctions(c.funcs)
			w := &bytes.Buffer{}
			err := h.Render(c.filename, w, c.vars)

			if (err != nil) != c.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if gotW := w.String(); gotW != c.wantW {
				t.Errorf("Render() gotW = %v, want %v", gotW, c.wantW)
				return
			}
		})
	}
}

func TestRenderer_LoadFilesMax(t *testing.T) {
	testWd, _ := filepath.Abs(fixtureDir)
	fixtures, _ := storage.NewLocalStorage(testWd)

	cases := []struct {
		name      string
		store     storage.Storage
		filenames []string
		vars      Variables
		funcs     template.FuncMap
		wantW     string
		wantErr   bool
	}{
		{
			"multiple-files",
			fixtures,
			[]string{"test-load-multiple-template-01", "test-load-multiple-template-02"},
			Variables{"A": 1, "B": 1},
			template.FuncMap{"add": func(a, b int) int { return a + b }},
			"render has template 2 content",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewTemplateManager(c.store, "", "tmpl")
			h.AddFunctions(c.funcs)
			w := &bytes.Buffer{}
			_, err := h.RenderFiles(w, c.vars, c.filenames...)

			if (err != nil) != c.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if gotW := w.String(); gotW != c.wantW {
				t.Errorf("Render() gotW = %v, want %v", gotW, c.wantW)
				return
			}
		})
	}
}
