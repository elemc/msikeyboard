package msikeyboard

/*
//////#cgo pkg-config: msikeyboard
#cgo LDFLAGS: -lmsikeyboard -L/usr/local/lib
#cgo CFLAGS: -I/usr/local/include
#include <msikeyboard/msikeyboard.h>
#include <stdlib.h>
*/
import "C"
import "fmt"

var (
	colors      = [...]string{"off", "red", "orange", "yellow", "green", "sky", "blue", "purple", "white"}
	modes       = [...]string{"normal", "gaming", "breathe", "demo", "wave"}
	intensities = [...]string{"high", "medium", "low", "light"}
)

// SideColorIntensity strcut hold color and intensity for region
type SideColorIntensity struct {
	Color     string
	Intensity string
}

// LEDSetting struct hold 3 regions and mode
type LEDSetting struct {
	Left   SideColorIntensity
	Middle SideColorIntensity
	Right  SideColorIntensity
	Mode   string
}

// Init function inits HID API and create new pointer for LEDSetting
func Init() (led *LEDSetting, err error) {
	var result C.int
	result = C.init_msi_keyboard()
	if result < 0 {
		return nil, fmt.Errorf("Error on initialization MSI keyboard")
	}
	led = new(LEDSetting)
	return
}

// Exit function free and destroy HID API
func Exit() (err error) {
	result := C.free_msi_keyboard()
	if result < 0 {
		return fmt.Errorf("Error on free MSI keyboard")
	}
	return
}

// Set function sets settings for keyboard
func (led *LEDSetting) Set() (err error) {
	err = led.Left.setColor("left")
	if err != nil {
		return
	}
	err = led.Middle.setColor("middle")
	if err != nil {
		return
	}
	err = led.Right.setColor("right")
	if err != nil {
		return
	}

	if led.Mode != "" {
		err = setMode(led.Mode)
	}

	return
}

// Check function checks names and settings
func (led *LEDSetting) Check() (err error) {
	err = led.Left.checkSide("left")
	if err != nil {
		return
	}
	err = led.Middle.checkSide("middle")
	if err != nil {
		return
	}
	err = led.Right.checkSide("right")
	if err != nil {
		return
	}

	if !checkName(modes[:], led.Mode) {
		return fmt.Errorf("unknown mode %s", led.Mode)
	}

	return
}

func (side *SideColorIntensity) setColor(region string) (err error) {
	if side.Color == "" {
		return
	}
	err = setColor(region, side.Color, side.Intensity)
	return
}

func (side *SideColorIntensity) checkSide(name string) error {
	if !checkName(colors[:], side.Color) {
		return fmt.Errorf("unknown color for %s region %s", name, side.Color)
	}
	if !checkName(intensities[:], side.Intensity) {
		return fmt.Errorf("unknown intensity for %s region %s", name, side.Intensity)
	}
	return nil
}

func checkName(list []string, name string) bool {
	if name == "" {
		return true
	}
	found := false
	for _, l := range list {
		if l == name {
			found = true
			break
		}
	}
	return found
}

// GetAllColors function returns all color names
func GetAllColors() []string {
	return colors[:]
}

// GetAllModes function returns all mode names
func GetAllModes() []string {
	return modes[:]
}

// GetAllIntensities function returns all intensity names
func GetAllIntensities() []string {
	return intensities[:]
}

//// cgo methods
func setColor(region, color, intensity string) (err error) {
	result := C.set_color_by_names(C.CString(region), C.CString(color), C.CString(intensity))
	if result < 0 {
		err = fmt.Errorf("error in set color (%s, %s, %s)", region, color, intensity)
	}
	return
}

func setMode(mode string) (err error) {
	result := C.set_mode_by_name(C.CString(mode))
	if result < 0 {
		err = fmt.Errorf("error in set mode %s", mode)
	}
	return
}
