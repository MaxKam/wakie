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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved computers",
	Long:  "List all saved computers in database. Will print out ID, Mac Address, and Alias.",
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
		)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "MAC Address", "Alias"})

		switch {
		case idNum != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `ID` = '%s';", idNum)
		case alias != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `Alias` = '%s';", alias)
		case macAddress != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `MAC_Address` = '%s';", macAddress)
		default:
			dbQueryString = "SELECT * FROM computers"
		}

		listDB, err := db.Query(dbQueryString)
		if err != nil {
			log.Fatalf("Unable to query database: %s", err)
		}

		defer listDB.Close()
		for listDB.Next() {
			err := listDB.Scan(&idFromDB, &macAddressFromDB, &aliasFromDB)
			if err != nil {
				log.Fatal(err)
			}
			t.AppendRow([]interface{}{idFromDB, macAddressFromDB, aliasFromDB})
			t.AppendSeparator()
		}

		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
