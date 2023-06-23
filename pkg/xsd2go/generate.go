package xsd2go

import (
	"fmt"
	"path/filepath"

	"github.com/gocomply/xsd2go/pkg/xsd"
)

func Convert(xsdFile string, outputDir string, outputFile string, goModulesPath string, goPackage string, nsPrefix string, tmplDir string) error {
	fmt.Printf("Processing '%s'\n", xsdFile)
	fmt.Printf("Cmd: gocomply_xsd2go convert "+
		"--xsd-file=%s "+
		"--output-dir=%s "+
		"--output-file=%s "+
		"--go-module=%s "+
		"--go-package=%s "+
		"--namespace-prefix=%s "+
		"--template-package=%s\n",
		xsdFile, outputDir, outputFile, goModulesPath, goPackage, nsPrefix, tmplDir,
	)

	templates, err := GetAllTemplates(tmplDir)
	if err != nil {
		return err
	}

	ws, err := xsd.NewWorkspace(goModulesPath, goPackage, nsPrefix, xsdFile, templates)
	if err != nil {
		return err
	}

	for _, sch := range ws.Cache {
		if sch.Empty() {
			continue
		}

		var schemaOutputDir = filepath.Join(outputDir, sch.GoPackageName())
		var schemaOutputFile = outputFile
		if schemaOutputFile == "" {
			schemaOutputFile = fmt.Sprintf("%s.go", sch.GoPackageName())
		}

		if err = GenerateTypes("element.tmpl", sch, schemaOutputDir, schemaOutputFile, tmplDir); err != nil {
			return err
		}
	}

	return nil
}
