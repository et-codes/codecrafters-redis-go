#!/bin/bash

# Configure GitHub identity
git config --global user.email "ericthornton43@gmail.com"
git config --global user.name "et-codes"

# Install Redis for testing
sudo apt install lsb-release curl gpg
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

sudo apt-get update
sudo apt-get install -y redis

# Install codecrafters CLI
curl https://codecrafters.io/install.sh | sh
