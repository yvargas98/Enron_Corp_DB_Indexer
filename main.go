package main

import (
	"Enron_Corp_DB_Indexer/indexer"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// err := indexer.IndexCreateOpenObserve()
	// if err != nil {
	// 	fmt.Printf("Error creating Enron Mail Index: %s\n", err)
	// 	return
	// }

	if len(os.Args) < 2 {
		fmt.Println("Enron Corp Directory DB Path is missing.")
		return
	}

	path := os.Args[1] + "/maildir/"

	// indexer.Indexer(path)

	id := 0
	user_list := indexer.List_all_folders(path)
	for _, user := range user_list {
		folders := indexer.List_all_folders(path + user)
		for _, folder := range folders {
			mail_files := indexer.List_files(path + user + "/" + folder + "/")
			for _, mail_file := range mail_files {
				//fmt.Println("Indexing: " + user + "/" + folder + "/" + mail_file)
				sys_file, _ := os.Open(path + user + "/" + folder + "/" + mail_file) //abre el archivo
				lines := bufio.NewScanner(sys_file)                                  //Lee el archivo línea por línea (https://golangdocs.com/reading-files-in-golang)
				id++                                                                 //cada vez que se invoque la función "parse_data" esta variable se pasa con un incremento de 1 para crear el ID de cada objeto en el JSON.
				indexer.PostDataToZincSearch(indexer.FormatData(lines, id))
				sys_file.Close() //cierra el archivo
			}
		}
	}
	fmt.Println("Indexing finished!!!!")
}
