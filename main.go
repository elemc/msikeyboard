package main

import (
	"flag"
	"fmt"
	"gomsikeyboard/msikeyboard"
	"log"
	"os"
	"strings"
)

const (
	blackColor = "black"
)

var (
	backlight bool
	color     string
	blink     int
	intensity string
	theme     string
	mode      string
)

func init() {
	flag.BoolVar(&backlight, "backlight", true, "enable backlight (set to true) or disable (set to false)")
	flag.StringVar(&color, "color", "white", fmt.Sprintf("set color: %s", strings.Join(msikeyboard.GetAllColors(), ", ")))
	//flag.IntVar(&blink, "blink", 0, "blink N ms")
	flag.StringVar(&intensity, "intensity", "high", "set intensity applies to all regions")
	flag.StringVar(&theme, "theme", "", fmt.Sprintf("set theme by name: %s", strings.Join(msikeyboard.GetNames(), ", ")))
	flag.StringVar(&mode, "mode", "normal", fmt.Sprintf("set mode: %s", strings.Join(msikeyboard.GetAllModes(), ", ")))
}

func main() {
	flag.Parse()
	dev, err := msikeyboard.GetDevice()
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer dev.Close()

	if _, ok := msikeyboard.Colors[color]; !ok {
		log.Fatalf("Color with name %s not found", color)
	}

	if theme != "" {
		dev.CurrentTheme, err = msikeyboard.GetTheme(theme)
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else if backlight {
		dev.CurrentTheme.Left.ColorName = color
		dev.CurrentTheme.Left.SecondaryName = blackColor
		dev.CurrentTheme.Middle.ColorName = color
		dev.CurrentTheme.Middle.SecondaryName = blackColor
		dev.CurrentTheme.Right.ColorName = color
		dev.CurrentTheme.Right.SecondaryName = blackColor
	} else {
		dev.CurrentTheme.Left.ColorName = blackColor
		dev.CurrentTheme.Left.SecondaryName = blackColor
		dev.CurrentTheme.Middle.ColorName = blackColor
		dev.CurrentTheme.Middle.SecondaryName = blackColor
		dev.CurrentTheme.Right.ColorName = blackColor
		dev.CurrentTheme.Right.SecondaryName = blackColor
	}
	dev.Intensity = intensity
	dev.Mode = mode

	err = dev.Set()
	if err != nil {
		log.Fatalf(err.Error())
	}
	os.Exit(0)
}
