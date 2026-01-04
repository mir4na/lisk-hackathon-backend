package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ShipmentVerification struct {
	ID                  uuid.UUID        `json:"id"`
	InvoiceID           uuid.UUID        `json:"invoice_id"`
	ContainerID         *string          `json:"container_id,omitempty"`
	BillOfLadingNumber  *string          `json:"bill_of_lading_number,omitempty"`
	Carrier             *string          `json:"carrier,omitempty"`
	OriginPort          *string          `json:"origin_port,omitempty"`
	DestinationPort     *string          `json:"destination_port,omitempty"`
	Status              string           `json:"status"`
	VerifiedAt          *time.Time       `json:"verified_at,omitempty"`
	APIResponse         json.RawMessage  `json:"api_response,omitempty"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
}

type VerifyShipmentRequest struct {
	InvoiceID          uuid.UUID `json:"invoice_id" binding:"required"`
	ContainerID        string    `json:"container_id"`
	BillOfLadingNumber string    `json:"bill_of_lading_number"`
}

type ShipmentTrackingResponse struct {
	ContainerID     string    `json:"container_id"`
	Status          string    `json:"status"`
	CurrentLocation string    `json:"current_location"`
	Vessel          string    `json:"vessel"`
	ETA             time.Time `json:"eta"`
	Events          []ShipmentEvent `json:"events"`
}

type ShipmentEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
}
