package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	// "gopkg.in/yaml.v2"
)

var saveDbPath string

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates SQLite Database and config file",
	Long: `Used to setup Wakie by creating a SQLite Database file at the specified path
and creates a config file. 
	
If config flag is not set, will default to saving the config file in ~/.config/wakie. If that
fails will attempt to save config file in same folder as the app.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVar(&saveDbPath, "saveDb", "~/.config/wakie/", "Path to folder where to save database file")
	setupCmd.MarkFlagRequired("saveDb")
}
