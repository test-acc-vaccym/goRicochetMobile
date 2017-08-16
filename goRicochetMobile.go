package goRicochetMobile

import (
	"github.com/s-rah/go-ricochet/application"
	"github.com/s-rah/go-ricochet/utils"
	"log"
	"time"
)

func GeneratePrivateKey() (string, error) {
	privateKey, err := utils.GeneratePrivateKey()
	if err != nil {
		return "", err
	}
	return utils.PrivateKeyToString(privateKey), nil
}

func EchoBot(privateKeyData string) {
	privateKey, err := utils.ParsePrivateKey([]byte(privateKeyData))

	if err != nil {
		log.Fatal("error parsing private key: %v", err)
	}

	echobot := new(application.RicochetApplication)

	l, err := application.SetupOnion("127.0.0.1:9051", "tcp4","", privateKey, 9878)
	//l, err := application.SetupOnion("/data/data/org.torproject.android/app_bin/control.txt", "unix","", privateKey, 9878)

	if err != nil {
		log.Fatalf("error setting up onion service: %v", err)
	}

	echobot.Init(privateKey, new(application.AcceptAllContactManager))
	echobot.OnChatMessage(func(rai *application.RicochetApplicationInstance, id uint32, timestamp time.Time, message string) {
		log.Printf("message from %v - %v", rai.RemoteHostname, message)
		rai.SendChatMessage(message)
	})
	log.Printf("echobot listening on %s", l.Addr().String())
	echobot.Run(l)
}
