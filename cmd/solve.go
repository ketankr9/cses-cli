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
	"bufio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// solveCmd represents the solve command
var solveCmd = &cobra.Command{
	Use:   "solve",
	Args: cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("show").Value.String() == "true" {
			ShowProblem(args[0])
		}

		SolveProblem(args[0], viper.GetString("lang"), checkAndGetEditor())
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)

	solveCmd.Flags().BoolP("show", "s", true, "print problem statement before solving")
	solveCmd.Flags().StringP("editor", "e", "", "use editor to solve the problem")
	err := viper.BindPFlag("editor", solveCmd.Flags().Lookup("editor"))

	cobra.CheckErr(err)
}

func checkAndGetEditor() string {
	editor := viper.GetString("editor")

	if editor == "" {
		scanner := bufio.NewScanner(os.Stdin)

		if !UpdateIfNew(scanner, &editor, "Editor") {
			panic("Editor still not configured")
		}
	}

	viper.Set("editor", editor)
	cobra.CheckErr(viper.WriteConfig())

	return editor
}
