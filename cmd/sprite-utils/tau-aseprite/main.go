package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabstv/primen/internal/spriteutils"
	"github.com/gabstv/primen/io/pb"
	"github.com/gabstv/primen/utils/aseprite"
	"github.com/golang/protobuf/proto"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "atlas",
			Usage: "Commands related to building atlases from aseprite sprite sheets",
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
						cli.BoolFlag{
							Name: "strict",
						},
						cli.IntFlag{
							Name:  "max-width",
							Value: 4096,
						},
						cli.IntFlag{
							Name:  "max-height",
							Value: 4096,
						},
						cli.IntFlag{
							Name:  "padding",
							Value: 0,
						},
						cli.StringFlag{
							Name:  "imageout",
							Usage: "Output the atlas image as a png file (optional) eg: atlas.png",
						},
						cli.BoolFlag{
							Name: "verbose, v",
						},
					},
					Action: cmdAtlasBuild,
				},
			},
		},
		cli.Command{
			Name:  "skel",
			Usage: "Generate skeletons",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "atlas-importer",
					ShortName: "ai",
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "workdir, wd",
			Usage: "Change the working directory before running commands",
		},
	}
	app.Before = func(c *cli.Context) error {
		if v := c.GlobalString("workdir"); v != "" {
			if err := os.Chdir(v); err != nil {
				return err
			}
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func vln(c *cli.Context, msg ...interface{}) {
	if c.Bool("verbose") {
		fmt.Println(msg...)
	}
}

func cmdAtlasBuild(c *cli.Context) error {
	isstrict := c.Bool("strict")
	tpl := &aseprite.AtlasImporter{}
	{
		tpln := c.String("template")
		tplb, err := ioutil.ReadFile(tpln)
		if err != nil {
			println("ioutil read template error " + tpln)
			wd, _ := os.Getwd()
			println(wd)
			return err
		}
		if err := json.Unmarshal(tplb, tpl); err != nil {
			println("json.Unmarshal template error " + tpln)
			wd, _ := os.Getwd()
			println(wd)
			return err
		}
		if v := c.StringSlice("sheet"); len(v) > 0 {
			//TODO: handle external sheets
			return errors.New("TODO: handle sheet flag")
		}
		if v := c.String("atlas"); v != "" {
			tpl.Output = v
		}
		if tpl.Output == "" {
			return errors.New("invalid atlas output")
		}
	}
	//
	idm := make(map[image.Image]*spriteutils.RectPackerNode)
	idmr := make(map[*spriteutils.RectPackerNode]image.Image)
	idmsprs := make(map[*spriteutils.RectPackerNode][]aseprite.FrameIO)
	pkr := &spriteutils.BinTreeRectPacker{}
	for si, v := range tpl.SpriteSheets {
		fb, err := ioutil.ReadFile(v)
		if err != nil {
			println("read spritesheet error " + v)
			wd, _ := os.Getwd()
			println(wd)
			return err
		}
		f, err := aseprite.Parse(fb)
		if err != nil {
			println("aseprite parse file err " + v)
			wd, _ := os.Getwd()
			println(wd)
			return err
		}
		if f.GetMetadata().Image == "" {
			if isstrict {
				return errors.New("The aseprite sheet needs the metadata image to work " + v)
			}
			println("The aseprite sheet needs the metadata image to work "+v, "skipping...")
			continue
		}
		imgb, err := ioutil.ReadFile(f.GetMetadata().Image)
		if err != nil {
			println("image " + f.GetMetadata().Image + " not found for " + v)
			if isstrict {
				return err
			}
			continue
		}
		img, _, err := image.Decode(bytes.NewReader(imgb))
		if err != nil {
			println("image " + f.GetMetadata().Image + " decode error for " + v)
			if isstrict {
				return err
			}
			continue
		}

		imbank := make(map[string]image.Image)
		imget := func(sheetname string, rect aseprite.FrameRect) image.Image {
			if vv, ok := imbank[sheetname+"_"+rect.String()]; ok {
				return vv
			}
			im := image.NewRGBA(image.Rect(0, 0, rect.W, rect.H))
			vln(c, "draw.Draw(im, rect.ToRect(), img, image.ZP, draw.Src)", rect.ToRect(), im.Bounds().String())
			//draw.Draw(im, rect.ToRect(), img, image.ZP, draw.Src)
			draw.Draw(im, im.Bounds(), img, rect.ToRect().Min, draw.Src)
			imbank[sheetname+"_"+rect.String()] = im
			return im
		}
		f.Walk(func(i aseprite.FrameInfo) bool {
			if frame, ok := tpl.SpriteWithFilename(i.Filename); ok {
				if frame.SheetName != "" && !strings.Contains(v, frame.SheetName) {
					println(i.Filename, "sheet name mismatch:", frame.SheetName, v)
					return true
				}
				if si != frame.SheetIndex && !strings.Contains(v, frame.SheetName) {
					println(i.Filename, "sheet name/index mismatch:", frame.Filename, si, v)
					return true
				}
				if frame.OutputName == "" {
					return true
				}
				// frame belongs to the correct sheet
				clipim := imget(v, i.Frame)
				if _, ok := idm[clipim]; !ok {
					node := pkr.AddRect(clipim.Bounds())
					idm[clipim] = node
					idmr[node] = clipim
				}
				if node, ok := idm[clipim]; ok {
					slc := idmsprs[node]
					if slc == nil {
						slc = make([]aseprite.FrameIO, 0)
					}
					slc = append(slc, frame)
					idmsprs[node] = slc
				}
			}
			return true
		})
	}
	//TODO: support border/margin etc
	atlases, err := pkr.Pack(context.TODO(), spriteutils.PackerInput{
		MaxWidth:  c.Int("max-width"),
		MaxHeight: c.Int("max-height"),
		Padding:   c.Int("padding"),
	})
	if err != nil {
		return err
	}
	outfdat := &pb.AtlasFile{
		Images:  make([][]byte, 0),
		Filters: make([]pb.ImageFilter, 0),
		Frames:  make(map[string]*pb.Frame),
	}
	for aindex, atlas := range atlases {
		outimg := image.NewRGBA(image.Rect(0, 0, atlas.Width, atlas.Height))
		for _, node := range atlas.Nodes {
			imclip := idmr[node]
			vln(c, "draw.Draw", node.R().String(), "imclip", "0,0", "draw.Over")
			draw.Draw(outimg, node.R(), imclip, image.ZP, draw.Over)
			imframes := idmsprs[node]
			for _, frame := range imframes {
				if _, ok := outfdat.Frames[frame.OutputName]; ok {
					println("warn: duplicated sprite name: " + frame.OutputName)
					if isstrict {
						return errors.New("warn: duplicated sprite name: " + frame.OutputName)
					}
				}
				fff := &pb.Frame{
					X:     uint32(node.X),
					Y:     uint32(node.Y),
					W:     uint32(node.Width),
					H:     uint32(node.Height),
					Image: uint32(aindex),
				}
				outfdat.Frames[frame.OutputName] = fff
				vln(c, "sprite", frame.OutputName, fff.String())
			}
		}
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, outimg); err != nil {
			return err
		}
		outfdat.Images = append(outfdat.Images, buf.Bytes())
		filter := pb.ImageFilter_DEFAULT
		switch tpl.ImageFilter {
		case "nearest", "NEAREST", "NN", "nn", "pixel", "PIXEL":
			filter = pb.ImageFilter_NEAREST
		case "linear", "LINEAR":
			filter = pb.ImageFilter_LINEAR
		}
		outfdat.Filters = append(outfdat.Filters, filter)
		if imout := c.String("imageout"); imout != "" {
			if strings.Contains(imout, "#") {
				imout = strings.ReplaceAll(imout, "#", strconv.Itoa(aindex))
			}
			if err := ioutil.WriteFile(imout, buf.Bytes(), 0644); err != nil {
				println("error wtriting image " + imout + " " + err.Error())
				if isstrict {
					return err
				}
			}
		}
	}
	fb, err := proto.Marshal(outfdat)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tpl.Output, fb, 0644)
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
		outtpl.SpriteSheets = append(outtpl.SpriteSheets, v)
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
		outp = strings.TrimSuffix(f1, ".json") + "-importer.json"
	}
	outtpl.Output = odat
	//
	b, err := json.MarshalIndent(outtpl, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outp, b, 0644)
}
