package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	RootCmd = &cobra.Command{
		Use:   "griddb-cloud-cli",
		Short: "A wrapper to making HTTP Requests to your GridDB Cloud Instance",
		Long: `A series of commands to help you manage your cloud-based DB.
Standouts include creating a container and graphing one using 'read graph' and 'create' respectfully`,
	}

	tokenManager *TokenManager
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	tokenManager = NewTokenManager()

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/griddb-cloud-cli/config.yaml)")
	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "israel imru <imru@fixstars.com>")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// in Linux, it's $HOME/.config
		// For Macos, it's ~/Library/Application\ Support
		configDir, err := os.UserConfigDir()
		cobra.CheckErr(err)
		appConfigDir := filepath.Join(configDir, "griddb-cloud-cli")
		err = os.MkdirAll(appConfigDir, 0755)
		cobra.CheckErr(err)

		viper.AddConfigPath(appConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config") // Looks for config.yaml
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		configDir, _ := os.UserConfigDir()
		fmt.Fprintln(os.Stderr, "Please enter credentials inside of `config.yaml` file in this directory: ", configDir+"/griddb-cloud-cli")
	}
}
