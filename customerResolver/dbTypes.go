package main

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	CustomerID          uuid.UUID `gorm:"primaryKey"`
	NameFirstName       string
	NameLastName        string
	CustomerEmail       string `gorm:"uniqueIndex;not null"`
	CustomerPhoneNumber uint8
	CustomerPostalCode  string
}

type Worker struct {
	WorkerId          uuid.UUID `gorm:"primaryKey"`
	WorkerFirstName   string
	WorkerLastName    string
	WorkerEmail       string
	WorkerIcNo        string
	WorkerPhoneNumber string
	Workerspeciality  string
}

type Location struct {
	LocationID            uuid.UUID `gorm:"primaryKey"`
	LocationLat           string
	LocationLng           string
	LocationCity          string
	LocationProvinceState string
	LocationZipCode       string
}

type JobRequest struct {
	JobRequestId          uuid.UUID `gorm:"primaryKey"`
	JobRequestCustomerId  uuid.UUID `gorm:"index"`
	JobRequestLocationId  uuid.UUID
	JobRequestSentAt      time.Time
	JobRequestCompletedAt time.Time
	JobRequestCity        string `gorm:"index"`
	JobRequestStatus      string
	JobRequestTitle       string
	JobRequestDescription string
	JobRequestTotalCost   float64
	JobRequestOffers      []Offer  `gorm:"foreignKey:OfferJobRequestId"`
	Customer              Customer `gorm:"foreignKey:JobRequestCustomerId"`
	JobRequestLocation    Location `gorm:"foreignKey:JobRequestLocationId"`
	JobRequestSpeciality  string   `gorm:"index"`
}

type Offer struct {
	OfferId           uuid.UUID `gorm:"primaryKey"`
	OfferJobRequestId uuid.UUID `gorm:"index"`
	OfferCustomerId   uuid.UUID
	OfferWorkerId     uuid.UUID `gorm:"index"`
	OfferPrice        float64
	OfferSentAt       time.Time
	OfferStatus       string
}

type Appointment struct {
	AppointmentId         uuid.UUID `gorm:"primaryKey"`
	AppointmentCustomerId uuid.UUID `gorm:"index"`
	AppointmentWorkerId   uuid.UUID `gorm:"index"`
	AppointmentOfferId    uuid.UUID
	AppointmentTime       time.Time
	AppointmentStatus     string
	Customer              Customer `gorm:"foreignKey:AppointmentCustomerId"`
	Worker                Worker   `gorm:"foreignKey:AppointmentWorkerId"`
	Offer                 Offer    `gorm:"foreignKey:AppointmentOfferId"`
}

// ======================= GQL TYPES  =======================
