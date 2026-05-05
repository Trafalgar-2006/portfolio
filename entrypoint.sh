#!/bin/sh
# entrypoint.sh — handles SSH host key persistence
set -e

mkdir -p .ssh

if [ -n "$SSH_HOST_KEY" ]; then
    # Decode the base64-encoded private key stored in Railway env var
    printf '%s' "$SSH_HOST_KEY" | base64 -d > .ssh/id_ed25519
    chmod 600 .ssh/id_ed25519
    echo "Using persistent SSH host key from environment variable"
else
    # Fallback: generate a new key (first deploy, or if var not set)
    if [ ! -f .ssh/id_ed25519 ]; then
        ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""
        echo "Generated new SSH host key (set SSH_HOST_KEY env var to persist it)"
    fi
fi

exec ./ssh-portfolio
