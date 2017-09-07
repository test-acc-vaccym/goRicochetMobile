package ODClient

import (
	"github.com/s-rah/go-ricochet/connection"
	"time"
	"log"
	"github.com/s-rah/go-ricochet/channels"
	"github.com/s-rah/go-ricochet/utils"
	"github.com/s-rah/go-ricochet"
)

type ODClient struct {
	connection.AutoConnectionHandler
	recvMessages chan string
	sendMessages chan string
	deviceName   string
	deviceLevel  int
	batteryLevel string
}

func (odClient *ODClient) Connect(privateKeyData string, serverAddr string) error {
	log.Println("ODCLient.Connect()")
	privateKey, err := utils.ParsePrivateKey([]byte(privateKeyData))
	if err != nil {
		log.Fatal("error parsing private key: %v", err)
	}

	odClient.Init(privateKey, serverAddr)

	odClient.RegisterChannelHandler("im.ricochet.contact.request", func() channels.Handler {
		contact := new(channels.ContactRequestChannel)
		contact.Handler = odClient
		return contact
	})

	odClient.RegisterChannelHandler("im.ricochet.chat", func() channels.Handler {
		chat := new(channels.ChatChannel)
		chat.Handler = odClient
		return chat
	})

	odClient.recvMessages = make(chan string)
	odClient.sendMessages = make(chan string)

	log.Println("ODClient connecting...")
	conn, err := goricochet.Open(serverAddr)
	if err != nil {
		log.Println("Error connecting %v", err)
		return err
	}
	log.Println("ODCleint connected!")
	_, err = connection.HandleOutboundConnection(conn).ProcessAuthAsClient(privateKey)
	if err != nil {
		log.Println("Error handling auth: %v", err)
		return err
	}
	log.Println("ODClient: Authenticated!")

	return nil
}

/************* Chat Channel Handler ********/

// ChatMessage passes the response to recvMessages.
func (odc *ODClient) ChatMessage(messageID uint32, when time.Time, message string) bool {
	log.Printf("Received Message: %s", message)
	odc.recvMessages <- message
	return true
}

// ChatMessageAck does nothing.
func (odc *ODClient) ChatMessageAck(messageID uint32) {
}

/************* Contact Channel Handler ********/

// GetContactDetails is purposely empty
func (odc *ODClient) GetContactDetails() (string, string) {
	return "", ""
}

// ContactRequest denies any contact request.
func (odc *ODClient) ContactRequest(name string, message string) string {
	return "Rejected"
}

// ContactRequestRejected purposly does nothing.
func (odc *ODClient) ContactRequestRejected() {
}

// ContactRequestAccepted purposly does nothing.
func (odc *ODClient) ContactRequestAccepted() {
}

// ContactRequestError purposly does nothing.
func (odc *ODClient) ContactRequestError() {
}