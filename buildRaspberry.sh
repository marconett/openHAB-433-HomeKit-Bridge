#!/bin/bash

env GOOS=linux GOARCH=arm GOARM=5 go build -v
scp thermostat marco@192.168.100.3:/home/marco
