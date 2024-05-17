package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TaskManagerEventLoop() {
	for {
		executeTasks()
		time.Sleep(10 * time.Second)
	}
}

func executeTasks() {
	var targets []Target
	db.Where("is_active = ?", true).Find(&targets)

	for _, target := range targets {
		executeTask(target)
	}
}

func executeTask(target Target) {
	log.Println("Fetching URL:", target.URL, time.Now())

	newHash, err := getPageHash(target.URL)
	if err != nil {
		log.Printf("Error checking for URL change %s: %v", target.URL, err)
		return
	}

	db.Create(&History{
		TargetID: target.ID,
		Hash:     newHash,
	})
}

func getPageHash(url string) (string, error) {
	content, err := fetchPage(url)
	if err != nil {
		return "", err
	}

	p := bytes.NewReader(content)
	doc, _ := goquery.NewDocumentFromReader(p)

	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})

	docBytes := []byte(doc.Text())
	currentHash := calculateHash(docBytes)
	return currentHash, nil
}

func fetchPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func calculateHash(content []byte) string {
	hash := md5.Sum(content)
	return fmt.Sprintf("%x", hash)
}
