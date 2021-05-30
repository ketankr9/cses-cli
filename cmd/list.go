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
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Example: "cses-cli list, cses-cli list -r",
	Short: "list all problem statements along with status",
	Long: `Without any flag it lists all the problems from the cache.
To invalidate/refresh/reload the problem list cache use flag '-r'
If you are authenticated it will display the submission status along with the list.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("refresh").Value.String() == "true" {
			UpdateProblemListCache(viper.GetString("cookie"))
		}
		ListProblems(viper.GetString("cookie"))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	listCmd.Flags().BoolP("refresh", "r", false, "hard refresh problems list cache")
}
