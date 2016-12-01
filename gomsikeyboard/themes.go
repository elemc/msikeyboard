package gomsikeyboard

import "fmt"

// Themes map of Theme
type Themes map[string]LEDSetting

var (
	defaultThemes Themes
)

func init() {
	defaultThemes = make(Themes)
	cool := LEDSetting{}
	cool.Left.Color = "green"
	cool.Middle.Color = "yellow"
	cool.Right.Color = "yellow"
	cool.Mode = "normal"
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
