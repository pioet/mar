/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show <id_or_tag>...",
	Short: "Display bookmarks specified by ID or Tag",
	Long:  "Display detailed information of bookmarks by ID or tag. ",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			bookmark, err := models.GetBookmarkByName(DB, arg)
			if err != nil {
				pterm.FgRed.Println("Error: failed to get bookmark. ")
				return
			}
			bookmark.PrintLong()
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
