package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	yaml "gopkg.in/yaml.v3"
)

// Miscellaneous constants.
const (
	configFilePath      = "config.yml"
	mqttProtocolVersion = 4 //v3.1.1
)

// Temperature constants.
const (
	temperatureFilePath         = "sensors/temperature/reading"
	tempConversionRate  float64 = 1000
)

type Config struct {
	AzureMQTTHost     string `yaml:"azure_mqtt_host"`
	AzureMQTTEndpoint string `yaml:"azure_mqtt_endpoint"`
	AzureDeviceName   string `yaml:"azure_device_name"`
	AzureDeviceKey    string `yaml:"azure_device_key"`
}

// getConfig reads the configuration file.
func getConfig() (Config, error) {
	execPath, err := os.Executable()
	if err != nil {
		return Config{}, fmt.Errorf("error getting executable path: %s", err)
	}

	configFileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", filepath.Dir(execPath), configFilePath))
	if err != nil {
		return Config{}, fmt.Errorf("error reading configuration file: %s", err)
	}

	c := Config{}

	err = yaml.Unmarshal([]byte(configFileContent), &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

// getTemperatureReading opens and reads the specified file,
// containing temperature readings from the sensor.
func getTemperatureReading() (float64, error) {
	execPath, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("error getting executable path: %s", err)
	}

	temperatureFileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", filepath.Dir(execPath), temperatureFilePath))
	if err != nil {
		return 0, fmt.Errorf("error reading temperature file: %s", err)
	}

	contentSplit := strings.Split(strings.TrimSpace(string(temperatureFileContent)), " ")

	tempVar := contentSplit[len(contentSplit)-1]
	tempVarSplit := strings.Split(tempVar, "=")

	tempString := tempVarSplit[len(tempVarSplit)-1]

	temp, err := strconv.ParseInt(tempString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing string to integer: %s", err)
	}

	return float64(temp) / tempConversionRate, nil
}

// generateAzureSASToken generates a SAS token for communication with Azure.
//
// The host must equal to the IoT Hub host on Azure.
// The device must equal to the device name added in the IoT Hub on Azure.
// The key must be provided in base64 format.
func generateAzureSASToken(host, device, key string) string {
	expires := time.Now().Add(1 * time.Hour).Unix()
	resource := url.QueryEscape(fmt.Sprintf("%s/devices/%s", host, device))
	toSign := resource + "\n" + strconv.Itoa(int(expires))

	keyDecoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Println(err, "decoding key")
	}

	hash := hmac.New(sha256.New, keyDecoded)
	hash.Reset()
	hash.Write([]byte(toSign))
	signed := url.QueryEscape(base64.StdEncoding.EncodeToString(hash.Sum(nil)))

	return "SharedAccessSignature sr=" + resource + "&sig=" + signed + "&se=" + strconv.Itoa(int(expires))
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err, "error getting configuration file")
	}

	clientOptions := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcps://%s:8883", cfg.AzureMQTTHost)).
		SetClientID(cfg.AzureDeviceName).
		SetUsername(fmt.Sprintf("%s/%s/?api-version=2021-04-12", cfg.AzureMQTTHost, cfg.AzureDeviceName)).
		SetPassword(generateAzureSASToken(cfg.AzureMQTTHost, cfg.AzureDeviceName, cfg.AzureDeviceKey)).
		SetProtocolVersion(mqttProtocolVersion).
		SetCleanSession(true)

	client := mqtt.NewClient(clientOptions)

	connectToken := client.Connect()

	<-connectToken.Done()

	if err := connectToken.Error(); err != nil {
		log.Println("error connecting MQTT client: ", err)

		return
	}

	tempFloat, err := getTemperatureReading()
	if err != nil {
		log.Println("error getting temperature reading: ", err)

		return
	}

	msg := fmt.Sprintf("%.1f", tempFloat)

	publishToken := client.Publish(
		cfg.AzureMQTTEndpoint,
		0,
		false,
		msg,
	)

	<-publishToken.Done()

	if err := publishToken.Error(); err != nil {
		log.Println("error publishing message: ", err)

		return
	}
}
