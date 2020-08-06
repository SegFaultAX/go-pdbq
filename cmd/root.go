package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "pdbq",
	Short: "PuppetDB Query Companion",
	Long:  `TODO`,
}

type (
	PDBQConfig struct {
		Clusters map[string]Cluster
	}

	Cluster struct {
		URL string
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	cfgFile, _ := rootCmd.PersistentFlags().GetString("config")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.SetConfigName(".pdbq")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func GetConfig() (*PDBQConfig, error) {
	var cfg PDBQConfig
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}
