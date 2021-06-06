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
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:     "show",
	Args:    cobra.ExactArgs(1),
	Short:   "show a problem statement by number",
	Example: "cses-cli show 1068",
	Long:    `Print the problem statement identified by the task number to stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		ShowProblem(args[0])
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
