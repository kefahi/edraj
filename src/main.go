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

	err = storage.PutDirMeta(folderPath, DirMeta{Id:"xyz", Categories:[]string{"A","B"}})
	if err != nil {
		log.Fatal(err)
	}

	err = storage.PutFileMeta(filePath, FileMeta{Id:"abc"})
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

	fileMeta, err := storage.GetFileMeta(folderPath+"myfile3")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fileMeta)

	dirMeta, err := storage.GetDirMeta(folderPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dirMeta)

  list, err = storage.ListDir(folderPath)
	fmt.Println(list)
	fmt.Println("=====")
	fmt.Println("done")

}
