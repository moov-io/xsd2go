# XSD2Go - Automatically generate golang xml parser based on XSD
This project is a fork of https://github.com/GoComply/xsd2go.

## Usage
Run this command with variable names `xsdFile`, `name`, `nsPrefix`, and `tmpl`:
```
moovio_xsd2go convert \
   --xsd-file=${xsdFile} \
   --output-dir=gen/${name} \
   --output-file=${name}.go \
   --go-package=${name} \
   --namespace-prefix=${nsPrefix} \
   --template-package=${tmpl}
```
