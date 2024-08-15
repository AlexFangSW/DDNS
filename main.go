package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

var (
	zoneID   string = os.Getenv("ZONE_ID")
	token    string = os.Getenv("TOKEN")
	recordID string = os.Getenv("RECORD_ID")
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

type ipiyReturn struct {
	IP string `json:"ip"`
}

func drainAndClose(body io.ReadCloser) error {
	_, err := io.Copy(io.Discard, body)
	if err != nil {
		return fmt.Errorf("drainAndClose: drain body failed: %w", err)
	}
	return body.Close()
}

func getCurrentIPv4() (currIP string, oErr error) {
	ret, err := httpClient.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", fmt.Errorf("getCurrentIPv4: ipify request failed: %w", err)
	}

	defer func() {
		oErr = errors.Join(oErr, drainAndClose(ret.Body))
	}()

	body, err := io.ReadAll(ret.Body)
	if err != nil {
		return "", fmt.Errorf("getCurrentIPv4: read request body failed: %w", err)
	}

	data := ipiyReturn{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("getCurrentIPv4: parse request body failed: %w", err)
	}

	log.Printf("current public ipv4: %q", data.IP)
	return data.IP, nil
}

func updateDNSRecord() error {
	// check necessary enviroment variables
	if zoneID == "" {
		return errors.New("updateDNSRecord: please provide ZONE_ID enviroment variable")
	}
	if token == "" {
		return errors.New("updateDNSRecord: please provide TOKEN enviroment variable")
	}

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return fmt.Errorf("updateDNSRecord: create new api instance failed: %w", err)
	}

	// get current public ip (ipv4)
	currIP, err := getCurrentIPv4()
	if err != nil {
		return fmt.Errorf("updateDNSRecord: get current ip failed: %w", err)
	}

	// check current ip
	identifier := cloudflare.ZoneIdentifier(zoneID)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	record, err := api.GetDNSRecord(ctxTimeout, identifier, recordID)
	if err != nil {
		return fmt.Errorf("updateDNSRecord: get dns record failed: %w", err)
	}

	log.Printf(
		"current record, type: %q, name: %q, content: %q",
		record.Type, record.Name, record.Content)

	if record.Type != "A" {
		return fmt.Errorf("must be an 'A' record, recived type %q", record.Type)
	}
	if record.Content == currIP {
		log.Println("same ip, no need to change")
		return nil
	}

	// update dns record
	recordParams := cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Content: currIP,
	}

	newRecord, err := api.UpdateDNSRecord(ctxTimeout, identifier, recordParams)
	if err != nil {
		return fmt.Errorf("updateDNSRecord: update dns record failed: %w", err)
	}

	log.Printf(
		"new record, type: %q, name: %q, content: %q",
		newRecord.Type, newRecord.Name, newRecord.Content)

	return nil
}

func main() {
	for {
		log.Println("start updating dns record")

		if err := updateDNSRecord(); err != nil {
			log.Fatal(err)
		}

		log.Println("finish updating dsn record")

		log.Println("sleeping for 5 minutes...")
		time.Sleep(5 * time.Minute)
	}
}
