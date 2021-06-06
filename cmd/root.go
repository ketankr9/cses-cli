/*
Copyright © 2021 UTSAV KRISHNAN github.com/ketankr9

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	cfgFile  string
	rootPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cses-cli",
	Short: "command line to solve probmens from https://cses.fi/problemset/",
	Long:  `command line to solve probmens from https://cses.fi/problemset/`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	rootPath = filepath.Join(home, ".cses")

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+rootPath+"/config.*)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// disabling it for now
	//viper.AutomaticEnv() // read in environment variables first to be overwritten by config

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in following directories with name "cses-cli" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.cses/")
		viper.AddConfigPath("/etc/cses/")

		viper.SetConfigName("config") // name of config file (without extension)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debug("Using config file:", viper.ConfigFileUsed())
	}
}
