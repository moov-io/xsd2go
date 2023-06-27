package xsd2go

import (
	"fmt"

	"github.com/gocomply/xsd2go/pkg/template"
	"github.com/gocomply/xsd2go/pkg/xsd"
)

type Params struct {
	XsdPath         string
	OutputDir       string
	OutputFile      string
	GoModuleImport  string
	TemplatePackage string
	TemplateName    string
	XmlnsOverrides  []string
}

func Convert(params Params) error {
	fmt.Printf("Processing '%s'\n", params.XsdPath)
	fmt.Printf("Cmd: gocomply_xsd2go convert "+
		"--xsd-file=%s "+
		"--output-dir=%s "+
		"--output-file=%s "+
		"--go-module-import=%s "+
		"--template-package=%s"+
		"--template-name=%s"+
		"--xmlns-override=%s"+
		"\n",
		params.XsdPath,
		params.OutputDir, params.OutputFile,
		params.GoModuleImport,
		params.TemplatePackage, params.TemplateName,
		params.XmlnsOverrides,
	)
	if params.TemplatePackage == "" {
		params.TemplatePackage = "template"
	}
	if params.TemplateName == "" {
		params.TemplateName = "types.tmpl"
	}

	ws, err := xsd.NewWorkspace(fmt.Sprintf("%s/%s", params.GoModuleImport, params.OutputDir), params.XsdPath, params.XmlnsOverrides)
	if err != nil {
		return err
	}

	for _, sch := range ws.Cache {
		if sch.Empty() {
			continue
		}

		var schemaOutputFile = params.OutputFile
		if schemaOutputFile == "" {
			schemaOutputFile = fmt.Sprintf("%s.go", sch.GoPackageName())
		}

		if err := template.GenerateTypes(sch, params.OutputDir, schemaOutputFile, params.TemplatePackage, params.TemplateName); err != nil {
			return err
		}
	}

	return nil
}
