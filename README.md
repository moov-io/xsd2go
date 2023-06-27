# XSD2Go - Automatically generate golang xml parser based on XSD
This project is a fork of https://github.com/GoComply/xsd2go.

## Usage
Run this command with variable names `xsdFile`, `name`, `nsPrefix`, and `tmpl`:
```
gocomply_xsd2go convert \
   --xsd-file=${xsdFile} \
   --go-module-import=${module} \
   --output-dir=gen/${name} \
   --output-file=${name}.go \
   --template-package=${tmpl} \
   --template-name=${name}.tmpl \
   --xmlns-override='my.namespace=your.namespace'
