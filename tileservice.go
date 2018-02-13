// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample storage demonstrates use of the cloud.google.com/go/storage package from App Engine flexible environment.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"os"
	"strings"
	"strconv"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)


var (
	storageClient *storage.Client
	// Set this in app.yaml when running in production.
	bucketName = os.Getenv("GCLOUD_STORAGE_BUCKET")
)

func main() {
	ctx := context.Background()

	var err error
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/map/", serveMapTile)
	http.HandleFunc("/", runinfo)
	//http.HandleFunc("/upload", uploadHandler)

	fmt.Println("server is starting up now")
	appengine.Main()
}

func runinfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "url "+r.URL.Path)
}

func serveMapTile(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	
	ctx := appengine.NewContext(r)
	bucket := storageClient.Bucket(bucketName)

	splits := strings.Split(r.URL.Path, "/")
	trf, _ := r.URL.Query()["transform"]
	
	if len(splits) != 8 {
		msg := fmt.Sprintf("Incorrect url format, expected /map/{mapname}/{style}/{level}/{col}/{row}/tile.png")
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	//Get splits
	mapname := splits[2]
	style := splits[3]
	level, _ := strconv.Atoi(splits[4])
	col, _ := strconv.Atoi(splits[5])
	row, _ := strconv.Atoi(splits[6])

	//Check that level col and row are ints

	//Get Transformation (if available)
	if(trf != nil && trf[0] != ""){
		trfFilename := mapname + "/" + style + "/" + trf[0] + ".json"		
		trfreader, err := bucket.Object(trfFilename).NewReader(ctx)
		if err != nil {
			trfFilename := mapname + "/" + trf[0] + ".json"		
			trfreader, err = bucket.Object(trfFilename).NewReader(ctx)
			if err != nil {
				msg := fmt.Sprintf("Could not get transformation: %v for %v", err, trfFilename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}
		defer trfreader.Close()

		transfContent, err := ioutil.ReadAll(trfreader)
		var trules map[string]interface{}
		err = json.Unmarshal(transfContent, &trules)
		checkErrors(w, err)

		//Transform
		level = level - int(trules["level"].(float64))
		col = col - (level + 1) * int(trules["col"].(float64))
		row = row - (level + 1) * int(trules["row"].(float64))
		
	}

	//Construct map path
	filename := mapname + "/" + style + "/" + strconv.Itoa(level) + "/" + strconv.Itoa(col) + "/" + strconv.Itoa(row) + ".png"

	//Get the file
	reader, err := bucket.Object(filename).NewReader(ctx)
	if err != nil {
		msg := fmt.Sprintf("Could not get file from store: %v for %v", err, filename)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer reader.Close()

	w.Header().Add("content-type", "image/png")

	//Copy the content
	_, err2 := io.Copy(w, reader)
	if err2 != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

}

func checkErrors(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(500)
		panic(err.Error())
	}
}

/*
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	ctx := appengine.NewContext(r)

	f, fh, err := r.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf("Could not get file: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer f.Close()

	sw := storageClient.Bucket(bucket).Object(fh.Filename).NewWriter(ctx)
	if _, err := io.Copy(sw, f); err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if err := sw.Close(); err != nil {
		msg := fmt.Sprintf("Could not put file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	u, _ := url.Parse("/" + bucket + "/" + sw.Attrs().Name)

	fmt.Fprintf(w, "Successful! URL: https://storage.googleapis.com%s", u.EscapedPath())
}
*/
