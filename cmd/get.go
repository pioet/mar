/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <id_or_tag>",
	Short: "Retrieve a bookmark's URI",
	Long:  `Retrieves the URI (Uniform Resource Identifier) for a specific bookmark`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bookmark, err := models.GetBookmarkByName(DB, args[0])
		if err != nil {
			pterm.FgRed.Println("Error: failed to get bookmark. ")
			return
		}
		pterm.DefaultBasicText.Println(bookmark.URI)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
