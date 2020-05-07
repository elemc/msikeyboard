package gomsikeyboard

import "testing"

func TestLEDSetting_Set(t *testing.T) {
	keyboard, err := Init()
	if err != nil {
		t.Fatalf("Unable to init keyboard: %s", err)
	}
	keyboard.Mode = "normal"
	keyboard.Regions["left"] = SideColorIntensity{
		Color:     "green",
		Intensity: "low",
	}
	keyboard.Regions["middle"] = SideColorIntensity{
		Color:     "purple",
		Intensity: "high",
	}
	keyboard.Regions["right"] = SideColorIntensity{
		Color:     "yellow",
		Intensity: "high",
	}
	/*keyboard.Regions["front_left"] = SideColorIntensity{
		Color:     "white",
		Intensity: "low",
	}
	keyboard.Regions["front_right"] = SideColorIntensity{
		Color:     "white",
		Intensity: "low",
	}*/
	err = keyboard.Set()
	if err != nil {
		t.Fatalf("Unable to set LED settings: %s", err)
	}
}
