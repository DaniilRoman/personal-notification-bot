package modules

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDbService struct {
	AppName    string
	TableName  string
	DynamoDB   *dynamodb.DynamoDB
}

type item struct {
	AppName    string `json:"app_name"`
	ItemName   string `json:"item_name"`
	ItemValue  string `json:"item_value"`
}

func NewDynamoDbService(accessKey, secretAccessKey, region string, endpointUrl *string) *DynamoDbService {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    endpointUrl,
		Credentials: credentials.NewStaticCredentials(accessKey, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatal("Failed to create AWS session:", err)
	}

	dynamoDB := dynamodb.New(session)

	return &DynamoDbService{
		AppName:    "person-notification-bot",
		TableName:  "items",
		DynamoDB:   dynamoDB,
	}
}

func (si *DynamoDbService) createItemTable() error {
	input := &dynamodb.CreateTableInput{
		TableName: &si.TableName,
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("app_name"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("item_name"), KeyType: aws.String("RANGE")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("app_name"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("item_name"), AttributeType: aws.String("S")},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}

	_, err := si.DynamoDB.CreateTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceInUseException {
				log.Println("Table", si.TableName, "already exists")
				return nil
			}
		}
		return err
	}

	log.Println("Table", si.TableName, "created successfully")
	return nil
}

func (si *DynamoDbService) getItem(itemName string) (string) {
	result, err := si.DynamoDB.GetItem(&dynamodb.GetItemInput{
		TableName: &si.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"app_name":  {S: aws.String(si.AppName)},
			"item_name": {S: aws.String(itemName)},
		},
	})

	if err != nil {
		log.Printf("operation failed: %w", err)
		return ""
	}

	if result.Item == nil {
		log.Printf("operation failed: %w", err)
		return ""
	}

	var storedItem item
	err = dynamodbattribute.UnmarshalMap(result.Item, &storedItem)
	if err != nil {
		log.Printf("operation failed: %w", err)
		return ""
	}

	return storedItem.ItemValue
}

func (si *DynamoDbService) saveItem(key, value string) {
	item := item{si.AppName, key, value}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Printf("operation failed: %w", err)
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &si.TableName,
	}

	_, err = si.DynamoDB.PutItem(input)
	if err != nil {
		log.Printf("operation failed: %w", err)
	}
}

func (item *DynamoDbService) GetActualItem(key, newValue string) string {
	prevValue := item.getItem(key)
	if prevValue != newValue {
		item.saveItem(key, newValue)
		return newValue
	} else {
		return prevValue
	}
}
