# milight-daemon
go daemon that controlls your milight lamps via http requests

## How to install

### precompiled
download the latest release from https://github.com/Nesurion/milight-daemon/releases

### from source
clone the repo by running:  
`git clone git@github.com:Nesurion/milight-daemon.git`  

compile the binary by running  
`make`

Note that in order to compile you need a working golang setup https://golang.org/doc/install.  
Go offers a an easy way to cross-compile your code (go >= 1.5.0).

To compile for a raspberry pi running debian run:  

raspberry pi version 1:  
`make pi`

raspberry pi veersion 2:  
`make pi2`

## Config
Milight-Daemon needs a config that contains information about your milight bridge.

A valid config looks like this:

```json
{
    "port": 8080,
    "bridge": "192.168.2.141"
}
```

`port`: milight-daemon listening port  
`bridge`: IP address of your milight bridge

**Important** the config must be located in the same folder as the milight-daemon binary. This will probaly change in the future.

## Flags
The following flags are supported:

```
Usage of ./milight-daemon:
  -mode string
    	Gin Mode (debug, release, test) (default "release")
```

get a list of all possible flags by running: `milight-daemon --help`

## API
The milight-daemon exposes the following enpoints:

**All Endpoints support the group parameter**:  
paramter: group = \[0, 1, 2, 3, 4\] (0 = all groups)  

### on
request method: `POST`  
endpoint: `/on`  
example: `http://localhost:8080/on?group=1`  

### off
request method: `POST`  
endpoint: `/off`
example: `http://localhost:8080/off?group=1`  

### rgb
request method: `POST`  
endpoint: `/rgb`  
paramter:
  - r = [1, 255]
  - g = [1, 255]
  - b = [1, 255]

example: `http://localhost:8080/rgb?group=1&r=255&g=100&b=200`  

### brightness
request method: `POST`  
endpoint: `/brightness`
paramter:
	- level = [1,100]
example: `http://localhost:8080/brightness?group=1&level=50`

### color
request method: `POST`  
endpoint: `/color`
paramter: 
  - color =
    - violet
    - blue
    - baby_blue
    - aqua
    - mint
    - seafoam_green
    - green
    - lime_green
    - yellow
    - yellow_orange
    - orange
    - red
    - pink
    - fusia
    - lilac
    - lavendar

example: `http://localhost:8080/color?group=1&color=blue`

### white
request method: `POST`  
endpoint: `/white`  
example: `http://localhost:8080/white?group=1`

### night
request method: `POST`  
endpoint: `/night`  
example: `http://localhost:8080/night?group=1`

### disco
request method: `POST`  
endpoint: `/disco`
paramter:
  - speed = \[up, down\]

*Note: speed is optional. To start disco mode you dont need the speed parameter.*  
example: `http://localhost:8080/disco?group=1`  
example: `http://localhost:8080/disco?group=1?speed=up` (speed=[up, down])  
