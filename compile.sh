#!/bin/bash

V=0.5.3
_os=('windows:win:.exe:386_-x86,amd64_-x64' 'linux:linux:x:386_-x86,amd64_-x64' 'darwin:mac:x:amd64_x')

for _s in ${_os[@]}
do
    unset _a
    _t=(${_s//:/ })
    _o=${_t[0]} # OS
    _n=${_t[1]} # Name
    _e=${_t[2]} # Ext
    _a=${_t[3]} # Arc
    _a=(${_a//,/ })

    if [ $_e == 'x' ]
    then
      _e=''
    fi
    for _i in ${_a[@]}
    do
        _is=(${_i//_/ })
        _ia=${_is[0]}
        _in=${_is[1]}
        if [ $_in == 'x' ]
        then
          _in=''
        fi
        rm -f ./bin/BXGo-$V-$_n$_in.zip
        #_nv=BXGo-$V-$_n$_in$_e
        _nv=BXGo-$V$_e
        GOOS=$_o GOARCH=$_ia go build -o $_nv ./app
        zip -r ./bin/BXGo-$V-$_n$_in.zip $_nv config.ini theme
        rm -f $_nv
    done
done
