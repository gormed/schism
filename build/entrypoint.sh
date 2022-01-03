#!/bin/bash

echo "Waiting for influxdb"
while [ -z "$(nc -vz "$INFLUXDB_HOST" "$INFLUXDB_PORT" &> /dev/null ; echo $?)" ] ; do echo '.'; sleep 1; done

exec "$@"