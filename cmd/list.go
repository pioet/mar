/*
Copyright © 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved bookmarks",
	Long: `Lists all saved bookmarks stored in the database.
By default, bookmarks are shown in a short, compact view.
Use the '--long' or '-l' flag to display detailed information. `,
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var bookmarks []models.Bookmark
		DB.Find(&bookmarks)
		// if all, _ := cmd.Flags().GetBool("all"); all {
		// 	DB.Find(&bookmarks) // find all bookmarks
		// } else {
		// 	DB.Where("Tag IS NOT NULL").Find(&bookmarks) // find bookmarks with tag
		// }

		// no bookmarks found.
		if len(bookmarks) == 0 {
			pterm.FgGreen.Println("No bookmark found.")
			return
		}

		if long, _ := cmd.Flags().GetBool("long"); long {
			for _, b := range bookmarks {
				b.PrintLong()
			}
		} else {
			for _, b := range bookmarks {
				b.PrintShort()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("long", "l", false, "Display in detailed format")
}
