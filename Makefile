generate_access_private_key:
	 openssl genrsa -out access-private.pem 2048

generate_refresh_private_key:
	openssl genrsa -out refresh-private.pem 2048

export_rsa_access_public_key:
	openssl rsa -in access-private.pem -outform PEM -pubout -out access-public.pem

export_rsa_refresh_public_key:
	 openssl rsa -in refresh-private.pem -outform PEM -pubout -out refresh-public.pem