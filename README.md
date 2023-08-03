# XSD2Go - Automatically generate golang xml parser based on XSD
This project is a fork of https://github.com/GoComply/xsd2go.

## Usage
Run this command with variable names `xsdFile`, `name`, `nsPrefix`, and `tmpl`:
```
moovio_xsd2go convert \
   ${xsdFile} \
   ${goModule} \
   ${outputDir} \
   --template-name=/templates/${name}.tmpl \
   --output-file=${output} \
   --xmlns-override='my.namespace=your.namespace'
```
Where the first parameter is the XSD file, the second parameter is the go module and the third parameter is the output folder. The remaining parameters are optional.

## Installation

```
go install github.com/moov-io/xsd2go/cli/moovio_xsd2go
```
