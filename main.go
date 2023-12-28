package main

import (
	"Enron_Corp_DB_Indexer/indexer"
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Enron Corp Directory DB Path is missing.")
		return
	}

	path := os.Args[1] + "/maildir/"

	id := 0
	userList := indexer.GetFolders(path)
	for _, user := range userList {
		folders := indexer.GetFolders(path + user)
		for _, folder := range folders {
			emailFiles := indexer.GetFiles(path + user + "/" + folder + "/")
			for _, mail_file := range emailFiles {
				sysFile, _ := os.Open(path + user + "/" + folder + "/" + mail_file)
				lines := bufio.NewScanner(sysFile)
				id++
				indexer.PostDataToZincSearch(indexer.FormatData(lines, id))
				sysFile.Close()
			}
		}
	}
	fmt.Println("Indexing finished!!!!")
}
