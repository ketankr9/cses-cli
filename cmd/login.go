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
	Example: "cses-cli login",
	Short: "authenticate to cses.fi server by username and password",
	Long: `This command retrieves the csrf token and signs it to the server,
reads the cookies and saves it in the config file along with the username
Note: password is not saved in the config file and is only used to sign the csrf token
`,
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
	ok := Login(csrf, cookie, username, PASSWORD)
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
