# RoryFlix

> **Turn an old laptop into your own personal media server.**

RoryFlix is a lightweight self-hosted media server written in Go that was designed to give an old laptop or desktop a second life. It allows you to stream your own videos to any device on your home network while also providing a modern streaming frontend capable of launching content from supported online sources.

Whether you're repurposing an old Linux Mint laptop connected to your TV or simply want a simple media server for your home, RoryFlix aims to be lightweight, easy to set up, and fun to tinker with.

---

# Features

## Local Media Library

* Stream your own video collection over your local network
* Automatically serves video files placed in the `videos/` folder
* Watch from phones, tablets, laptops, desktops, or smart TVs using a web browser

---

## Streaming Frontend

RoryFlix includes a Netflix-style frontend that can search for movies and television shows using the OMDb API.

Supported streaming providers include:

* MoviesAPI
* VSembed (formerly VidSrc)

The frontend supports:

* Movie search
* TV episode search
* Watch history
* Favorites
* Playlists

---

## Remote Control

One of RoryFlix's biggest features is its built-in remote.

Using another device on your network, you can:

* Search movies
* Search TV shows
* Launch media
* Pause/Play
* Toggle fullscreen
* Hide the mouse cursor
* Refresh the viewer
* Open or close the Viewer window
* Adjust volume
* Move the mouse
* Click
* Adjust cursor step size

This allows an old laptop connected to your television to behave similarly to a Roku or other streaming device.

---

## Viewer

The Viewer is a fullscreen Firefox window running on the server itself.

The remote communicates with the server and tells Firefox what to display, making RoryFlix ideal for:

* Living room TVs
* Bedroom TVs
* Projectors
* Media PCs

---

# Requirements

RoryFlix was built on **Linux Mint** and is primarily intended to run there.

It should also work on:

* Ubuntu
* Most Debian-based Linux distributions

Windows and macOS are not officially supported at this time because the remote currently depends on `xdotool` for mouse and keyboard control.

---

# Dependencies

RoryFlix requires:

* Go (Golang)
* Firefox
* xdotool
* An OMDb API key

---

# Installation

## Clone the repository

```bash
git clone https://github.com/YOUR_USERNAME/roryflix.git

cd roryflix
```

---

## Install dependencies

Make the installer executable:

```bash
chmod +x install.sh
```

Run it:

```bash
./install.sh
```

The installer will install:

* Go
* Firefox
* xdotool
* Create videos folder

and create an empty

```text
apikey.txt
```

file.

---

# Getting an OMDb API Key

RoryFlix uses the OMDb API to convert movie and TV show names into IMDb IDs.

You can obtain a free API key here:

https://www.omdbapi.com/apikey.aspx

Once you receive your key, place **only the API key** inside

```text
apikey.txt
```

Example:

```text
12345678
```

---

# Building RoryFlix

Compile the server:

```bash
go build server.go
```

---

# Running RoryFlix

Start the server:

```bash
./server
```

The server will print something similar to:

```text
Service running on 192.168.1.25:8080
```

---

# Opening RoryFlix

On the server:

```text
http://localhost:8080
```

From another device:

```text
http://SERVER-IP:8080
```

Example:

```text
http://192.168.1.25:8080
```

---

# Firewall

It is recommended you use the default gui firewall tool that comes with linux mint.
Press windows or super key and Search "firewall"
In the linux mint firewall tool set incoming to allow.

If other devices still cannot access RoryFlix, your firewall may still be blocking connections.

If so try this:

On Linux systems using UFW:

```bash
sudo ufw allow 8080/tcp
```

Then verify:

```bash
sudo ufw status
```

---

# Using RoryFlix

## Local Videos

Place your videos into:

```text
videos/
```

Supported formats include:

* mp4
* webm
* ogg

They will automatically appear under the **Videos** page.

---

## Streaming

Use the search page to:

* Search Movies
* Search TV Shows

RoryFlix uses the OMDb API to locate the correct IMDb ID before opening the selected content using your chosen provider.

---

## Remote

Open the Remote page from another device.

The Remote allows you to:

* Search for content
* Launch content
* Pause playback
* Toggle fullscreen
* Move the mouse
* Click
* Control volume
* Hide the mouse
* Refresh Firefox
* Open or close the Viewer

---

# install.sh

The included installer installs all required dependencies automatically.

```bash
#!/bin/bash

set -e

echo "=================================="
echo "      RoryFlix Installer"
echo "=================================="

sudo apt update

sudo apt install -y \
    golang \
    firefox \
    xdotool \
    git

if [ ! -f apikey.txt ]; then
    touch apikey.txt
fi

echo
echo "Installation complete!"
echo
echo "Get an OMDb API key:"
echo "https://www.omdbapi.com/apikey.aspx"
echo
echo "Paste your key into:"
echo "apikey.txt"
echo
echo "Build:"
echo "go build server.go"
echo
echo "Run:"
echo "./server"
```

---

# Project Goals

RoryFlix was created with a few simple goals:

* Give old laptops a second life
* Provide a lightweight self-hosted media server
* Make streaming simple
* Keep everything open source
* Be easy to modify and extend
* Have a fun Netflix-inspired interface

---

# Platform Support

| Platform   | Status            |
| ---------- | ----------------- |
| Linux Mint | ✅ Fully Supported |
| Ubuntu     | ✅ Should Work     |
| Debian     | ✅ Likely Works    |
| Windows    | ⚠ Untested        |
| macOS      | ⚠ Untested        |

The Remote currently depends on **xdotool**, which is the main reason Linux is recommended. There is a windows version of xdotool. Intheory everything could work in windows. I often wondered if xdotool/Roryflix would work in WSL (windows subsystem for linux) feel free to port it over yourself and make a windows fork.

---

# Disclaimer

RoryFlix does **not** host, provide, or distribute copyrighted content.

The local media server only serves files that you place inside your own `videos` directory.

Streaming functionality simply opens content from external providers selected by the user. Users are responsible for ensuring their use complies with the laws and terms of service applicable in their jurisdiction.

---

# License

Software licensing is not punk rock!

---

Enjoy giving your old laptop a second life fighting the system.
