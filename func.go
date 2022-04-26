package repository

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

// 拷贝template 包helper 方法
func readFileFS(fsys fs.FS) func(string) (string, []byte, error) {
	return func(file string) (name string, b []byte, err error) {
		name = path.Base(file)
		b, err = fs.ReadFile(fsys, file)
		return
	}
}

// parseFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func parseFiles(t *template.Template, readFile func(string) (string, []byte, error), filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		name, b, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func GlobDirectory(dir string, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		pattern = fmt.Sprintf("%s/*%s", dir, pattern)
		return filepath.Glob(pattern)
	}
	var matches []string
	reg := getPattern(pattern)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			err := errors.Errorf("dir:%s filepath.Walk info is nil", dir)
			return err
		}
		if !info.IsDir() {
			path = strings.ReplaceAll(path, "\\", "/")
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func GlobFS(fsys fs.FS, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		return fs.Glob(fsys, pattern)
	}
	var matches []string
	reg := getPattern(pattern)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func getPattern(pattern string) *regexp.Regexp {
	regStr := strings.TrimLeft(pattern, ".")
	regStr = strings.ReplaceAll(regStr, ".", "\\.")
	regStr = strings.ReplaceAll(regStr, "**", ".*")
	reg := regexp.MustCompile(regStr)
	return reg
}

//TrimSpaces  去除开头结尾的非有效字符
func TrimSpaces(s string) string {
	return strings.Trim(s, "\r\n\t\v\f ")
}

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
