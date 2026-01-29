#!/bin/bash
set -e

REPO_URL="https://github.com/madss-bin/pgsync.git"
LOGO_FILE="branding/logo.txt"

C_BLUE='\033[38;5;39m'
C_PURPLE='\033[38;5;135m'
C_PINK='\033[38;5;213m'
C_GREEN='\033[38;5;82m'
C_GREY='\033[38;5;240m'
NC='\033[0m'

hide_cursor() { echo -ne "\033[?25l"; }
show_cursor() { echo -ne "\033[?25h"; }
cleanup() {
    show_cursor
    if [[ -n "$TEMP_DIR" ]]; then
        rm -rf "$TEMP_DIR"
    fi
}
trap cleanup EXIT

show_logo() {
    local logo="${1:-$LOGO_FILE}"
    [[ ! -f "$logo" ]] && return

    if command -v tte &> /dev/null; then
        HAS_TTE=1
    else
        HAS_TTE=0
    fi

    echo
    if [[ $HAS_TTE -eq 1 ]]; then
        cat "$logo" | tte \
            --frame-rate 60 beams \
            --beam-row-symbols "▂" "▁" "_" \
            --beam-column-symbols "▌" "▍" "▎" "▏" \
            --beam-delay 3 \
            --beam-row-speed-range 30-120 \
            --beam-column-speed-range 18-30 \
            --beam-gradient-stops FFEB3B FFB74D FF8A80 \
            --beam-gradient-steps 2 6 \
            --beam-gradient-frames 2 \
            --final-gradient-stops FFEB3B FFB74D FF8A80 F48FB1 EC407A \
            --final-gradient-steps 12 \
            --final-gradient-frames 2 \
            --final-gradient-direction vertical \
            --final-wipe-speed 3 2>/dev/null || { echo -e "${C_BLUE}" && cat "$logo" && echo -e "${NC}"; }
    else
        echo -e "${C_BLUE}"
        cat "$logo"
        echo -e "${NC}"
    fi
    echo
}

run_step() {
    local desc="$1"
    shift
    local cmds=("$@")
    local total_steps=${#cmds[@]}
    
    echo -e "${C_PURPLE}:: ${C_BLUE}$desc${NC}"

    for ((i=0; i<total_steps; i++)); do
        local cmd="${cmds[$i]}"
        local step_num=$((i+1))
        local percent=$(( step_num * 100 / total_steps ))
        
        local width=40
        local filled=$(( percent * width / 100 ))
        local empty=$(( width - filled ))
        local bar=$(printf "%0.s━" $(seq 1 $filled))
        local space=$(printf "%0.s━" $(seq 1 $empty))
        
        echo -ne "\033[1A\033[K"
        echo -e "${C_GREY}> $cmd${NC}"
        echo -ne "${C_GREEN}▕${C_PINK}${bar}${C_GREY}${space}${C_GREEN}▏ ${C_PINK}${percent}%${NC}\r"
        if ! eval "$cmd" > /dev/null 2>&1; then
             echo -e "\n${C_PINK}Command failed: $cmd${NC}"
             exit 1
        fi
    done
    
    echo -e "${C_GREEN}▕${C_PINK}$(printf "%0.s━" $(seq 1 40))${C_GREEN}▏ ${C_PINK}100%${NC}"
}

if [ ! -d ".git" ]; then
    echo -e "${C_BLUE}Detected non-git environment. Cloning repository...${NC}"
    TEMP_DIR=$(mktemp -d)
    git clone --depth 1 "$REPO_URL" "$TEMP_DIR" > /dev/null 2>&1
    cd "$TEMP_DIR"
fi

hide_cursor
show_logo

if [ -f /etc/arch-release ]; then
    DISTRO="arch"
elif [ -f /etc/fedora-release ]; then
    DISTRO="fedora"
elif [ -f /etc/debian_version ]; then
    DISTRO="ubuntu"
else
    DISTRO="unknown"
fi

STEPS_DEPS=()
if ! command -v go &> /dev/null; then
   case $DISTRO in
        "arch") STEPS_DEPS+=("sudo pacman -S --noconfirm go") ;;
        "fedora") STEPS_DEPS+=("sudo dnf install -y golang") ;;
        "ubuntu") 
            STEPS_DEPS+=("sudo apt update")
            STEPS_DEPS+=("sudo apt install -y golang-go") 
            ;;
        *)
            echo -e "${C_PINK}Please install Go manually.${NC}"
            exit 1
            ;;
    esac
fi

if [ ${#STEPS_DEPS[@]} -gt 0 ]; then
    run_step "Installing Dependencies" "${STEPS_DEPS[@]}"
fi

CMD_BUILD="go build -o bin/pgsync ."
CMD_INSTALL="sudo cp bin/pgsync /usr/local/bin/pgsync"
CMD_CHMOD="sudo chmod +x /usr/local/bin/pgsync"

run_step "Building and Installing" "$CMD_BUILD" "$CMD_INSTALL" "$CMD_CHMOD"

echo
echo -e "${C_GREEN}✓ pgsync installed successfully!${NC}"
echo -e "Run ${C_BLUE}pgsync${NC} to start."
