package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/moov-io/xsd2go/pkg"
	"github.com/moov-io/xsd2go/pkg/xsd"
)

func GenerateTypes(schema *xsd.Schema, outputDir string, outputFile string, tmplDir string, templateName string) error {
	t, err := newTemplate(tmplDir, templateName)
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

	_, err = exec.Command("goimports", "-w", f.Name()).Output()
	if err != nil {
		return err
	}

	return nil
}

func newTemplate(tmplDir string, templateName string) (*template.Template, error) {
	in, err := getFile(filepath.Join(tmplDir, templateName))
	if err != nil {
		return nil, err
	}
	defer in.Close()

	tempText, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	return template.New(templateName).Funcs(template.FuncMap{}).Parse(string(tempText))
}

func getFile(tmplPath string) (fs.File, error) {
	file, err := pkg.Templates.Open(tmplPath)
	if err == nil {
		return file, err
	}
	return os.Open(tmplPath)
}
