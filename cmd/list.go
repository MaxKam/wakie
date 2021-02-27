package cmd

import (
	"database/sql"
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
			id         int64
			macAddress string
			alias      string
		)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "MAC Address", "Alias"})

		listDB, err := db.Query("SELECT * FROM computers")
		if err != nil {
			log.Fatalf("Unable to query database: %s", err)
		}

		defer listDB.Close()
		for listDB.Next() {
			err := listDB.Scan(&id, &macAddress, &alias)
			if err != nil {
				log.Fatal(err)
			}
			t.AppendRow([]interface{}{id, macAddress, alias})
			t.AppendSeparator()
		}

		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
