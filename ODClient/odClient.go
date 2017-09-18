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
	connection *connection.Connection
	recvMessages chan string
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
		log.Println("handler for contact.request")
		contact := new(channels.ContactRequestChannel)
		contact.Handler = odClient
		return contact
	})

	odClient.RegisterChannelHandler("im.ricochet.chat", func() channels.Handler {
		log.Println("handler for chat")
		chat := new(channels.ChatChannel)
		chat.Handler = odClient
		return chat
	})

	odClient.recvMessages = make(chan string)

	log.Println("ODClient connecting...")
	odClient.connection, err = goricochet.Open(serverAddr)
	if err != nil {
		log.Println("Error connecting %v", err)
		return err
	}
	log.Println("ODCleint connected!")
	log.Println("starting auth...")
	known, err := connection.HandleOutboundConnection(odClient.connection).ProcessAuthAsClient(privateKey)
	if err != nil {
		log.Println("Error handling auth: %v", err)
		return err
	}

	log.Println("go Process")
	// TODO: end with breakChannel
	go odClient.connection.Process(odClient)

	if !known {
		err := odClient.connection.RequestOpenChannel("im.ricochet.contact.request", odClient)
		if err != nil {
			log.Printf("could not contact %s", err)
		}
	}

	log.Println("ODClient: Authenticated")

	log.Println("go")

	log.Println("RequestOpenChanel chat")
	err = odClient.connection.RequestOpenChannel("im.ricochet.chat", odClient)
	if err != nil {
		log.Println("Error: " + err.Error())
	}

	//log.Println("sending greeting message")
	//odClient.SendMessage("hello from the client")

	return nil
}

/*func (odClient *ODClient) RequestContact() {
	odClient.connection.Do(func() error {
		channel := odClient.connection.Channel("im.ricochet.contact.request", channels.Outbound)
		if channel != nil {
			contactRequestChannel, ok := (*channel.Handler).(*channels.ContactRequestChannel)
			if ok {
				//contactRequestChannel.Handler
			}
		}else {
			log.Println("ERROR: failed to find chat channel")
		}
		return nil
	})
}*/

func (odClient *ODClient) SendMessage(message string) {
	odClient.connection.Do(func() error {
		channel := odClient.connection.Channel("im.ricochet.chat", channels.Outbound)
		if channel != nil {
			chatchannel, ok := (*channel.Handler).(*channels.ChatChannel)
			if ok {
				chatchannel.SendMessage(message)
			}
		}else {
			log.Println("ERROR: failed to find chat channel")
		}
		return nil
	})
}

func (odClient *ODClient) GetMessage() string {
	message := <-odClient.recvMessages
	return message
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
	return "AndroidOD Client", ""
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