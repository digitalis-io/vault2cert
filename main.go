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
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/digitalis-io/vault2cert/cmd"
	"github.com/digitalis-io/vault2cert/pkg/certs"
)

func writeToFile(filename string, contents string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(contents + "\n")
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func main() {
	cmd.Execute()
	ctx := context.Background()

	issuedCerts, err := certs.IssueCert(ctx, cmd.Role, cmd.CommonName, cmd.Mount)
	if err != nil {
		panic(err)
	}
	if cmd.SavePath == "" {
		fmt.Println(issuedCerts.Data.Certificate)
	} else {
		certFile := filepath.Join(cmd.SavePath, cmd.CommonName+".crt")
		keyFile := filepath.Join(cmd.SavePath, cmd.CommonName+".key")
		caFile := filepath.Join(cmd.SavePath, cmd.CommonName+".ca")
		writeToFile(certFile, issuedCerts.Data.Certificate)
		writeToFile(keyFile, issuedCerts.Data.PrivateKey)
		writeToFile(caFile, issuedCerts.Data.IssuingCa)

	}

	if cmd.JksPath != "" {
		certs.WriteToJks(issuedCerts.Data.Certificate, issuedCerts.Data.PrivateKey, issuedCerts.Data.IssuingCa, cmd.JksPath, cmd.JksPassword)
	}
}
