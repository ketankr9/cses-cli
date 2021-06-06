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
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:     "submit",
	Args:    cobra.ExactArgs(1),
	Example: "cses-cli submit 1068.Weird-Algorithm.cpp",
	Short:   "submit the solution file and fetch the status",
	Long: `If the solution is ACCEPTED and github is configured, then it will push the solution to github,
with limitation that the repo has to be public, since it uses github generated token`,

	Run: func(cmd *cobra.Command, args []string) {
		ghConfig := &GithubConfig{}

		err := viper.UnmarshalKey("github", ghConfig)
		cobra.CheckErr(err)

		SubmitSolution(args[0], viper.GetString("cookie"), viper.GetString("csrf"), ghConfig)
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
