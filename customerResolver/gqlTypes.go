package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type GqlAppointment struct {
	ID         uuid.UUID            `json:"id"`
	CustomerID uuid.UUID            `json:"customerId"`
	WorkerID   uuid.UUID            `json:"workerId"`
	OfferID    uuid.UUID            `json:"offerId"`
	Time       time.Time            `json:"time"`
	Status     GqlAppointmentStatus `json:"status"`
}

type GqlCreateCustomerInput struct {
	FName         string  `json:"fName"`
	LName         string  `json:"lName"`
	Email         string  `json:"email"`
	PhoneNo       string  `json:"phoneNo"`
	PostalZipCode *string `json:"postalZipCode"`
}

type GqlCreateJobRequestInput struct {
	CustomerID  string                  `json:"customerId"`
	Location    *GqlCreateLocationInput `json:"location"`
	Status      GqlJobStatus            `json:"status"`
	Title       string                  `json:"title"`
	City        string                  `json:"city"`
	Description string                  `json:"description"`
}

type GqlCreateLocationInput struct {
	CustomerID string `json:"customerId"`
	Lng        int    `json:"lng"`
	Lat        int    `json:"lat"`
	State      string `json:"state"`
	City       string `json:"city"`
	Address    string `json:"address"`
}

type GqlCreateOfferInput struct {
	CustomerID string  `json:"customerId"`
	WorkerID   string  `json:"workerId"`
	JobID      string  `json:"jobId"`
	Price      float64 `json:"price"`
}

type GqlCreateWorkerInput struct {
	FName      string              `json:"fName"`
	LName      string              `json:"lName"`
	Email      string              `json:"email"`
	IcNo       string              `json:"icNo"`
	PhoneNo    string              `json:"phoneNo"`
	Speciality GqlWorkerSpeciality `json:"speciality"`
}

type GqlCustomer struct {
	ID            string            `json:"id"`
	FName         string            `json:"fName"`
	LName         string            `json:"lName"`
	Email         string            `json:"email"`
	PhoneNo       string            `json:"phoneNo"`
	PostalZipCode *string           `json:"postalZipCode"`
	JobRequests   []*GqlJobRequest  `json:"jobRequests"`
	Appointments  []*GqlAppointment `json:"appointments"`
	Offers        []*GqlOffer       `json:"offers"`
}

