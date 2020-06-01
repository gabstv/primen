package main

import (
	"errors"

	"github.com/gabstv/primen/internal/aseprite"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "aseprimen"
	app.Authors = []cli.Author{
		{
			Name:  "Gabriel Ochsenhofer",
			Email: "gabriel.ochsenhofer <*at*> gmail [*dot*] com",
		},
	}
	app.Copyright = "2020 Gabriel Ochsenhofer"
	app.Description = "A set of utilities to import content from Aseprite."

	app.Commands = []cli.Command{
		{
			Name:      "tplgen",
			ShortName: "gen",
			Usage:     "Generate an import template to import Aseprite sheets",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type, t",
					Usage: "Import type: slices | frame_tags | frames",
					Value: "default",
				},
				cli.StringSliceFlag{
					Name:  "input, i",
					Usage: "Aseprite input file(s) (location of json files)",
				},
			},
			Action: cmdTplGen(),
		},
	}
}

func cmdTplGen() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		typestr := aseprite.AtlasImportStrategy(c.String("type"))
		switch typestr {
		case aseprite.Slices, aseprite.Frames, aseprite.FrameTags, aseprite.Default:
			// ok
		default:
			return errors.New("invalid type")
		}
		inputFiles := c.StringSlice("input")
		outfile := aseprite.AtlasImporterGroup{}
		println(inputFiles[0])
		println(outfile.Output)
		return nil
	}
}
