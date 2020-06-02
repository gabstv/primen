package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gabstv/primen/internal/aseprite"
	"github.com/urfave/cli"
)

const (
	flagType        = "type"
	flagSheet       = "sheet"
	flagImageFilter = "image-filter"
	flagMaxWidth    = "max-width"
	flagMaxHeight   = "max-height"
	flagPadding     = "padding"
	flagOutput      = "output"
	flagOverwrite   = "overwrite"
	flagFPS         = "fps"
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
					Name:  flagType + ", t",
					Usage: "Import type: slices | frame_tags | frames",
					Value: "default",
				},
				cli.StringSliceFlag{
					Name:  flagSheet + ", s",
					Usage: "Aseprite sheet file(s) (location of json files)",
				},
				cli.StringFlag{
					Name:  flagImageFilter + ", f",
					Usage: "Image filter: default | linear | nearest (alias: pixel, nn)",
					Value: "default",
				},
				cli.UintFlag{
					Name:  flagMaxWidth,
					Usage: "Max atlas width",
					Value: 4096,
				},
				cli.UintFlag{
					Name:  flagPadding,
					Usage: "Padding between sprites",
					Value: 0,
				},
				cli.StringFlag{
					Name:  flagOutput + ", o",
					Usage: "Output file",
				},
				cli.BoolFlag{
					Name: flagOverwrite,
				},
				cli.IntFlag{
					Name:  flagFPS,
					Usage: "Animations FPS",
					Value: 24,
				},
			},
			Action: cmdTplGen(),
		},
	}
}

func cmdTplGen() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		typestr := aseprite.AtlasImportStrategy(c.String(flagType))
		switch typestr {
		case aseprite.Slices, aseprite.Frames, aseprite.FrameTags, aseprite.Default:
			// ok
		default:
			return errors.New("invalid type")
		}
		inputFiles := c.StringSlice(flagSheet)
		outfile := aseprite.AtlasImporterGroup{
			ImageFilter: c.String(flagImageFilter),
			MaxWidth:    int(c.Uint(flagMaxWidth)),
			MaxHeight:   int(c.Uint(flagMaxHeight)),
			Padding:     int(c.Uint(flagPadding)),
		}
		f, err := getOutput(c, flagOutput, flagOverwrite)
		if err != nil {
			return err
		}
		defer f.Close()
		if len(inputFiles) < 1 {
			// generate generic template
			tpl := aseprite.AtlasImporter{
				ImportStrategy: typestr,
				Frames: []aseprite.FrameIO{
					{
						Filename:   "example.png",
						OutputName: "example",
					},
				},
				AsepriteSheet: "asepritesheet.json",
			}
			switch typestr {
			case aseprite.Default, aseprite.Frames:
				outfile.Templates = []aseprite.AtlasImporter{
					tpl,
				}
				d, _ := json.MarshalIndent(outfile, "", "    ")
				rdr := bytes.NewReader(d)
				if _, err := io.Copy(f, rdr); err != nil {
					return err
				}
				return nil
			}
		}
		outfile.Templates = make([]aseprite.AtlasImporter, 0)
		for _, asefn := range inputFiles {
			asebytes, err := ioutil.ReadFile(asefn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error reading file "+asefn+": "+err.Error())
				continue
			}
			inf, err := aseprite.Parse(asebytes)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error parsing file "+asefn+": "+err.Error())
				continue
			}
			tpl := aseprite.AtlasImporter{
				AsepriteSheet:  asefn,
				Animations:     make([]aseprite.AnimationIO, 0),
				FrameTags:      make([]aseprite.FrameTagIO, 0),
				Frames:         make([]aseprite.FrameIO, 0),
				ImportStrategy: typestr,
				Slices:         make([]aseprite.SliceIO, 0),
			}
			if (typestr == aseprite.Default && len(inf.Meta.Slices) == 0) || typestr != aseprite.Slices {
				for _, v := range inf.Frames {
					tpl.Frames = append(tpl.Frames, aseprite.FrameIO{
						Filename:   v.Filename,
						OutputName: v.Filename,
					})
				}
				for _, v := range inf.Meta.FrameTags {
					tpl.FrameTags = append(tpl.FrameTags, aseprite.FrameTagIO{
						Name:          v.Name,
						OutputPattern: v.Name + "#",
					})
					if v.From != v.To {
						tpl.Animations = append(tpl.Animations, aseprite.AnimationIO{
							FrameTag:   v.Name,
							ClipMode:   string(v.Direction),
							FPS:        c.Int(flagFPS),
							OutputName: v.Name,
						})
					}
				}
			}
			if typestr == aseprite.Default || typestr == aseprite.Slices {
				for _, v := range inf.Meta.Slices {
					tpl.Slices = append(tpl.Slices, aseprite.SliceIO{
						Name:          v.Name,
						OutputPattern: v.Name + "#",
					})
					minframe := -1
					maxframe := -1
					for _, vkey := range v.Keys {
						if minframe == -1 || minframe > vkey.Frame {
							minframe = vkey.Frame
						}
						if vkey.Frame > maxframe {
							maxframe = vkey.Frame
						}
					}
					if minframe != maxframe {
						tpl.Animations = append(tpl.Animations, aseprite.AnimationIO{
							ClipMode:   "forward",
							FPS:        c.Int(flagFPS),
							Slice:      v.Name,
							OutputName: v.Name,
						})
					}
				}
			}
			outfile.Templates = append(outfile.Templates, tpl)
		}
		d, _ := json.MarshalIndent(outfile, "", "    ")
		rdr := bytes.NewReader(d)
		if _, err := io.Copy(f, rdr); err != nil {
			return err
		}
		return nil
	}
}

func getOutput(c *cli.Context, flagname, flagoverwrite string) (w io.WriteCloser, err error) {
	outn := c.String(flagname)
	if outn == "" {
		outn = c.Args().First()
	}
	if outn == "" {
		return os.Stdout, nil
	}
	if fi, _ := os.Stat(outn); fi != nil {
		if fi.IsDir() {
			return nil, errors.New(outn + " is a directory")
		}
		if !c.Bool(flagoverwrite) {
			return nil, errors.New("Output file exists. Use flag --" + flagOverwrite + " to overwrite.")
		}
	}
	return os.Create(outn)
}
