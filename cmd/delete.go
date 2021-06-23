package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3" // Driver for sql

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete saved MAC Address from database",
	Long: `Delete saved MAC Address from database.
	Can select record to delete by ID, MAC Address, or Alias.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Error opening database file. %s", err)
		}

		var (
			idFromDB         int64
			macAddressFromDB string
			ipAddressFromDB  string
			aliasFromDB      string
			userResponse     string
			dbQueryString    string
			dbDeleteString   string
		)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "MAC Address", "IP Address", "Alias"})

		switch {
		case idNum != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `ID` = '%s';", idNum)
			dbDeleteString = fmt.Sprintf("DELETE FROM computers WHERE ID = '%s'", idNum)
		case alias != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `Alias` = '%s';", alias)
			dbDeleteString = fmt.Sprintf("DELETE FROM computers WHERE Alias = '%s'", alias)
		case macAddress != "":
			dbQueryString = fmt.Sprintf("SELECT * FROM computers WHERE `MAC_Address` = '%s';", macAddress)
			dbDeleteString = fmt.Sprintf("DELETE FROM computers WHERE MAC_Address = '%s'", macAddress)
		default:
			fmt.Println("Please specify an ID(-i), MAC Address(-m), or Alias(-a). To see list of all records in database, run 'wakie list'")
			os.Exit(1)
		}

		// Finds entry in database by user inputed flag, prints entry on screen, and asks user to confirm deletion

		listDB := db.QueryRow(dbQueryString)
		if err != nil {
			log.Fatalf("Unable to query database: %s", err)
		}
		err = listDB.Scan(&idFromDB, &macAddressFromDB, &ipAddressFromDB, &aliasFromDB)
		if err != nil {
			log.Fatal(err)
		}
		t.AppendRow([]interface{}{idFromDB, macAddressFromDB, ipAddressFromDB, aliasFromDB})
		t.AppendSeparator()
		fmt.Println("Entry to be deleted:")
		t.Render()
		fmt.Print(" *** Are you sure you want to delete this entry? [Yes|No] ")
		fmt.Scanf("%s", &userResponse)

		switch userResponse {
		case "Yes", "Y", "yes", "y":
			dbResult, err := db.Exec(dbDeleteString)
			if err != nil {
				log.Fatalf("Error deleting entry from database: %s", err)
			}
			rowsAffected, _ := dbResult.RowsAffected()
			rowsAffectedString := strconv.FormatInt(rowsAffected, 10)
			fmt.Printf("Deleted rows: %s\n", rowsAffectedString)
		default:
			fmt.Println("Canceling delete.")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
