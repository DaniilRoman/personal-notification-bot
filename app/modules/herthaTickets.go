package modules

import (
	"fmt"
	"log"
	utils "main/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HerthaTickets(dynamodb *utils.DynamoDbService) *HerthaTicketsData {
	res, err := herthaTickets()
	if err != nil {
		log.Printf("Error in Hertha Tickets module: %s", err)
	}
	res.Data = dynamodb.GetValueIfChanged("hertha_tickets", res.Data)
	return res
}

const herthaTicketsUrl = "https://ticket-onlineshop.com/ols/hbsctk/en/tk/"

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

type HerthaTicketsData struct {
	Data string
}

func (d *HerthaTicketsData) String() string {
	if d == nil || d.Data == "" {
		return ""
	}
	return fmt.Sprintf("[Hertha Berlin tickets](%s):\n%s", herthaTicketsUrl, d.Data)
}
