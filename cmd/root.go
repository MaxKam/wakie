package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3" // Driver for sql

	"github.com/linde12/gowol"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var dbPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wakie",
	Short: "Wake-on-LAN utility",
	Long:  `Utility for sending Magic Packets for Wake-on-LAN, as well as managing an address book of computer MAC Addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		var macAddress string
		var err error
		if cmd.Flags().Changed("alias") {
			aliasValue, err := cmd.Flags().GetString("alias")
			if err != nil {
				log.Fatalf("Error getting flag value")
			}
			macAddress = queryMAC("alias", aliasValue)
			sendPacket(macAddress)
		} else if cmd.Flags().Changed("id") {
			aliasValue, err := cmd.Flags().GetString("id")
			if err != nil {
				log.Fatalf("Error getting flag value")
			}
			macAddress = queryMAC("ID", aliasValue)
			sendPacket(macAddress)
		} else if cmd.Flags().Changed("mac") {
			macAddress, err = cmd.Flags().GetString("mac")
			cobra.CheckErr(err)
			sendPacket(macAddress)
		} else {
			fmt.Println("Welcome to Wakie. Please specify a MAC Address or --help flag for list of commands.")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	fmt.Println("Wakie:")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/wakie/wakie.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("alias", "a", "", "Alias of saved MAC address to send magic packet to")
	rootCmd.Flags().StringP("id", "i", "", "ID of saved MAC address to send magic packet to")
	rootCmd.Flags().StringP("mac", "m", "", "Manually entered MAC address to send magic packet to")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		configPath := home + "/.config/wakie"

		// Search config in home directory with name ".wakie" (without extension).
		viper.AddConfigPath(configPath)
		viper.AddConfigPath(".")
		viper.SetConfigName("wakie")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, " - Using config file:", viper.ConfigFileUsed())
	}

	dbPath = viper.GetString("db.dbPath")
}

func queryMAC(flagName, flagValue string) string {
	fmt.Printf(" - Getting MAC Address with %s of %s\n", flagName, flagValue)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database file. %s", err)
	}

	var idColumn string
	var macColumn string
	var aliasColumn string

	querySQLStmt := fmt.Sprintf("SELECT * FROM computers WHERE `%s` = '%s';", flagName, flagValue)

	queryResult := db.QueryRow(querySQLStmt)
	err = queryResult.Scan(&idColumn, &macColumn, &aliasColumn)
	if err != nil {
		log.Fatalf("Unable to find MAC with %s of %s", flagName, flagValue)
	}

	return macColumn
}

func sendPacket(macAddress string) {
	packet, err := gowol.NewMagicPacket(macAddress)
	if err != nil {
		log.Fatalln(err)
	}
	packet.Send("255.255.255.255")
	fmt.Printf(" - Magic Packet sent to %s\n", macAddress)
}
