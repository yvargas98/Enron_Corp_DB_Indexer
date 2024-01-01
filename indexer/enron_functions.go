package indexer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ECEmail struct {
	ID                        int    `json:"ID"`
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

func GetFolders(folder_name string) []string {
	files, err := os.ReadDir(folder_name)
	if err != nil {
		log.Fatal(err)
	}
	var folders []string

	for _, file := range files {
		filePath := filepath.Join(folder_name, file.Name())

		if file.IsDir() {
			_, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening directory %s: %s\n", filePath, err)
				continue
			}

			folders = append(folders, file.Name())
		}
	}
	return folders
}

func GetFiles(folder_name string) []string {
	files, err := os.ReadDir(folder_name)
	if err != nil {
		log.Fatal(err)
	}
	var fileNames []string
	for _, file := range files {
		filePath := filepath.Join(folder_name, file.Name())

		if file.IsDir() {
			continue
		}

		if file.Name() != ".DS_Store" {
			sysFile, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", filePath, err)
				continue
			}
			sysFile.Close()
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func FormatData(data_lines *bufio.Scanner, id int) ECEmail {
	var data ECEmail
	for data_lines.Scan() {
		data.ID = id
		if strings.Contains(data_lines.Text(), "Message-ID:") {
			data.Message_ID = data_lines.Text()[11:]
		} else if strings.Contains(data_lines.Text(), "Date:") {
			data.Date = data_lines.Text()[5:]
		} else if strings.Contains(data_lines.Text(), "From:") {
			data.From = data_lines.Text()[5:]
		} else if strings.Contains(data_lines.Text(), "To:") {
			data.To = data_lines.Text()[3:]
		} else if strings.Contains(data_lines.Text(), "Subject:") {
			data.Subject = data_lines.Text()[8:]
		} else if strings.Contains(data_lines.Text(), "Cc:") {
			data.Cc = data_lines.Text()[3:]
		} else if strings.Contains(data_lines.Text(), "Mime-Version:") {
			data.Mime_version = data_lines.Text()[9:]
		} else if strings.Contains(data_lines.Text(), "Content-Type:") {
			data.Content_Type = data_lines.Text()[9:]
		} else if strings.Contains(data_lines.Text(), "Content-Transfer-Encoding:") {
			data.Content_Transfer_Encoding = data_lines.Text()[9:]
		} else if strings.Contains(data_lines.Text(), "X-From:") {
			data.X_from = data_lines.Text()[9:]
		} else if strings.Contains(data_lines.Text(), "X-To:") {
			data.X_to = data_lines.Text()[9:]
		} else if strings.Contains(data_lines.Text(), "X-cc:") {
			data.X_cc = data_lines.Text()[6:]
		} else if strings.Contains(data_lines.Text(), "X-bcc:") {
			data.X_bcc = data_lines.Text()[6:]
		} else if strings.Contains(data_lines.Text(), "X-Folder:") {
			data.X_folder = data_lines.Text()[9:]
		} else if strings.Contains(data_lines.Text(), "X-Origin:") {
			data.X_origin = data_lines.Text()[9:]
		} else if strings.Contains(data_lines.Text(), "X-FileName:") {
			data.X_filename = data_lines.Text()[9:]
		} else {
			data.Content = data.Content + data_lines.Text()
		}
	}
	return data
}

func PostDataToZincSearch(data ECEmail) {
	jsonData, _ := json.MarshalIndent(data, "", "   ")
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
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
