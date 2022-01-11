package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DBerror error

type CustomerHandler struct {
}

type CognitoIdentity struct {
	Sub      string             `json:"sub"`
	Username string             `json:"username"`
	Claims   *map[string]string `json:"claims"`
	SourceIp []string           `json:"sourceIp"`
}

//Info carries meta-data about the request, such as the fields selected by the client, parent type name and field
type Info struct {
	FieldName           string                 `json:"fieldName"`
	ParentTypeName      string                 `json:"parentTypeName"`
	Variables           map[string]interface{} `json:"variables"`
	SelectionSetList    []string               `json:"selectionSetList"`
	SelectionSetGraphQL string                 `json:"selectionSetGraphQL"`
}

//Request is the parent type that gets passed by graphQl Api Server
type Request struct {
	Arguments map[string]interface{} `json:"arguments"`
	Source    map[string]interface{} `json:"source"`
	Identity  CognitoIdentity        `json:"identity"`
	Info      Info                   `json:"info"`
}

//========== Arguments type for specific field operations ===============

//DealsByInquiryIdArguments has the arguments passed by user with graphQl call
type CustomerByIdArguments struct {
	CustomerId uuid.UUID `json:"customerId"`
}

type CreateCustomerArguments struct {
	CreateCustomerInput GqlCreateCustomerInput `json:"createCustomerInput"`
}

type UpdateCustomerArguments struct {
	UpateCustomerInput UpdateCustomerInput `json:"updateCustomerInput"`
}
type DeleteCustomerArguments struct {
	Id uuid.UUID `json:"id"`
}

func (c *CustomerHandler) CreateCustomer(req *Request) (GqlCustomer, error) {
	var gqlCstmr GqlCustomer
	// var dbAdm autoMatesOrm.Customer

	args := CreateCustomerArguments{}
	argsByts, err := json.Marshal(req.Arguments)
	if err != nil {
		return gqlCstmr, err
	}
	err = json.Unmarshal(argsByts, &args)
	if err != nil {
		return gqlCstmr, err
	}

	dbCstmr, err := InsertCustomerToDB(&args)
	if err != nil {
		return gqlCstmr, err
	}

	gqlCstmr, err = ConvertDBCustomerToGqlCustomer(dbCstmr)
	if err != nil {
		return gqlCstmr, err
	}

	return gqlCstmr, nil
}

func (c *CustomerHandler) UpdateCustomer(req *Request) (GqlCustomer, error) {
	var gqlCstmr GqlCustomer
	// var dbAdm autoMatesOrm.Customer

	args := CreateCustomerArguments{}
	argsByts, err := json.Marshal(req.Arguments)
	if err != nil {
		return gqlCstmr, err
	}
	err = json.Unmarshal(argsByts, &args)
	if err != nil {
		return gqlCstmr, err
	}

	dbCstmr, err := InsertCustomerToDB(&args)
	if err != nil {
		return gqlCstmr, err
	}

	gqlCstmr, err = ConvertDBCustomerToGqlCustomer(dbCstmr)
	if err != nil {
		return gqlCstmr, err
	}

	//implement insert admin logic

	//raise an event to event eventBridge

	fmt.Println("my customerId", gqlCstmr)

	return gqlCstmr, nil
}

func (c *CustomerHandler) DeleteCustomer(req *Request) (GqlCustomer, error) {
	var gqlCstmr GqlCustomer
	var dbCstmr Customer

	args := DeleteCustomerArguments{}
	argsByts, err := json.Marshal(req.Arguments)
	if err != nil {
		return gqlCstmr, err
	}
	err = json.Unmarshal(argsByts, &args)
	if err != nil {
		return gqlCstmr, err
	}

	err = DB.Where("customer_id = ?", args.Id).Delete(&dbCstmr).Error
	if err != nil {
		return gqlCstmr, err
	}

	return gqlCstmr, nil
}

