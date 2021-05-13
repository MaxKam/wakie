package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3" // Driver for sql
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
		if cfgFile == "" {
			cfgFile = "~/.config/wakie"
		}

		if saveDbPath == "$HOME/.config/wakie/" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Unable to get users home dir: %s", err)
			}

			saveDbPath = homeDir + "/.config/wakie/wakie.db"
		}
		file := dbFileExists(saveDbPath)

		if file {
			log.Fatal("Error - Database already exists at the supplied path")
		}

		createDB(saveDbPath)

	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVar(&saveDbPath, "saveDb", "$HOME/.config/wakie/", "Path to folder where to save database file")
}

func dbFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func createDB(path string) {
	createTableStmt := `CREATE TABLE 'computers' 
	('ID' INTEGER PRIMARY KEY AUTOINCREMENT, 
    'Alias' STRING NULL UNIQUE, 
    'MAC_Address' STRING NULL UNIQUE);`

	_, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create db file. Please make sure directory path is correct: %s", err)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Unable to open db file: %s", err)
	}

	_, err = db.Exec(createTableStmt)
	if err != nil {
		log.Fatalf("Unable to create table in db file: %s", err)
	}

	fmt.Println(" - database file and table successfully created.")

}
