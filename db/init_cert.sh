#!/bin/bash

# postgres user has UID and GID 999, while user from host has 1000:1000
chown 999:999 /home/ca/server.crt
chown 999:999 /home/ca/server.key