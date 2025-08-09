/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm <id_or_tag>...",
	Short: "Remove bookmarks specified by ID or Tag",
	Long:  `Remove bookmarks specified by ID or Tag. You can find backups in the location specified in your configuration file (database.dustbin).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			bookmark, err := models.GetBookmarkByName(DB, arg)
			if err != nil {
				pterm.FgRed.Println("Error: failed to find bookmark. ")
				return
			}

			// result := DB.Delete(&models.Bookmark{}, bookmark.ID)
			result := DB.Unscoped().Delete(&models.Bookmark{}, bookmark.ID)
			if result.Error != nil {
				pterm.FgRed.Println("Error: failed to remove bookmark. ")
				return
			}
			if result.RowsAffected == 0 {
				pterm.FgRed.Printf("No bookmark found with [%s].\n", arg)
			} else {
				pterm.FgWhite.Printf("Bookmark with <%s> removed successfully.\n", arg)
			}
			dustbinPath := viper.GetString("database.dustbin")
			bookmark.BackupToDustbin(dustbinPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
