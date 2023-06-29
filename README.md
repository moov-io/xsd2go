# XSD2Go - Automatically generate golang xml parser based on XSD
This project is a fork of https://github.com/GoComply/xsd2go.

## Usage
Run this command with variable names `xsdFile`, `name`, `nsPrefix`, and `tmpl`:
```
moovio_xsd2go convert \
   ${xsdFile} \
   ${module} \
   gen \
   --output-file=${name}.go \
   --template-package=${tmpl} \
   --template-name=${name}.tmpl \
   --xmlns-override='my.namespace=your.namespace'
```
Where the first parameter is the XSD file, the second parameter is the go module and the third parameter is the output folder. The remaining parameters are optional.
