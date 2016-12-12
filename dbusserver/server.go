package dbusserver

import (
	"fmt"

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
	select {}
}

// Stop function stops a DBus server
func (server *DBusServer) Stop() {
	if server.Conn != nil {
		server.Conn.Close()
	}
}

// Set is a dbus method Set, set color and intensity to region and return OK
// or error message
func (server *DBusServer) Set(region, color, intensity string) (string, *dbus.Error) {
	fmt.Printf("Region: %s\n\tColor: %s\n\tIntensity: %s\n", region, color, intensity)
	return "OK", nil
}

// SetMode is a dbus method SetMode, set to keyboard and return OK
// or error message
func (server *DBusServer) SetMode(mode string) (string, *dbus.Error) {
	fmt.Printf("Mode: %s\n", mode)
	return "OK", nil
}
