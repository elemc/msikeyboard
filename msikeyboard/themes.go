package msikeyboard

import "fmt"

// Region struct set primary and secondary colors
type Region struct {
	ColorName     string
	SecondaryName string
}

// Theme is a struct for theme
type Theme struct {
	Name   string
	Left   Region
	Middle Region
	Right  Region
}

// Themes map of Theme
type Themes map[string]Theme

var (
	defaultThemes Themes
)

func init() {
	defaultThemes = make(Themes)
	cool := Theme{}
	cool.Name = "cool"
	cool.Left.ColorName = "green"
	cool.Middle.ColorName = "yellow"
	cool.Right.ColorName = "yellow"
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
func GetTheme(name string) (theme Theme, err error) {
	theme, ok := defaultThemes[name]
	if !ok {
		return theme, fmt.Errorf("theme with giving name %s not found", name)
	}
	return
}
