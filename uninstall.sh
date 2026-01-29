#!/bin/bash
set -e

# Colors
C_BLUE='\033[38;5;39m'
C_PURPLE='\033[38;5;135m'
C_PINK='\033[38;5;213m'
C_GREEN='\033[38;5;82m'
C_GREY='\033[38;5;240m'
C_RED='\033[38;5;196m'
NC='\033[0m'

echo
echo -e "${C_BLUE}pgsync Uninstaller${NC}"
echo -e "${C_GREY}===================${NC}"
echo

BIN_PATH="/usr/local/bin/pgsync"
CONFIG_DIR="$HOME/.pgsync"

# 1. Remove Binary
if [ -f "$BIN_PATH" ]; then
    echo -e "${C_PURPLE}Found pgsync binary at: ${NC}$BIN_PATH"
    echo -e "${C_GREY}Removing binary...${NC}"
    
    if [ -w $(dirname "$BIN_PATH") ]; then
        rm "$BIN_PATH"
    else
        sudo rm "$BIN_PATH"
    fi
    echo -e "${C_GREEN}✓ Binary removed.${NC}"
else
    echo -e "${C_PURPLE}Binary not found at $BIN_PATH (skipping)${NC}"
fi

echo

# 2. Remove Config/History
if [ -d "$CONFIG_DIR" ]; then
    echo -e "${C_PURPLE}Found configuration/history at: ${NC}$CONFIG_DIR"
    echo -ne "${C_PINK}Remove history and configuration? [y/N]: ${NC}"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]]; then
        rm -rf "$CONFIG_DIR"
        echo -e "${C_GREEN}✓ Configuration directory removed.${NC}"
    else
        echo -e "${C_GREY}Configuration kept.${NC}"
    fi
fi

echo
echo -e "${C_GREEN}Uninstallation successful.${NC}"
echo
