package dbusserver

import (
	"fmt"
	"log"
	"strings"

	"github.com/elemc/msikeyboard/gomsikeyboard"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	dBusName = "name.elemc.msikeyboard"
	dBusPath = "/name/elemc/msikeyboard"
	intro    = `<node>
    <interface name="` + dBusName + `">
        <method name="Set">
            <arg direction="in" type="s"/>
            <arg direction="in" type="s"/>
            <arg direction="in" type="s"/>
            <arg direction="out" type="s"/>
        </method>
        <method name="SetMode">
            <arg direction="in" type="s"/>
            <arg direction="out" type="s"/>
        </method>
    </interface>` + introspect.IntrospectDataString + `</node>`
)

// DBusServer struct for D-Bus server
type DBusServer struct {
	Conn *dbus.Conn
	led  gomsikeyboard.LEDSetting
}

// Start function start DBus server
func (server *DBusServer) Start() (err error) {
	server.Conn, err = dbus.SessionBus()
	if err != nil {
		return
	}
	defer server.Stop()
	reply, err := server.Conn.RequestName(dBusName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("name %s already taken", dBusName)
	}
	err = server.Conn.Export(server, dBusPath, dBusName)
	if err != nil {
		return
	}
	err = server.Conn.Export(server, dBusPath, dBusName)
	if err != nil {
		return
	}
	err = server.Conn.Export(introspect.Introspectable(intro), dBusPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return
	}
	log.Printf("starting D-Bus server")
	select {}
}

// Stop function stops a DBus server
func (server *DBusServer) Stop() {
	if server.Conn != nil {
		server.Conn.Close()
	}
	log.Printf("stopping D-Bus server")
}

// Set is a dbus method Set, set color and intensity to region and return OK
// or error message
func (server *DBusServer) Set(region, color, intensity string) (result string, dbusErr *dbus.Error) {
	result = "OK"
	var err error
	if strings.ToLower(region) == "all" {
		server.led.Left.Color = color
		server.led.Left.Intensity = intensity
		server.led.Middle.Color = color
		server.led.Middle.Intensity = intensity
		server.led.Right.Color = color
		server.led.Right.Intensity = intensity
	} else {
		switch strings.ToLower(region) {
		case "left":
			server.led.Left.Color = color
			server.led.Left.Intensity = intensity
		case "middle":
			server.led.Middle.Color = color
			server.led.Middle.Intensity = intensity
		case "right":
			server.led.Right.Color = color
			server.led.Right.Intensity = intensity
		}
	}

	err = server.led.Check()
	if err != nil {
		result = err.Error()
		dbusErr = &dbus.ErrMsgInvalidArg
		return
	}

	err = server.led.Set()
	if err != nil {
		result = err.Error()
		return
	}
	return
}

// SetMode is a dbus method SetMode, set to keyboard and return OK
// or error message
func (server *DBusServer) SetMode(mode string) (result string, dbusErr *dbus.Error) {
	server.led.Mode = mode

	result = "OK"
	var err error

	err = server.led.Check()
	if err != nil {
		result = err.Error()
		dbusErr = &dbus.ErrMsgInvalidArg
		return
	}

	err = server.led.Set()
	if err != nil {
		result = err.Error()
		return
	}
	return
}
