package modules

import (
	"log"
    utils "main/utils"
)

func MobileNumberNotification() string {
    mobileNumberData, err := mobileNumberNotification()
    if err != nil {
       log.Printf("Error in mobile number module: %s", err)
    }
    return mobileNumberData
}

func mobileNumberNotification() (string, error) {
	tele2Url := "https://nnov.tele2.ru/api/shop/products/numbers/bundles/1/groups?query=9524596234&exclude&siteId=siteNNOV"
	isFound := false
	res := mobileNumberResponse{}
	err := utils.DoGet(tele2Url, &res)
	if err != nil {
		return "", err
	 }

	for _, group := range res.Data {
		for _, bound := range group.Groups {
			if len(bound.Bundles) != 0 {
				isFound = true
			}
		} 
	}

	if isFound {
		return "Number was found in Tele2: https://nnov.tele2.ru/shop/number?pageParams=type%3Dchoose%26price%3D0%26search_num%3D9524596234", nil
	} else {
		return "", nil
	}
}

type mobileNumberResponse struct {
    Data []groupResponse `json:"data"`
}

type groupResponse struct {
    Groups []bundleResponse `json:"bundleGroups"`
}

type bundleResponse struct {
	Bundles []any `json:"bundles"`
}


