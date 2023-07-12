package xsd2go

import (
	"fmt"

	"github.com/moov-io/xsd2go/pkg/template"
	"github.com/moov-io/xsd2go/pkg/xsd"
)

func Convert(xsdPath, goModule, outputDir string, xmlnsOverrides []string, templateName string, outputFile string) error {
	fmt.Printf("Processing '%s'\n", xsdPath)
	fmt.Printf("Cmd: moovio_xsd2go convert "+
		"%s %s %s --xmlns-override=%s --template-name=%s --output-file=%s \n",
		xsdPath,
		goModule,
		outputDir,
		xmlnsOverrides,
		templateName,
		outputFile,
	)

	ws, err := xsd.NewWorkspace(fmt.Sprintf("%s/%s", goModule, outputDir), xsdPath, xmlnsOverrides)
	if err != nil {
		return err
	}

	for _, sch := range ws.Cache {
		if sch.Empty() {
			continue
		}

		// TODO JB
		// Set the name of the generated file. This is where the original program wants to use "models.go".
		// This is for imported XSDs which generate to their own file, so a parameter can't work here.
		var generatedOutputFile = fmt.Sprintf("%s.go", sch.GoPackageName())
		if outputFile != "" {
			generatedOutputFile = outputFile
		}

		if err := template.GenerateTypes(sch, outputDir, generatedOutputFile, templateName); err != nil {
			return err
		}
	}

	return nil
}
