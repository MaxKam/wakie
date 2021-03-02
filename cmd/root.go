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
var idNum string
var macAddress string
var alias string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wakie",
	Short: "Wake-on-LAN utility",
	Long:  `Utility for sending Magic Packets for Wake-on-LAN, as well as managing an address book of computer MAC Addresses.`,
	Run: func(cmd *cobra.Command, args []string) {

		switch {
		case idNum != "":
			macAddress = queryMAC("ID", idNum)
		case alias != "":
			macAddress = queryMAC("Alias", alias)
		case macAddress != "":
			// Nothing to do here since MAC address is passed in directly. Still need a case for it
			// so that default case isn't executed, which is for nothing being passed in.
		default:
			fmt.Println("Welcome to Wakie. Please specify a MAC Address or --help flag for list of commands.")
			os.Exit(1)
		}

		sendPacket(macAddress)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.config/wakie/wakie.yaml)")
	rootCmd.PersistentFlags().StringVarP(&idNum, "id", "i", "", "ID of saved MAC address")
	rootCmd.PersistentFlags().StringVarP(&macAddress, "mac", "m", "", "Manually entered MAC address")
	rootCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "Alias of saved MAC address")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

// Looks up MAC address by ID or Alias and returns address as a string
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

// Sends the magic packet to the specified address
func sendPacket(targetAddress string) {
	packet, err := gowol.NewMagicPacket(targetAddress)
	if err != nil {
		log.Fatalln(err)
	}
	packet.Send("255.255.255.255")
	fmt.Printf(" - Magic Packet sent to %s\n", targetAddress)
}
