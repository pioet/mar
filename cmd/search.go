/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"strings"

	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <keywords>...",
	Short: "Search for bookmarks by multiple keywords",
	Long:  `Search for bookmarks that contain all specified keywords (case-insensitive) in their title, URI, comment, or tag. `,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := DB
		// Build the query dynamically for each keyword
		for _, keyword := range args {
			// Use a sub-query to group the OR conditions for each keyword
			// This ensures the AND logic for multiple keywords works correctly
			query = query.Where(
				DB.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(keyword)+"%").
					Or("LOWER(uri) LIKE ?", "%"+strings.ToLower(keyword)+"%").
					Or("LOWER(comment) LIKE ?", "%"+strings.ToLower(keyword)+"%").
					Or("LOWER(tag) LIKE ?", "%"+strings.ToLower(keyword)+"%"),
			)
		}

		var bookmarks []models.Bookmark
		result := query.Find(&bookmarks)

		if result.Error != nil {
			pterm.FgRed.Printf("Error during search: %v\n", result.Error)
			return
		}
		// Check if any bookmarks were found
		if len(bookmarks) == 0 {
			pterm.FgYellow.Println("No bookmarks found matching all keywords.")
			return
		}
		// Print the results
		for _, bookmark := range bookmarks {
			bookmark.PrintLong()
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
