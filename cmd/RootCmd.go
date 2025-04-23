package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	RootCmd = &cobra.Command{
		Use:   "griddb-cli",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $DIRECTORY/config/config.json)")
	RootCmd.PersistentFlags().StringP("author", "a", "israel imru", "author name for copyright attribution")
	RootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "israel imru <imru@fixstars.com>")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("./config")
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
