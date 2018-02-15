package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

var (
	RootPath  = "/tmp/edraj/content"
	TrashPath = "/tmp/edraj/trash"
	TStorage  = Storage{RootPath, TrashPath}

	TDirMeta = DirMeta{
		ID:          "1",
		OwnerID:     "1000",
		Permissions: []string{"read", "write", "execute"},
		Tags:        []string{"home", "test"},
		Categories:  []string{"Private", "Family"},
	}

	TFileMeta = FileMeta{
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
	TStorage.CanonicalPath("/tmp/edraj")

	// exists(RootPath)
	// exists(TrashPath)
}

func TestValidDir(t *testing.T) {
	fmt.Printf("Testing method ValidDir...\n")
	TStorage.ValidDir("/Dir", true)

	state, _ := exists("/tmp/edraj/content/Dir/") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	}

}

func TestPutDirMeta(t *testing.T) {
	fmt.Printf("Testing method PutDirMeta...\n")
	TStorage.PutDirMeta("/Dir", TDirMeta)

	state, _ := exists("/tmp/edraj/content/Dir/.meta.json") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	}
}

// Issiues found while making PutFileMeta:
// The function creates the metafile even when the target file isnt in the path given
func TestPutFileMeta(t *testing.T) {
	fmt.Printf("Testing method PutFileMeta...\n")
	TStorage.PutFileMeta("/Dir/test.todo", TFileMeta)

	state, _ := exists("/tmp/edraj/content/Dir/.test.todo.meta.json") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	}
	// exists("/tmp/edraj/content/Dir/test.todo")
}

func TestGetFileMeta(t *testing.T) {
	fmt.Printf("Testing method GetFileMeta...\n")
	object, _ := TStorage.GetFileMeta("/Dir/test.todo")

	if object.ID != TFileMeta.ID || object.OwnerID != TFileMeta.OwnerID || object.ContentType != TFileMeta.ContentType || object.AuthorID != TFileMeta.AuthorID || object.Signature != TFileMeta.Signature || object.Payload != TFileMeta.Payload || object.Checksum != TFileMeta.Checksum || object.Schema != TFileMeta.Schema {
		t.Error("GetFileMeta returned corrupt FileMeta")
	} // Comparing the two objects did'nt work
	//  invalid operation: object != TFileMeta (struct containing []string cannot be compared)
	// so I used the long method. Maybe there is an easier way to do it
	// It will get longer if i worked on checking []string like "Permissions". so i execluded them for the moment.
}

func TestGetDirMeta(t *testing.T) {
	fmt.Printf("Testing method GetDirMeta...\n")
	object, _ := TStorage.GetDirMeta("/Dir")

	state, _ := exists("/tmp/edraj/content/Dir")
	if state == false {
		t.Error("Doesn't Exist!")
	}

	if object.ID != TDirMeta.ID || object.OwnerID != TFileMeta.OwnerID {
		t.Error("GetDirMeta returned corrupt TDirMeta")
	} // Comparing the two objects did'nt work
	//  invalid operation: object != TDirMeta (struct containing []string cannot be compared)
	// so I used the long method. Maybe there is an easier way to do it
	// It will get longer if i worked on checking []string like "Permissions". so i execluded them for the moment.
}

// TestListDir is still under construction
func TestListDir(t *testing.T) {
	fmt.Printf("Testing method ListDir...\n")

	/*list*/
	_, err := TStorage.ListDir("/Dir")
	if err != nil {
		t.Error("Method ListDir returned error!")
	}
	// fmt.Printf("%v\n", list)
}

func TestDeleteFile(t *testing.T) {
	fmt.Println("Testing method DeleteFile...")

	exec.Command("sh", "-c", "touch /tmp/edraj/content/Dir/DelFileTEST").Output()
	time.Sleep(100 * time.Millisecond)

	state, _ := exists("/tmp/edraj/content/Dir/DelFileTEST")
	if state == false {
		t.Error("Test File failed to be created!\n")
	}
	TStorage.DeleteFile("/Dir/DelFileTEST")

	state, _ = exists("/tmp/edraj/content/Dir/DelFileTEST")
	if state == true {
		t.Error("Test File Was not removed!\n")
	}
	state, _ = exists("/tmp/edraj/trash/Dir/DelFileTEST")
	if state == false {
		t.Error("Deleted File is not in trash!\n")
	}
}

