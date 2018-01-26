package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

var (
	storage = Storage{RootPath: "/tmp/edraj/content", TrashPath: "/tmp/edraj/trash"}

	dirMeta = DirMeta{
		ID:          "1",
		OwnerID:     "1000",
		Permissions: []string{"read", "write", "execute"},
		Tags:        []string{"home", "test"},
		Categories:  []string{"Private", "Family"},
	}

	fileMeta = FileMeta{
		ID:          "2",
		OwnerID:     "1000",
		Permissions: []string{"read", "write", "execute"},
		Tags:        []string{"home", "test"},
		Categories:  []string{"Private", "Family"},
		ContentType: "Text",
		AuthorID:    "1",
		Signature:   "signiture", // Replace with appropriate test value
		Payload:     "Hello World",
		Checksum:    "93623ac7d9badb95b01f74ceb2d17702f142e692",
		Schema:      "MySchema",
	}
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
	storage.ValidDir("/Dir", true)

	state, _ := exists("/tmp/edraj/content/Dir/") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	}

}

func TestPutDirMeta(t *testing.T) {
	fmt.Printf("Testing method PutDirMeta...\n")
	storage.PutDirMeta("/Dir", dirMeta)

	state, _ := exists("/tmp/edraj/content/Dir/.meta.json") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	}
}

// Issiues found while making PutFileMeta:
// The function creates the metafile even when the target file isnt in the path given
func TestPutFileMeta(t *testing.T) {
	fmt.Printf("Testing method PutFileMeta...\n")
	storage.PutFileMeta("/Dir/test.todo", fileMeta)

	state, _ := exists("/tmp/edraj/content/Dir/.test.todo.meta.json") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	}
	// exists("/tmp/edraj/content/Dir/test.todo")
}

func TestGetFileMeta(t *testing.T) {
	fmt.Printf("Testing method GetFileMeta...\n")
	object, _ := storage.GetFileMeta("/Dir/test.todo")

	if object.ID != fileMeta.ID || object.OwnerID != fileMeta.OwnerID || object.ContentType != fileMeta.ContentType || object.AuthorID != fileMeta.AuthorID || object.Signature != fileMeta.Signature || object.Payload != fileMeta.Payload || object.Checksum != fileMeta.Checksum || object.Schema != fileMeta.Schema {
		t.Error("GetFileMeta returned corrupt fileMeta")
	} // Comparing the two objects did'nt work
	//  invalid operation: object != fileMeta (struct containing []string cannot be compared)
	// so I used the long method. Maybe there is an easier way to do it
	// It will get longer if i worked on checking []string like "Permissions". so i execluded them for the moment.
}

func TestGetDirMeta(t *testing.T) {
	fmt.Printf("Testing method GetDirMeta...\n")
	object, _ := storage.GetDirMeta("/Dir")

	state, _ := exists("/tmp/edraj/content/Dir")
	if state == false {
		t.Error("Doesn't Exist!")
	}

	if object.ID != dirMeta.ID || object.OwnerID != fileMeta.OwnerID {
		t.Error("GetDirMeta returned corrupt dirMeta")
	} // Comparing the two objects did'nt work
	//  invalid operation: object != dirMeta (struct containing []string cannot be compared)
	// so I used the long method. Maybe there is an easier way to do it
	// It will get longer if i worked on checking []string like "Permissions". so i execluded them for the moment.
}

// TestListDir is still under construction
func TestListDir(t *testing.T) {
	fmt.Printf("Testing method ListDir...\n")

	list, _ := storage.ListDir("/Dir") //
	//syslist := exec.Command("find", "/tmp/edraj/content/Dir")
	// Under Construction
	// for x := 0; x > 10; x++ {
	//		fmt.Println([x]list)
	//	}
	fmt.Printf("%v\n", list)
}

func TestDeleteFile(t *testing.T) {
	fmt.Printf("Testing method DeletFile...\n")

	exec.Command("sh", "-c", "touch /tmp/edraj/content/Dir/DelFileTEST").Output()
	time.Sleep(100 * time.Millisecond)

	state, _ := exists("/tmp/edraj/content/Dir/DelFileTEST")
	if state == false {
		t.Error("Doesn't Exist!")
	}

	storage.DeleteFile("/Dir/DelFileTEST")
	state, _ = exists("/tmp/edraj/content/Dir/DelFileTEST")
	if state == true {
		t.Error("File Still Exists!")
	}

	state, _ = exists("/tmp/edraj/trash/Dir/DelFileTEST")
	if state == false {
		t.Error("File Doesn't Exist in Trash!")
	}
}

// Deleting Dirs Doesnt work for some reason
func TestDeleteDir(t *testing.T) {
	fmt.Printf("Testing method DeletDir...\n")

	exec.Command("sh", "-c", "mkdir  /tmp/edraj/content/Dir/DelDirTEST").Output()
	time.Sleep(100 * time.Millisecond)

	storage.DeleteDir("/Dir/DelDirTEST/")
	state, _ := exists("/tmp/edraj/content/Dir/DelDirTEST")
	if state == true {
		t.Error("Dir Still Exists!")
	}

	state, _ = exists("/tmp/edraj/trash/Dir/DelFileTEST")
	if state == false {
		t.Error("File Doesn't Exist in Trash!")
	}
}

func TestMoveFile(t *testing.T) {
	fmt.Printf("Testing method MoveFile...\n")

	exec.Command("sh", "-c", "touch /tmp/edraj/content/Dir/MvFileTEST").Output()
	time.Sleep(100 * time.Millisecond)

	storage.MoveFile("/Dir/MvFileTEST", "/MvFileTEST")
	state, _ := exists("/tmp/edraj/content/Dir/MvFileTEST")
	if state == true {
		t.Error("File Still Exists!")
	}
	state, _ = exists("/tmp/edraj/content/MvFileTEST")
	if state == false {
		t.Error("File Doesn't Exists!")
	}
}

// Moving Dirs Doesnt work for some reason
func TestMoveDir(t *testing.T) {
	fmt.Printf("Testing method MoveDir...\n")

	exec.Command("sh", "-c", "mkdir /tmp/edraj/content/Dir/MvDirTEST").Output()
	time.Sleep(100 * time.Millisecond)

	storage.MoveDir("/Dir/MvDirTEST", "/MvDirTEST")
	state, _ := exists("/tmp/edraj/content/Dir/MvDirTEST")
	if state == true {
		t.Error("Dir Still Exists!")
	}
	state, _ = exists("/tmp/edraj/content/MvDirTEST")
	if state == false {
		t.Error("Dir Doesn't Exists!")
	}
}
