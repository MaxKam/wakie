package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

		// Check if database file already exists at specified path and creates file if not
		if saveDbPath == "$HOME/.config/wakie" {
			saveDbPath = homeDir + "/.config/wakie"
		}
		fullDbPath := saveDbPath + "/wakie.db"

		file := fileExists(fullDbPath)
		if file {
			log.Fatal("Error - Database already exists at the supplied path")
		}

		fileCreateResult, err := createFile(fullDbPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(" - Database %s\n", fileCreateResult)

		err = createDbTable(fullDbPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(" - Table created in database")

		// Writes out config file with path to database file once database has been created
		if cfgFile == "" {
			cfgFile = homeDir + "/.config/wakie/wakie.yaml"
		}
		viper.Set("db.dbPath", fullDbPath)
		err = viper.WriteConfigAs(cfgFile)
		if err != nil {
			log.Fatalf("Error writing config file: %s", err)
		}
		fmt.Printf(" - config file saved at %s\n", cfgFile)
		fmt.Println(" - Wakie is now setup and ready for use :-)")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVar(&saveDbPath, "saveDb", "$HOME/.config/wakie", "Path to folder where to save database file")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func createFile(filePath string) (string, error) {
	_, err := os.Create(filePath)
	if err != nil {
		errorMsg := "Failed to create file: " + err.Error()
		return "", errors.New(errorMsg)
	}
	return "File created at: " + filePath, nil
}

func createDbTable(dbPath string) error {
	createTableStmt := `CREATE TABLE 'computers' 
	('ID' INTEGER PRIMARY KEY AUTOINCREMENT, 
    'MAC_Address' STRING NULL UNIQUE,
	'Alias' STRING NULL UNIQUE);`

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		errorMsg := "Unable to open db file: " + err.Error()
		return errors.New(errorMsg)
	}

	_, err = db.Exec(createTableStmt)
	if err != nil {
		errorMsg := "Unable to create table in db file: " + err.Error()
		return errors.New(errorMsg)
	}

	return nil

}
