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

func Chunk[T any](slice []T, n uint64) <-chan []T {
	if n == 0 {
		panic("n can`t be less than 1")
	}

	channel := make(chan []T, 1)

	go func() {
		defer close(channel)
		for i := uint64(0); i < uint64(len(slice)); i += n {
			// Clamp the last chunk to the slice bound as necessary.
			end := min(n, uint64(len(slice[i:])))

			// Set the capacity of each chunk so that appending to a chunk does
			// not modify the original slice.
			channel <- slice[i : i+end : i+end]
		}
	}()

	return channel
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
	var outputFileIndex int

	// get entries
	entries := getEntries(200)

	// create EPUBs
	chunks := Chunk(entries, 20)

	for chunk := range chunks {
		filename := fmt.Sprintf("Wallabag %v.epub", outputFileIndex)
		fmt.Printf("Creating file: %s\n", filename)

		// Create a new EPUB
		e, err := epub.NewEpub(filename)
		if err != nil {
			log.Println("Error initializing Epub:", err)
		}

		// Set the author
		e.SetAuthor("Wallabag")

		// Add articles
		for _, entry := range chunk {
			fmt.Printf("Adding: %s\n", entry.Title)

			// Add section
			title := entry.Title
			content := fmt.Sprintf("<h1>%s</h1>%s", title, entry.Content)

			_, err = e.AddSection(content, entry.Title, "", "")
			if err != nil {
				log.Println("Error adding article", err)
			}
		}

		// write epub
		fmt.Println("Embedding images...")
		e.EmbedImages() // this has to stay here

		err = e.Write(fmt.Sprintf("output/%s", filename))
		if err != nil {
			fmt.Println("Error creating EPUB")
		} else {
			fmt.Println("Successfully created EPUB")
		}

		outputFileIndex++
	}
}
