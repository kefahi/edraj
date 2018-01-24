package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	rootPath := "/tmp/edraj/content"
	trashPath := "/tmp/edraj/trash"
	storage := Storage{rootPath, trashPath}
	folderPath := "this/is/a/test/"
	filePath := folderPath + "myfile1"
	var err error
  var list []string

	err = storage.PutFileMeta(filePath, `{"key":"some random text"}`)
	if err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader("Hello world")
  err = storage.PutFilePayload(filePath, reader)
	if err != nil {
		log.Fatal(err)
	}
  list, err = storage.ListDir(folderPath)
	fmt.Println(list)

	err = storage.MoveFile(filePath, folderPath+"myfile3")
	if err != nil {
		log.Fatal(err)
	}
  list, err = storage.ListDir(folderPath)
	fmt.Println(list)
	fmt.Println("=====")
	fmt.Println("done")

}
