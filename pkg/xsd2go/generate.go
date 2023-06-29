package xsd2go

import (
	"fmt"

	"github.com/moov-io/xsd2go/pkg/template"
	"github.com/moov-io/xsd2go/pkg/xsd"
)

func Convert(xsdPath, goModule, outputDir string, xmlnsOverrides []string, outputFile, templatePackage, templateName string) error {
	fmt.Printf("Processing '%s'\n", xsdPath)
	fmt.Printf("Cmd: moovio_xsd2go convert "+
		"%s %s %s "+
		"--xmlns-override=%s "+
		"--output-file=%s "+
		"--template-package=%s "+
		"--template-name=%s "+
		"\n",
		xsdPath,
		goModule,
		outputDir,
		xmlnsOverrides,
		outputFile,
		templatePackage,
		templateName,
	)
	if templatePackage == "" {
		templatePackage = "/pkg/template"
	}
	if templateName == "" {
		templateName = "types.tmpl"
	}

	ws, err := xsd.NewWorkspace(fmt.Sprintf("%s/%s", goModule, outputDir), xsdPath, xmlnsOverrides)
	if err != nil {
		return err
	}

	for _, sch := range ws.Cache {
		if sch.Empty() {
			continue
		}

		var schemaOutputFile = outputFile
		// TODO JB: this is where the original program wants to use "models.go"
		if schemaOutputFile == "" {
			schemaOutputFile = fmt.Sprintf("%s.go", sch.GoPackageName())
		}

		if err := template.GenerateTypes(sch, outputDir, schemaOutputFile, templatePackage, templateName); err != nil {
			return err
		}
	}

	return nil
}
