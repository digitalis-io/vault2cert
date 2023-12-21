package certs

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/pavlo-v-chernykh/keystore-go/v4"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func IssueCert(ctx context.Context, role string, commonName string, mountPoint string) (*vault.Response[schema.PkiIssueWithRoleResponse], error) {
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress(GetEnv("VAULT_ADDR", "http://127.0.0.1:8200")),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// authenticate with a root token (insecure)
	if err := client.SetToken(GetEnv("VAULT_TOKEN", "root")); err != nil {
		log.Fatal(err)
	}

	req := schema.PkiIssueWithRoleRequest{
		CommonName: commonName,
	}

	headers := http.Header{}
	if GetEnv("CLOUDFLARE_TOKEN", "") != "" {
		headers.Set("cf-access-token", GetEnv("CLOUDFLARE_TOKEN", ""))
	}

	certs, err := client.Secrets.PkiIssueWithRole(ctx, role, req,
		vault.WithCustomHeaders(headers),
		vault.WithMountPath(mountPoint),
	)

	return certs, err
}

func writeKeyStore(ks keystore.KeyStore, filename string, password []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = ks.Store(f, password)
	if err != nil {
		return err
	}
	return nil
}

func WriteToJks(cert string, key string, ca string, path string, pass string) error {
	password := []byte(pass)
	defer zeroing(password)

	ks1 := keystore.New()

	pkeIn := keystore.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:   []byte(key),
		CertificateChain: []keystore.Certificate{
			{
				Type:    "X509",
				Content: []byte(cert),
			},
		},
	}

	if err := ks1.SetPrivateKeyEntry("alias", pkeIn, password); err != nil {
		return err
	}

	err := writeKeyStore(ks1, path, password)
	if err != nil {
		return err
	}
	return nil
}

func zeroing(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}
