/*
Copyright © 2024 https://github.com/gatariee

*/
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "winton-client",
    Short: "A CLI client for Winton C2",
    Long:  `winton-client is a CLI client for Winton, if run with no flags- anonymous (default) will be used, you will likely not be able to authenticate to your teamserver.
	`,
    Run: func(cmd *cobra.Command, args []string) {
		s := NewSession()
		s.Start()
    },
}

func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        viper.AddConfigPath("$HOME")
        viper.SetConfigName("config")
    }

    viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("[!] No config file specified, using defaults.")
		viper.SetDefault("operator", "anonymous")
		viper.SetDefault("teamserver.ip", "127.0.0.1")
		viper.SetDefault("teamserver.port", 80)
	}
}
