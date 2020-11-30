# ===
# === ENV variables
# ===

# ===
# === function statements
# ===
arch() {
    case $1 in
    aarch64)
        echo arm64
        ;;
    x86_64)
        echo amd64
        ;;
    amd64) # FreeBSD.
        echo amd64
        ;;
    arm64) # FreeBSD.
        echo arm64
        ;;
    *)
        ;;
    esac
}

# ===
# === main code start here
# ===
# params checks
if [ ! "$1" = "" ]; then
    echo "ERROR: need to specify a target architecture"
    exit
fi
version="$(curl -fsSLI -o /dev/null -w "%{url_effective}" https://github.com/kubeedge/kubeedge/releases/latest)"
version="${version#https://github.com/kubeedge/kubeedge/releases/tag/}"
version="${version#v}"
ARCH="$(arch $1)"
if [ ${ARCH} == "" ]; then
    echo "ERROR: architecture error"
    exit
fi
curl -fOL https://github.com/kubeedge/kubeedge/releases/download/v$version/kubeedge-v${version}-linux-${ARCH}.tar.gz
tar -zxvf kubeedge-v${version}-linux-${ARCH}.tar.gz
rm kubeedge-v${version}-linux-${ARCH}.tar.gz
mv kubeedge-v${version}-linux-${ARCH} kubeedgeRuntime-${ARCH}
tar -zxcf kubeedgeRuntime-${ARCH}.tar.gz ./kubeedgeRuntime-${ARCH}
chmod +x kubeedgeRuntime-${ARCH}/cloud/cloudcore/cloudcore
cp kubeedgeRuntime-${ARCH}/cloud/cloudcore/cloudcore /usr/local/bin
chmod +x kubeedgeRuntime-${ARCH}/edge/edgecore
cp kubeedgeRuntime-${ARCH}/edge/edgecore /usr/local/bin