package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	tsumekiri "github.com/ariarijp/tsumekiri/internal/pkg"
	"github.com/imroc/req"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func getFileList(srv *drive.Service, qStr string) *drive.FileList {
	r, err := srv.Files.List().PageSize(1000).
		Fields("nextPageToken, files(id,name,size,hasThumbnail,thumbnailLink,contentHints/thumbnail)").
		Q(qStr).
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	return r
}

func mkdir(name *string) {
	_, err := os.Stat(*name)
	if err != nil {
		os.Mkdir(*name, 0755)
	}
}

func downloadFiles(fl *drive.FileList) {
	downloadFolder := "thumbnails"
	mkdir(&downloadFolder)

	if len(fl.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, f := range fl.Files {
			if f.HasThumbnail {
				res, _ := req.Get(f.ThumbnailLink)
				res.ToFile(filepath.Join(downloadFolder, f.Name+".jpg"))
			}
		}
	}
}

func Run() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := tsumekiri.GetClient(config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	fl := getFileList(srv, os.Args[1])
	downloadFiles(fl)
}
