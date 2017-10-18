package od

import (
	"github.com/dballard/goRicochetMobile/od/odClient"
	"log"
	"strconv"
)

var (
	// Downsampling array from https://git.mascherari.press/oniondildonics/client/src/master/main.go
	// moddified
	levelArr = []int{1, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6}
	client *odClient.ODClient = nil
)

func ODClientConnect(privateKey string, serverAddr string) error {
	log.Println("ODClientConnect()")
	client = new(odClient.ODClient)
	err := client.Connect(privateKey, serverAddr)
	return err
}

func ODClientDisconnect() {
	log.Println("ODClientDisconnect()")
	client.Disconnect()
	client = nil
}

func GetDeviceName() string {
	client.SendMessage("/name")
	name := client.GetMessage()
	return name
}

func GetBatteryLevel() string {
	client.SendMessage("/battery")
	batteryLevel := client.GetMessage()
	return batteryLevel
}

func GetVibeLevel() int {
	client.SendMessage("/level")
	level, err := strconv.Atoi(client.GetMessage())
	if err != nil {
		// TODO: don't swallow errors
		return 0
	}
	return levelArr[level] // not bounds checking...
}

func SetVibeLevel(newVibeLevel int) {
	client.SendMessage("/level " + strconv.Itoa(newVibeLevel))
}
