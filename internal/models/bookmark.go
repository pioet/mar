/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package models

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// "bytes"

	// toml "github.com/pelletier/go-toml/v2"
	"github.com/pterm/pterm"
	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	Bid     uint    `gorm:"column:bid;unique;not null"`
	Title   string  `gorm:"column:title;not null"`
	URI     string  `gorm:"column:uri;unique;not null"`
	Comment string  `gorm:"column:comment;not null"`
	Tag     *string `gorm:"column:tag;uniqueIndex"`
}

func (b *Bookmark) PrintLong() {
	// bid
	pterm.NewStyle(pterm.FgLightGreen, pterm.Bold).Print(b.Bid, ". ")
	// tag
	if b.Tag != nil {
		pterm.FgGray.Print("@")
		pterm.NewStyle(pterm.FgLightRed, pterm.Bold).Print(*b.Tag, " ")
	}
	// title
	pterm.NewStyle(pterm.FgLightGreen, pterm.Bold).Println(b.Title)
	// uri
	spaceWidth := len(strconv.Itoa(int(b.Bid))) + 2
	pterm.FgGray.Print(fmt.Sprintf("%-*s", spaceWidth, ">"))
	pterm.FgLightBlue.Println(b.URI)
	// comment
	if strings.TrimSpace(b.Comment) != "" {
		pterm.FgGray.Print(fmt.Sprintf("%-*s", spaceWidth, "#"))
		pterm.FgWhite.Println(b.Comment)
	}
	// new line to separate the bookmark.
	pterm.FgWhite.Println()
}

const aliasWidth int16 = 16

func (b *Bookmark) PrintShort() {
	// tag
	if b.Tag != nil {
		pterm.FgGray.Print("@")
		pterm.FgLightRed.Printf("%-*s", aliasWidth, *b.Tag)
	} else {
		pterm.DefaultBasicText.Printf("%-*s", aliasWidth+1, "")
	}
	// title
	pterm.FgGray.Print(" -> ")
	pterm.NewStyle(pterm.FgLightGreen, pterm.Bold).Println(b.Bid, ". ", b.Title)
}

func (b *Bookmark) GetPlainText() string {
	plainTextPattern := "Title='%s'\nURI='%s'\nComment='%s'\nTag='%s'\n\n"
	content := ""
	if b.Tag != nil {
		content = fmt.Sprintf(plainTextPattern, b.Title, b.URI, b.Comment, *b.Tag)
	} else {
		content = fmt.Sprintf(plainTextPattern, b.Title, b.URI, b.Comment, "")
	}
	return content
}

// GetNextBid returns the next bid for a new bookmark.
func GetNextBid(db *gorm.DB) uint {
	var maxBid uint
	err := db.Model(&Bookmark{}).Select("COALESCE(MAX(bid), 0)").Row().Scan(&maxBid)
	if err != nil {
		return 1
	}
	return maxBid + 1
}

func GetBookmarkByBid(db *gorm.DB, bid uint) (*Bookmark, error) {
	var bm Bookmark
	if err := db.Where("bid = ?", bid).First(&bm).Error; err != nil {
		return nil, err
	}
	return &bm, nil
}

func GetBookmarkByTag(db *gorm.DB, tag string) (*Bookmark, error) {
	var bm Bookmark
	if err := db.Where("tag = ?", tag).First(&bm).Error; err != nil {
		return nil, err
	}
	return &bm, nil
}

func GetBookmarkByURI(db *gorm.DB, uri string) (*Bookmark, error) {
	var bm Bookmark
	if err := db.Where("uri = ?", uri).First(&bm).Error; err != nil {
		return nil, err
	}
	return &bm, nil
}

func GetBookmarkByName(db *gorm.DB, name string) (*Bookmark, error) {
	if bid_int, err := strconv.Atoi(name); err == nil {
		return GetBookmarkByBid(db, uint(bid_int))
	} else {
		return GetBookmarkByTag(db, name)
	}
}

func (b *Bookmark) BackupToDustbin(dustbinFilePath string) {
	// Open the file in append mode.
	// os.O_APPEND: append to the file when writing.
	// os.O_CREATE: create the file if it doesn't exist.
	// os.O_WRONLY: open the file write-only.
	// 0644: file permissions (read/write for owner, read-only for others).
	file, err := os.OpenFile(dustbinFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		pterm.FgRed.Println("Error: opening dustbin file.")
		return // Exit the function if we can't open the file
	}
	defer file.Close() // Ensure the file is closed when the function exits

	// Write the formatted string to the file.
	if _, err := file.WriteString(b.GetPlainText()); err != nil {
		pterm.FgRed.Println("Error: writing dustbin file.")
	}
}

func (b *Bookmark) CreateDb(db *gorm.DB) error {
	if err := db.Create(b).Error; err != nil {
		return fmt.Errorf("error: failed to create bookmark in database: %w", err)
	}
	return nil
}

func (b *Bookmark) GetHTML() string {
	// Chrome bookmark time format: Unix microseconds
	// We'll just use seconds for simplicity
	timeStr := strconv.FormatInt(b.CreatedAt.Unix()*1000000, 10)
	title := htmlEscape(b.Title)
	uri := htmlEscape(b.URI)

	return fmt.Sprintf(`    <DT><A HREF="%s" ADD_DATE="%s">%s</A>
`, uri, timeStr, title)
}

// 简单HTML转义
func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		`&`, "&amp;",
		`"`, "&quot;",
		`<`, "&lt;",
		`>`, "&gt;",
	)
	return replacer.Replace(s)
}