type GqlJobRequest struct {
	ID          string              `json:"id"`
	CustomerID  string              `json:"customerId"`
	Location    *GqlLocation        `json:"location"`
	City        string              `json:"city"`
	Status      GqlJobStatus        `json:"status"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	TotalCost   float64             `json:"totalCost"`
	Offers      []*GqlOffer         `json:"offers"`
	SentAt      time.Time           `json:"sentAt"`
	CompletedAt *time.Time          `json:"completedAt"`
	Speciality  GqlWorkerSpeciality `json:"speciality"`
}

type GqlLocation struct {
	ID         string `json:"id"`
	CustomerID string `json:"customerId"`
	Lng        int    `json:"lng"`
	Lat        int    `json:"lat"`
	State      string `json:"state"`
	City       string `json:"city"`
	Address    string `json:"address"`
}

type GqlOffer struct {
	ID         string         `json:"id"`
	CustomerID string         `json:"customerId"`
	WorkerID   string         `json:"workerId"`
	JobID      string         `json:"jobId"`
	Price      float64        `json:"price"`
	SentAt     time.Time      `json:"sentAt"`
	Status     GqlOfferStatus `json:"status"`
	SuggestedTime time.Time `json:"suggestedTime"`
}

type GqlPost struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type UpdateCustomerInput struct {
	ID            string  `json:"id"`
	FName         *string `json:"fName"`
	LName         *string `json:"lName"`
	Email         *string `json:"email"`
	PhoneNo       *string `json:"phoneNo"`
	PostalZipCode *string `json:"postalZipCode"`
}

type UpdateJobRequestInput struct {
	ID     string        `json:"id"`
	Status *GqlJobStatus `json:"status"`
}

type UpdateOfferInput struct {
	ID    string   `json:"id"`
	Price *float64 `json:"price"`
}

type UpdateWorkerInput struct {
	ID         string               `json:"id"`
	FName      *string              `json:"fName"`
	LName      *string              `json:"lName"`
	Email      *string              `json:"email"`
	IcNo       *int                 `json:"icNo"`
	PhoneNo    *int                 `json:"phoneNo"`
	Speciality *GqlWorkerSpeciality `json:"speciality"`
}

type GqlWorker struct {
	ID           string              `json:"id"`
	FName        string              `json:"fName"`
	LName        string              `json:"lName"`
	Email        string              `json:"email"`
	IcNo         string              `json:"icNo"`
	PhoneNo      string              `json:"phoneNo"`
	Speciality   GqlWorkerSpeciality `json:"speciality"`
	Offers       []*GqlOffer         `json:"offers"`
	Appointments []*GqlAppointment   `json:"appointments"`
}

type GqlAppointmentStatus string

const (
	AppointmentStatusUpcoming          GqlAppointmentStatus = "UPCOMING"
	AppointmentStatusCustomerCancelled GqlAppointmentStatus = "CUSTOMER_CANCELLED"
	AppointmentStatusWorkerCanceled    GqlAppointmentStatus = "WORKER_CANCELED"
	AppointmentStatusCompleted         GqlAppointmentStatus = "COMPLETED"
	AppointmentStatusCustomerNoShowUp  GqlAppointmentStatus = "CUSTOMER_NO_SHOW_UP"
	AppointmentStatusWorerNoShowUp     GqlAppointmentStatus = "WORER_NO_SHOW_UP"
)

var AllGqlAppointmentStatus = []GqlAppointmentStatus{
	AppointmentStatusUpcoming,
	AppointmentStatusCustomerCancelled,
	AppointmentStatusWorkerCanceled,
	AppointmentStatusCompleted,
	AppointmentStatusCustomerNoShowUp,
	AppointmentStatusWorerNoShowUp,
}

func (e GqlAppointmentStatus) IsValid() bool {
	switch e {
	case AppointmentStatusUpcoming, AppointmentStatusCustomerCancelled, AppointmentStatusWorkerCanceled, AppointmentStatusCompleted, AppointmentStatusCustomerNoShowUp, AppointmentStatusWorerNoShowUp:
		return true
	}
	return false
}

func (e GqlAppointmentStatus) String() string {
	return string(e)
}

func (e *GqlAppointmentStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GqlAppointmentStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AppointmentStatus", str)
	}
	return nil
}

func (e GqlAppointmentStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GqlJobStatus string

const (
	JobStatusCreated           GqlJobStatus = "CREATED"
	JobStatusAccepted          GqlJobStatus = "ACCEPTED"
	JobStatusCustomerCanceled  GqlJobStatus = "CUSTOMER_CANCELED"
	JobStatusWorkerCalnceled   GqlJobStatus = "WORKER_CALNCELED"
	JobStatusCompleted         GqlJobStatus = "COMPLETED"
	JobStatusClientDidntShowUp GqlJobStatus = "CLIENT_DIDNT_SHOW_UP"
	JobStatusWorkerDidntShowUp GqlJobStatus = "WORKER_DIDNT_SHOW_UP"
)

var AllGqlJobStatus = []GqlJobStatus{
	JobStatusCreated,
	JobStatusAccepted,
	JobStatusCustomerCanceled,
	JobStatusWorkerCalnceled,
	JobStatusCompleted,
	JobStatusClientDidntShowUp,
	JobStatusWorkerDidntShowUp,
}

func (e GqlJobStatus) IsValid() bool {
	switch e {
	case JobStatusCreated, JobStatusAccepted, JobStatusCustomerCanceled, JobStatusWorkerCalnceled, JobStatusCompleted, JobStatusClientDidntShowUp, JobStatusWorkerDidntShowUp:
		return true
	}
	return false
}

func (e GqlJobStatus) String() string {
	return string(e)
}

func (e *GqlJobStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GqlJobStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid JobStatus", str)
	}
	return nil
}

func (e GqlJobStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GqlOfferStatus string

const (
	OfferStatusSent      GqlOfferStatus = "SENT"
	OfferStatusAccepted  GqlOfferStatus = "ACCEPTED"
	OfferStatusCompelted GqlOfferStatus = "COMPELTED"
	OfferStatusCanceled  GqlOfferStatus = "CANCELED"
)

var AllGqlOfferStatus = []GqlOfferStatus{
	OfferStatusSent,
	OfferStatusAccepted,
	OfferStatusCompelted,
	OfferStatusCanceled,
}

func (e GqlOfferStatus) IsValid() bool {
	switch e {
	case OfferStatusSent, OfferStatusAccepted, OfferStatusCompelted, OfferStatusCanceled:
		return true
	}
	return false
}

func (e GqlOfferStatus) String() string {
	return string(e)
}

func (e *GqlOfferStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GqlOfferStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OfferStatus", str)
	}
	return nil
}

func (e GqlOfferStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GqlWorkerSpeciality string

const (
	WorkerSpecialityHandyman   GqlWorkerSpeciality = "HANDYMAN"
	WorkerSpecialityDriver     GqlWorkerSpeciality = "DRIVER"
	WorkerSpecialityAirconspec GqlWorkerSpeciality = "AIRCONSPEC"
	WorkerSpecialityPlumber    GqlWorkerSpeciality = "PLUMBER"
)

var AllGqlWorkerSpeciality = []GqlWorkerSpeciality{
	WorkerSpecialityHandyman,
	WorkerSpecialityDriver,
	WorkerSpecialityAirconspec,
	WorkerSpecialityPlumber,
}

func (e GqlWorkerSpeciality) IsValid() bool {
	switch e {
	case WorkerSpecialityHandyman, WorkerSpecialityDriver, WorkerSpecialityAirconspec, WorkerSpecialityPlumber:
		return true
	}
	return false
}

func (e GqlWorkerSpeciality) String() string {
	return string(e)
}

func (e *GqlWorkerSpeciality) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GqlWorkerSpeciality(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid WorkerSpeciality", str)
	}
	return nil
}

func (e GqlWorkerSpeciality) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