func (c *CustomerHandler) CustomerMuationHandler(req *Request) (GqlCustomer, error) {
	gqlCstmr := GqlCustomer{}
	var err error

	switch req.Info.FieldName {
	case "createCustomer":
		gqlCstmr, err = c.CreateCustomer(req)
		if err != nil {
			return gqlCstmr, err
		}
	case "updateCustomer":
		gqlCstmr, err = c.UpdateCustomer(req)
		if err != nil {
			return gqlCstmr, err
		}
	case "deleteCustomer":
		gqlCstmr, err = c.DeleteCustomer(req)
		if err != nil {
			return gqlCstmr, err
		}
	default:
		errMessage := fmt.Sprintf("handler error: field name %v is unkown", req.Info.FieldName)
		err = errors.New(errMessage)
		return gqlCstmr, err
	}

	return gqlCstmr, nil

}

func (c *CustomerHandler) CustomerQueryHandler(req *Request) (GqlCustomer, error) {
	var gqlCstmr GqlCustomer
	var dbCstmr Customer

	args := CustomerByIdArguments{}
	argsByts, err := json.Marshal(req.Arguments)
	if err != nil {
		return gqlCstmr, err
	}
	err = json.Unmarshal(argsByts, &args)
	if err != nil {
		return gqlCstmr, err
	}

	err = DB.First(&dbCstmr, "customer_id = ?", args.CustomerId).Error
	if err != nil {
		return gqlCstmr, err
	}

	gqlCstmr, err = ConvertDBCustomerToGqlCustomer(dbCstmr)
	if err != nil {
		return gqlCstmr, err
	}

	return gqlCstmr, nil
}

func (c *CustomerHandler) CustomerResolver(ctx context.Context, req Request) (*GqlCustomer, error) {

	cstmr := GqlCustomer{}
	var err error

	switch req.Info.ParentTypeName {
	case "Query":
		cstmr, err = c.CustomerQueryHandler(&req)
		if err != nil {
			return nil, err
		}
	case "Mutation":
		cstmr, err = c.CustomerMuationHandler(&req)
		if err != nil {
			return nil, err
		}
	default:
		errMessage := fmt.Sprintf("handler error: type %v is unkown", req.Info.ParentTypeName)
		err = errors.New(errMessage)
		return nil, err
	}

	return &cstmr, nil

}

func main() {

	DB, DBerror = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=database-1.cdqfiq1lbrhl.ap-southeast-1.rds.amazonaws.com user=fyp2 password=Hnaas0916421820_97 dbname=fyp2 port=5432 sslmode=disable search_path=naasfyp",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if DBerror != nil {
		panic("Error opening connection with db")
	}

	fmt.Println("i the customer resolver go invoked")
	c := CustomerHandler{}

	lambda.Start(c.CustomerResolver)

}

// ====================== helpers ======================
func ConvertDBCustomerToGqlCustomer(dbCstmr Customer) (GqlCustomer, error) {

	// gqlAddrs := GqlLocation{

	// }

	gqlCstmr := GqlCustomer{
		ID:            dbCstmr.CustomerID.String(),
		FName:         dbCstmr.NameFirstName,
		LName:         dbCstmr.NameLastName,
		Email:         dbCstmr.CustomerEmail,
		PhoneNo:       string(dbCstmr.CustomerPhoneNumber),
		PostalZipCode: &dbCstmr.CustomerPostalCode,
		// Appointments: ,
		// JobRequests: ,

	}

	fmt.Println("ConvertDBCustomerToGqlCustomer end", gqlCstmr)

	return gqlCstmr, nil
}

func InsertCustomerToDB(args *CreateCustomerArguments) (Customer, error) {

	dbCstmr := Customer{
		CustomerID:    uuid.New(),
		NameFirstName: args.CreateCustomerInput.FName,
		NameLastName:  args.CreateCustomerInput.LName,
		CustomerEmail: args.CreateCustomerInput.Email,
		// CustomerPhoneNumber: int16(args.createCustomerInput.PhoneNo),
		CustomerPostalCode: *args.CreateCustomerInput.PostalZipCode,
	}

	fmt.Println("customerId from insert function", dbCstmr.CustomerID)

	err := DB.Create(&dbCstmr).Error
	if err != nil {
		return dbCstmr, err
	}

	return dbCstmr, nil
}
