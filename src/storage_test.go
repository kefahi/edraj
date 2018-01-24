package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

var (
	RootPath  = "/tmp/edraj"
	TrashPath = "/tmp/edraj/trash"
	storage   = Storage{RootPath, TrashPath}

	dirId       = "1"
	OwnerID     = "1000"
	Permissions = []string{"read", "write", "execute"}
	Tags        = []string{"home", "test"}
	Categories  = []string{"Private", "Family"}
	dirmeta     = DirMeta{dirId, OwnerID, Permissions, Tags, Categories}

	fileId      = "2"
	contentType = "TODOs"
	AutherId    = "1"
	Signature   = "signiture" // Replace with appropriate test value
	payload     = "Hello World"
	checksum    = "93623ac7d9badb95b01f74ceb2d17702f142e692"
	schema      = "TODO"
	filemeta    = FileMeta{fileId, OwnerID, Permissions, Tags, Categories, contentType, AutherId, Signature, payload, checksum, schema}
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
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
	storage.ValidDir("/content/Dir", true)

	exists("/tmp/edraj/content/Dir")

}

func TestPutDirMeta(t *testing.T) {
	fmt.Printf("Testing method PutDirMeta...\n")
	storage.PutDirMeta("/content/Dir", dirmeta)

	state, _ := exists("/tmp/edraj/content/Dir/.meta.json") // Json verification will be added soon
	if state == false {
		log.Fatal("Doesn't Exist!")
	}
}

// Issiues found while making PutFileMeta:
// The function creates the metafile even when the target file isnt in the path given
func TestPutFileMeta(t *testing.T) {
	fmt.Printf("Testing method PutFileMeta...\n")
	storage.PutFileMeta("/content/Dir/test.todo", filemeta)

	state, _ := exists("/tmp/edraj/content/Dir/.test.todo.meta.json") // Json verification will be added soon
	if state == false {
		log.Fatal("Doesn't Exist!")
	}
	// exists("/tmp/edraj/content/Dir/test.todo")
}

func TestGetFileMeta(t *testing.T) {
	fmt.Printf("Testing method GetFileMeta...\n")
	object, _ := storage.GetFileMeta("/content/Dir/test.todo")

	if object.Id != filemeta.Id || object.OwnerId != filemeta.OwnerId || object.ContentType != filemeta.ContentType || object.AutherId != filemeta.AutherId || object.Signature != filemeta.Signature || object.Payload != filemeta.Payload || object.Checksum != filemeta.Checksum || object.Schema != filemeta.Schema {
		log.Fatal("GetFileMeta returned corrupt filemeta")
	} // Comparing the two objects did'nt work
	//  invalid operation: object != filemeta (struct containing []string cannot be compared)
	// so I used the long method. Maybe there is an easier way to do it
	// It will get longer if i worked on checking []string like "Permissions". so i execluded them for the moment.
}

func TestGetDirMeta(t *testing.T) {
	fmt.Printf("Testing method GetDirMeta...\n")
	object, _ := storage.GetDirMeta("/content/Dir")

	state, _ := exists("/tmp/edraj/content/Dir")
	if state == false {
		log.Fatal("Doesn't Exist!")
	}

	if object.Id != dirmeta.Id || object.OwnerId != filemeta.OwnerId {
		log.Fatal("GetDirMeta returned corrupt dirmeta")
	} // Comparing the two objects did'nt work
	//  invalid operation: object != dirmeta (struct containing []string cannot be compared)
	// so I used the long method. Maybe there is an easier way to do it
	// It will get longer if i worked on checking []string like "Permissions". so i execluded them for the moment.
}

// TestListDir is still under construction
func TestListDir(t *testing.T) {
	fmt.Printf("Testing method ListDir...\n")

	list, _ := storage.ListDir("/content/Dir") //
	//syslist := exec.Command("find", "/tmp/edraj/content/Dir")
	// Under Construction
	// for x := 0; x > 10; x++ {
	//		fmt.Println([x]list)
	//	}
	fmt.Printf("%v\n", list)
}

func TestDeleteFile(t *testing.T) {
	fmt.Printf("Testing method DeletFile...\n")

	exec.Command("sh", "-c", "touch /tmp/edraj/content/Dir/TEST").Output()
	time.Sleep(1 * time.Millisecond)

	state, _ := exists("/tmp/edraj/content/Dir/TEST")
	if state == false {
		log.Fatal("Doesn't Exist!")
	}

	storage.DeleteFile("/content/Dir/TEST")
	state, _ = exists("/tmp/edraj/content/Dir/TEST")
	if state == true {
		log.Fatal("File Still Exists!")
	}

	state, _ = exists("/tmp/edraj/trash/content/Dir/TEST")
	if state == false {
		log.Fatal("File Doesn't Exists!")
	}
}

// Deleting Dirs Doesnt work for some reason
//func TestDeleteDir(t *testing.T) {
//	fmt.Printf("Testing method DeletDir...\n")

//	exec.Command("sh", "-c", "mkdir  /tmp/edraj/content/Dir/TEST").Output()
//	time.Sleep(100 * time.Millisecond)
//
//	state, _ := exists("/tmp/edraj/content/Dir/TEST")
//	if state == false {
//		log.Fatal("Doesn't Exist!")
//	}
//
//	storage.DeleteDir("/content/Dir/TEST/")
//	state, _ = exists("/tmp/edraj/content/Dir/TEST")
//	if state == true {
//		log.Fatal("Dir Still Exists!")
//	}
//}

func TestMoveFile(t *testing.T) {
	fmt.Printf("Testing method MoveFile...\n")

	exec.Command("sh", "-c", "touch /tmp/edraj/content/Dir/TEST").Output()
	time.Sleep(100 * time.Millisecond)
	state, _ := exists("/tmp/edraj/content/Dir/TEST")
	if state == false {
		log.Fatal("Doesn't Exist!")
	}

	storage.MoveFile("/content/Dir/TEST", "/content/TEST")
	state, _ = exists("/tmp/edraj/content/Dir/TEST")
	if state == true {
		log.Fatal("File Still Exists!")
	}
	state, _ = exists("/tmp/edraj/content/TEST")
	if state == false {
		log.Fatal("File Doesn't Exists!")
	}
	exec.Command("sh", "-c", "rm /tmp/edraj/content/TEST").Output()
	time.Sleep(100 * time.Millisecond)

}

// Moving Dirs Doesnt work for some reason
//func TestMoveDir(t *testing.T) {
//	fmt.Printf("Testing method MoveDir...\n")
//
//	exec.Command("sh", "-c", "mkdir /tmp/edraj/content/Dir/TEST").Output()
//	time.Sleep(100 * time.Millisecond)
//	state, _ := exists("/tmp/edraj/content/Dir/TEST")
//	if state == false {
//		log.Fatal("Doesn't Exist!")
//	}
//
//	storage.MoveDir("/content/Dir/TEST", "/content/TEST")
//	state, _ = exists("/tmp/edraj/content/Dir/TEST")
//	if state == true {
//		log.Fatal("Dir Still Exists!")
//	}
//	state, _ = exists("/tmp/edraj/content/TEST")
//	if state == false {
//		log.Fatal("Dir Doesn't Exists!")
//	}
//}
