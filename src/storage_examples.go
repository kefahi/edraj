package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func examplesMain() {
	rootPath := "/tmp/edraj/content"
	trashPath := "/tmp/edraj/trash"
	storage := Storage{rootPath, trashPath}
	folderPath := "this/is/a/test/"
	filePath := folderPath + "myfile1"
	var err error
	var list []string

	err = storage.PutDirMeta(folderPath, DirMeta{ID: "xyz", Categories: []string{"A", "B"}})
	if err != nil {
		log.Fatal(err)
	}

	err = storage.PutFileMeta(filePath, FileMeta{ID: "abc"})
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

	log.Println("Dir", rootPath, "exists?", dirExists(rootPath))
	log.Println("File", rootPath+"/"+filePath, "exists?", fileExists(rootPath+"/"+filePath))

	err = storage.MoveFile(filePath, folderPath+"myfile3")
	if err != nil {
		log.Fatal(err)
	}

	fileMeta, err := storage.GetFileMeta(folderPath + "myfile3")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fileMeta)

	dirMeta, err := storage.GetDirMeta(folderPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dirMeta)

	err = storage.MoveDir("this", "that")
	if err != nil {
		log.Fatal(err)
	}

	list, err = storage.ListDir("/that/is/a/test")
	fmt.Println(list)

	err = storage.DeleteDir("/that/is/a/test")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(filepath.Glob(trashPath + "/that/is/a/test/*"))
	fmt.Println("=====")
	fmt.Println("done")

}
