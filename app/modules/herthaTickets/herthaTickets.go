package hertha

import (
	"fmt"
	"log"
	"main/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const herthaTicketsUrl = "https://ticket-onlineshop.com/ols/hbsctk/en/tk/"

func HerthaTickets(dynamodb *utils.DynamoDbService) *HerthaTicketsData {
	res, err := herthaTickets()
	if err != nil {
		log.Printf("Error in Hertha Tickets module: %s", err)
	}
	if res == nil {
		return nil
	}
	res.Data = dynamodb.GetValueIfChanged("hertha_tickets", res.Data)
	return res
}

func herthaTickets() (*HerthaTicketsData, error) {
	doc, err := utils.GetDoc(herthaTicketsUrl)
	if err != nil {
		return nil, err
	}

	tickets := ""
	parentDiv := doc.Find(".event-card__headings")
	parentDiv.Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		title = strings.ReplaceAll(title, "  ", "")
		title = strings.ReplaceAll(title, "\n", "")
		title = strings.ReplaceAll(title, "Hertha BSC", "Hertha BSC - ")
		tickets += fmt.Sprintf("%s\n", title)
	})

	return &HerthaTicketsData{tickets}, nil
}
