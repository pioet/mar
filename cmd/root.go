/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/glebarez/sqlite"
	"github.com/pioet/mar/internal/models"
	"github.com/pioet/mar/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var rootCmd = &cobra.Command{
	Use:     "mar",
	Short:   "Manage and navigate your bookmarks effifently",
	Long:    `Mar is a CLI to manage and navigate bookmarks effifently. For more details: https://github.com/pioet/mar`,
	Args:    cobra.MaximumNArgs(1),
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		bookmark, err := models.GetBookmarkByName(DB, args[0])
		if err != nil {
			pterm.FgRed.Println("Error: The specified bookmark was not found. ")
			return
		}
		utils.OpenURI(bookmark.URI)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initDb)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	appdataDir := filepath.Join(home, ".bookmark-cli")
	if appdataDir == "" {
		pterm.FgRed.Println("Error: Application data directory is not set. ")
		os.Exit(1)
	}

	viper.AddConfigPath(appdataDir)
	viper.SetConfigType("toml")
	viper.SetConfigName(".config.toml")
	configPath := filepath.Join(appdataDir, ".config.toml")

	// set default values for configuration
	viper.SetDefault("database.path", filepath.Join(appdataDir, "bookmarks.db"))
	viper.SetDefault("database.dustbin", filepath.Join(appdataDir, "dustbin.txt"))

	// read config file
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Case 1: Config file not found.
			// Ensure the directory exists before trying to write the file.
			if _, dirErr := os.Stat(appdataDir); os.IsNotExist(dirErr) {
				pterm.FgGreen.Printf("Creating application data directory: %s\n", appdataDir)
				if mkdirErr := os.MkdirAll(appdataDir, 0755); mkdirErr != nil {
					pterm.FgRed.Printf("Failed to create config directory %s: %v\n", appdataDir, mkdirErr)
					os.Exit(1)
				}
			}
			// Write the default config.
			if writeErr := viper.SafeWriteConfigAs(configPath); writeErr != nil {
				pterm.FgRed.Printf("Failed to create default config file at %s: %v\n", configPath, writeErr)
				os.Exit(1)
			}
		} else {
			// Case 2: Other types of errors (e.g., parsing error, permission issues).
			pterm.FgRed.Printf("Error reading config file: %v.", err)
			os.Exit(1)
		}
	}
}

var DB *gorm.DB // global database connection

// initDb initializes the database connection and runs migrations.
func initDb() {
	dbPath := viper.GetString("database.path")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		pterm.FgRed.Printf("Failed to connect to database at %s: %v\n", dbPath, err)
		cobra.CheckErr(err)
	}
	err = db.AutoMigrate(&models.Bookmark{})
	if err != nil {
		pterm.FgRed.Printf("Failed to migrate database: %v\n", err)
		cobra.CheckErr(err)
	}
	DB = db
}

func isSubCommand(parentCmd *cobra.Command, name string) bool {
	for _, cmd := range parentCmd.Commands() {
		if cmd.Name() == name {
			return true
		}
	}
	return false
}

// check whether the tag is valid.
func IsValidTag(tagStr string) error {
	// check is the tag is as same as the sub-command name
	if isSubCommand(rootCmd, tagStr) {
		return fmt.Errorf("error: tag '%s' is a reserved command name. ", tagStr)
	}
	// check the first letter of the tag
	firstChar := rune(tagStr[0])
	if !((firstChar >= 'a' && firstChar <= 'z') ||
		(firstChar >= 'A' && firstChar <= 'Z') ||
		(firstChar == '_')) {
		return fmt.Errorf("error: tag '%s' should begin with a letter or underscore. ", tagStr)
	}

	// check the letter contents
	if !regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*$`).MatchString(tagStr) {
		return fmt.Errorf("error: tag '% s' can contain only letters, numbers, and underscores. ", tagStr)
	}

	return nil
}
