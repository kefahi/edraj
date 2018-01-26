package main

import (
	"encoding/json"
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

// FileMeta : the structure for File meta details
type FileMeta struct {
	ID          string   `json:"id"`
	OwnerID     string   `json:"owner"`
	Permissions []string `json:"permissions"`
	Tags        []string `json:"tags"`
	Categories  []string `json:"categories"`
	ContentType string   `json:"content-type"`
	AuthorID    string   `json:"author"`
	Signature   string   `json:"signature"`
	Payload     string   `json:"payload"`
	Checksum    string   `json:"checksum"`
	Schema      string   `json:"schema,omitempty"`
}

// DirMeta : the structure for Directory meta details
type DirMeta struct {
	ID          string   `json:"id"`
	OwnerID     string   `json:"owner"`
	Permissions []string `json:"permissions"`
	Tags        []string `json:"tags"`
	Categories  []string `json:"categories"`
}

// CanonicalPath : Make sure path is confined to its root directory
func (s *Storage) CanonicalPath(entryPath string) (string, error) {
	canonicalPath := path.Clean(s.RootPath + "/" + entryPath)
	var err error
	if !path.IsAbs(canonicalPath) || !strings.HasPrefix(canonicalPath, s.RootPath) {
		err = errors.New("Invalid input string: " + canonicalPath)
	}
	return canonicalPath, err
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return (err == nil && info.Mode().IsRegular() && (info.Mode().Perm()&0600 == 0600))
}

func dirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return (err == nil && info.IsDir() && (info.Mode().Perm()&0700 == 0700))
}

// ValidDir : Make sure the provided path meets the ValidDir requirements; creates Directories if needed
func (s *Storage) ValidDir(dirPath string, createIfMissing bool) (string, error) {
	canonicalPath, err := s.CanonicalPath(dirPath)
	if err != nil {
		return canonicalPath, err // errors.New("Invalid path:" + canonicalPath)
	}
	info, err := os.Stat(canonicalPath)
	if err != nil {
		if createIfMissing && os.IsNotExist(err) {
			return canonicalPath, os.MkdirAll(s.RootPath+"/"+dirPath, os.ModeDir|0755)
		}
		return canonicalPath, err
	} else if !info.IsDir() || (info.Mode()&0700 != 0700) {
		return canonicalPath, errors.New("Invalid path: " + canonicalPath)
	}
	return canonicalPath, nil
}

// PutDirMeta : Assigns meta details to a directory (the meta file name is .meta.json)
func (s *Storage) PutDirMeta(dirPath string, dirMeta DirMeta) error {
	canonicalPath, err := s.ValidDir(dirPath, true)
	if err != nil {
		return err
	}

	// TODO check json L3
	data, err := json.Marshal(dirMeta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(canonicalPath+"/.meta.json", data, 0644)
}

// PutFileMeta : Assigns meta details to a file (adds "."+fileName+".meta.json")
func (s *Storage) PutFileMeta(filePath string, fileMeta FileMeta) error {
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
	data, err := json.Marshal(fileMeta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(canonicalPath+"/."+fileName+".meta.json", data, 0644)
}

// GetFileMeta : return the File's meta data
func (s *Storage) GetFileMeta(filePath string) (FileMeta, error) {
	dirPath := path.Dir(filePath)
	canonicalDirPath, err := s.ValidDir(dirPath, false)
	if err != nil {
		return FileMeta{}, err
	}

	fileName := path.Base(filePath)
	canonicalFilePath := canonicalDirPath + "/." + fileName + ".meta.json"

	fileMeta := FileMeta{}
	info, err := os.Stat(canonicalFilePath)
	if err != nil || info.IsDir() /*|| ( !info.IsDir() && (info.Mode()&0600 == 0600))*/ {
		return fileMeta, err
	}
	data, err := ioutil.ReadFile(canonicalFilePath)
	if err == nil {
		// TODO Do json L4 level checking
		json.Unmarshal(data, &fileMeta)
	}
	return fileMeta, err
}

// GetDirMeta : return the Directory's meta data
func (s *Storage) GetDirMeta(dirPath string) (DirMeta, error) {
	canonicalDirPath, err := s.ValidDir(dirPath, false)
	if err != nil {
		return DirMeta{}, err
	}

	canonicalFilePath := canonicalDirPath + "/.meta.json"

	dirMeta := DirMeta{}
	info, err := os.Stat(canonicalFilePath)
	if err != nil || info.IsDir() /*&& !info.IsDir() && (info.Mode()&0600 == 0600)*/ {
		return dirMeta, err
	}
	data, err := ioutil.ReadFile(canonicalFilePath)
	if err == nil {
		// TODO Do json L4 level checking
		json.Unmarshal(data, &dirMeta)
	}
	return dirMeta, err
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
	_, err = os.Stat(trashDirPath)
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

	parentPath := path.Dir(dirPath)
	trashDirPath := s.TrashPath + parentPath
	if !dirExists(trashDirPath) {
		err = os.MkdirAll(trashDirPath, os.ModeDir|0755)
		if err != nil {
			return err
		}
		log.Println("Created missing trashdir", trashDirPath)
	} // TODO add an else if the full target path exists to determine what to do (I thin we should simply delete the old conflicting)

	return os.Rename(canonicalDirPath, trashDirPath+"/"+path.Base(dirPath))
}

// MoveFile : Move a directory within the RootPath
func (s *Storage) MoveFile(fromPath string, toPath string) error {
	canonicalFromPath, err := s.ValidDir(path.Dir(fromPath), false)
	if err != nil {
		return err
	}
	canonicalToPath, err := s.ValidDir(path.Dir(toPath), true)
	if err != nil {
		return err
	}

	fromFileName := path.Base(fromPath)
	toFileName := path.Base(toPath)
	err = os.Rename(canonicalFromPath+"/"+fromFileName, canonicalToPath+"/"+toFileName)
	if err != nil {
		return err
	}
	return os.Rename(canonicalFromPath+"/."+fromFileName+".meta.json", canonicalToPath+"/."+toFileName+".meta.json")
}

// MoveDir : Move a directory within the RootPath
func (s *Storage) MoveDir(fromPath string, toPath string) error {
	canonicalFromPath, err := s.ValidDir(fromPath, false)
	if err != nil {
		return err
	}
	parentToPath := path.Dir(toPath)
	folderName := path.Base(toPath)
	canonicalToPath, err := s.ValidDir(parentToPath, true)
	if err != nil {
		return err
	}

	return os.Rename(canonicalFromPath, canonicalToPath+"/"+folderName)
}

// PutFilePayload : takes the io.Reader to retrieve the data that would be persisted to the payload file
func (s *Storage) PutFilePayload(filePath string, reader io.Reader) error {
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
	/*bytesWritten*/
	_, err = io.Copy(newFile, reader)
	if err != nil {
		return err
	}
	// log.Printf("Copied %d bytes.", bytesWritten)

	// Commit the file contents
	// Flushes memory to disk
	err = newFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// GetFilePayload : return the io.Reader that would serve the file payload
func (s *Storage) GetFilePayload(filePath string) (io.Reader, error) {
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
