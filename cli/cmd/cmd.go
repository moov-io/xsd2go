package cmd

import (
	"github.com/gocomply/xsd2go/pkg/xsd2go"
	"github.com/urfave/cli"
	"os"
)

// Execute ...
func Execute() error {
	app := cli.NewApp()
	app.Name = "GoComply XSD2Go"
	app.Usage = "Automatically generate golang xml parser based on XSD"
	app.Commands = []cli.Command{
		convert,
	}

	return app.Run(os.Args)
}

var convert = cli.Command{
	Name:      "convert",
	Usage:     "convert XSD to golang code to parse xml files generated by given xsd",
	ArgsUsage: "XSD-FILE GO-MODULE-IMPORT OUTPUT-DIR",
	Before: func(c *cli.Context) error {
		if c.NArg() != 3 {
			return cli.NewExitError("Exactly 3 arguments are required", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		metaschemaDir, goModule, outputDir := c.Args()[0], c.Args()[1], c.Args()[2]
		err := xsd2go.Convert(metaschemaDir, goModule, outputDir)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},

}
