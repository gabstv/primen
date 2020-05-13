package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabstv/tau/utils/aseprite"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cli.Command{
			Name: "atlas",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "build",
					ShortName: "b",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "template, t",
							Usage:    "Template file used to generate the atlas (use skel to generate)",
							Required: true,
						},
						cli.StringSliceFlag{
							Name:  "sheet, s",
							Usage: "Provide extra sprite sheet files (aseprite json)",
						},
						cli.StringFlag{
							Name:  "atlas, o",
							Usage: "atlas output file (xyz.dat)",
						},
					},
					Action: cmdAtlasBuild,
				},
			},
		},
		cli.Command{
			Name: "skel",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "atlas-importer",
					ShortName: "atlasi",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "output, o",
						},
						cli.StringSliceFlag{
							Name:  "input, i",
							Usage: "Input Aseprite Sprite Sheet (JSON)",
						},
						cli.StringFlag{
							Name:  "atlasout",
							Usage: "Atlas output path",
							Value: "atlas.dat",
						},
					},
					Action: cmdAtlasSkel,
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func cmdAtlasBuild(c *cli.Context) error {
	//tpln := c.String("template")
	return errors.New("TODO")
}

func cmdAtlasSkel(c *cli.Context) error {
	outp := c.String("output")
	inpt := c.StringSlice("input")
	odat := c.String("atlasout")
	if len(inpt) == 0 {
		outtpl := aseprite.AtlasImporter{
			Sprites: []aseprite.FrameIO{
				aseprite.FrameIO{
					Filename:   "input name 1",
					OutputName: "sprite name 1",
					SheetIndex: 0,
				},
				aseprite.FrameIO{
					Filename:   "input name 2",
					OutputName: "sprite name 2",
					SheetIndex: 0,
				},
			},
			SpriteSheets: []string{
				"spritesheet.json",
			},
			Output: odat,
		}
		of := "atlas-importer.json"
		if outp != "" {
			of = outp
		}
		b, err := json.MarshalIndent(outtpl, "", "    ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(of, b, 0644)
	}
	outtpl := &aseprite.AtlasImporter{
		Sprites:      make([]aseprite.FrameIO, 0),
		SpriteSheets: make([]string, 0),
	}
	for _, v := range inpt {
		d, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		f, err := aseprite.Parse(d)
		if err != nil {
			return nil
		}
		xi := len(outtpl.SpriteSheets)
		outtpl.SpriteSheets = append(outtpl.SpriteSheets, f.GetMetadata().Image)
		f.Walk(func(i aseprite.FrameInfo) bool {
			outtpl.Sprites = append(outtpl.Sprites, aseprite.FrameIO{
				Filename:   i.Filename,
				SheetIndex: xi,
				SheetName:  v,
				OutputName: i.Filename,
			})
			return true
		})
	}
	if outp == "" {
		_, f1 := filepath.Split(inpt[0])
		outp = strings.TrimSuffix(".json") + "-importer.json"
	}
	outtpl.Output = odat
	//
	b, err := json.MarshalIndent(outtpl, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outp, b, 0644)
}
