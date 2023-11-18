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
	fmt.Print(parentDiv)
	parentDiv.Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		title = strings.Replace(title, "  ", "", -1)
		title = strings.Replace(title, "\n", "", -1)
		title = strings.Replace(title, "Hertha BSC", "Hertha BSC - ", -1)
		tickets += fmt.Sprintf("%s\n", title)
	})


	return tickets, nil
}


