// Copyright 2023 Shaun Mitchell

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// verbose states if the -v or --verbose flag as been set for debug purposes
// current version, is set a build time
var (
	verbose bool
	Version string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arete",
	Short: "Arete makes deploying GCP infrastructure easier",
	Long: `Arete is a wrapper that makes deploying solutions onto Google Cloud Platform easier.

It utilizes Googles Config Connector and Config Controller to deploy declaritive resources into your environment
with as little changes as required.`,
}

// rootCmd represents the base command when called without any subcommands
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print out the current version of arete",
	Run: func(cmd *cobra.Command, args[]string) {
		fmt.Printf("%s\n", Version)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}


// Init the CLI and add global flags
func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(versionCmd)
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}


