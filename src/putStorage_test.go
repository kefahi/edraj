package main

import (
	"testing"
	"time"
)


func TestPutDirMeta(t *testing.T) {
	fmt.Printf("Testing method PutDirMeta...\n")
	storage.PutDirMeta("/tmp/putStorage", dirMeta)

	state, _ := exists("/tmp/putStorage/.putStorage.json") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	} else if state == true {
		  permissions, err := os.Stat("/tmp/putStorage")
	    exec.Command("sh", "-c", "rmdir /tmp/putStorage").Output()
			if permissions == 0644 {
			} else {
			  t.Error("Permissions incorrect")
			}
	}
}



func TestPutDirMeta(t *testing.T) {
	fmt.Printf("Testing method PutDirMeta...\n")
	exec.Command("sh", "-c", "touch /tmp/putStorage.txt").Output()
	time.Sleep( 50 * time.Millisecond )

	storage.PutDirMeta("/tmp/putStorage.txt", dirMeta)

	state, _ := exists("/tmp/putStorage") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	} else {
	     exec.Command("sh", "-c", "rm /tmp/putStorage.txt").Output()
	  }
}


func TestPutFileMeta(t *testing.T) {
	fmt.Printf("Testing method PutFileMeta...\n")

	storage.PutFileMeta("/tmp/putStorage.txt", dirMeta)

	state, _ := exists("/tmp/putStorage") // Json verification will be added soon
	if state == false {
		t.Error("Doesn't Exist!")
	}
}


