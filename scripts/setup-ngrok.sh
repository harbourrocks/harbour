#!/bin/sh

sudo apt install -y unzip wget

DOWNLOAD_FILE=/tmp/ngrok-stable-linux-amd64.zip

DOWNLOAD_PATH=https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip

wget $DOWNLOAD_PATH -P /tmp

unzip $DOWNLOAD_FILE -d ./scripts

rm $DOWNLOAD_FILE

chmod +x ./scripts/ngrok
