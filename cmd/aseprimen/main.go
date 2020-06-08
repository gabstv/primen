package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gabstv/primen/internal/aseprite"
	"github.com/gabstv/primen/internal/atlaspacker"
	"github.com/golang/protobuf/proto"
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
			ShortName: "t",
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
					Name:  flagMaxHeight,
					Usage: "Max atlas height",
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
				cli.StringFlag{
					Name:  "atlasout",
					Usage: "Atlas output file (eg: atlas.dat)",
				},
			},
			Action: cmdTplGen(),
		},
		{
			Name:      "build",
			ShortName: "b",
			Usage:     "build an atlas using template file(s)",
			Action:    cmdBuild(),
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "margin-left",
					Usage: "Atlas margin left (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "margin-right",
					Usage: "Atlas margin right (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "margin-top",
					Usage: "Atlas margin top (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "margin-bottom",
					Usage: "Atlas margin bottom (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "padding, p",
					Usage: "Atlas padding (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "fixed-width",
					Usage: "Atlas fixed width (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "fixed-height",
					Usage: "Atlas fixed height (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "max-width",
					Usage: "Max atlas width (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "max-height",
					Usage: "Max atlas height (pixels)",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "count",
					Usage: "Max atlas count",
					Value: 0,
				},
				cli.BoolFlag{
					Name: "debug",
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func cmdTplGen() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		inputFiles := c.StringSlice(flagSheet)
		outfile := aseprite.AtlasImporterGroup{
			ImageFilter: c.String(flagImageFilter),
			MaxWidth:    int(c.Uint(flagMaxWidth)),
			MaxHeight:   int(c.Uint(flagMaxHeight)),
			Padding:     int(c.Uint(flagPadding)),
			Output:      c.String("atlasout"),
		}
		f, err := getOutput(c, flagOutput, flagOverwrite)
		if err != nil {
			return err
		}
		defer f.Close()
		if len(inputFiles) < 1 {
			// generate generic template
			tpl := aseprite.AtlasImporter{
				Frames: []aseprite.FrameIO{
					{
						Filename: "sprite1",
						Pivot:    aseprite.Vec2{2, 3},
					},
					{
						Filename: "sprite2",
						Pivot:    aseprite.Vec2{2, 3},
					},
				},
				AsepriteSheet: "asepritesheet.json",
			}
			outfile.Templates = []aseprite.AtlasImporter{
				tpl,
			}
			outfile.Animations = []aseprite.Animation{
				{
					Name: "person",
					Clips: []aseprite.AnimationClip{
						{
							Name:     "idle",
							ClipMode: "loop",
							FPS:      12,
							Frames: []string{
								"sprite1",
								"sprite2",
							},
						},
					},
				},
			}
			d, _ := json.MarshalIndent(outfile, "", "    ")
			rdr := bytes.NewReader(d)
			if _, err := io.Copy(f, rdr); err != nil {
				return err
			}
			return nil
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
				AsepriteSheet:         asefn,
				Frames:                make([]aseprite.FrameIO, 0),
				ExportUndefinedFrames: true,
			}
			for _, v := range inf.Frames {
				tpl.Frames = append(tpl.Frames, aseprite.FrameIO{
					Filename: v.Filename,
				})
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

func cmdBuild() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if !c.Args().Present() {
			return errors.New("no templates files specified")
		}
		for i, tpl := 0, c.Args().Get(0); tpl != ""; i, tpl = i+1, c.Args().Get(i+1) {
			tplb, err := ioutil.ReadFile(tpl)
			if err != nil {
				return fmt.Errorf("error reading template file %w", err)
			}
			g := &aseprite.AtlasImporterGroup{}
			if err := json.Unmarshal(tplb, g); err != nil {
				return fmt.Errorf("error parsing template file %w", err)
			}
			ctx := context.Background()
			src, err := getSources(c, g)
			if err != nil {
				return fmt.Errorf("error reading template file sources %w", err)
			}
			pbfile, err := aseprite.Import(ctx, aseprite.ImportInput{
				Template: g,
				PackerI: atlaspacker.PackerInput{
					MarginLeft:   c.Int("margin-left"),
					MarginRight:  c.Int("margin-right"),
					MarginTop:    c.Int("margin-top"),
					MarginBottom: c.Int("margin-bottom"),
					Padding:      c.Int("padding"),
					FixedWidth:   c.Int("fixed-width"),
					FixedHeight:  c.Int("fixed-height"),
					MaxWidth:     c.Int("max-width"),
					MaxHeight:    c.Int("max-height"),
					Count:        c.Int("count"),
					Debug:        c.Bool("debug"),
				},
				Source: src,
			})
			if err != nil {
				return fmt.Errorf("error creating atlas file for template '%s': %w", tpl, err)
			}
			outn := g.Output
			if outn == "" {
				outn = tpl + ".atlas.dat"
			}
			b, err := proto.Marshal(pbfile)
			if err != nil {
				return fmt.Errorf("error marshalling atlas file for template '%s': %w", tpl, err)
			}
			if err := ioutil.WriteFile(outn, b, 0644); err != nil {
				return fmt.Errorf("error saving atlas file '%s' for template '%s': %w", outn, tpl, err)
			}
		}
		return nil
	}
}

func getSources(c *cli.Context, g *aseprite.AtlasImporterGroup) ([]aseprite.AsepriteInput, error) {
	out := make([]aseprite.AsepriteInput, 0)
	for _, f := range g.Templates {
		asej, err := ioutil.ReadFile(f.AsepriteSheet)
		if err != nil {
			return nil, err
		}
		asef, err := aseprite.Parse(asej)
		if err != nil {
			return nil, err
		}
		inp := aseprite.AsepriteInput{
			Filename:  f.AsepriteSheet,
			FrameData: asef,
		}
		imgb, err := ioutil.ReadFile(asef.Meta.Image)
		if err != nil {
			return nil, err
		}
		inp.ImageData = imgb
		out = append(out, inp)
	}
	return out, nil
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
