package modules

import (
	"fmt"
	"log"
	utils "main/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HerthaTickets() string {
	res, err := herthaTickets()
    if err != nil {
		log.Printf("Error in Hertha Tickets module: %s", err)
	 }
	return res
}

func herthaTickets() (string, error) {
 	url := "https://ticket-onlineshop.com/ols/hbsctk/en/tk/"
	doc, err := utils.GetDoc(url)
	if err != nil {
		return "", err
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


	return tickets, nil
}


