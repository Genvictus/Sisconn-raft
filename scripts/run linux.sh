#!/bin/bash

gnome-terminal -- bash -c "go run ./cmd/server/ --port 6964; exec bash"
gnome-terminal -- bash -c "go run ./cmd/server/ --port 6965; exec bash"
gnome-terminal -- bash -c "go run ./cmd/server/ --port 6966; exec bash"
gnome-terminal -- bash -c "go run ./cmd/server/ --port 6967; exec bash"
gnome-terminal -- bash -c "go run ./cmd/server/ --port 6968; exec bash"
gnome-terminal -- bash -c "go run ./cmd/server/ --port 6969; exec bash"
