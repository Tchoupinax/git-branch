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

VERSION="0.0.11"

info "Welcome to git-branch!"

info "Downloading version v$VERSION..."
ARCH=$(dpkg --print-architecture)

if [[ $ARCH = "amd64" ]]; then
  curl https://github.com/Tchoupinax/git-branch/archive/refs/tags/v0.0.11.tar.gz \
    -Lo /tmp/git-branch-v0.0.11.tar.gz 2> /dev/null
else
  fatal "Architecture $ARCH is not supported by this script for the moment, please be free to open an issue."
  exit 1
fi

info "Installing..."
tar -xvzf /tmp/git-branch-v0.0.11.tar.gz
rm -f /tmp/git-branch-v0.0.11.tar.gz

cd git-branch-0.0.11
info "building..."
go build -o git-branch

warn "Asking you password for installing in your path [/usr/local/bin]"
sudo mv git-branch /usr/local/bin/gb
cd ..
rm -rf git-branch-0.0.11

info "Successfully installed!"
