/*
Copyright © 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/pioet/mar/internal/models"
)

var addCmd = &cobra.Command{
	Use:   "add <URI> [tag]",
	Short: "Adds a new bookmark",
	Long:  `Adds a new bookmark to your collection. You'll need to provide the bookmark's URI (Uniform Resource Identifier), and you can optionally add a tag to help you access it later. `,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// validate the URI.
		// case 1: webpage.
		// case 2: local file. Then transfer to absolute path.
		uriVar, err := models.NewUri(args[0])
		if err != nil {
			pterm.FgRed.Println(err)
			return
		}
		uri := uriVar.UriOutput

		// set the tag if provided
		var tagInput string
		if len(args) > 1 {
			tagInput = args[1]
		}

		// flags-tile
		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			title, _ = uriVar.GetTitle() // automatically fetch the title from the URI
		}
		// flags-comment
		comment, _ := cmd.Flags().GetString("comment")

		// create the bookmark object
		bookmark := models.Bookmark{
			Bid:   models.GetNextBid(DB),
			Title: title,
			URI:   uri,
		}
		if comment != "" {
			bookmark.Comment = comment
		}
		if tagInput != "" {
			err := IsValidTag(tagInput)
			if err != nil {
				pterm.FgRed.Println(err)
				return
			}
			bookmark.Tag = &tagInput
		}

		if err = bookmark.CreateDb(DB); err != nil {
			pterm.FgRed.Println(err)
			return
		} else {
			pterm.FgWhite.Println("Bookmark added:", bookmark.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("title", "t", "", "Specify a title for the bookmark. (Default: automatically fetched from the URI)")
	addCmd.Flags().StringP("comment", "c", "", "Add additional information to the bookmark. (Default: none)")
}
