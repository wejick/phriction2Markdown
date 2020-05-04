package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const connstring = "root:giovani@tcp(localhost:32779)/phriction"
const basePath = "./exported/"

type phriction struct {
	DBConnstring string `json:"connstring"`

	DB *sql.DB
}

type phIndex struct {
	ID   string `db:"id"`
	Slug string `db:"slug"`
}

type phDocument struct {
	DocID   string
	Slug    string
	Title   string
	Content string

	DateCreated int
	LastUpdate  int
}

func (ph *phriction) GetAllIndex() (index []phIndex) {
	query := "SELECT id,slug FROM phriction_document ORDER BY depth ASC;"
	rows, err := ph.DB.Query(query)
	if err != nil {
		log.Println("couldn't get index", err)
	}
	defer rows.Close()

	for rows.Next() {
		var idx phIndex
		err = rows.Scan(&idx.ID, &idx.Slug)
		if err != nil {
			log.Println("couldn't get index", err)
			break
		}

		index = append(index, idx)
	}

	return
}

func (ph *phriction) GetDocumentByID(id string) (document phDocument, err error) {
	query := "SELECT documentID, title, slug, content, dateCreated, dateModified FROM phriction_content WHERE documentID = " + id + " ORDER BY version DESC LIMIT 1"
	row := ph.DB.QueryRow(query)

	err = row.Scan(&document.DocID, &document.Title, &document.Slug, &document.Content, &document.DateCreated, &document.LastUpdate)

	return
}

func initPH() (ph phriction, err error) {
	ph = phriction{
		DBConnstring: connstring,
	}

	db, err := sql.Open("mysql", ph.DBConnstring)
	ph.DB = db

	return
}

func createDir(path string) (err error) {
	os.MkdirAll(basePath+path, os.ModePerm)
	return
}

func sanitizePath(path string) (sanitizedPath string) {
	sanitizedPath = strings.ReplaceAll(path, " ", "_")
	return
}

func sanitizeFileName(title string) (sanitizedFilename string) {
	sanitizedFilename = strings.ReplaceAll(title, " ", "_")
	sanitizedFilename = strings.ReplaceAll(title, "/", "_")
	return
}

func writeDoc(doc phDocument) (err error) {
	sanitizedSlug := sanitizePath(doc.Slug)
	sanitizedTitle := sanitizeFileName(doc.Title)

	err = createDir(sanitizedSlug)
	if err != nil {
		log.Println("couldn't create dire", err)
	}
	f, err := os.Create(basePath + sanitizedSlug + "/" + sanitizedTitle + ".md")
	if err != nil {
		log.Println("couldn't write doc", doc.Slug, err)
		return
	}
	defer f.Close()

	mdfized := mdifize(doc.Content)
	f.WriteString(mdfized)

	return
}

func main() {
	ph, err := initPH()
	if err != nil {
		log.Fatalln("couldn't connect to db", err)
	}
	defer ph.DB.Close()

	index := ph.GetAllIndex()

	for _, idx := range index {
		doc, err := ph.GetDocumentByID(idx.ID)
		if err != nil {
			log.Println("couldn't get doc, id : ", idx, err)
			continue
		}
		writeDoc(doc)
	}
}
