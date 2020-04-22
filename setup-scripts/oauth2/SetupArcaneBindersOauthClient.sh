#!/bin/sh

echo "Creating Arcane Binders OAuth2.0 client..."

docker run --rm -it --network mtg-organizer_default \
		-e HYDRA_ADMIN_URL="http://hydra:4445" \
		oryd/hydra:latest \
		clients create --skip-tls-verify \
		--id ArcaneBinders \
		--grant-types implicit \
		--response-types token,id_token \
		--scope openid,test \
		--callbacks http://192.168.50.185:3000/auth_callback
