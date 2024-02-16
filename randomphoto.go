package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"regexp"

	"github.com/kris-nova/logger"
)

var photosWithFaces []Photo

func viewHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	if nil == photosWithFaces {
		log.Println("Getting faces")
		photosWithFaces, err = getPhotosWithFaces()
		if err != nil {
			log.Fatal("No faces got")
		} else {
			log.Printf("Got %d faces\n", len(photosWithFaces))
		}
	}
	offset := rand.IntN(len(photosWithFaces))
	photo := photosWithFaces[offset]
	resp, err := getPhoto(fmt.Sprintf("photos/%s/dl", photo.UID))

	if err != nil {
		log.Fatalln("Fail")
	}

	w.Write(resp)
}

var validPath = regexp.MustCompile("^/photo")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func get(request string) (*http.Response, error) {
	baseUrl := "http://10.0.20.10:2345/"
	resp, err := http.Get(baseUrl + "api/v1/" + request)
	return resp, err
}

func getPhoto(request string) ([]byte, error) {
	resp, err := get(request)
	if err == nil {
		body, err := io.ReadAll(resp.Body)
		return body, err
	}
	return nil, err
}

type Photo struct {
	UID string `json:"UID"`
}

func getQuery(query string) ([]Photo, error) {
	resp, err := get(fmt.Sprintf("photos?q=%s&count=10000", query))
	if err == nil {
		body, err := io.ReadAll(resp.Body)

		var photos []Photo
		err = json.Unmarshal(body, &photos)
		if err == nil {
			return photos, err
		}
	}
	return nil, err
}

func getPhotosWithFaces() ([]Photo, error) {
	return getQuery("faces:2")
}

func getPhotosWithPeople() ([]Photo, error) {
	return getQuery("person:karl|eli|belinda|isabella")
}

func main() {
	logger.Level = 4

	http.HandleFunc("/photo/", makeHandler(viewHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
