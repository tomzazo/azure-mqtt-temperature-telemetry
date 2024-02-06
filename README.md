## Azure IoT Hub temperature telemetry

An example of a Raspberry Pi being used to read temperature using a DS18B20 temperature sensor and report those readings to an Azure IoT Hub _without any SDKs_.

### Pre-requisites

- Raspberry Pi
  - Raspbian/Raspberry Pi OS with the one-wire interface enabled
- DS18B20 temperature sensor
- IoT Hub setup on Azure

### Required secrets

A file with secrets must first be supplied at root named `config.yml`

Example of `config.yml`:

```
azure_mqtt_host: "[Azure IoT Hub hostname]"
azure_mqtt_endpoint: "devices/[IoT Hub device ID]/messages/events/"
azure_device_name: "[IoT Hub device ID]"
azure_device_key: "[IoT Hub device key]"
```

### Usage

1.) Change Ansible inventory to point to your Raspberry Pi device.

2.) Run `just setup` at root.

Raspberry Pi will read the temperature in Celsius every 30 seconds and then report it to Azure every minute.
