#!/bin/bash
set -e

# httpd
echo "[*] Starting httpd"
httpd-foreground -c "LoadModule cgid_module modules/mod_cgid.so"
