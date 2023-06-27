package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gocomply/xsd2go/pkg/xsd2go"
)

var tests = []struct {
	xsdFile           string
	outputFile        string
	goModuleImport    string
	xmlnsOverrides    []string
	expectedOutputDir string
	expectedFiles     []string
}{
	{
		xsdFile:           "complex.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "complex",
		expectedFiles:     []string{"complex.xsd.out"},
	},
	{
		xsdFile:           "cpe-naming_2.3.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "cpe_naming_2_3",
		expectedFiles:     []string{"cpe-naming_2.3.xsd.out"},
	},
	{
		xsdFile:           "restriction.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "restriction",
		expectedFiles:     []string{"restriction.xsd.out"},
	},
	{
		xsdFile:           "simple.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "simple",
		expectedFiles:     []string{"simple.xsd.out"},
	},
	{
		xsdFile:           "simple-8859-1.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "simple_8859_1",
		expectedFiles:     []string{"simple-8859-1.xsd.out"},
	},
	{
		xsdFile:           "swid-2015-extensions-1.0.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "swid_2015_extensions_1_0",
		expectedFiles:     []string{"swid-2015-extensions-1.0.xsd.out"},
	},
	{
		xsdFile:           "xmldsig-core-schema.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "xmldsig_core_schema",
		expectedFiles:     []string{"xmldsig-core-schema.xsd.out"},
	},
	{
		xsdFile:           "incl.xsd",
		outputFile:        "models.go",
		goModuleImport:    "user.com/private",
		expectedOutputDir: "incl",
		expectedFiles:     []string{"incl.xsd.out"},
	},
}

func TestSanity(t *testing.T) {
	var dname = "gen"
	defer os.RemoveAll(dname)

	xsdPath := "xsd-examples/valid/"
	expectedPath := "xsd-examples/valid"

	for indx := range tests {
		t.Run(tests[indx].xsdFile, func(t *testing.T) {
			xsdFile := path.Join(xsdPath, tests[indx].xsdFile)

			err := xsd2go.Convert(xsd2go.Params{
				XsdPath:         xsdFile,
				OutputDir:       dname,
				OutputFile:      tests[indx].outputFile,
				GoModuleImport:  tests[indx].goModuleImport,
				TemplatePackage: "rtp",
				TemplateName:    "element.tmpl",
				XmlnsOverrides:  tests[indx].xmlnsOverrides,
			})
			require.NoError(t, err)

			outputDir := path.Join(dname, tests[indx].expectedOutputDir)
			golangFiles, err := filepath.Glob(outputDir + "/*")
			require.NoError(t, err)
			assert.Equal(t, len(tests[indx].expectedFiles), len(golangFiles), "Expected to find %v generated files in %s but found %v", len(tests[indx].expectedFiles), outputDir, len(golangFiles))

			for indx2 := range tests[indx].expectedFiles {
				if indx2 < len(golangFiles) {
					actual, err := os.ReadFile(golangFiles[indx2])
					require.NoError(t, err)

					expected, err := os.ReadFile(path.Join(expectedPath, tests[indx].expectedFiles[indx2]))
					require.NoError(t, err)

					t.Logf("Comparing %s to %s", golangFiles[indx2], tests[indx].expectedFiles[indx2])
					assert.Equal(t, string(expected), string(actual))
				}
			}
		})
	}
}
