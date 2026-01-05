// zonefile_export.go demonstrates how to export a DNS zone file from NS1.
//
// This example shows the complete workflow for exporting a zone file:
//  1. Initiate the zone file export
//  2. Poll for export completion status
//  3. Download the completed zone file
//  4. Writes the zonefile to cwd, then prints the first 1500 characters to stdout
//
// Usage:
//
//	export NS1_APIKEY="your-api-key-here"
//	go run zonefile_export.go "example.com"
//
// The zone file will be saved to a file named after the zone (e.g., example.com.txt).
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
)

var client *api.Client

// Helper that initializes rest api client from environment variable.
func init() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		log.Fatal("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client = api.NewClient(doer, api.SetAPIKey(k))
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run zonefile_export.go <zone-name>")
	}

	zoneName := os.Args[1]

	// Initiate the zone file export
	log.Printf("Initiating zone file export for %s...\n", zoneName)
	exportStatus, resp, err := client.Zones.ExportZonefile(zoneName)

	// print http response data
	jsonBytes, _ := json.MarshalIndent(exportStatus, "", "  ")
	log.Printf("Received %s; Body: %s\n", resp.Status, jsonBytes)

	if err != nil {
		if errors.Is(err, api.ErrZoneMissing) {
			log.Fatalf("Zone %s not found", zoneName)
		}
		log.Fatal(err)
	}

	msg := fmt.Sprintf("Export initiated. Status: %s", exportStatus.Status)
	if exportStatus.Message != "" {
		msg = fmt.Sprintf("%s; Message: %s", msg, exportStatus.Message)
	}
	log.Println(msg)

	// Poll for export completion
	log.Println("Polling for export completion...")
	maxAttempts := 30
	pollInterval := 2 * time.Second

	for i := 0; i < maxAttempts; i++ {
		exportStatus, _, err = client.Zones.GetExportZonefileStatus(zoneName)
		if err != nil {
			if errors.Is(err, api.ErrZoneMissing) {
				log.Fatalf("No export found for zone %s", zoneName)
			}
			log.Fatal(err)
		}

		msg = fmt.Sprintf("Attempt %d/%d - Status: %s", i+1, maxAttempts, exportStatus.Status)
		if exportStatus.Message != "" {
			msg = fmt.Sprintf("%s - %s", msg, exportStatus.Message)
		}
		log.Print(msg)

		if exportStatus.Status == "COMPLETED" {
			log.Printf("Export completed at: %s\n", exportStatus.GeneratedAt)
			break
		}
		if exportStatus.Status == "FAILED" {
			log.Fatalf("Export failed: %s", exportStatus.Message)
		}
		if slices.Contains([]string{"QUEUED", "GENERATING"}, exportStatus.Status) {
			// Continue polling
			time.Sleep(pollInterval)
			continue
		}
		log.Fatalf("Unknown exportStatus: %s", exportStatus.Status)
	}

	if exportStatus.Status != "COMPLETED" {
		log.Fatal("Export did not complete within the expected time")
	}

	// Download the zone file
	log.Println("Downloading zone file...")
	buf, resp, err := client.Zones.DownloadZonefile(zoneName)
	if err != nil {
		if errors.Is(err, api.ErrZoneMissing) {
			log.Fatalf("Zone file not found for %s", zoneName)
		}
		log.Fatal(err)
	}

	// Get filename from Content-Disposition header if available
	filename := fmt.Sprintf("%s.txt", zoneName)
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
		log.Printf("Content-Disposition: %s\n", contentDisposition)
		// You could parse the filename from the header here if needed
	}

	// Save to file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Zone file saved to: %s\n", filename)
	log.Printf("File size: %d bytes\n", buf.Len())

	log.Println("\nFirst 1500 characters of zone file:")
	content := buf.String()
	if len(content) > 1500 {
		fmt.Println(content[:1500] + "...")
	} else {
		fmt.Println(content)
	}
}
