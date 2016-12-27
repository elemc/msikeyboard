#!/bin/sh

commit=`grep -r "%global commit" msikeyboard.spec | awk '{print $3}'`
short_commit=`echo ${commit:0:7}`
repo="msikeyboard"

wget -c https://github.com/elemc/${repo}/archive/${commit}/${repo}-${short_commit}.tar.gz
