/*
Copyright © 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"strings"

	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <id_or_tag>",
	Short: "Modify an existing bookmark",
	Long:  `Interactively modify an existing bookmark's details (title, URI, tag, and comment). `,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bookmark, err := models.GetBookmarkByName(DB, args[0])
		if err != nil {
			pterm.FgRed.Println("Error: The specified bookmark was not found. ")
			return
		}
		editBookmarkInteractively(bookmark)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	// editCmd.Flags().BoolP("editor", "e", false, "Use default editor to edit. ")
}

// editBookmarkInteractively edit bookmark interactively
func editBookmarkInteractively(bookmark *models.Bookmark) {
	// interactively edit bookmark details
	bookmark.PrintLong()

	// choose the field to edit
	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions([]string{"Edit Title", "Edit URI", "Edit Comment", "Edit Tag"}).Show()
	switch selectedOption {
	case "Edit Title":
		titleInput, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("New Title").Show()
		bookmark.Title = strings.TrimSpace(titleInput)
	case "Edit URI":
		uriInput, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("New URI").Show()
		bookmark.URI = strings.TrimSpace(uriInput)
	case "Edit Comment":
		commentInput, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("New Comment").Show()
		bookmark.Comment = strings.TrimSpace(commentInput)
	case "Edit Tag":
		tagInput, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("New Tag (remain blank to delete)").Show()
		tagInput = strings.TrimSpace(tagInput)
		// Tag(pointer) can not be empty
		if tagInput == "" {
			bookmark.Tag = nil
		} else {
			bookmark.Tag = &tagInput
		}
	}

	// save the bookmark
	if err := DB.Save(bookmark).Error; err != nil {
		pterm.FgRed.Println("Error: Failed to update bookmark. ")
		return
	}
	pterm.FgWhite.Printf("Bookmark (ID: %d) successfully updated.\n", bookmark.Bid)
}
