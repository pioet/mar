/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// fetch the title of webpage
func GetWebTitle(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // 10 sec timeout
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}
	// add headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	// 4. Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	// check HTTP code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d %s", resp.StatusCode, resp.Status)
	}
	// using goquery load HTML doc
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML doc: %w", err)
	}
	// use CSS selecor to find <title>
	title := doc.Find("title").Text()
	if title == "" {
		return "", fmt.Errorf("failed to find title. ")
	}
	return title, nil
}
