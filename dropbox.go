package main

import (
	"github.com/pearkes/Dropbox-Go/dropbox"
	"io/ioutil"
	"log"
	"os"
)

func newDropboxUrl() string {
	s := dropbox.Session{
		AppKey:      DROPBOX_KEY,
		AppSecret:   DROPBOX_SECRET,
		AccessType:  "app_folder",
		RedirectUri: os.Getenv("DROPBOX_CALLBACK"),
	}
	url := dropbox.GenerateAuthorizeUrl(s.AppKey, &dropbox.Parameters{RedirectUri: s.RedirectUri})
	return url
}

// Fills a dropbox with the default stuff
func fillDropbox(token string) {
	// The session to write with
	s := dropbox.Session{
		AppKey:     DROPBOX_KEY,
		AppSecret:  DROPBOX_SECRET,
		AccessType: "app_folder",
		Token:      token,
	}

	files := []string{"01-east-river_01.jpg", "01-east-river_02.md", "02-code_01.jpg", "02-code_02.md", "03-coffee_01.jpg", "03-coffee_02.md", "first-victory-theme.css", "first-victory-theme.js"}

	for _, f := range files {
		u := dropbox.Uri{
			Root: "sandbox",
			Path: "/" + f,
		}
		content, err := ioutil.ReadFile("files/" + f)
		_, err = dropbox.UploadFile(s, content, u, nil)
		if err != nil {
			log.Printf("error uploading default file: %s", err.Error())
		}
	}
}
