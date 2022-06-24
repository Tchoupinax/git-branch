# --- helper functions for logs ---
function info() {
  echo '[INFO] ' "$@"
}
function warn() {
  echo '[WARN] ' "$@" >&2
}
function fatal() {
  echo '[ERROR] ' "$@" >&2
  exit 1
}

VERSION=0.0.10

info "Welcome to git-branch!"

info "Downloading version v$VERSION..."
ARCH=$(dpkg --print-architecture)

if [[ $ARCH = "amd64" ]]; then
  curl https://github.com/Tchoupinax/git-branch/releases/download/v0.0.10/git-branch_0.0.10_linux_amd64.tar.gz \
    -Lo /tmp/git-branch_0.0.10_linux_amd64.tar.gz 2> /dev/null
else
  fatal "Architecture $ARCH is not supported by this script for the moment, please be free to open an issue."
  exit 1
fi

info "Installing..."
tar -xf /tmp/git-branch_0.0.10_linux_amd64.tar.gz
rm -f /tmp/git-branch_0.0.10_linux_amd64.tar.gz

warn "Asking you password for installing in your path [/usr/local/bin]"
sudo mv git-branch /usr/local/bin/gb
rm -f git-branch

info "Successfully installed!"
