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

	mode        string
	theme       string
	flagregions map[string]*string
	left        string
	middle      string
	right       string
	all         string

	addr      string
	startHTTP bool

	enabelDBus bool
)

func init() {
	flagregions = make(map[string]*string)
	colors := strings.Join(gomsikeyboard.GetAllColors(), ", ")
	modes := strings.Join(gomsikeyboard.GetAllModes(), ", ")
	intensities := strings.Join(gomsikeyboard.GetAllIntensities(), ", ")

	regions := gomsikeyboard.GetAllRegions()
	for _, region := range regions {
		flagregions[region] = new(string)
		flag.StringVar(flagregions[region], region, "",
			fmt.Sprintf("color:intensity for %s keyboard region (colors: %s) (intensities: %s)", region, colors, intensities))
	}

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
		side := gomsikeyboard.SideColorIntensity{}
		side.Color, side.Intensity = getColorIntensity(all)
		for _, region := range gomsikeyboard.GetAllRegions() {
			led.Regions[region] = side
		}
	} else {
		for region, params := range flagregions {
			if *params == "" {
				continue
			}
			side := gomsikeyboard.SideColorIntensity{}
			side.Color, side.Intensity = getColorIntensity(all)
			led.Regions[region] = side
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
