set -e 
package="selfhttps"
tmpdir="/tmp/${package}"
version=$(curl -fsSL https://raw.githubusercontent.com/ahuigo/selfhttps/main/version)
version=${version#v}
function is_installed() {
  if [ -f /usr/local/bin/selfhttps ]; then
    echo "selfhttps is installed"
    return 0
  else
    return 1
  fi
}

function getOs() {
    # return darwin linux windows
    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux*) os="linux";;
        Darwin*) os="darwin";;
        CYGWIN*) os="windows";;
        MINGW*) os="windows";;
        *) os="unknown";;
    esac
    echo "${os}"
}
function getArch() {
    # return arm64 amd64 386
    unameOut="$(uname -m)"
    case "${unameOut}" in
        "x86_64") arch="amd64";;
        "i386") arch="386";;
        "arm64") arch="arm64";;
        *) arch="unknown";;
    esac
    echo "${arch}"
}

function install() {
    os=$(getOs)
    arch=$(getArch)
    if [ "$os" != "linux" ] && [ "$os" != "darwin" ]; then
        echo "Not support os: $os"
        return 1
    fi
    if [ "$arch" == "unknown" ]; then
        echo "Not support arch: $arch"
        return 1
    fi
    url="https://github.com/ahuigo/selfhttps/releases/download/v${version}/selfhttps_${version}_${os}_${arch}.tar.gz"
    echo "Downloading $url"

    mkdir -p $tmpdir
    { wget -O $tmpdir/a.tar.gz $url || curl -L -C - -o $tmpdir/a.tar.gz $url; } && tar -zxvf $tmpdir/a.tar.gz -C $tmpdir 
    { mv "$tmpdir/$package" /usr/local/bin/ || sudo mv "$tmpdir/$package" /usr/local/bin/; } && rm -rf $tmpdir
    echo "Install /usr/local/bin/$package success!"
}

function main() {
  if is_installed; then
    echo "selfhttps is installed"
    return 0
  fi
  echo "Installing selfhttps..."
  install
}
main