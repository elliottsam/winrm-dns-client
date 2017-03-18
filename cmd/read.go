// Copyright © 2017 Sam Elliott <me@sam-e.co.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/elliottsam/winrm-dns-client/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Reads a DNS record from the specified zone",
	Long: `Reads either single A or CNAME records or all A and CNAME records
from a Microsoft DNS Zone
`,
	Run: func(cmd *cobra.Command, args []string) {
		if dnsZone == "" {
			fmt.Println("DnsZone is a required parameter")
			os.Exit(1)
		}

		rec := dns.Record{
			Dnszone: dnsZone,
			Name:    name,
		}
		ClientConfig := dns.GenerateClient(viper.GetString("servername"), viper.GetString("username"), viper.GetString("password"))
		ClientConfig.ConfigureWinRMClient()

		resp, err := dns.ReadRecord(&ClientConfig, rec)
		if err != nil {
			log.Fatal("Error:", err)
		}
		dns.OutputTable(resp)
	},
}

func init() {

	RootCmd.AddCommand(readCmd)

	readCmd.PersistentFlags().StringVarP(&dnsZone, "DnsZone", "d", "", "DNS Zone to read against, this is required")
	readCmd.PersistentFlags().StringVarP(&name, "Name", "n", "", "Name of record to lookup")
	readCmd.MarkPersistentFlagRequired("DnsZone")

}
