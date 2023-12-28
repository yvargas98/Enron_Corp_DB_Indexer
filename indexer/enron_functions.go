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

const (
	ZincSearchUrl = "http://localhost:4080/api/enron_corp/_doc"
	ZSusername    = "admin"
	ZSpassword    = "Complexpass#123"
)

func List_all_folders(folder_name string) []string { //recibe como parámetro el folder "maildir".
	files, err := os.ReadDir(folder_name) //"ioutil.ReadDir" extrae todos los subfolders y los guarda en "files"
	if err != nil {
		log.Fatal(err)
	}
	var folders []string

	for _, file := range files {
		filePath := filepath.Join(folder_name, file.Name())

		// Verificar si es un directorio
		if file.IsDir() {
			// Ignorar directorios que no se pueden abrir
			_, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening directory %s: %s\n", filePath, err)
				continue
			}

			// Si llegamos aquí, el directorio es válido
			folders = append(folders, file.Name())
		}
	}
	return folders
}

// Lista cada uno de los archivos o correos
func List_files(folder_name string) []string {
	files, err := os.ReadDir(folder_name) //https://golang.cafe/blog/how-to-list-files-in-a-directory-in-go.html
	if err != nil {
		log.Fatal(err)
	}
	var files_names []string //array donde se guardarán los nombres de los archivos contenidos en las subcarpetas.
	for _, file := range files {
		filePath := filepath.Join(folder_name, file.Name())

		// Verificar si es un directorio
		if file.IsDir() {
			// Ignorar directorios
			continue
		}

		// Ignorar archivos con el nombre ".DS_Store"
		if file.Name() != ".DS_Store" {
			// Verificar explícitamente si el archivo se puede abrir
			sysFile, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", filePath, err)
				continue
			}
			sysFile.Close()

			// Si llegamos aquí, el archivo es válido
			files_names = append(files_names, file.Name())
		}
	}
	return files_names
}

func FormatData(data_lines *bufio.Scanner, id int) ECEmail { //parse_data
	var data ECEmail
	for data_lines.Scan() {
		data.ID = id
		if strings.Contains(data_lines.Text(), "Message-ID:") {
			data.Message_ID = data_lines.Text()[11:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 11)
		} else if strings.Contains(data_lines.Text(), "Date:") {
			data.Date = data_lines.Text()[5:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 5)
		} else if strings.Contains(data_lines.Text(), "From:") {
			data.From = data_lines.Text()[5:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 6)
		} else if strings.Contains(data_lines.Text(), "To:") {
			data.To = data_lines.Text()[3:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 3)
		} else if strings.Contains(data_lines.Text(), "Subject:") {
			data.Subject = data_lines.Text()[8:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 8)
		} else if strings.Contains(data_lines.Text(), "Cc:") {
			data.Cc = data_lines.Text()[3:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 3)
		} else if strings.Contains(data_lines.Text(), "Mime-Version:") {
			data.Mime_version = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else if strings.Contains(data_lines.Text(), "Content-Type:") {
			data.Content_Type = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else if strings.Contains(data_lines.Text(), "Content-Transfer-Encoding:") {
			data.Content_Transfer_Encoding = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else if strings.Contains(data_lines.Text(), "X-From:") {
			data.X_from = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else if strings.Contains(data_lines.Text(), "X-To:") {
			data.X_to = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else if strings.Contains(data_lines.Text(), "X-cc:") {
			data.X_cc = data_lines.Text()[6:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 6)
		} else if strings.Contains(data_lines.Text(), "X-bcc:") {
			data.X_bcc = data_lines.Text()[6:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 6)
		} else if strings.Contains(data_lines.Text(), "X-Folder:") {
			data.X_folder = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else if strings.Contains(data_lines.Text(), "X-Origin:") {
			data.X_origin = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else if strings.Contains(data_lines.Text(), "X-FileName:") {
			data.X_filename = data_lines.Text()[9:]
			fmt.Println("data_lines.Text():", data_lines.Text())
			fmt.Println("len(data_lines.Text()):", len(data_lines.Text()))
			fmt.Println("índice utilizado:", 9)
		} else {
			data.Content = data.Content + data_lines.Text()
		}
	}
	return data
}

func PostDataToZincSearch(data ECEmail) { //index_data
	jsonData, _ := json.MarshalIndent(data, "", "   ")
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
