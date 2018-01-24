package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	RootPath  = "/tmp/edraj/content"
	TrashPath = "/tmp/edraj/trash"
	storage   = Storage{RootPath, TrashPath}

	Id          = "1"
	OwnerID     = "1000"
	Permissions = []string{"read", "write", "execute"}
	Tags        = []string{"home", "test"}
	Categories  = []string{"Private", "Family"}
	dirmeta     = DirMeta{Id, OwnerID, Permissions, Tags, Categories}
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		log.Fatal("\n\tFile/Dir wasnt Found/Created!\n\n")
	}
	return true, err
}

func TestCanonicalPath(t *testing.T) {
	fmt.Printf("Testing method CanonicalPath...\n")
	storage.CanonicalPath("/tmp/edraj")

	// exists(RootPath)
	// exists(TrashPath)
}

func TestValidDir(t *testing.T) {
	fmt.Printf("Testing method ValidDir...\n")
	storage.ValidDir("/tmp/edraj/Dir", true)

	exists("/tmp/edraj/content/tmp/edraj/Dir")

}

func TestPutDirMet(t *testing.T) {
	fmt.Printf("Testing method DirMeti...\n")
	storage.PutDirMeta("/tmp/edraj/Dir", dirmeta)

	exists("/tmp/edraj/content/tmp/edraj/Dir/.meta.json")
}
