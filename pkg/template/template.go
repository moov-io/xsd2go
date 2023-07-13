package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/markbates/pkger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/moov-io/xsd2go/pkg/xsd"
)

func GenerateTypes(schema *xsd.Schema, outputDir string, outputFile string, templateName string) error {
	t, err := newTemplate(templateName)
	if err != nil {
		return err
	}

	packageName := schema.GoPackageName()
	dir := filepath.Join(outputDir, packageName)
	err = os.MkdirAll(dir, os.FileMode(0722))
	if err != nil {
		return err
	}
	goFile := filepath.Clean(filepath.Join(dir, outputFile))
	fmt.Printf("\tGenerating '%s'\n", goFile)
	f, err := os.Create(goFile)
	if err != nil {
		return fmt.Errorf("Could not create '%s': %v", goFile, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, schema); err != nil {
		return fmt.Errorf("Could not execute template: %v", err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Could not gofmt output file\nError was: '%v'\nFile was:\n%s\n", err, buf.String())
	}

	_, err = f.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func newTemplate(templateName string) (*template.Template, error) {
	in, err := getFile(templateName)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	tempText, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	return template.New(templateName).Funcs(template.FuncMap{
		"title": cases.Title(language.AmericanEnglish).String,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"split": strings.Split,
	}).Parse(string(tempText))
}

// getFile returns a fs.File either using pkger or the OS. This allows for templates outside the packaged program to be used.
func getFile(templateName string) (fs.File, error) {
	file, err := pkger.Open(templateName)
	if err == nil {
		return file, err
	}
	return os.Open(templateName)
}
