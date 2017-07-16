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
	Username    string
}

func SendMarketItems(marketItems []string, config ClientConfig, locationId string) {
	client := &http.Client{}

	injestRequest := InjestRequest{
		Marketitems: marketItems,
		Locationid:  locationId,
		Username:    config.Username,
	}

	data, err := json.Marshal(injestRequest)
	//	log.Printf("%v", string(data))

	if err != nil {
		log.Printf("Error while marshalling payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", config.IngestUrl, bytes.NewBuffer([]byte(string(data))))
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

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Printf("Got bad response code: %v", resp.StatusCode)
		return
	}

	log.Printf("Sent market payload with %v entries.", len(marketItems))

	resp.Body.Close()
}
