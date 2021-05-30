/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !doLogin() {
			cobra.CompErrorln("failed to login")
		}

		UpdateProblemListCache(viper.GetString("cookie"))
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func doLogin() bool {

	var wg sync.WaitGroup
	var cookie string
	var csrf string

	wg.Add(1)
	go func() {
		cookie, csrf = NewCookieCsrf()
		wg.Done()
	}()

	username := viper.GetString("username")
	scanner := bufio.NewScanner(os.Stdin)
	UpdateIfNew(scanner, &username, "username")

	fmt.Print("password: ")
	scanner.Scan()
	PASSWORD := scanner.Text()

	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Prefix = "Logging "

	s.Start()
	wg.Wait()
	ok := login(csrf, cookie, username, PASSWORD)
	s.Stop()

	if !ok {
		return false
	}

	viper.Set("username", username)
	viper.Set("csrf", csrf)
	viper.Set("cookie", cookie)

	cobra.CheckErr(viper.WriteConfig())

	return true
}

func login(csrf string, cookie string, username string, pass string) bool {

	params := "csrf_token=" + csrf + "&nick=" + username + "&pass=" + pass
	if LoginRequest(params, cookie) {
		return true
	}
	return false
}
