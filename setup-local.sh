#!/bin/bash

export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN=root

vault server -dev -dev-root-token-id=root > /tmp/vault.log 2>&1 &
sleep 5

vault secrets enable pki
vault secrets tune -max-lease-ttl=87600h pki
vault write -field=certificate pki/root/generate/internal \
     common_name="example.com" \
     issuer_name="root-2023" \
     ttl=87600h > /tmp/root_2023_ca.crt

vault write pki/roles/test allow_any_name=true

vault write pki/config/urls \
     issuing_certificates="$VAULT_ADDR/v1/pki/ca" \
     crl_distribution_points="$VAULT_ADDR/v1/pki/crl"

