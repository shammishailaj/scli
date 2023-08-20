#!/bin/bash

# Script to install flutter on Ubuntu
wget https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.10.6-stable.tar.xz

# Create a new directory
mkdir flutter

# Move the downloaded file to the new directory
mv flutter_linux_3.10.6-stable.tar.xz flutter/

# Change to the new directory
cd flutter

# Extract the downloaded file
tar -xvf flutter_linux_3.10.6-stable.tar.xz

# Add the flutter directory to the path
export PATH="$PATH:`pwd`/flutter/bin"

# Pre-download development binaries
flutter precache

# Run flutter doctor
flutter doctor
