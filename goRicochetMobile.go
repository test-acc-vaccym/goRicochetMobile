package goRicochetMobile

import (
	"github.com/s-rah/go-ricochet/application"
	"github.com/s-rah/go-ricochet/utils"
	"log"
	"net/http"
	"time"
	"github.com/yawning/bulb/utils/pkcs1"
	"crypto/rsa"
)



func GeneratePrivateKey() (string, error) {
	privateKey, err := utils.GeneratePrivateKey()
	if err != nil {
		return "", err
	}
	return utils.PrivateKeyToString(privateKey), nil
}

func GetOnionAddress(privateKey string) string {
	pk, _ := utils.ParsePrivateKey([]byte(privateKey))
	pubKey := rsa.PublicKey(pk.PublicKey)
	addr, err := pkcs1.OnionAddr(&pubKey)
	if err != nil || addr == "" {
		return ""
	}
	return addr
}



/******** Testing by standing up an echobot ******/

func TestNet() (ok bool, ex error) {
	_, err := http.Get("http://golang.org/")
	if err != nil {
		return false, err
	}
	return true, nil
}

func EchoBot(privateKeyData string)  {
	privateKey, err := utils.ParsePrivateKey([]byte(privateKeyData))
	if err != nil {
		log.Fatal("error parsing private key: %v", err)
	}

	log.Println("Setup onion hidden service via tor control...")
	l, err := application.SetupOnion("127.0.0.1:9051", "tcp4","", privateKey, 9878)
	if err != nil {
		log.Fatalf("error setting up onion service: %v", err)
	}

	echobot := new(application.RicochetApplication)
	echobot.Init(privateKey, new(application.AcceptAllContactManager))

	echobot.OnChatMessage(func(rai *application.RicochetApplicationInstance, id uint32, timestamp time.Time, message string) {
		log.Printf("message from %v - %v", rai.RemoteHostname, message)
		rai.SendChatMessage(message)
	})
	log.Printf("echobot started on %s", l.Addr().String())
	echobot.Run(l)
}