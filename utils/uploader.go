package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type InjestRequest struct {
	Marketitems []string
	Locationid  string
	Username    string
}

func SendMarketItems(marketItems []string, ingestUrl string, locationId string) {
	client := &http.Client{}

	username := "unknown user"
	if user, err := ioutil.ReadFile("C:\\Users\\Public\\Documents\\username.txt"); err == nil {
		username = string(user)
	}
	
	if user, err := ioutil.ReadFile("/media/username.txt"); err == nil {
		username = string(user)
	}

	injestRequest := InjestRequest{
		Marketitems: marketItems,
		Locationid:  locationId,
		Username:    username
	}

	data, err := json.Marshal(injestRequest)
	log.Printf("%s", data)

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
