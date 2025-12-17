#!/bin/bash

set -e

CERT_PATH="${CERT_PATH:-h3/examples/server.cert}"

SPKI=`openssl x509 -inform der -in "$CERT_PATH" -pubkey -noout \
  | openssl pkey -pubin -outform der \
  | openssl dgst -sha256 -binary \
  | openssl enc -base64`

echo "Got cert key $SPKI"

echo "Opening chromium"

case `uname` in
    (*Linux*)  chromium --origin-to-force-quic-on=127.0.0.1:4433 --ignore-certificate-errors-spki-list="$SPKI" --enable-logging --v=1 ;;
    (*Darwin*)  open -a "Google Chrome" --args --origin-to-force-quic-on=127.0.0.1:4433 --ignore-certificate-errors-spki-list="$SPKI" --enable-logging --v=1 ;;
esac

## Logs are stored to ~/Library/Application Support/Google/Chrome/chrome_debug.log