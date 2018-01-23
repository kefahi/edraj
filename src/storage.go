package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Storage : An adapter to access meta-data enabled storage pools.
// The functions only operate on files within the pool, extra check are put in the code to guarantee that.
// Paths must be 1. Absolute, 2. End with a slash
type Storage struct {
	RootPath  string
	TrashPath string
}

// TODO define the meta data structurs and return/accept them in the meta functions

// FileMeta : the structure that defines the fields of a file meta data
type FileMeta struct {
}

// DirMeta : the structure that defines the fields of a directory meta data
type DirMeta struct {
}

// CanonicalPath : Make sure path is confined to its root directory
func (s *Storage) CanonicalPath(entryPath string) (string, error) {
	canonicalPath := path.Clean(s.RootPath + "/" + entryPath)
	var err error
	if !path.IsAbs(canonicalPath) || !strings.HasPrefix(canonicalPath, s.RootPath+entryPath) {
		err := errors.New("Invalid input string" + canonicalPath)
	}
	return canonicalPath, err
}

// ValidDir : Make sure the provided path meets the ValidDir requirements; creates Directories if needed
func (s *Storage) ValidDir(dirPath string, createIfMissing bool) (string, error) {
	canonicalPath, err := s.CanonicalPath(dirPath)
	if err != nil {
		return canonicalPath, errors.New("Invalid path:" + canonicalPath)
	}
	info, err := os.Stat(canonicalPath)
	if err != nil {
		if createIfMissing && os.IsNotExist(err) {
			return canonicalPath, os.MkdirAll(s.RootPath+dirPath, os.ModeDir|0755)
		}
		return canonicalPath, err
	} else if !info.IsDir() || (info.Mode()&0700 != 0700) {
		return canonicalPath, errors.New("Invalid path:" + canonicalPath)
	}
	return canonicalPath, nil
}

// PutDirMeta : Assigns meta details to a directory (the meta file name is .meta.json)
func (s *Storage) PutDirMeta(dirPath string, meta string) error {
	canonicalPath, err := s.ValidDir(dirPath, true)
	if err != nil {
		return err
	}

	// TODO check json L3
	return ioutil.WriteFile(canonicalPath+"/.meta.json", []byte(meta), 0644)
}

// PutFileMeta : Assigns meta details to a file (adds "."+fileName+".meta.json")
func (s *Storage) PutFileMeta(filePath string, meta string) error {
	// TODO check L4 json.
	// L1: the file is valid json
	// L2: the file is valid json schema
	// L3: the data is checked
	// L4: the external data is checked. e.g. checksup of related payloads, links to acl groups/users

	dirPath := path.Dir(filePath)
	canonicalPath, err := s.ValidDir(dirPath, true)
	if err != nil {
		return err
	}

	fileName := path.Base(filePath)
	return ioutil.WriteFile(canonicalPath+"/."+fileName+"/.meta.json", []byte(meta), 0644)
}

// GetFileMeta : return the File's meta data
func (s *Storage) GetFileMeta(filePath string) (string, error) {
	dirPath := path.Dir(filePath)
	canonicalDirPath, err := s.ValidDir(dirPath, false)
	if err != nil {
		return "", err
	}

	fileName := path.Base(filePath)
	canonicalFilePath := canonicalDirPath + "/." + fileName + "/.meta.json"

	info, err := os.Stat(canonicalFilePath)
	if err != nil && !info.IsDir() && (info.Mode()&0600 == 0600) {
		content, err := ioutil.ReadFile(canonicalFilePath)
		// TODO Do json L4 level checking
		if err == nil {
			return string(content), nil
		}
	}
	return "", err
}

// ListDir : Return list of childern (files and folders) under the dirPath. .meta.json (meta files) should be excluded
func (s *Storage) ListDir(dirPath string) ([]string, error) {
	canonicalPath, err := s.ValidDir(dirPath, false)
	if err != nil {
		return nil, err
	}
	// TODO make sure all .meta.json and .XYZ.meta.json are excluded from the returned list
	return filepath.Glob(canonicalPath + "/*")
}

// DeleteFile : Moves a file along with its meta data to trash
func (s *Storage) DeleteFile(filePath string) error {
	dirPath := path.Dir(filePath)
	canonicalDirPath, err := s.ValidDir(dirPath, false)
	if err != nil {
		return err
	}

	fileName := path.Base(filePath)
	canonicalFileMetaPath := canonicalDirPath + "/." + fileName + "/.meta.json"
	canonicalFilePath := canonicalDirPath + "/" + fileName

	trashDirPath := s.TrashPath + "/" + dirPath
	info, err := os.Stat(trashDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(s.TrashPath+dirPath, os.ModeDir|0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err1 := os.Rename(canonicalFileMetaPath, s.TrashPath+"/."+filePath+"/.meta.json")
	err2 := os.Rename(canonicalFilePath, s.TrashPath+"/"+filePath)
	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}
	return nil
}

// DeleteDir : Moves a file along with its meta data to trash
func (s *Storage) DeleteDir(dirPath string) error {
	canonicalDirPath, err := s.ValidDir(dirPath, false)
	if err != nil {
		return err
	}

	trashDirPath := s.TrashPath + "/" + path.Dir(dirPath)
	info, err := os.Stat(trashDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(s.TrashPath+dirPath, os.ModeDir|0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = os.Rename(canonicalDirPath, trashDirPath)
	if err != nil {
		return err
	}
	return nil
}

// MoveDir : Move a directory within the RootPath
func (s *Storage) MoveDir(fromPath string, toPath string) error {
	canonicalFromPath, err := s.ValidDir(fromPath, false)
	if err != nil {
		return err
	}
	canonicalToPath, err := s.ValidDir(fromPath, true)
	if err != nil {
		return err
	}

	return os.Rename(canonicalFromPath, canonicalToPath)
}

// PutPayload : takes the io.Reader to retrieve the data that would be persisted to the payload file
func (s *Storage) PutPayload(filePath string, reader io.Reader) error {
	canonicalDirPath, err := s.ValidDir(path.Dir(filePath), false)
	if err != nil {
		return err
	}
	// TODO check that the payload file has the appropriate meta file

	fileName := path.Base(filePath)
	newFile, err := os.Create(canonicalDirPath + "/" + fileName)
	if err != nil {
		return err
	}

	defer newFile.Close()

	// Copy the bytes to destination from source
	bytesWritten, err := io.Copy(newFile, reader)
	if err != nil {
		return err
	}
	log.Printf("Copied %d bytes.", bytesWritten)

	// Commit the file contents
	// Flushes memory to disk
	err = newFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// GetPayload : return the io.Reader that would serve the file payload
func (s *Storage) GetPayload(filePath string) (io.Reader, error) {
	canonicalDirPath, err := s.ValidDir(path.Dir(filePath), false)
	if err != nil {
		return nil, err
	}
	fileName := path.Base(filePath)
	info, err := os.Stat(canonicalDirPath + "/" + fileName)
	if err != nil {
		return nil, err
	}
	if info.IsDir() || (info.Mode()&0600 != 0600) {
		return nil, errors.New("Unable to read file: " + canonicalDirPath + fileName)
	}
	// TODO check that appropriate meta file exists
	return os.Open(canonicalDirPath + "/" + fileName)
}
