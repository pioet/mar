/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reassign sequential ID to all bookmarks",
	Long:  `Reassign (reindex) the IDs of all bookmarks starting from 1. `,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// get all bookmarks.
		var bookmarks []models.Bookmark
		DB.Find(&bookmarks)

		// reindex the ID of bookmarks.
		for currentIndex, bookmark := range bookmarks {
			bookmark.Bid = uint(currentIndex + 1)
			DB.Save(&bookmark)
		}
		pterm.FgWhite.Println("Reindexed bookmarks successfully.")
	},
}

func init() {
	rootCmd.AddCommand(reindexCmd)
}
