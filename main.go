package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elemc/msikeyboard/dbusserver"
	"github.com/elemc/msikeyboard/gomsikeyboard"
	"github.com/elemc/msikeyboard/httpserver"
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

	addr      string
	startHTTP bool

	enabelDBus bool
)

func init() {
	colors := strings.Join(gomsikeyboard.GetAllColors(), ", ")
	modes := strings.Join(gomsikeyboard.GetAllModes(), ", ")
	intensities := strings.Join(gomsikeyboard.GetAllIntensities(), ", ")
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
		strings.Join(gomsikeyboard.GetNames(), ", ")))
	flag.StringVar(&addr, "addr", "localhost:9097",
		"set http server address as host:port for listen")
	flag.BoolVar(&startHTTP, "start-http", true, "set start http server or no")
	flag.BoolVar(&enabelDBus, "start-dbus", true, "set start dbus server or non")
}

func main() {
	flag.Parse()
	gomsikeyboard.Init()
	defer gomsikeyboard.Exit()

	led := gomsikeyboard.LEDSetting{}
	var err error

	if theme != "" {
		led, err = gomsikeyboard.GetTheme(theme)
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

	if startHTTP {
		// start http server
		server := httpserver.Server{}
		server.Addr = addr
		defer server.Close()
		if enabelDBus {
			go func() {
				err = server.Start()

			}()
		} else {
			err = server.Start()
		}
		if err != nil {
			log.Printf("EXIT: %s", err)
			os.Exit(1)
		}
	}
	if enabelDBus {
		dbus := dbusserver.DBusServer{}
		err = dbus.Start()
		if err != nil {
			log.Fatalf("Error in running dbus server: %s", err)
		}
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
