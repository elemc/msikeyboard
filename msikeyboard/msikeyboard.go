package msikeyboard

import (
	"fmt"
	"log"

	"github.com/GeertJohan/go.hid"
)

// Device is a main struct
type Device struct {
	d            *hid.Device
	CurrentTheme Theme
	Intensity    string
	Mode         string
}

// Codes type for match number by name
type Codes map[string]int

const (
	// VID is a Vendor ID for MSI keyboard
	VID = 0x1770
	// PID is a Product ID for MSI keyboards
	PID = 0xff00
)

var (
	// Regions is a regions map
	Regions Codes
	// Colors is a colors map
	Colors Codes
	// Levels is a levels map
	Levels Codes
	// Modes is a modes map
	Modes Codes
)

func init() {
	Regions = make(Codes)
	Regions["left"] = 1
	Regions["middle"] = 2
	Regions["right"] = 3

	Colors = make(Codes)
	Colors["black"] = 0
	Colors["red"] = 1
	Colors["orange"] = 2
	Colors["yellow"] = 3
	Colors["green"] = 4
	Colors["cyan"] = 5
	Colors["blue"] = 6
	Colors["purple"] = 7
	Colors["white"] = 8

	Levels = make(Codes)
	Levels["light"] = 3
	Levels["low"] = 2
	Levels["med"] = 1
	Levels["high"] = 0

	Modes = make(Codes)
	Modes["normal"] = 1
	Modes["gaming"] = 2
	Modes["breathe"] = 3
	Modes["demo"] = 4
	Modes["wave"] = 5
}

func getCodesKeys(data Codes) (result []string) {
	for k := range data {
		result = append(result, k)
	}
	return result
}

// GetAllColors function returns color list
func GetAllColors() []string {
	return getCodesKeys(Colors)
}

// GetAllModes function returns all modes
func GetAllModes() []string {
	return getCodesKeys(Modes)
}

// GetDevice function getting keyboard device pointer
func GetDevice() (device *Device, err error) {
	d := new(Device)

	dil, err := hid.Enumerate(VID, PID)
	if err != nil {
		return nil, err
	}
	if len(dil) == 0 {
		return nil, fmt.Errorf("device %4x:%4x not found", VID, PID)
	}
	di := dil[0]

	d.d, err = di.Device()
	if err != nil {
		return nil, err
	}

	log.Printf("Path: %s", di.Path)
	log.Printf("Vendor ID: %x", di.VendorId)
	log.Printf("Product ID: %x", di.ProductId)
	log.Printf("Serial number: %x", di.SerialNumber)
	log.Printf("Manufacturer: %s", di.Manufacturer)
	log.Printf("Product: %s", di.Product)

	return
}

// Close function close a device
func (d *Device) Close() {
	if d.d != nil {
		d.d.Close()
	}
}

// Set functuion sets mode and colors
func (d *Device) Set() (err error) {

	err = d.SetColor("left", d.CurrentTheme.Left.ColorName, d.Intensity)
	if err != nil {
		return err
	}
	err = d.SetColor("middle", d.CurrentTheme.Middle.ColorName, d.Intensity)
	if err != nil {
		return err
	}
	err = d.SetColor("right", d.CurrentTheme.Right.ColorName, d.Intensity)
	if err != nil {
		return err
	}

	err = d.SetMode()
	return
}

// SetColor function sets color to region with color and intensity
func (d *Device) SetColor(region, color, intensity string) (err error) {
	iRegion, ok := Regions[region]
	if !ok {
		return fmt.Errorf("Unknown region: %s", region)
	}
	iColor, ok := Colors[color]
	if !ok {
		return fmt.Errorf("Unknown color: %s", color)
	}
	iIntensity, ok := Levels[intensity]
	if !ok {
		return fmt.Errorf("Unknown intensity: %s", region)
	}

	data := make([]byte, 8)
	data[0] = 1
	data[1] = 2
	data[2] = 66
	data[3] = byte(iRegion)
	data[4] = byte(iColor)
	data[5] = byte(iIntensity)
	data[6] = 0
	data[7] = 236
	_, err = d.d.SendFeatureReport(data)

	return
}

// SetMode function sets mode to keyboard
func (d *Device) SetMode() (err error) {
	period := 2

	err = d.setMode(0, d.CurrentTheme.Left, period)
	if err != nil {
		return err
	}
	err = d.setMode(3, d.CurrentTheme.Left, period)
	if err != nil {
		return err
	}
	err = d.setMode(6, d.CurrentTheme.Left, period)
	if err != nil {
		return err
	}

	err = d.commit()

	return
}

func (d *Device) setMode(begin int, r Region, period int) (err error) {
	err = d.sendData(begin+1, Colors[r.ColorName], Levels[d.Intensity], 0)
	if err != nil {
		return err
	}
	err = d.sendData(begin+2, Colors[r.SecondaryName], Levels[d.Intensity], 0)
	if err != nil {
		return err
	}
	err = d.sendData(begin+3, period, period, period)
	if err != nil {
		return err
	}
	return
}

func (d *Device) sendData(section, color, intensity, period int) (err error) {
	data := make([]byte, 8)
	data[0] = 1
	data[1] = 2
	data[2] = 67
	data[3] = byte(section)
	data[4] = byte(color)
	data[5] = byte(intensity)
	data[6] = byte(period)
	data[7] = 236

	_, err = d.d.SendFeatureReport(data)
	return
}

func (d *Device) commit() (err error) {
	data := make([]byte, 8)
	data[0] = 1
	data[1] = 2
	data[2] = 65
	data[3] = byte(Modes[d.Mode])
	data[4] = 0
	data[5] = 0
	data[6] = 0
	data[7] = 236

	_, err = d.d.SendFeatureReport(data)
	return

}
