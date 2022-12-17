package main

import (
	"crypto/tls"
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttClientInit create mqtt client config
func MqttClientInit(server, clientID, username, password string) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientID).SetCleanSession(true)
	if username != "" {
		opts.SetUsername(username)
		if password != "" {
			opts.SetPassword(password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	opts.SetTLSConfig(tlsConfig)
	return opts
}

// MqttConnect function felicitates the MQTT connection
func MqttConnect() error {
	// Initiate the MQTT connection
	ClientOpts := MqttClientInit("tcp://127.0.0.1:1884", "eventbus", "", "")
	Client = MQTT.NewClient(ClientOpts)
	if TokenClient := Client.Connect(); TokenClient.Wait() && TokenClient.Error() != nil {
		return fmt.Errorf("client.Connect() Error is %s" + TokenClient.Error().Error())
	}
	return nil
}

func PublishMqtt(topic, message string) error {
	TokenClient := Client.Publish(topic, 0, false, message)
	if TokenClient.Wait() && TokenClient.Error() != nil {
		return fmt.Errorf("client.publish() Error in topic %s. reason: %s. ", topic, TokenClient.Error().Error())
	}
	log.Printf("publish topic %s message %s", topic, message)
	return nil
}
