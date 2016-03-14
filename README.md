# openHAB-433-HomeKit-Bridge

Work in Progress. 433 is currently not optional.

## Intro

We have a bunch of [RWE Smarthome devices](https://www.rwe-smarthome.de/web/cms/de/2935200/home/geraete-apps/geraete/) at home. The RWE app to control those devices is really sluggish and I like Apples idea to try and combine all those different smarthome solutions with a single protocol.

RWE Smarthome doesn't support HomeKit, but I found that [openHAB](https://github.com/openhab/openhab) has a binding that integrates RWE devices. openHAB itself makes it possible to control it's devices using a REST-API, so all I had to do, was to create an interface between openHAB and a HomeKit bridge implementation.
There's [OpenHAB HomeKit Bridge](https://github.com/htreu/OpenHAB-HomeKit-Bridge), which does exactly that, but it was too hard to implement what I wanted to do. So I implemented my own REST-Client in Go and went with the awesome HomeKit implementation [HomeControl](https://github.com/brutella/hc).

I focused on openHAB thermostats, but it's really is easy to extend thanks to HomeControl.

I also have a couple of cheap 433 MHz RF outlets and have been using a Raspberry Pi to control those outlets with [Raspberry Pi Remote](https://github.com/xkonni/raspberry-remote) for a while. I wanted to get those into HomeKit as well. Their send-command implementation is included with this project.

## Building
* Install openHAB
* [Enable and configure the rwe-smarthome-binding in openHAB](https://github.com/openhab/openhab/wiki/RWE-Smarthome-Binding)
* Setup items and sitemap (see *openHAB config*)
* [Install Go](https://golang.org/doc/install)
* [Setup Go workspace](https://golang.org/doc/code.html#Organization)
* `cd $GOPATH/src`
* Clone this repo `git clone https://github.com/marconett/openHAB-433-HomeKit-Bridge; cd openHAB-433-HomeKit-Bridge`
* Install dependencies `go get`
* Build it `go build` or `env GOOS=linux GOARCH=arm GOARM=5 go build -v` for Raspberry
* Start it `./openHAB-433-HomeKit-Bridge`

### 433 (not optional yet)
* [Setup Raspberry Pi Remote](https://github.com/xkonni/raspberry-remote#setup)
* `sudo gpio export 17 out` to [Enable User Mode](https://github.com/xkonni/raspberry-remote#user-mode)
* `make` the send command
* move the *send* binary into the same folder as the openHAB-433-HomeKit-Bridge binary

## CLI Reference
```
Usage of:
  -host string
        OpenHAB Host (default "localhost")
  -name string
        HomeKit Bridge Name (default "LBI_Bridge")
  -pin string
        HomeKit Bridge PIN (default "32191123")
  -port string
        OpenHAB Port (default "8080")
  -sitemap string
        OpenHAB Sitemap (default "default")
```

## openHAB config

**NOTE**: The application uses string-matching of the item names to differentiate them. The underscore is important. The substring before the underscore determines the name of the room the thermostat is in. The substring after the underscore tells the application which element of the thermostat it is. It matches case-insensitive for *currenthum*, *settemp*, *currenttemp* and *setmode*.

### .items
```
Number Marco_CurrentHumidity "Luftfeuchtigkeit Marco [%.1f %%]" <temperature> {rwe="id=DEVICEID,param=humidity"}
Number Marco_CurrentTemp "Temperatur Marco [%.1f °C]" <temperature> {rwe="id=DEVICEID,param=temperature"}
Number Marco_SetTemp "Solltemperatur Marco [%.1f °C]" <temperature> {rwe="id=DEVICEID,param=settemperature"}
Switch Marco_SetMode "Thermostat Modus Marco" <temperature> {rwe="id=DEVICEID,param=operationmodeauto"}

Number  Thomas_CurrentHumidity "Luftfeuchtigkeit Thomas [%.1f %%]" <temperature> (rwe) {rwe="id=DEVICEID,param=humidity"}
Number  Thomas_CurrentTemp "Temperatur Thomas [%.1f °C]" <temperature> (rwe) {rwe="id=DEVICEID,param=temperature"}
Number  Thomas_SetTemp "Solltemperatur Thomas [%.1f °C]" <temperature> (rwe) {rwe="id=DEVICEID,param=settemperature"}
Switch  Thomas_SetMode "Thermostat Modus Thomas" <temperature> (rwe) {rwe="id=DEVICEID,param=operationmodeauto"}

Number  Buero_CurrentHumidity "Luftfeuchtigkeit Buero [%.1f %%]" <temperature> (rwe) {rwe="id=DEVICEID,param=humidity"}
Number  Buero_CurrentTemp "Temperatur Buero [%.1f °C]" <temperature> (rwe) {rwe="id=DEVICEID,param=temperature"}
Number  Buero_SetTemp "Solltemperatur Buero [%.1f °C]" <temperature> (rwe) {rwe="id=DEVICEID,param=settemperature"}
Switch  Buero_SetMode "Thermostat Modus Buero" <temperature> (rwe) {rwe="id=DEVICEID,param=operationmodeauto"}
```

### .sitemap
```
sitemap default label="HomeKit" {

  Text item=Marco_CurrentTemp
  Text item=Marco_CurrentHumidity
  Text item=Marco_SetTemp
  Switch item=Marco_SetMode

  Text item=Thomas_CurrentTemp
  Text item=Thomas_CurrentHumidity
  Text item=Thomas_SetTemp
  Switch item=Thomas_SetMode

  Text item=Buero_CurrentTemp
  Text item=Buero_CurrentHumidity
  Text item=Buero_SetTemp
  Switch item=Buero_SetMode
}
```
