package union

import (
	"log"
	"main/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const unionBerlinUrl = "https://tickets.union-zeughaus.de/unveu/heimspiele_2.htm"

func UnionBerlinTickets(dynamodb *utils.DynamoDbService) *UnionBerlinTicketsData {
	res, err := unionBerlinTickets()
	if err != nil {
		log.Printf("Error in Union Berlin Tickets module: %s", err)
	}
	res.data = dynamodb.GetValueIfChanged("union_berlin_tickets", res.data)
	return res
}

func unionBerlinTickets() (*UnionBerlinTicketsData, error) {
	doc, err := utils.GetDoc(unionBerlinUrl)
	if err != nil {
		return nil, err
	}

	tickets := ""
	parentDiv := doc.Find(".ticket.listitem.gamehome")
	parentDiv.Each(func(i int, s *goquery.Selection) {
		elementsStrings := []string{}
		s.Find("h2").Each(func(i int, s *goquery.Selection) {
			elementsStrings = append(elementsStrings, s.Text())
		})
		tickets += strings.Join(elementsStrings, " - ") + "\n"
	})

	return &UnionBerlinTicketsData{tickets}, nil
}
