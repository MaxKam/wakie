package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3" // Driver for sql

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing record in the database",
	Long: `Update an existing record in the database.
	Can search for record to update by ID, MAC Address, or Alias.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Error opening database file. %s", err)
		}

		var (
			idFromDB         int64
			macAddressFromDB string
			aliasFromDB      string
			dbQueryString    string
			dbUpdateString   string
		)

		switch {
		case idNum != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `ID` = '%s';", idNum)
		case alias != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `Alias` = '%s';", alias)
		case macAddress != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `MAC_Address` = '%s';", macAddress)
		default:
			fmt.Println("Please specify an ID(-i), MAC Address(-m), or Alias(-a). To see list of all records in database, run 'wakie list'")
			os.Exit(1)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "MAC Address", "Alias"})
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
