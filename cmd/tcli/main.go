package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

//go:generate sql2var -I newsystem.go.tpl -O gen.go

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "new",
			ShortName: "n",
			Action:    cmdNew,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "package, p",
					EnvVar: "GOPACKAGE",
				},
				cli.StringFlag{
					Name: "component, c",
				},
				cli.IntFlag{
					Name: "priority",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		if e, ok := err.(*cli.ExitError); ok {
			os.Exit(e.ExitCode())
		}
		os.Exit(1)
	}
}

type tags struct {
	Component string
	Package   string
	Priority  int
}

var tplfns = template.FuncMap{
	"tolower": strings.ToLower,
}

func cmdNew(c *cli.Context) error {
	outd := map[string]interface{}{
		"Tags": tags{
			Component: c.String("component"),
			Package:   c.String("package"),
		},
	}
	tpl := template.Must(template.New("").Funcs(tplfns).Parse(newSystemTpl))
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, outd); err != nil {
		println(err.Error())
		return err
	}

	ffn := ""

	if v := c.Args().First(); v != "" {
		dir, fname := path.Split(v)
		if dir != "" {
			if _, err := os.Stat(dir); err != nil {
				if !os.IsNotExist(err) {
					return cli.NewExitError("could not read dir "+dir, 3)
				}
				if err := os.MkdirAll(path.Clean(dir), 0744); err != nil {
					return cli.NewExitError("could not create dir "+dir, 3)
				}
			}
		}
		if !strings.HasSuffix(fname, ".go") {
			ffn = dir + fname + ".go"
		} else {
			ffn = dir + fname
		}
	} else if envx := os.Getenv("GOFILE"); envx != "" {
		ffn = envx[:len(envx)-len(path.Ext(envx))]
	} else {
		return cli.NewExitError("specify an output file", 5)
	}
	return ioutil.WriteFile(ffn, buf.Bytes(), 0744)
}
