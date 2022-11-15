package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type AssetResonse struct {
	Url    string
	Id     int64
	Assets []*Asset
}

type Asset struct {
	Url  string
	Id   int64
	Name string
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("please provide asset name")
	}

	assetNeed := make(map[string]struct{})
	for _, arg := range args {
		assetNeed[arg] = struct{}{}
	}

	client := http.DefaultClient
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if token == "" {
		log.Fatalln("missing GITHUB_ACCESS_TOKEN")
	}
	owner := os.Getenv("GITHUB_OWNER")
	if owner == "" {
		log.Fatalln("missing GITHUB_OWNER")
	}
	repo := os.Getenv("GITHUB_REPO")
	if repo == "" {
		log.Fatalln("missing GITHUB_REPO")
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	req.Header.Set("Authorization", "token "+token)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var assetResponse AssetResonse
	if err := json.Unmarshal(body, &assetResponse); err != nil {
		log.Fatalln(err)
	}
	var wg sync.WaitGroup
	for _, asset := range assetResponse.Assets {
		if _, ok := assetNeed[asset.Name]; !ok {
			continue
		}

		url := asset.Url
		name := asset.Name
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			req.Header.Set("Authorization", "token "+token)
			req.Header.Set("Accept", "application/octet-stream")

			res, err := client.Do(req)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			defer res.Body.Close()

			out, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
			if _, err := io.Copy(out, res.Body); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
		}()
	}
	wg.Wait()
}
