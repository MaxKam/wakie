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

		flagValue, err := cmd.Flags().GetString("alias")

		sqlStatement, err := db.Prepare("INSERT INTO 'main'.'computers'('MAC_Address', 'Alias') VALUES(?, ?);")

		insertEntry, err := sqlStatement.Exec(formattedMAC, flagValue)
		if err != nil {
			fmt.Printf("Error inserting data: %s", err)
		}

		fmt.Println(insertEntry)
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	saveCmd.Flags().StringP("alias", "a", "", "Alias for computer (required)")
	saveCmd.MarkFlagRequired("alias")
}
