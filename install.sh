#!/bin/bash

set -e

echo "=================================="
echo "      RoryFlix Installer"
echo "=================================="
echo

# Ensure apt exists
if ! command -v apt >/dev/null 2>&1; then
    echo "This installer currently supports Debian/Ubuntu/Linux Mint."
    exit 1
fi

echo "Updating package lists..."
sudo apt update

echo
echo "Installing dependencies..."
sudo apt install -y \
    golang \
    firefox \
    xdotool \
    git

echo
echo "Creating apikey.txt if it doesn't exist..."

if [ ! -f apikey.txt ]; then
    touch apikey.txt
    echo "Created apikey.txt"
else
    echo "apikey.txt already exists."
fi

echo
echo "=================================="
echo "Installation complete!"
echo "=================================="
echo
echo "Next steps:"
echo
echo "1. Obtain a free OMDb API key:"
echo "   https://www.omdbapi.com/apikey.aspx"
echo
echo "2. Paste ONLY your API key into:"
echo "   apikey.txt"
echo
echo "3. Build RoryFlix:"
echo "   go build server.go"
echo
echo "4. Run RoryFlix:"
echo "   ./server"
echo
echo "5. Open a browser to:"
echo "   http://localhost:8080"
echo
echo "To access RoryFlix from other devices,"
echo "you may need to allow TCP port 8080 through"
echo "your firewall."
