package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type InjestRequest struct {
	Marketitems []string
	Locationid  string
	Username    string
}

func SendMarketItems(marketItems []string, config ClientConfig, locationId string) {
	publicItems := make([]string, 0)
	playerItems := make([]string, 0)

	buyerString := fmt.Sprintf("\"BuyerName\":\"%s\"", config.Username)
	sellerString := fmt.Sprintf("\"SellerName\":\"%s\"", config.Username)

	for _, item := range marketItems {
		buyMatched, _ := regexp.MatchString(buyerString, item)
		sellMatched, _ := regexp.MatchString(sellerString, item)
		if buyMatched || sellMatched {
			playerItems = append(playerItems, item)
		} else {
			publicItems = append(publicItems, item)
		}
	}

	if len(publicItems) > 0 {
		SendMarketItemsToEndpoint(publicItems, config.Username, locationId, config.MarketIngestUrl)
	}

	if len(playerItems) > 0 {
		SendMarketItemsToEndpoint(playerItems, config.Username, locationId, config.PlayerIngestUrl)
	}
}

func SendMarketItemsToEndpoint(marketItems []string, username string, locationId string, url string) {
	injestRequest := InjestRequest{
		Marketitems: marketItems,
		Locationid:  locationId,
		Username:    username,
	}

	data, err := json.Marshal(injestRequest)
	//	log.Printf("%v", string(data))

	if err != nil {
		log.Printf("Error while marshalling payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(string(data))))
	if err != nil {
		log.Printf("Error while create new reqest: %v", err)
		return
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error while sending market data: %v", err)
		return
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Printf("Got bad response code: %v", resp.StatusCode)
		return
	}

	log.Printf("Sent market payload with %v entries.", len(marketItems))

	resp.Body.Close()
}
