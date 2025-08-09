/*
Copyright © 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "Import bookmarks from a plaintext file",
	Long:    `Imports bookmarks from a specified plaintext file and saves them to the database. `,
	Args:    cobra.NoArgs,
	Example: `mar import --input bookmarks.txt`,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, err := cmd.Flags().GetString("input")
		if err != nil {
			pterm.FgRed.Printfln("--input is required")
			return
		}

		file, err := os.Open(inputPath)
		if err != nil {
			pterm.FgRed.Printfln("Error opening file at %s: %v", inputPath, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var newBookmarks []models.Bookmark
		var currentBookmark models.Bookmark

		for scanner.Scan() {
			line := scanner.Text()

			if strings.TrimSpace(line) == "" {
				if currentBookmark.Title != "" {
					newBookmarks = append(newBookmarks, currentBookmark)
				}
				currentBookmark = models.Bookmark{}
				continue
			}

			parts := strings.SplitN(line, "='", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(parts[1], "'")

			switch key {
			case "Title":
				currentBookmark.Title = value
			case "URI":
				currentBookmark.URI = value
			case "Comment":
				currentBookmark.Comment = value
			case "Tag":
				if value != "" {
					currentBookmark.Tag = &value
				} else {
					currentBookmark.Tag = nil
				}
			}
		}

		if currentBookmark.Title != "" {
			newBookmarks = append(newBookmarks, currentBookmark)
		}

		if len(newBookmarks) == 0 {
			pterm.FgYellow.Println("No bookmarks found in the file.")
			return
		}

		pterm.FgGreen.Printfln("Found %d bookmarks to import...", len(newBookmarks))

		var importedCount int
		var failedImports []string

		for _, bookmark := range newBookmarks {
			bookmark.Bid = models.GetNextBid(DB)
			result := DB.Create(&bookmark)
			if result.Error != nil {
				failedImports = append(failedImports, bookmark.Title)
				pterm.FgRed.Printfln("Failed to import bookmark '%s': %v", bookmark.Title, result.Error)
			} else {
				importedCount++
			}
		}

		pterm.FgGreen.Printfln("Successfully imported %d bookmarks.", importedCount)

		if len(failedImports) > 0 {
			pterm.FgYellow.Printfln("Failed to import the following bookmarks:")
			for _, title := range failedImports {
				pterm.FgYellow.Printfln("- %s", title)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("input", "i", "", "Path to the input file from which bookmarks will be imported")
	importCmd.MarkFlagRequired("input")
}
