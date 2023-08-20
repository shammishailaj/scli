#!/bin/bash

# Install Dart SDK and CLI tools on Ubuntu
# https://github.com/dart-lang/sdk#install

# Update APT Cache
sudo apt-get update

# Install dependency https transport for apt
sudo apt-get install apt-transport-https

# Download the GPG key for the dart repo and store it in a keyring file
wget -qO- https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo gpg --dearmor -o /usr/share/keyrings/dart.gpg

# Add the dart repo to APT sources
echo 'deb [signed-by=/usr/share/keyrings/dart.gpg arch=amd64] https://storage.googleapis.com/download.dartlang.org/linux/debian stable main' | sudo tee /etc/apt/sources.list.d/dart_stable.list

# Update APT Cache
sudo apt-get update

# Install dart
sudo apt-get install dart
