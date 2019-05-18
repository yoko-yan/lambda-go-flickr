package main

import (
	"bytes"
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	ENDPOINT = "https://api.flickr.com/services/rest/"
)

// レスポンスJSONデータ用構造体
type APIResponse struct {
	Photos struct {
		Page    int    `json:"page"`
		Pages   int    `json:"pages"`
		Perpage int    `json:"perpage"`
		Total   string `json:"total"`
		Photo   []struct {
			ID       string `json:"id"`
			Owner    string `json:"owner"`
			Secret   string `json:"secret"`
			Server   string `json:"server"`
			Farm     int    `json:"farm"`
			Title    string `json:"title"`
			Ispublic int    `json:"ispublic"`
			Isfriend int    `json:"isfriend"`
			Isfamily int    `json:"isfamily"`
		} `json:"photo"`
	} `json:"photos"`
	Stat string `json:"stat"`
}

// 返却用の構造体
type ReturnResponse struct {
	Status int                   `json:"status"`
	Photos []ReturnResponsePhoto `json:"photo"`
}

// 返却用の構造体のItem
type ReturnResponsePhoto struct {
	URL string `json:"url"`
}

// 写真情報の構造体
type PhotoInfo struct {
	ID       string `json:"id"`
	Owner    string `json:"owner"`
	Secret   string `json:"secret"`
	Server   string `json:"server"`
	Farm     int    `json:"farm"`
	Title    string `json:"title"`
	Ispublic int    `json:"ispublic"`
	Isfriend int    `json:"isfriend"`
	Isfamily int    `json:"isfamily"`
	URL      string `json:"url"`
}

func crearteJson(photos map[string]PhotoInfo) ReturnResponse {

	var data = ReturnResponse{}
	data.Status = 200
	data.Photos = []ReturnResponsePhoto{}
	for _, photoInfo := range photos {
		var photo ReturnResponsePhoto
		photo.URL = photoInfo.URL
		data.Photos = append(data.Photos, photo)
	}
	return data
}

func buildQuery(q map[string]string) string {
	queries := make([]string, 0)
	for k, v := range q {
		qq := fmt.Sprintf("%s=%s", k, v)
		queries = append(queries, qq)
	}
	return strings.Join(queries, "&")
}

type MyEvent struct {
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
	PerPage string `json:"perpage"`
	Page    string `json:"page"`
}

func hello(ctx context.Context, params MyEvent) (interface{}, error) {

	lat := "35.658587"
	lon := "139.745433"
	per_page := "15"
	page := "1"

	if params.Lat != "" {
		lat = params.Lat
	}

	if params.Lon != "" {
		lon = params.Lon
	}

	if params.PerPage != "" {
		per_page = params.PerPage
	}

	if params.Page != "" {
		page = params.Page
	}

	// URL生成
	q := map[string]string{
		"method":         "flickr.photos.search",
		"api_key":        os.Getenv("FLICKR_API_KEY"),
		"lat":            lat,
		"lon":            lon,
		"format":         "json",
		"nojsoncallback": "1",
		"per_page":       per_page,
		"page":           page,
		"radius":         "1",
	}
	url := fmt.Sprintf("%s?%s", ENDPOINT, buildQuery(q))
	// URLを叩いてデータを取得
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 取得したデータをJSONデコード
	var data APIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	// 取得したデータを整形して出力する
	var photos map[string]PhotoInfo = map[string]PhotoInfo{}
	for _, photo := range data.Photos.Photo {
		var photoInfo PhotoInfo
		photoInfo.ID = photo.ID
		photoInfo.Owner = photo.Owner
		photoInfo.Secret = photo.Secret
		photoInfo.Server = photo.Server
		photoInfo.Farm = photo.Farm
		photoInfo.Title = photo.Title
		photoInfo.Ispublic = photo.Ispublic
		photoInfo.Isfriend = photo.Isfriend
		photoInfo.Isfamily = photo.Isfamily
		// 写真URLを組み立てる
		var photoUrl = fmt.Sprintf("https://farm%d.staticflickr.com/%s/%s_%s_m.jpg", photo.Farm, photo.Server, photo.ID, photo.Secret)
		photoInfo.URL = photoUrl
		// 	fmt.Printf("https://farm%d.staticflickr.com/%s/%s_%s.jpg", photo.Farm, photo.Server, photo.ID, photo.Secret)
		// 	fmt.Printf("%s: [%d][%s][%s]\n", photo.ID, photo.Farm, photo.Server, photo.Secret)
		photos[photo.ID] = photoInfo
	}

	var data2 = crearteJson(photos)

	// jsonエンコード
	outputJson, err := json.Marshal(&data2)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	buf.Write(outputJson)
	// println(buf.String())

	return data2.Photos, nil
}

func main() {
	lambda.Start(hello)
}
