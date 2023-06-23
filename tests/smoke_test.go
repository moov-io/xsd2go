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
	xsdFile         string
	outputDir       string
	outputFile      string
	goModule        string
	goPackage       string
	namespacePrefix string
	expectedFiles   []string
}{
	{
		xsdFile:         "complex.xsd",
		outputDir:       "complex",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "simple_schema",
		namespacePrefix: "complex",
		expectedFiles:   []string{"complex.go.out"},
	},
	{
		xsdFile:         "cpe-naming_2.3.xsd",
		outputDir:       "cpe_naming_2_3",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "simple_schema",
		namespacePrefix: "",
		expectedFiles:   []string{"cpe-naming_2.3.go.out"},
	},
	{
		xsdFile:         "restriction.xsd",
		outputDir:       "restriction",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "simple_schema",
		namespacePrefix: "",
		expectedFiles:   []string{"restriction.go.out"},
	},
	{
		xsdFile:         "simple.xsd",
		outputDir:       "simple",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "simple_schema",
		namespacePrefix: "",
		expectedFiles:   []string{"simple.go.out"},
	},
	{
		xsdFile:         "simple-8859-1.xsd",
		outputDir:       "simple_8859_1",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "simple_schema",
		namespacePrefix: "",
		expectedFiles:   []string{"simple-8859-1.go.out"},
	},
	{
		xsdFile:         "swid-2015-extensions-1.0.xsd",
		outputDir:       "swid_2015_extensions_1_0",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "swid_2015_extensions_1_0",
		namespacePrefix: "",
		expectedFiles:   []string{"swid-2015-extensions-1.0.go.out"},
	},
	{
		xsdFile:         "xmldsig-core-schema.xsd",
		outputDir:       "xmldsig_core_schema",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "xml_dsig",
		namespacePrefix: "",
		expectedFiles:   []string{"xmldsig-core-schema.go.out"},
	},
	{
		xsdFile:         "incl.xsd",
		outputDir:       "incl",
		outputFile:      "models.go",
		goModule:        "user.com/private",
		goPackage:       "incl",
		namespacePrefix: "",
		expectedFiles:   []string{"incl.go.out"},
	},
}

func TestSanity(t *testing.T) {
	dname, err := os.MkdirTemp("", "xsd2go_tests_")
	assert.Nil(t, err)
	defer os.RemoveAll(dname)

	xsdPath := "xsd-examples/xsd/"
	expectedPath := "xsd-examples/assertions"

	for indx := range tests {
		t.Run(tests[indx].xsdFile, func(t *testing.T) {
			outputDir := path.Join(dname, tests[indx].outputDir)
			xsdFile := path.Join(xsdPath, tests[indx].xsdFile)

			err = xsd2go.Convert(
				xsdFile,
				outputDir,
				tests[indx].outputFile,
				tests[indx].goModule,
				tests[indx].goPackage,
				tests[indx].namespacePrefix,
				"rtp",
			)
			require.NoError(t, err)

			golangFiles, err := filepath.Glob(outputDir + "/*/*")
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