/*
func TestDeleteDir(t *testing.T) {
	fmt.Println("Testing method DeleteDir...")

	exec.Command("sh", "-c", "mkdir /tmp/edraj/content/Dir/DelDirTEST").Output()
	time.Sleep(100 * time.Millisecond)

	state, _ := exists("/tmp/edraj/content/Dir/DelDirTEST")
	if state == false {
		t.Error("Test Dir failed to be created!\n")
	}
	TStorage.DeleteFile("/Dir/DelDirTEST")

	state, _ = exists("/tmp/edraj/content/Dir/DelDirTEST")
	if state == true {
		t.Error("Test Dir Was not removed!\n")
	}
	state, _ = exists("/tmp/edraj/trash/Dir/DelDirTEST")
	if state == false {
		t.Error("Deleted Dir is not in trash!\n")
	}
}
*/
/*
func TestMoveFile(t *testing.T) {
	fmt.Println("Testing method MoveFile...")

	exec.Command("sh", "-c", "touch /tmp/edraj/content/Dir/mvFileTEST3").Output()
	time.Sleep(100 * time.Millisecond)

	state, _ := exists("/tmp/edraj/content/Dir/mvFileTEST3")
	if state == false {
		t.Error("Test File failed to be created!\n")
	}
	TStorage.MoveFile("/Dir/DelFileTEST3", "/")

	state, _ = exists("/tmp/edraj/content/Dir/mvFileTEST3")
	if state == true {
		t.Error("Test File Was not moved!\n")
	}
	state, _ = exists("/tmp/edraj/mvFileTEST3")
	if state == false {
		t.Error("Moved File is not in moved location!\n")
	}
}
*/
/*
func TestMoveDir(t *testing.T) {
	fmt.Println("Testing method MoveDir...")

	exec.Command("sh", "-c", "mkdir /tmp/edraj/content/Dir/mvDirTEST3").Output()
	time.Sleep(100 * time.Millisecond)

	state, _ := exists("/tmp/edraj/content/Dir/mvDirTEST3")
	if state == false {
		t.Error("Test Dir failed to be created!\n")
	}
	TStorage.MoveDir("/Dir/mvDirTEST3", "/")

	state, _ = exists("/tmp/edraj/content/Dir/mvDirTEST3")
	if state == true {
		t.Error("Test Dir Was not moved!\n")
	}
	state, _ = exists("/tmp/edraj/mvDirTEST3")
	if state == false {
		t.Error("Moved Dir is not in moved location!\n")
	}
}
*/
// func TestDeletingWithoutForwardSlash (t *testing.T) {
//
//         fmt.Println("Testing method DeleteFile without Forward Slash at the beggining ...\n")
//
//         exec.Command("sh", "-c", "touch /tmp/edraj/content/Dir/DelFileTEST2").Output()
//         time.Sleep(100 * time.Millisecond)
//
//         state, _ := exists("/tmp/edraj/content/Dir/DelFileTEST2")
//         if state == false {
//                 t.Error("Test File failed to be created!\n")
//         }
//         TStorage.DeleteFile("Dir/DelFileTEST2")
//
// 	state, _ = exists("/tmp/edraj/trash")
//         if state == false {
//                 t.Error(TrashPath + "doesnt Exist!\n")
//         }
//
//         state, _ = exists("/tmp/edraj/content/Dir/DelFileTEST2")
//         if state == true {
//                 t.Error("Test File Was not removed!\n")
//         }
//         state, _ = exists("/tmp/edraj/trash/DelFileTEST2")
//         if state == false {
//                 t.Error("Deleted File is not in trash!\n")
//         }
// }
//
// func TestDeleteDirWithoutForwardSlash(t *testing.T)  {
//         fmt.Println("Testing method DeleteDir...\n")
//
//         exec.Command("sh", "-c", "mkdir /tmp/edraj/content/Dir/DelDirTEST2").Output()
//         time.Sleep(100 * time.Millisecond)
//
//         state, _ := exists("/tmp/edraj/content/Dir/DelDirTEST2")
//         if state == false {
//                 t.Error("Test Dir failed to be created!\n")
//         }
//         TStorage.DeleteFile("Dir/DelDirTEST2")
//
//         state, _ = exists("/tmp/edraj/content/Dir/DelDirTEST2")
//         if state == true {
//                 t.Error("Test Dir Was not removed!\n")
//         }
//         state, _ = exists("/tmp/edraj/trash/DelDirTEST2")
//         if state == false {
//                 t.Error("Deleted Dir is not in trash!\n")
//         }
// }
