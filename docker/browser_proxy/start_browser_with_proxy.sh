#!/bin/bash

echo "Starting nginx reverse proxy for Chrome debugging..."
nginx -c /etc/nginx/chrome-debug-proxy.conf -g "daemon on;"

echo "Starting Kasm browser enviroment..."
exec /dockerstartup/vnc_startup.sh
