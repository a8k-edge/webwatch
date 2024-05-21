package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aryann/difflib"
	"gorm.io/gorm"
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

	resp, err := http.Get(target.URL)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", target.URL, err)
		return
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})

	hash := calculateHash([]byte(doc.Text()))
	isChanged := false
	var history History
	err = db.Order("created_at desc").Where("target_id == ?", target.ID).First(&history).Error
	var diffBytes []byte
	prev := make([]string, 0)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		isChanged = true
	} else if err != nil {
		log.Println("Error querying target:", err)
		return
	} else {
		isChanged = history.Hash != hash
		preDiffs := make([]difflib.DiffRecord, 0)
		if err := json.Unmarshal([]byte(history.Diff), &preDiffs); err != nil {
			log.Println("Error unmarshal prev history diff:", err)
			return
		}
		for i := range preDiffs {
			if preDiffs[i].Delta == difflib.LeftOnly {
				continue
			}
			prev = append(prev, preDiffs[i].Payload)
		}
	}

	diffBytes, err = json.Marshal(difflib.Diff(prev, strings.Split(doc.Text(), "\n")))
	if err != nil {
		log.Println("Error marshal history diff:", err)
		return
	}

	db.Create(&History{
		TargetID:   target.ID,
		Hash:       hash,
		IsChanged:  isChanged,
		StatusCode: resp.StatusCode,
		Diff:       string(diffBytes),
	})
}

func calculateHash(content []byte) string {
	hash := md5.Sum(content)
	return fmt.Sprintf("%x", hash)
}
