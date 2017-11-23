package bbqapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

var connected = false
var mqttClient mqtt.Client

// Temp is the JSON structure for the temperatures
type Temp struct {
	Index int     `json:"index"`
	Temp  float64 `json:"temp"`
}

// Temps array
type Temps []Temp

// LastMessage contains the current IoT message received
var LastMessage Temps

// Temps container to unmarshall
//type Temps []Temp

// MainLoop handles ongoing process
func MainLoop() {
	mqttClient = IotConnect()
	if mqttClient.IsConnected() {
		defer mqttClient.Disconnect(0)
	}
	router := gin.Default()
	router.GET("/", defaultTarget)

	router.GET("/temp", tempTarget)
	readIotTemp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)

}

func defaultTarget(c *gin.Context) {
	c.String(http.StatusOK, "BBQ App \n\nTry  GET /temp")
}

func tempTarget(c *gin.Context) {
	fmt.Println(LastMessage)
	c.JSON(http.StatusOK, LastMessage)
}

// IotConnect func
// Organization ID "sbdt88"
// URL: sbdt88.messaging.internetofthings.ibmcloud.com
// https://orgId.messaging.internetofthings.ibmcloud.com:443/api/v0002
// Device Type "RaspberryPi"
// Device ID "bbq-raspi"
// Authentication Method "use-token-auth"
// Authentication Token "XSZWm(3SzU2norg15W"
func IotConnect() mqtt.Client {
	if connected {
		return mqttClient
	}
	appID := "bbqapp"
	orgID := "sbdt88"
	apiKey := "a-sbdt88-kw8hgl76vz"
	authToken := "Djc@1fIhP)qE(V39bQ"
	pwd := authToken
	host := "ssl://sbdt88.messaging.internetofthings.ibmcloud.com:8883" // tcp:1883 is disabled -- MUST be ssl
	// https://developer.ibm.com/answers/questions/264069/how-do-you-publish-to-a-mqtt-topic-in-iot-using-ja.html
	connectString := fmt.Sprintf("a:%s:%s", orgID, appID)
	cOpt := mqtt.NewClientOptions()
	cOpt.SetClientID(connectString)
	cOpt.SetUsername(apiKey)
	cOpt.SetPassword(pwd)
	cOpt.AddBroker(host)
	fmt.Println("cOpt : ", cOpt)
	mqttClient = mqtt.NewClient(cOpt)
	tok := mqttClient.Connect()
	if tok.Wait() && tok.Error() != nil {
		fmt.Println("Error: ", tok.Error())
	}
	fmt.Println("mqttClient = ", mqttClient.IsConnected())
	connected = true
	return mqttClient
}

func readIotTemp() {
	//iot-2/type/{device type}/id/{device id}/evt/{event type}/fmt/{format type}
	fmt.Println("reading temp to IoT")
	// login first
	mq := IotConnect()
	// https://developer.ibm.com/answers/questions/264069/how-do-you-publish-to-a-mqtt-topic-in-iot-using-ja.html
	devType := "RaspberryPi"
	devID := "bbq-raspi"
	topic := fmt.Sprintf("iot-2/type/%s/id/%s/evt/status/fmt/json", devType, devID)
	fmt.Println("topic: ", topic)

	var msgRcvd = func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
		var data Temps
		err := json.Unmarshal(message.Payload(), &data)
		if err != nil {
			fmt.Println("Error unmarshalling: ", err)
		}
		LastMessage = data
	}

	if token := mq.Subscribe(topic, 0, msgRcvd); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}
