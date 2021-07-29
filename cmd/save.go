package cmd

import (
	"fmt"
	"log"
	"net"

	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Driver for sql
	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save 'MAC Address'",
	Short: "Saves a computer's MAC address, along with an alias and IP address for that MAC",
	Long: `Saves a computer's MAC address, giving it a unique ID, along with an alias and IP address for that MAC address.
Wakie can use alias, ID or IP address to lookup saved MAC Address.`,
	Run: func(cmd *cobra.Command, args []string) {
		// First validate that MAC address is in the correct format
		formattedMAC, err := net.ParseMAC(args[0])
		if err != nil {
			log.Fatalf("Error validating MAC address. Please double check entered MAC address. \n %s", err)
		}

		// Also validate that IP address is in the correct format
		formattedIP := net.ParseIP(ipAddress)
		if formattedIP == nil {
			log.Fatal("Error validating IP address. Please double check entered IP address.")
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Error opening database file. %s", err)
		}

		insertSQLStmt, err := db.Prepare("INSERT INTO 'main'.'computers'('MAC_Address', 'IP_Address', 'Alias') VALUES(?, ?, ?);")
		if err != nil {
			cobra.CheckErr(err)
		}

		insertEntry, err := insertSQLStmt.Exec(formattedMAC.String(), formattedIP.String(), alias)
		if err != nil {
			insertSQLStmt.Close()
			if err.Error() == "UNIQUE constraint failed: computers.Alias" {
				log.Fatalln("Error saving: Computer alias already exists in the database.")
			} else if err.Error() == "UNIQUE constraint failed: computers.MAC_Address" {
				log.Fatalln("Error saving: MAC Address already exists in the database.")
			} else {
				log.Fatalf("Error saving: %s", err)
			}

		}

		dbSaveID, err := insertEntry.LastInsertId()
		if err != nil {
			insertSQLStmt.Close()
			fmt.Println("Unable to get status of save to database operation")
		}

		fmt.Printf(" - MAC address saved to database with ID: %d\n", dbSaveID)
		insertSQLStmt.Close()

	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
