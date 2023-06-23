package xsd

import (
	"fmt"
	"os"
	"path/filepath"
)

type Workspace struct {
	Cache             map[string]*Schema
	GoModulesPath     string
	templateOverrides map[string]Override
}

func NewWorkspace(goModulesPath string, goPackage string, nsPrefix string, xsdFile string, templates map[string]Override) (*Workspace, error) {
	ws := Workspace{
		Cache:             map[string]*Schema{},
		GoModulesPath:     goModulesPath,
		templateOverrides: templates,
	}
	var err error
	if err != nil {
		return nil, err
	}

	_, err = ws.loadXsd(goPackage, nsPrefix, xsdFile, true)
	if err != nil {
		return nil, err
	}
	return &ws, ws.compile()
}

func (ws *Workspace) loadXsd(goPackage string, nsPrefix string, xsdPath string, cache bool) (*Schema, error) {
	cached, found := ws.Cache[xsdPath]
	if found {
		return cached, nil
	}
	fmt.Println("\tParsing:", xsdPath)

	xsdPathClean := filepath.Clean(xsdPath)
	f, err := os.Open(xsdPathClean)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	schema, err := parseSchema(f)
	if err != nil {
		return nil, err
	}

	schema.NsPrefix = nsPrefix
	schema.ModulesPath = ws.GoModulesPath
	schema.filePath = xsdPath
	schema.TemplateOverrides = ws.templateOverrides
	schema.goPackageNameOverride = goPackage
	// Won't cache included schemas - we need to append contents to the current schema.
	if cache {
		ws.Cache[xsdPath] = schema
	}

	dir := filepath.Dir(xsdPath)

	for idx := range schema.Includes {
		si := schema.Includes[idx]
		includeNsPrefix := schema.xmlnsPrefixByXmlns(schema.Includes[idx].Namespace)
		if err = si.load(ws, includeNsPrefix, goPackage, dir); err != nil {
			return nil, err
		}

		isch := si.IncludedSchema
		schema.Imports = append(isch.Imports, schema.Imports...)
		schema.Elements = append(isch.Elements, schema.Elements...)
		schema.Attributes = append(isch.Attributes, schema.Attributes...)
		schema.AttributeGroups = append(isch.AttributeGroups, schema.AttributeGroups...)
		schema.ComplexTypes = append(isch.ComplexTypes, schema.ComplexTypes...)
		schema.SimpleTypes = append(isch.SimpleTypes, schema.SimpleTypes...)
		schema.inlinedElements = append(isch.inlinedElements, schema.inlinedElements...)
		for key, sch := range isch.importedModules {
			schema.importedModules[key] = sch
		}
	}

	for idx := range schema.Imports {
		importNsPrefix := schema.xmlnsPrefixByXmlns(schema.Imports[idx].Namespace)
		if err = schema.Imports[idx].load(ws, importNsPrefix, goPackage, dir); err != nil {
			return nil, err
		}
	}
	schema.compile()
	return schema, nil
}

func (ws *Workspace) compile() error {
	uniqPkgNames := map[string]string{}

	for _, schema := range ws.Cache {
		goPackageName := schema.GoPackageName()
		prevXmlns, ok := uniqPkgNames[goPackageName]
		if ok {
			return fmt.Errorf("Malformed workspace. Multiple XSD files refer to itself with xmlns shorthand: '%s':\n - %s\n - %s\nWhile this is xsd in XSD it is impractical for golang code generation.\nConsider providing --xmlns-override=%s=mygopackage", goPackageName, prevXmlns, schema.TargetNamespace, schema.TargetNamespace)
		}
		uniqPkgNames[goPackageName] = schema.TargetNamespace
	}

	return nil
}
