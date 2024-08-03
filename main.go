package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Strubbl/wallabago/v9"
	"github.com/go-shiori/go-epub"
	_ "github.com/joho/godotenv/autoload"
)

func getEntries(length int) []wallabago.Item {
	entries, err := wallabago.GetEntries(
		wallabago.APICall,
		0, 0, "", "", 1, length, "", 0, -1, "", "")
	if err != nil {
		log.Println("Cannot obtain articles from Wallabag")
	}

	return entries.Embedded.Items
}

func init() {
	wallabagConfig := wallabago.WallabagConfig{
		WallabagURL:  os.Getenv("WALLABAG_URL"),
		ClientID:     os.Getenv("WALLABAG_CLIENT_ID"),
		ClientSecret: os.Getenv("WALLABAG_CLIENT_SECRET"),
		UserName:     os.Getenv("WALLABAG_USERNAME"),
		UserPassword: os.Getenv("WALLABAG_PASSWORD"),
	}
	wallabago.SetConfig(wallabagConfig)

	// create data dir
	wd, _ := os.Getwd()
	_ = os.MkdirAll(filepath.Join(wd, "output"), os.ModePerm)
}

func main() {
	// get entries
	entries := getEntries(2) // debug
	fmt.Println(len(entries))

	// Create a new EPUB
	e, err := epub.NewEpub("My title")
	if err != nil {
		log.Println(err)
	}

	// Set the author
	e.SetAuthor("Wallabag")

	// add articles
	for _, entry := range entries {
		fmt.Println(entry.Title)

		// Add a section
		_, err = e.AddSection(entry.Content, entry.Title, "", "")
		if err != nil {
			log.Println(err)
		}
	}

	// write epub
	fmt.Println("Embedding images...")
	e.EmbedImages() // this has to stay here

	err = e.Write("output/My EPUB.epub")
	if err != nil {
		fmt.Println("Error creating EPUB")
	} else {
		fmt.Println("Successfully created EPUB")
	}
}
