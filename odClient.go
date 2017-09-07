package goRicochetMobile

import (
	"github.com/s-rah/go-ricochet/connection"
	"time"
	"log"
)

type ODClient struct {
	connection.AutoConnectionHandler
	messages chan string
	deviceName string
	deviceLevel int
	batteryLevel string
}

/************* Chat Channel Handler ********/

// ChatMessage passes the response to messages.
func (odc *ODClient) ChatMessage(messageID uint32, when time.Time, message string) bool {
	log.Printf("Received Message: %s", message)
	odc.messages <- message
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