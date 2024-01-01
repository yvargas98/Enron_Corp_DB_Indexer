package indexer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ECEmail struct {
	ID                        int    `json:"id"`
	Message_ID                string `json:"message_id"`
	Date                      string `json:"date"`
	From                      string `json:"from"`
	To                        string `json:"to"`
	Subject                   string `json:"subject"`
	Cc                        string `json:"cc"`
	Mime_version              string `json:"mime_version"`
	Content_Type              string `json:"content_type"`
	Content_Transfer_Encoding string `json:"content_transfer_encoding"`
	Bcc                       string `json:"bcc"`
	X_from                    string `json:"x_from"`
	X_to                      string `json:"x_to"`
	X_cc                      string `json:"x_cc"`
	X_bcc                     string `json:"x_bcc"`
	X_folder                  string `json:"x_folder"`
	X_origin                  string `json:"x_origin"`
	X_filename                string `json:"x_filename"`
	Content                   string `json:"content"`
}

func GetFolders(folderName string) []string {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	var folders []string

	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}
	return folders
}

func GetFiles(folderName string) []string {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	var fileNames []string
	for _, file := range files {
		if file.IsDir() == false && file.Name() != ".DS_Store" {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func FormatData(data_lines *bufio.Scanner, id int) ECEmail {
	var data ECEmail
	data.ID = id

	for data_lines.Scan() {
		line := data_lines.Text()
		switch {
		case strings.Contains(line, "Message-ID:"):
			data.Message_ID = line[11:]
		case strings.Contains(line, "Date:"):
			data.Date = line[5:]
		case strings.Contains(line, "From:"):
			data.From = line[5:]
		case strings.Contains(line, "To:"):
			data.To = line[3:]
		case strings.Contains(line, "Subject:"):
			data.Subject = line[8:]
		case strings.Contains(line, "Cc:"):
			data.Cc = line[3:]
		case strings.Contains(line, "Mime-Version:"):
			data.Mime_version = line[9:]
		case strings.Contains(line, "Content-Type:"):
			data.Content_Type = line[9:]
		case strings.Contains(line, "Content-Transfer-Encoding:"):
			data.Content_Transfer_Encoding = line[9:]
		case strings.Contains(line, "X-From:"):
			data.X_from = line[9:]
		case strings.Contains(line, "X-To:"):
			data.X_to = line[9:]
		case strings.Contains(line, "X-cc:"):
			data.X_cc = line[6:]
		case strings.Contains(line, "X-bcc:"):
			data.X_bcc = line[6:]
		case strings.Contains(line, "X-Folder:"):
			data.X_folder = line[9:]
		case strings.Contains(line, "X-Origin:"):
			data.X_origin = line[9:]
		case strings.Contains(line, "X-FileName:"):
			data.X_filename = line[9:]
		default:
			data.Content += line
		}
	}
	return data
}

func PostDataToOpenObserve(data ECEmail) error {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %s", err)
	}

	ZincSearchUrl := os.Getenv("INDEXER_URL")
	ZSusername := os.Getenv("SEARCH_SERVER_USERNAME")
	ZSpassword := os.Getenv("SEARCH_SERVER_PASSWORD")

	req, err := http.NewRequest(http.MethodPost, ZincSearchUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error reading request.", err)
	}

	req.SetBasicAuth(ZSusername, ZSpassword)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
