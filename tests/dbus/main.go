package main

import "github.com/godbus/dbus"

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	obj := conn.Object("name.elemc.msikeyboard", "/name/elemc/msikeyboard")
	call := obj.Call("name.elemc.msikeyboard.SetMode", 0, "normal")
	if call.Err != nil {
		panic(call.Err)
	}
	call = obj.Call("name.elemc.msikeyboard.Set", 0, "middle", "blue", "low")
	if call.Err != nil {
		panic(call.Err)
	}

}
