#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/upasspdf"
mkdir -p "$CONFIG_DIR/data" "$CONFIG_DIR/logs"

if [ ! -f "$CONFIG_DIR/settings.json" ] && [ -f "config/example/settings.example.json" ]; then
  cp "config/example/settings.example.json" "$CONFIG_DIR/settings.json"
fi

if [ ! -f "$CONFIG_DIR/.env" ] && [ -f "config/example/.env.example" ]; then
  cp "config/example/.env.example" "$CONFIG_DIR/.env"
  chmod 600 "$CONFIG_DIR/.env"
fi

mkdir -p tmp
ln -sfn "$CONFIG_DIR" tmp/config
