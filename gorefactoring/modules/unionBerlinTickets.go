package modules

import (
	"log"
	utils "main/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func UnionBerlinTickets(dynamodb *utils.DynamoDbService) *UnionBerlinTicketsData {
	res, err := unionBerlinTickets()
    if err != nil {
		log.Printf("Error in Union Berlin Tickets module: %s", err)
	}
	res.data = dynamodb.GetActualItem("union_berlin_tickets", res.data)
	return res
}

func unionBerlinTickets() (*UnionBerlinTicketsData, error) {
 	url := "https://tickets.union-zeughaus.de/unveu/heimspiele_2.htm"
	doc, err := utils.GetDoc(url)
	if err != nil {
		return nil, err
	}
 
	tickets := "[Union Berlin tickets]("+url+"):\n"
	parentDiv := doc.Find(".ticket.listitem.gamehome")
	parentDiv.Each(func(i int, s *goquery.Selection) {
		elementsStrings := []string {}
		s.Find("h2").Each(func(i int, s *goquery.Selection) {
			elementsStrings = append(elementsStrings, s.Text())
		})
		tickets += strings.Join(elementsStrings, " - ") + "\n"
	})


	return &UnionBerlinTicketsData{tickets}, nil
}


type UnionBerlinTicketsData struct {
	data string
}

func (d *UnionBerlinTicketsData) String() string {
	if d == nil {
		return ""
	}
    return d.data
}
