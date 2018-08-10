package gomsikeyboard

/*
#cgo pkg-config: msikeyboard
#include <msikeyboard/msikeyboard.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// SideColorIntensity strcut hold color and intensity for region
type SideColorIntensity struct {
	Color     string
	Intensity string
}

// LEDSetting struct hold 3 regions and mode
type LEDSetting struct {
	Regions map[string]SideColorIntensity
	Mode    string
}

// Init function inits HID API and create new pointer for LEDSetting
func Init() (led *LEDSetting, err error) {
	var result C.int
	result = C.init_msi_keyboard()
	if result < 0 {
		return nil, fmt.Errorf("Error on initialization MSI keyboard")
	}
	led = new(LEDSetting)
	led.Regions = make(map[string]SideColorIntensity)
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
	for region, side := range led.Regions {
		err = side.setColor(region)
		if err != nil {
			return
		}
	}
	if led.Mode != "" {
		err = setMode(led.Mode)
	}

	return
}

// Check function checks names and settings
func (led *LEDSetting) Check() (err error) {
	for region, side := range led.Regions {
		err = side.checkSide(region)
		if err != nil {
			return
		}
	}

	if !checkName(GetAllModes(), led.Mode) {
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
	if !checkName(GetAllColors(), side.Color) {
		return fmt.Errorf("unknown color for %s region %s", name, side.Color)
	}
	if !checkName(GetAllIntensities(), side.Intensity) {
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

func getStringSliceForFunc(list **C.char, size C.size_t) []string {
	var result []string
	slice := (*[1 << 31]*C.char)(unsafe.Pointer(list))[:size:size]

	//for i := 0; i < int(size); i++ {
	for _, l := range slice {
		item := l
		sitem := C.GoString(item)
		result = append(result, sitem)
	}
	C.free(unsafe.Pointer(list))

	return result
}

// GetAllColors function returns all color names
func GetAllColors() []string {
	var size C.size_t
	list := C.get_colors(&size)
	return getStringSliceForFunc(list, size)
}

// GetAllRegions function returns all region names
func GetAllRegions() []string {
	var size C.size_t
	list := C.get_regions(&size)
	return getStringSliceForFunc(list, size)
}

// GetAllModes function returns all mode names
func GetAllModes() []string {
	var size C.size_t
	list := C.get_modes(&size)
	return getStringSliceForFunc(list, size)
}

// GetAllIntensities function returns all intensity names
func GetAllIntensities() []string {
	var size C.size_t
	list := C.get_intensities(&size)
	return getStringSliceForFunc(list, size)
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
