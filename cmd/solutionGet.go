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
	"github.com/shaunmitchellve/arete/internal/cmdsolution"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var branch, subFolder string

// solutionCmd represents the create command
var solutionGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Download the solution into your cache dir and add it to the list of known solutions",
	Example: ` arete solution get https://github.com/user/repo`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sl := cmdsolution.SolutionsList{}

		if err := sl.GetCoreSolutions(); err != nil {
			return err
		}

		if err := sl.GetSolution(args[0], branch, subFolder); err != nil {
			return err
		}

		log.Info().Msg("Solution Added")
		return nil
	},
}

// init the command and add flags
func init() {
	solutionCmd.AddCommand(solutionGetCmd)

	solutionGetCmd.Flags().StringVar(&branch, "branch", "main", "If the solution.yaml file is in a different branch from the default")

	solutionGetCmd.Flags().StringVar(&subFolder, "sub-folder", "/", "If the solution.yaml file is not in the root of the repo then provide the path here.")
}