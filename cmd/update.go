package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3" // Driver for sql

	"github.com/spf13/cobra"
)

var (
	idFromDB         int64
	macAddressFromDB string
	ipAddressFromDB  string
	aliasFromDB      string
	dbQueryString    string
	newAliasValue    string
	newMacAddrValue  string
	newIPAddrValue   string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing record in the database",
	Long: `Update an existing record in the database.
	Can search for record to update by ID, MAC Address, or Alias but not by IP Address since Wakie allows for duplicate IP's. 
	Can update Alias, IP Address or MAC Address, or all at the same time.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Error opening database file. %s", err)
		}

		// Create sql query string based on user inputed flag to get record.
		// All update requests will specify record by ID.
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

		// Create table that will be printed out, showing old and updated record
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"ID", "MAC Address", "IP Address", "Alias"})

		// Query db for old record
		listDB := db.QueryRow(dbQueryString)
		if err != nil {
			log.Fatalf("Unable to query database: %s", err)
		}
		err = listDB.Scan(&idFromDB, &macAddressFromDB, &ipAddressFromDB, &aliasFromDB)
		if err != nil {
			log.Fatal(err)
		}
		t.AppendRow(table.Row{"Old record:"})
		t.AppendRow([]interface{}{idFromDB, macAddressFromDB, ipAddressFromDB, aliasFromDB})
		t.AppendSeparator()

		// Update Alias if updateAlias flag is set
		if newAliasValue != "" {
			updateStmt, err := db.Prepare("UPDATE computers SET Alias=? where ID=?")
			cobra.CheckErr(err)

			_, err = updateStmt.Exec(newAliasValue, idFromDB)
			cobra.CheckErr(err)

			updateStmt.Close()

		}

		// Update MAC address if updateMac flag is set
		if newMacAddrValue != "" {
			// First confirm valid mac address
			newMacAddrValue, err := net.ParseMAC(newMacAddrValue)
			if err != nil {
				log.Fatalf("Updated MAC Address is not valid. Please double check the entered address. %s", err)
			}

			updateStmt, err := db.Prepare("UPDATE computers SET MAC_Address=? where ID=?")
			cobra.CheckErr(err)

			_, err = updateStmt.Exec(newMacAddrValue.String(), idFromDB)
			cobra.CheckErr(err)

			updateStmt.Close()
		}

		// Update IP address if updateIP flag is set
		if newIPAddrValue != "" {
			// First confirm valid mac address
			newIPAddrValue := net.ParseIP(newIPAddrValue)
			if newIPAddrValue == nil {
				log.Fatalf("Updated IP Address is not valid. Please double check the entered address. %s", err)
			}

			updateStmt, err := db.Prepare("UPDATE computers SET IP_Address=? where ID=?")
			cobra.CheckErr(err)

			_, err = updateStmt.Exec(newIPAddrValue.String(), idFromDB)
			cobra.CheckErr(err)

			updateStmt.Close()
		}

		// Query DB again for updated record and print out table
		listDB = db.QueryRow(fmt.Sprintf("SELECT * FROM computers WHERE `ID` = '%d';", idFromDB))
		if err != nil {
			log.Fatalf("Unable to query database: %s", err)
		}
		err = listDB.Scan(&idFromDB, &macAddressFromDB, &ipAddressFromDB, &aliasFromDB)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(" - Record has been updated:")
		t.AppendRow(table.Row{"Updated record:"})
		t.AppendRow([]interface{}{idFromDB, macAddressFromDB, ipAddressFromDB, aliasFromDB})
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&newAliasValue, "updateAlias", "", "Specify updated Alias")
	updateCmd.Flags().StringVar(&newMacAddrValue, "updateMac", "", "Specify updated MAC Address")
	updateCmd.Flags().StringVar(&newIPAddrValue, "updateIP", "", "Specify updated IP Address")
}
