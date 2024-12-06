package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	out, err := exec.Command("sensors").Output()
	if err != nil {
		log.Fatal(err)
	}
	result := strings.FieldsFunc(string(out), func(r rune) bool {
		return r == '\n'
	})
	var leadingInt = regexp.MustCompile(`[-+][0-9,.]*`)

	field1 := "Tctl:"
	field2 := "Sensor 1:"

	//field = "in0:"
	str1 := ""
	str2 := ""
	for _, d := range result {
		if strings.HasPrefix(d, field1) {
			str1 = d
			continue
		}
		if strings.HasPrefix(d, field2) {
			str2 = d
			break
		}
	}

	str1 = getCPU(str1, leadingInt)
	str2 = getNVME(str2, leadingInt)

	// cpu_temp, err := strconv.ParseFloat(str1, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// nvme_temp, err := strconv.ParseFloat(str2, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cpu_temp := str1
	nvme_temp := str2

	//fmt.Println(cpu_temp, nvme_temp)

	client := connect()
	defer client.Disconnect(250)

	// Topics to publih for MQTT
	TOPIC_TO_PUBLISH1 := "homeassistant/sensor/proxmox_system/cpu_temp"
	TOPIC_TO_PUBLISH2 := "homeassistant/sensor/proxmox_system/nvme_temp"

	// Publish a message with topic
	publish(client, TOPIC_TO_PUBLISH1, cpu_temp)
	publish(client, TOPIC_TO_PUBLISH2, nvme_temp)

	// Keep the connection alive
	time.Sleep(1 * time.Second)
}

func getCPU(s string, r *regexp.Regexp) string {
	//s = "Tctl:         +38.8°C"

	s = r.FindString(s)
	return s
}

func getNVME(s string, r *regexp.Regexp) string {
	//s = "Sensor 1:     +44.9°C  (low  = -273.1°C, high = +65261.8°C)"

	s = r.FindString(s)
	return s
}

func connect() mqtt.Client {
	var broker = "192.168.1.45"
	var port = 1883
	var clientId = "VM_system"

	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	options.SetUsername("Username")
	options.SetPassword("Password")
	options.SetClientID(clientId)

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	//fmt.Println("Connected to MQTT Broker:", broker)
	return client
}

func subscribe(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), msg.Payload())
	})

	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	//log.Printf("Subscribed to topic: %s\n", topic)
}

func publish(client mqtt.Client, topic string, message string) {
	token := client.Publish(topic, 0, false, message)
	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Printf("Published message: %s\n", message)
}
