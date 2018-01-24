package main

import (
	"testing"
)

var (
	RootPath = "/tmp/edraj/content"
	TrashPath = "/tmp/edraj/trash"
)

func TestCanonicalPath(t *testing.T) {
	storage := Storage{RootPath, TrashPath}
	storage.CanonicalPath("/tmp/edraj")
}

