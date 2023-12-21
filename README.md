# Vault2Cert

It requests a SSL certificate from a HashiCorp vault server and stores it to either

## Usage

```bash
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN=root

go run main.go --mount pki --common-name hello.example.com --role=test --write-to /tmp
```

This will write the PEM files to /tmp

```bash
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN=root

go run main.go --mount pki --common-name hello.example.com --role=test --jks /tmp/hello.jks --jkspassword=changeme
```

Same as before but store the keys in JKS

