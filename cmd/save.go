package cmd

import (
	"fmt"

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
		db, err := sql.Open("sqlite3", "/home/max/.config/wakie/wakie.db")
		if err != nil {
			fmt.Printf("Error opening database file. %s", err)
		}

		flagValue, err := cmd.Flags().GetString("alias")

		sqlStatement, err := db.Prepare("INSERT INTO 'main'.'computers'('MAC_Address', 'Alias') VALUES(?, ?);")

		insertEntry, err := sqlStatement.Exec(args[0], flagValue)
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
