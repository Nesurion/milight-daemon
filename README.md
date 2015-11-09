# milight-daemon
go daemon that controlls your milight lamps via http requests

## How to install

### precompiled
download the latest release from https://github.com/Nesurion/milight-daemon/releases

### from source
clone the repo by running: `git clone git@github.com:Nesurion/milight-daemon.git`
compile the binary by running `make`.
Note that in order to compile you need a working golang setup https://golang.org/doc/install.
Go offers a an easy way to cross-compile your code (go >= 1.5.0). To compile for a raspberry pi running debian run:
`GOOS=linux GOARCH=arm GOARM=5 make`

## Config
Milight-Daemon needs a config that contains information about your milight bridge.
A valid config looks like this:

```json
{
    "port": 8080,
    "bridge": "192.168.2.141"
}
```

port: milight-daemon listening port
bridge: IP address of your milight bridge

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

### on
endpoint: `/on`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/on?group=1`

### off
endpoint: `/off`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/off?group=1`

### rgb
endpoint: `/rgb`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/rgb?group=1&r=255&g=100&b=200` (r,g,b = [1,255])

### brightness
endpoint: `/brightness`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/brightness?group=1&level=50` (level = [1,100])

### color
endpoint: `/color`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/color?group=1&color=blue`
	color:
		violet
		blue
		baby_blue
		aqua
		mint
		seafoam_green
		green
		lime_green
		yellow
		yellow_orange
		orange
		red
		pink
		fusia
		lilac
		lavendar

### white
endpoint: `/white`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/white?group=1`

### night
endpoint: `/night`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/night?group=1`

### disco
endpoint: `/disco`
paramter: group = 0,1,2,3,4 (0 = all groups)
example: `http://localhost:8080/disco?group=1`
example: `http://localhost:8080/disco?group=1?speed=up` (speed=[up, down])
