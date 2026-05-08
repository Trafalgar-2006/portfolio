#!/usr/bin/env bash
# install.sh — drops a `mohith` command on your system
# Usage: curl -fsSL https://raw.githubusercontent.com/Trafalgar-2006/portfolio/master/install.sh | bash

set -e

CMD_NAME="mohith"
SSH_HOST="mohith.is-a.dev"
SSH_PORT="41074"
INSTALL_DIR="/usr/local/bin"

# ── colours ──────────────────────────────────────────────────────────────────
BOLD="\033[1m"
GREEN="\033[1;32m"
CYAN="\033[1;36m"
YELLOW="\033[1;33m"
RESET="\033[0m"

echo ""
echo -e "${CYAN}${BOLD}  Installing '${CMD_NAME}' command...${RESET}"
echo ""

# ── write the tiny wrapper ────────────────────────────────────────────────────
WRAPPER=$(mktemp)
cat > "$WRAPPER" << EOF
#!/usr/bin/env bash
# mohith — opens Mohith's SSH portfolio
exec ssh ${SSH_HOST} -p ${SSH_PORT} "\$@"
EOF
chmod +x "$WRAPPER"

# ── move to install dir (sudo if needed) ─────────────────────────────────────
TARGET="$INSTALL_DIR/$CMD_NAME"

if mv "$WRAPPER" "$TARGET" 2>/dev/null; then
  :
elif sudo mv "$WRAPPER" "$TARGET" 2>/dev/null; then
  :
else
  # Fallback: install to ~/.local/bin and remind user to add it to PATH
  LOCAL_BIN="$HOME/.local/bin"
  mkdir -p "$LOCAL_BIN"
  mv "$WRAPPER" "$LOCAL_BIN/$CMD_NAME"
  TARGET="$LOCAL_BIN/$CMD_NAME"

  # Add to shell rc if not already present
  for RC in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
    if [ -f "$RC" ] && ! grep -q "$LOCAL_BIN" "$RC" 2>/dev/null; then
      echo "export PATH=\"\$PATH:$LOCAL_BIN\"" >> "$RC"
    fi
  done

  echo -e "${YELLOW}  Installed to $TARGET${RESET}"
  echo -e "${YELLOW}  Restart your shell (or run: export PATH=\"\$PATH:$LOCAL_BIN\")${RESET}"
  echo ""
fi

# ── done ─────────────────────────────────────────────────────────────────────
echo -e "${GREEN}${BOLD}  ✓ Done!${RESET}  Just run:  ${BOLD}${CMD_NAME}${RESET}"
echo ""
echo -e "  (or connect directly:  ${CYAN}ssh mohith.is-a.dev -p ${SSH_PORT}${RESET})"
echo ""
