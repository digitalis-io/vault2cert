/*
Copyright © 2023 Sergio Rua <sergio.rua@digitalis.io>

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
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/digitalis-io/vault2cert/pkg/certs"
	"github.com/spf13/cobra"
)

var Mount string
var CommonName string
var Role string
var JksPath string
var JksPassword string
var SavePath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github.com/digitalis-io/vault2cert",
	Short: "Utility to retrieve a SSL cert from HashiCorp Vault",
	Long:  `Use this tool to request a cert from Vault and store it to in either in PEM or JKS format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if CommonName == "" {
			return fmt.Errorf("common name is mandatory")
		}
		if Role == "" {
			return fmt.Errorf("role is mandatory")
		}
		if JksPassword == "" && JksPath != "" {
			JksPassword = generatePassword(12)
			fmt.Printf("JKS password: %s\n", JksPassword)
		}
		return nil
	},
}

func generatePassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password string
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return ""
		}
		password += string(chars[num.Int64()])
	}
	return password
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&Mount, "mount", certs.GetEnv("VAULT_PKI_MOUNT", "pki"), "Path for the PKI mount, default pki")
	rootCmd.Flags().StringVar(&CommonName, "common-name", "", "SSL cert common name, ie, something.example.com")
	rootCmd.Flags().StringVar(&Role, "role", certs.GetEnv("VAULT_PKI_ROLE", ""), "Vault role to call for issuying the cert")
	rootCmd.Flags().StringVar(&JksPath, "jks", "", "Write keys to JKS keystore")
	rootCmd.Flags().StringVar(&JksPassword, "jkspassword", "", "Password for the JKS store")
	rootCmd.Flags().StringVar(&SavePath, "write-to", "", "Path where to write out the certificate and key")
}
