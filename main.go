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
	blackColor            = "off"
	defaultColorIntensity = "red:high"
)

var (
	// backlight bool
	// color     string
	// blink     int
	// intensity string

	mode   string
	theme  string
	left   string
	middle string
	right  string
	all    string
)

func init() {
	colors := strings.Join(msikeyboard.GetAllColors(), ", ")
	modes := strings.Join(msikeyboard.GetAllModes(), ", ")
	intensities := strings.Join(msikeyboard.GetAllIntensities(), ", ")
	flag.StringVar(&left, "left", "",
		fmt.Sprintf("color:intensity for left keyboard region (colors: %s) (intensities: %s)", colors, intensities))
	flag.StringVar(&middle, "middle", "",
		fmt.Sprintf("color:intensity for middle keyboard region (colors: %s) (intensities: %s)", colors, intensities))
	flag.StringVar(&right, "right", "",
		fmt.Sprintf("color:intensity for right keyboard region (colors: %s) (intensities: %s)", colors, intensities))
	flag.StringVar(&all, "all", "",
		fmt.Sprintf("color:intensity for all keyboard regions (colors: %s) (intensities: %s)", colors, intensities))
	flag.StringVar(&mode, "mode", "", fmt.Sprintf("set mode: %s", modes))
	flag.StringVar(&theme, "theme", "", fmt.Sprintf("set theme by name: %s",
		strings.Join(msikeyboard.GetNames(), ", ")))
}

func main() {
	flag.Parse()
	msikeyboard.Init()
	defer msikeyboard.Exit()

	led := msikeyboard.LEDSetting{}
	var err error

	if theme != "" {
		led, err = msikeyboard.GetTheme(theme)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	if all != "" {
		color, intensity := getColorIntensity(all)
		led.Left.Color = color
		led.Left.Intensity = intensity
		led.Middle.Color = color
		led.Middle.Intensity = intensity
		led.Right.Color = color
		led.Right.Intensity = intensity
	} else {
		if left != "" {
			led.Left.Color, led.Left.Intensity = getColorIntensity(left)
		}
		if middle != "" {
			led.Middle.Color, led.Middle.Intensity = getColorIntensity(middle)
		}
		if right != "" {
			led.Right.Color, led.Right.Intensity = getColorIntensity(right)
		}
	}
	if mode != "" {
		led.Mode = mode
	}

	err = led.Check()
	if err != nil {
		log.Fatalf("Error in parsing parameters: %s", err)
	}

	err = led.Set()
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	os.Exit(0)
}

func getColorIntensity(arg string) (color, intensity string) {
	paramList := strings.Split(arg, ":")
	if len(paramList) < 2 {
		color = arg
		intensity = "high"
		return
	}

	color = paramList[0]
	intensity = paramList[1]
	return
}
