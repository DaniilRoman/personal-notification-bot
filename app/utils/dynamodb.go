package utils

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDbService struct {
	DynamoDB *dynamodb.DynamoDB

	appName       string
	itemTableName string

	blogsStatsTableName  string
	popularWordsStatName string
}

type item struct {
	AppName   string `json:"app_name"`
	ItemName  string `json:"item_name"`
	ItemValue string `json:"item_value"`
}

type blogsStats struct {
	StatName  string `json:"stat_name"`
	Date      string `json:"date"`
	StatValue string `json:"stat_value"`
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
		appName:              "person-notification-bot",
		itemTableName:        "items",
		DynamoDB:             dynamoDB,
		blogsStatsTableName:  "blogs_stats",
		popularWordsStatName: "popular_words",
	}
}

func (item *DynamoDbService) GetValueIfChanged(key, newValue string) string {
	prevValue := item.GetItem(key)
	if strings.TrimSpace(prevValue) != strings.TrimSpace(newValue) {
		item.SaveItem(key, newValue)
		return newValue
	} else {
		return ""
	}
}

func (si *DynamoDbService) createItemTable() error {
	input := &dynamodb.CreateTableInput{
		TableName: &si.itemTableName,
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
				log.Println("Table", si.itemTableName, "already exists")
				return nil
			}
		}
		return err
	}

	log.Println("Table", si.itemTableName, "created successfully")
	return nil
}

func (si *DynamoDbService) createBlogsStatsTable() error {
	input := &dynamodb.CreateTableInput{
		TableName: &si.blogsStatsTableName,
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("stat_name"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("date"), KeyType: aws.String("RANGE")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("stat_name"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("date"), AttributeType: aws.String("S")},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	}

	_, err := si.DynamoDB.CreateTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceInUseException {
				log.Println("Table", si.blogsStatsTableName, "already exists")
				return nil
			}
		}
		return err
	}

	log.Println("Table", si.blogsStatsTableName, "created successfully")
	return nil
}

func (si *DynamoDbService) GetItem(itemName string) string {
	result, err := si.DynamoDB.GetItem(&dynamodb.GetItemInput{
		TableName: &si.itemTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"app_name":  {S: aws.String(si.appName)},
			"item_name": {S: aws.String(itemName)},
		},
	})

	if err != nil {
		log.Printf("get item operation failed: %w", err)
		return ""
	}

	if result.Item == nil {
		log.Printf("Item is nil with the name: %s", itemName)
		return ""
	}

	var storedItem item
	err = dynamodbattribute.UnmarshalMap(result.Item, &storedItem)
	if err != nil {
		log.Printf("unmarshal item operation failed: %w", err)
		return ""
	}

	return storedItem.ItemValue
}

func (si *DynamoDbService) SaveItem(key, value string) {
	item := item{si.appName, key, value}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Printf("unmarshal item operation failed: %w", err)
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &si.itemTableName,
	}

	_, err = si.DynamoDB.PutItem(input)
	if err != nil {
		log.Printf("save item operation failed: %w", err)
	}
}

func (si *DynamoDbService) SavePopularWords(key, value string) {
	stat := blogsStats{si.popularWordsStatName, key, value}
	av, err := dynamodbattribute.MarshalMap(stat)
	if err != nil {
		log.Printf("unmarshal stat operation failed: %w", err)
		return
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &si.blogsStatsTableName,
	}

	_, err = si.DynamoDB.PutItem(input)
	if err != nil {
		log.Printf("save stat operation failed: %w", err)
	}
}

func (service *DynamoDbService) AppendPopularWords(key, value string) {
	currentValue := service.GetItem(key)
	value = strings.Join([]string{value, currentValue}, ",")
	service.SavePopularWords(key, value)
}

func (service *DynamoDbService) GetStatsFromPrevDays(lastDays []string) string {
	keys := []map[string]*dynamodb.AttributeValue{}
	for _, lastDay := range lastDays {
		value := map[string]*dynamodb.AttributeValue{
			"stat_name": {
				S: aws.String(service.popularWordsStatName),
			},
			"date": {
				S: aws.String(lastDay),
			},
		}
		keys = append(keys, value)
	}

	items := map[string]*dynamodb.KeysAndAttributes{
		service.blogsStatsTableName: {Keys: keys},
	}
	result, err := service.DynamoDB.BatchGetItem(
		&dynamodb.BatchGetItemInput{RequestItems: items},
	)

	if err != nil {
		log.Printf("operation failed: %w", err)
		return ""
	}

	words := []string{}
	for _, items := range result.Responses {
        for _, item := range items {
			var blogsStats blogsStats
			err = dynamodbattribute.UnmarshalMap(item, &blogsStats)
			if err != nil {
				log.Printf("operation failed: %w", err)
			} else {
				words = append(words, blogsStats.StatValue)
			}
		}	
	}

	return strings.Join(words, ",")
}

func (si *DynamoDbService) GetBlogsStat(date string) string {
	result, err := si.DynamoDB.GetItem(&dynamodb.GetItemInput{
		TableName: &si.blogsStatsTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"stat_name":  {S: aws.String(si.popularWordsStatName)},
			"date": {S: aws.String(date)},
		},
	})

	if err != nil {
		log.Printf("get blog stats operation failed: %w", err)
		return ""
	}

	if result.Item == nil {
		log.Printf("Blog stat is nil for a date: %s", date)
		return ""
	}

	var stats blogsStats
	err = dynamodbattribute.UnmarshalMap(result.Item, &stats)
	if err != nil {
		log.Printf("operation failed: %w", err)
		return ""
	}

	return stats.StatValue
}