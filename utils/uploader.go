package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type InjestRequest struct {
	Marketitems []string
	Locationid  string
}

func SendMarketItems(marketItems []string, ingestUrl string, locationId string) {
	client := &http.Client{}
	log.Printf("data location %v", locationId)
	injestRequest := InjestRequest{
		Marketitems: marketItems,
		Locationid:  locationId,
	}

	data, err := json.Marshal(injestRequest)
	//	log.Printf("%s", data)
	if err != nil {
		log.Printf("Error while marshalling payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", ingestUrl, bytes.NewBuffer([]byte(string(data))))
	if err != nil {
		log.Printf("Error while create new reqest: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error while sending market data: %v", err)
		return
	}

	if resp.StatusCode != 201 {
		log.Printf("Got bad response code: %v", resp.StatusCode)
		return
	}

	log.Printf("Sent market payload with %v entries.", len(marketItems))

	defer resp.Body.Close()
}
