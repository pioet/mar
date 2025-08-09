/*
Copyright © 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"os"

	"github.com/pioet/mar/internal/models"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all bookmarks to a file",
	Long:  `Exports all bookmarks from the database to a specified file, supporting text and html formats. `,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var bookmarks []models.Bookmark
		DB.Find(&bookmarks)

		outputPath, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")

		file, err := os.Create(outputPath)
		if err != nil {
			pterm.FgRed.Printfln("Error creating file at %s", outputPath)
			return
		}
		defer file.Close()

		switch format {
		case "", "txt":
			for _, bookmark := range bookmarks {
				_, err := file.WriteString(bookmark.GetPlainText())
				if err != nil {
					pterm.FgRed.Println("Error writing bookmark to file")
					return
				}
			}
			pterm.FgGreen.Println("Exported bookmarks in plain text format.")

		case "html":
			// Chrome bookmark file header
			file.WriteString(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
`)
			for _, bookmark := range bookmarks {
				file.WriteString(bookmark.GetHTML())
			}
			file.WriteString("</DL><p>\n")
			pterm.FgGreen.Println("Exported bookmarks in HTML format.")

		default:
			pterm.FgRed.Printfln("Unsupported format: %s", format)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("output", "o", "", "Path to the output file where bookmarks will be exported")
	exportCmd.MarkFlagRequired("output")
	exportCmd.Flags().StringP("format", "f", "txt", "Output file format: txt or html")
}
