package xsd2go

import (
	"fmt"
	"path"
	"strings"

	"github.com/moov-io/xsd2go/pkg/template"
	"github.com/moov-io/xsd2go/pkg/xsd"
)

func Convert(xsdPath, goModule, outputDir string, xmlnsOverrides []string, templateName string, outputFile string) error {
	if strings.Contains(templateName, ".tgo") {
		outputFile = path.Base(templateName)
		outputFile = strings.ReplaceAll(outputFile, ".tgo", ".go")
	} else if strings.Contains(templateName, ".go.tmpl") {
		outputFile = path.Base(templateName)
		outputFile = strings.ReplaceAll(outputFile, ".go.tmpl", ".go")
	}

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

		if err := template.GenerateTypes(sch, outputDir, outputFile, templateName); err != nil {
			return err
		}
	}

	return nil
}
