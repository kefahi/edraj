package main

import (
	//"syscall"
	"fmt"
	//"io/ioutil"
	"log"
	//"os"
	//"path/filepath"
)

/*
func ListDir(path string) {
	fmt.Println("--------")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}*/

/*
func Fatal(err *error) {
	if err != nil {
		log.Fatal(err)
	}
}*/

func main() {
	rootPath := "/tmp/edraj/content"
	trashPath := "/tmp/edraj/trash"
	storage := Storage{rootPath, trashPath}
	folderPath := "this/is/a/test/"
	//fileName := "myfile1"
	var err error
	//full_path := root_path + folderPath
	/*err := syscall.Chroot(root_path)
	if err != nil {
		log.Fatal(err)
	}*/
	//if _, err := os.Stat(full_path); err == nil {
	/*if storage.ValidDir(full_path) {
		os.RemoveAll(full_path)
		fmt.Printf("folder previously existed\n")
	}
	os.MkdirAll(full_path, os.ModeDir|0755)*/
	//err := storage.PutMeta(folderPath, file_name, `{"key":"some random text"}`)
	if err != nil {
		log.Fatal(err)
	}
	//os.Link(full_path+file_name, full_path+"myfile2")
	//ListDir(full_path)
	fmt.Println(storage.ListDir(folderPath))
	//os.RemoveAll(full_path + file_name)
	//ListDir(full_path)
	//os.Rename(full_path+"myfile2", full_path+"myfile3")
	err = storage.MoveFile(folderPath+"myfile1", folderPath+"myfile3")
	if err != nil {
		log.Fatal(err)
	}
	//ListDir(full_path)
	fmt.Println(storage.ListDir(folderPath))
	files, err := storage.ListDir(folderPath) //filepath.Glob(full_path+"*")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("=====")
	fmt.Println(files)
	fmt.Println("done")

}
