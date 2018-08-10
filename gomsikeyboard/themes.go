package gomsikeyboard

import (
	"fmt"
	"log"
)

// Themes map of Theme
type Themes map[string]*LEDSetting

var (
	defaultThemes Themes
)

func init() {
	defaultThemes = make(Themes)
	cool, err := Init()
	if err != nil {
		log.Printf("Unable to initialize default theme: %s", err)
		return
	}
	cool.Mode = "normal"
	green := SideColorIntensity{Color: "green"}
	yellow := SideColorIntensity{Color: "green"}
	cool.Regions["left"] = green
	cool.Regions["middle"] = yellow
	cool.Regions["right"] = yellow
	defaultThemes["cool"] = cool
}

// GetNames function get names of default themes
func GetNames() (result []string) {
	for k := range defaultThemes {
		result = append(result, k)
	}
	return
}

// GetTheme function for return theme from default themes by name
func GetTheme(name string) (theme LEDSetting, err error) {
	theme, ok := defaultThemes[name]
	if !ok {
		return theme, fmt.Errorf("theme with giving name %s not found", name)
	}
	return
}
