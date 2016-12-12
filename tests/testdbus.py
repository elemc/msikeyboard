#!/usr/bin/env python3

import dbus

dbusName = "name.elemc.msikeyboard"
dbusPath = "/name/elemc/msikeyboard"

bus = dbus.SessionBus()
msi = bus.get_object(dbusName, dbusPath)

set_mode = msi.get_dbus_method("SetMode", dbusName)
result = set_mode("normal")
if result != "OK":
    print("error in SetMode method: " + result)
else:
    print("SetMode execute successfull")
_set = msi.get_dbus_method("Set", dbusName)
result = _set("left", "red", "high")
if result != "OK":
    print("error in SetMode method: " + result)
else:
    print("SetMode execute successfull")
