#!bin/sh

cat /sys/devices/w1_bus_master1/28-*/w1_slave > /home/pi/deploy/sensors/temperature/reading

