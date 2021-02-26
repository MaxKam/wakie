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
	Short: "Saves a computer's MAC address, along with an alias for that address",
	Long: `Saves a computer's MAC address, giving it a unique ID, along with an alias for that address.
Wakie can use either alias or ID to lookup saved MAC Address`,
	Run: func(cmd *cobra.Command, args []string) {
		// First validate that MAC address is in the correct format
		formattedMAC, err := net.ParseMAC(args[0])
		if err != nil {
			log.Fatalf("Error validating MAC address. Please double check the address entered. \n %s", err)
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("Error opening database file. %s", err)
		}

		aliasFlagValue, err := cmd.Flags().GetString("alias")

		insertSQLStmt, err := db.Prepare("INSERT INTO 'main'.'computers'('MAC_Address', 'Alias') VALUES(?, ?);")

		insertEntry, err := insertSQLStmt.Exec(formattedMAC.String(), aliasFlagValue)
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

		fmt.Println(fmt.Sprintf("MAC address saved to database with ID: %d", dbSaveID))
		insertSQLStmt.Close()

	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringP("alias", "a", "", "Alias for computer (required)")
	saveCmd.MarkFlagRequired("alias")
}
