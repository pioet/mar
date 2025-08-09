/*
Copyright © 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all saved bookmarks",
	Long:  `Delete all saved bookmarks. Before deletion, all bookmarks will be backed up to the location specified in your configuration file (database.dustbin). `,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Use pterm's built-in confirmation for a cleaner check.
		confirm, _ := pterm.DefaultInteractiveConfirm.Show("Are you sure you want to continue?")
		if !confirm {
			pterm.FgYellow.Println("Action aborted.")
			return
		}
		// find all bookmarks
		var bookmarks []models.Bookmark
		DB.Find(&bookmarks)
		dustbinPath := viper.GetString("database.dustbin")
		for _, b := range bookmarks {
			b.BackupToDustbin(dustbinPath)
		}

		// Add a "Where" clause to explicitly indicate you want to delete all records.
		result := DB.Unscoped().Where("1 = 1").Delete(&models.Bookmark{})

		// Handle potential database errors.
		if result.Error != nil {
			pterm.FgRed.Printf("Error: Failed to clear bookmarks. Details: %v\n", result.Error)
			return
		}

		// Provide user feedback based on the outcome.
		if result.RowsAffected == 0 {
			pterm.FgYellow.Println("No bookmarks were found to clear.")
		} else {
			pterm.FgGreen.Printf("Successfully cleared %d bookmark(s).\n", result.RowsAffected)
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
