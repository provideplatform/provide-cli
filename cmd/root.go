package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

var networkID string
var applicationID string

const (
	// Note: Viper downcases key names, so hyphenating for better readability.
	authTokenConfigKey = "auth-token"
	apiTokenConfigKey  = "api-token"
)

var rootCmd = &cobra.Command{
	Use:   "prvd",
	Short: "Provide command-line interface",
	Long: `The Provide CLI exposes network and dapp developer tools for
seamless interaction with the microservices which power the 
provide.services PaaS.

Run with the --help flag to see available options`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.provide-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".provide-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".provide-cli")

		configPath := fmt.Sprintf("%s/.provide-cli.yaml", home)
		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			if os.IsNotExist(err) {
				err = viper.WriteConfigAs(configPath)
			}
		}

		viper.RegisterAlias(authTokenConfigKey, "token")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		os.Chmod(viper.ConfigFileUsed(), 0600)

		if verbose {
			fmt.Println("Using configuration:", viper.ConfigFileUsed())
		}
	}
}

func requireAPIToken() string {
	token := viper.Get(authTokenConfigKey)
	if token == nil {
		log.Printf("Authorized API token required in prvd configuration; have you authenticated or otherwise configured an API token?")
		os.Exit(1)
	}
	return token.(string)
}
