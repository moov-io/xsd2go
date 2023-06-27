package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gocomply/xsd2go/pkg/xsd"
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
	if err := t.ExecuteTemplate(&buf, templateName, schema); err != nil {
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
	t := template.New(templateName).Funcs(template.FuncMap{
		// Allow any template ending to be included inline. The main template will call this function at a specific point.
		"InclCType": func(tmplName string, data *xsd.ComplexType) (string, error) {
			tmplName = fmt.Sprintf("%s%s", tmplName, xsd.TemplateTypeInclude)
			return includeTemplate(tmplDir, tmplName, data)
		},
		"InclSType": func(tmplName string, data *xsd.SimpleType) (string, error) {
			tmplName = fmt.Sprintf("%s%s", tmplName, xsd.TemplateTypeInclude)
			return includeTemplate(tmplDir, tmplName, data)
		},
		"InclEType": func(tmplName string, data *xsd.Element) (string, error) {
			tmplName = fmt.Sprintf("%s%s", tmplName, xsd.TemplateTypeInclude)
			return includeTemplate(tmplDir, tmplName, data)
		},
		"InclElem": func(tmplName string, data *xsd.Element) (string, error) {
			tmplName = fmt.Sprintf("%s%s", tmplName, xsd.TemplateTypeElement)
			return includeTemplate(tmplDir, tmplName, data)
		},
	})
	return parseTemplate(t, tmplDir, templateName)
}

func includeTemplate(tmplDir string, tmplName string, data any) (string, error) {
	t2 := template.New(tmplName)
	t2, tmplErr := parseTemplate(t2, tmplDir, tmplName)
	if tmplErr != nil {
		return "", tmplErr
	}
	var tmplBuf bytes.Buffer
	tmplErr = t2.Execute(&tmplBuf, data)
	return tmplBuf.String(), tmplErr
}

func parseTemplate(t *template.Template, tmplDir string, templateName string) (*template.Template, error) {
	in, err := Templates.Open(tmplDir + "/" + templateName)
	if err != nil {
		return t, err
	}
	defer in.Close()

	tempText, err := ioutil.ReadAll(in)
	if err != nil {
		return t, err
	}
	t, err = t.Parse(string(tempText))
	if err != nil {
		return t, err
	}
	return t, nil
}

func GetAllTemplates(tmplDir string) (map[string]xsd.Override, error) {
	dir, err := Templates.ReadDir(tmplDir)
	if err != nil {
		return nil, err
	}

	templates := make(map[string]xsd.Override)
	for indx := range dir {
		name, found := strings.CutSuffix(dir[indx].Name(), xsd.TemplateTypeInclude)
		if found {
			if _, ok := templates[name]; !ok {
				templates[name] = xsd.Override{TemplateName: name}
			}
			override := templates[name]
			override.IsIncl = true
			templates[name] = override
		}
		name, found = strings.CutSuffix(dir[indx].Name(), xsd.TemplateTypeElement)
		if found {
			if _, ok := templates[name]; !ok {
				templates[name] = xsd.Override{TemplateName: name}
			}
			override := templates[name]
			override.IsElem = true
			templates[name] = override
		}
	}

	return templates, nil
}
