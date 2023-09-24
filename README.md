## Azure IoT Hub temperature telemetry

An example of Raspberry Pi being used to read temperature using a DS18B20 temperature sensor and report those readings to Azure IoT Hub.

### Pre-requisites

- Raspberry Pi
  - Raspberry Pi OS with one-wire interface enabled
- DS18B20 temperature sensor
- IoT Hub setup on Azure

### Required secrets

Two files with secrets must first be supplied at root: `.env-secret` and `config.yml`

Example of `.env-secret` (required by `direnv`):

```
TARGET_IP=[rPi ipv4 address]
TARGET_USER=[rPi user]
TARGET_DIRECTORY=[project target directory on rPi]
```

Example of `config.yml`:

```
azure_mqtt_host: "[Azure IoT Hub hostname]"
azure_mqtt_endpoint: "devices/[IoT Hub device ID]/messages/events/[IoT Hub endpoint name]"
azure_device_name: "[IoT Hub device ID]"
azure_device_key: "[IoT Hub device key]"
```

### Usage

1.) Run `make deploy` from root.

2.) Register systemd services from `systemd/system` in Raspberry Pi OS.

Raspberry Pi will read the temperature in Celsius each 30 seconds and then report it to Azure each minute.
