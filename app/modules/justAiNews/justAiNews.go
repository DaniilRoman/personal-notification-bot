package word

import (
	"fmt"
	"log"
	utils "main/utils"
	"strconv"
)

const JustAiNewsPageKey = "JustAiNewsPage"

func JustAiNews(dynamodb *utils.DynamoDbService) *JustAiNewsData {
	res, err := justAiNews(dynamodb)
    if err != nil {
		log.Printf("Error in JustAI news module: %s", err)
	 }
	return res
}

func constructJustAiUrl(currentPage int, dynamodb *utils.DynamoDbService) string {
	const justAiNewsURL = "https://just-ai.com/rus/mail"
	return fmt.Sprintf("%s%d", justAiNewsURL, currentPage)
}

func justAiNews(dynamodb *utils.DynamoDbService) (*JustAiNewsData, error) {
	var currentPage, err = getCurrentPage(dynamodb)
	if err != nil {
		return nil, err
	}
	var justAiNewsURL = constructJustAiUrl(currentPage, dynamodb)

	statusCode, err := utils.GetStatusCode(justAiNewsURL)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		log.Printf("Status code for JustAI news is not 200. The code is: %d", statusCode)
		return nil, nil
	}

	incrementPage(currentPage, dynamodb)

	res := fmt.Sprintf("[New JustAI news've come](%s)\n", justAiNewsURL)
	return &JustAiNewsData{res}, nil
}

func incrementPage(value int, dynamodb *utils.DynamoDbService) {
	dynamodb.SaveItem(JustAiNewsPageKey, strconv.Itoa(value))
}

func getCurrentPage(dynamodb *utils.DynamoDbService) (int, error) {
	pageStr := dynamodb.GetItem(JustAiNewsPageKey)
	return strconv.Atoi(pageStr)
}
