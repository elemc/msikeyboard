msikeyboard
===========
Is a CLI for change color, intensity and mode on MSI keyboards

gomsikeyboard
=============
Is a go library for libmsikeyboard ( https://github.com/elemc/libmsikeyboard )

Build
-----
For build run:
$ go get -u github.com/elemc/msikeyboard

REST API
--------

URL arguments:
* theme - is a theme name
* all-color - is a all (left, middle, right) color name
* all-intensity - is a all (left, middle, right) intensity
* left-color - is a left color name
* left-intensity - is a left intensity
* middle-color - is a middle color name
* middle-intensity - is a middle intensity
* right-color - is a right color name
* right-intensity - is a right intensity
* mode - is a mode

Example URLs:
* http://localhost:9097/?theme=cool
* http://localhost:9097/set?left-color=white&left-intensity=low&middle-color=blue&middle-intensity=light&right-color=red&right-intensity=high&mode=gaming
* http://localhost:9097/set?all-color=white&all-intensity=high
