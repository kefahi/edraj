package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo server
type Mongo struct {
	ServerAddress  string
	DatabaseName   string
	CollectionName string
	session        *mgo.Session
	database       *mgo.Database
	collection     *mgo.Collection
}

// Connect to mongodb
func (m *Mongo) Connect() {
	var err error
	m.session, err = mgo.Dial(m.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	m.database = m.session.DB(m.DatabaseName)
	m.collection = m.database.C(m.CollectionName)
}

// FindAll matches
func (m *Mongo) FindAll() (interface{}, error) {
	var result interface{}
	err := m.collection.Find(bson.M{}).All(&result)
	return result, err
}

// FindByID for a single document
func (m *Mongo) FindByID(id string) (interface{}, error) {
	var result interface{}
	err := m.collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return result, err

}

// Insert object
func (m *Mongo) Insert(docs ...interface{}) error {
	return m.collection.Insert(docs)
}

// Delete object
func (m *Mongo) Delete(selector interface{}) error {
	return m.collection.Remove(selector)
}

// Update object
func (m *Mongo) Update(id string, selector interface{}) error {
	return m.collection.Update(id, selector)
}

/*

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
*/
