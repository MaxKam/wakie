package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3" // Driver for sql

	"github.com/linde12/gowol"
	"github.com/spf13/viper"
)

var cfgFile string
var dbPath string
var idNum string
var macAddress string
var ipAddress string
var alias string
var homeDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "1.6.5",
	Use:     "wakie",
	Short:   "Wake-on-LAN utility",
	Long:    `Utility for sending Magic Packets for Wake-on-LAN, as well as managing an address book of computer MAC Addresses.`,
	Run: func(cmd *cobra.Command, args []string) {

		switch {
		case idNum != "":
			macAddress, ipAddress = queryMAC("ID", idNum)
		case alias != "":
			macAddress, ipAddress = queryMAC("Alias", alias)
		case macAddress != "":
			// Nothing to do here since MAC address is passed in directly. Still need a case for it
			// so that default case isn't executed, which is for no flags being passed in.
		default:
			fmt.Println(color.Ize(color.Bold, "Welcome to Wakie. Please specify a MAC Address or --help flag for list of commands."))
			os.Exit(1)
		}

		sendPacket(macAddress, ipAddress)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	fmt.Println(color.Ize(color.Bold, "Wakie:"))
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.config/wakie/wakie.yaml or app folder)")
	rootCmd.PersistentFlags().StringVarP(&idNum, "id", "i", "", "ID of saved MAC address")
	rootCmd.PersistentFlags().StringVarP(&macAddress, "mac", "m", "", "Manually entered MAC address")
	rootCmd.PersistentFlags().StringVar(&ipAddress, "ip", "255.255.255.255", "IP address of MAC address. Needed in case host is connected to multiple networks.")
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
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			log.Fatalf("Unable to get users home dir: %s", err)
		}
		configPath := homeDir + "/.config/wakie"

		// Search config in home directory with name ".wakie" (without extension).
		viper.AddConfigPath(configPath)
		viper.AddConfigPath(".")
		viper.SetConfigName("wakie")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, " - Using config file:", viper.ConfigFileUsed())
	}

	dbPath = viper.GetString("db.dbPath")
}

// queryMAC looks up MAC address by ID or Alias and returns MAC and IP address as a string
func queryMAC(flagName, flagValue string) (string, string) {
	fmt.Printf(" - Getting MAC Address with %s of %s\n", flagName, flagValue)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database file. %s", err)
	}

	var idColumn string
	var macColumn string
	var ipColumn string
	var aliasColumn string

	querySQLStmt := fmt.Sprintf("SELECT * FROM computers WHERE `%s` = '%s';", flagName, flagValue)

	queryResult := db.QueryRow(querySQLStmt)
	err = queryResult.Scan(&idColumn, &macColumn, &ipColumn, &aliasColumn)
	if err != nil {
		log.Fatalf("Unable to find MAC with %s of %s", flagName, flagValue)
	}

	return macColumn, ipColumn
}

// sendPacket sends the magic packet to the specified address
func sendPacket(targetMacAddress, targetIPAddress string) {
	packet, err := gowol.NewMagicPacket(targetMacAddress)
	if err != nil {
		log.Fatalln(err)
	}

	err = packet.Send(targetIPAddress)
	if err != nil {
		log.Fatal("Error sending magic packet")
	}

	fmt.Printf(" - Magic Packet sent to %s\n", targetMacAddress)
}
