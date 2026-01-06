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
//	go run zonefile_export.go -stream "example.com"  # Use streaming mode for large files
//	go run zonefile_export.go -chunk "example.com"  # Use chunk to process data as it comes in
//
// The zone file will be saved to a file named after the zone (e.g., example.com.txt).
// Use -stream flag to download directly to disk without loading into memory (more efficient for large zones).
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client = api.NewClient(doer, api.SetAPIKey(k))
}

func downloadChunk(zoneName string) error {
	log.Println("Downloading zone file (chunking mode)...")
	fileName := fmt.Sprintf("%s.txt", zoneName)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	reader, resp, err := client.Zones.DownloadZonefileReader(zoneName)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // IMPORTANT: Must close when done

	noLines := 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = writer.WriteString(line)
		if err != nil {
			return err
		}

		noLines += 1

		if noLines%10 == 0 {
			log.Printf("Processed %d lines\n", noLines)
		}
	}
	return nil
}

func downloadStream(zoneName string) error {
	// Streaming mode: download directly to disk without loading into memory
	log.Println("Downloading zone file (streaming mode)...")

	fileName := fmt.Sprintf("%s.txt", zoneName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", fileName, err)
	}
	defer file.Close()

	resp, err := client.Zones.DownloadZonefileWriter(zoneName, file)
	if err != nil {
		return err
	}

	// Get filename from Content-Disposition header if available
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
		log.Printf("Content-Disposition: %s\n", contentDisposition)
	}

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", fileName, err)
	}

	log.Printf("Zone file saved to: %s\n", fileName)
	log.Printf("File size: %d bytes\n", fileInfo.Size())

	return nil
}

func downloadBytes(zoneName string) error {
	// Standard mode: load into memory then write to disk
	log.Println("Downloading zone file (standard mode)...")

	buf, resp, err := client.Zones.DownloadZonefile(zoneName)
	if err != nil {
		if errors.Is(err, api.ErrZoneMissing) {
			log.Fatalf("Zone file not found for %s", zoneName)
		}
		return fmt.Errorf("failed to download zone file (standard mode): %v", err)
	}

	// Get filename from Content-Disposition header if available
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
		log.Printf("Content-Disposition: %s\n", contentDisposition)
		// You could parse the filename from the header here if needed
	}

	// Save to file
	fileName := fmt.Sprintf("%s.txt", zoneName)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", fileName, err)
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write zone file to %s: %v", fileName, err)
	}

	log.Printf("Zone file saved to: %s\n", fileName)
	log.Printf("File size: %d bytes\n", buf.Len())
	return nil
}

func main() {
	useStream := flag.Bool("stream", false, "Use streaming mode to download directly to disk (more memory efficient)")
	useChunk := flag.Bool("chunk", false, "Use chunking mode to process data as it streams")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Usage: go run zonefile_export.go [-stream] <zone-name>")
	}

	zoneName := flag.Arg(0)

	log.Printf("Initiating zone file export for %s...\n", zoneName)
	exportStatus, resp, err := client.Zones.ExportZonefile(zoneName)

	jsonBytes, err := json.MarshalIndent(exportStatus, "", "  ")
	if err != nil {
		log.Fatalf("Error processing JSON response: %v", err)
	}
	log.Printf("Received %s; Body: %s\n", resp.Status, jsonBytes)

	msg := fmt.Sprintf("Export initiated. Status: %s", exportStatus.Status)
	if exportStatus.Message != "" {
		msg = fmt.Sprintf("%s; Message: %s", msg, exportStatus.Message)
	}
	log.Println(msg)

	log.Println("Polling for export completion...")
	maxAttempts := 30
	pollInterval := 2 * time.Second

	for i := 0; i < maxAttempts; i++ {
		exportStatus, _, err = client.Zones.GetExportZonefileStatus(zoneName)
		if err != nil {
			if errors.Is(err, api.ErrZoneMissing) {
				log.Fatalf("No export found for zone %s", zoneName)
			}
			log.Fatalf("Failed to get export status (attempt %d/%d): %v", i+1, maxAttempts, err)
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

	switch {
	case *useChunk:
		err := downloadChunk(zoneName)
		if err != nil {
			log.Fatalf("Failed to download file via chunking: %v", err)
		}
	case *useStream:
		err := downloadStream(zoneName)
		if err != nil {
			log.Fatalf("Failed to download file via streaming: %v", err)
		}
	default:
		err := downloadBytes(zoneName)
		if err != nil {
			log.Fatalf("Failed to download file via bytes: %v", err)
		}
	}

	log.Println("First 1500 characters of zone file:")

	filename := fmt.Sprintf("%s.txt", zoneName)

	previewFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open file for preview: %v\n", err)
	} else {
		defer previewFile.Close()

		preview := make([]byte, 1500)
		n, err := previewFile.Read(preview)
		if err != nil && err != io.EOF {
			log.Printf("Warning: Could not read file preview: %v\n", err)
		} else {
			if n >= 1500 {
				fmt.Println(string(preview) + "...")
			} else {
				fmt.Println(string(preview[:n]))
			}
		}
	}
}
