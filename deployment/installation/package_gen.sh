#! /bin/bash
# ===
# === ENV variables
# ===
ARCH= # will be init latter
PACKAGE_NAME="package"
PACKAGE_PATH="/tmp/apulisedge/${PACKAGE_NAME}/"
LOG_DIR=/var/log/apulisedge
LOG_FILE="${LOG_FILE_PATH}/package_gen.log"
CLOUD_DOMAIN="apulis.cn"
INSTALL_SCRIPT_PATH="${PACKAGE_PATH}/scripts"


# ===
# === funtion statements
# ===
getArch()
{
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

LOG_INFO()
{
    message="$(date +%Y-%m-%dT%H:%M:%S) | INFO | $*"
    echo "${message}"
    echo "${message}" >> "${LOG_FILE}" 2>&1
}

LOG_ERROR()
{
    message="$(date +%Y-%m-%dT%H:%M:%S) | ERROR | $*"
    echo "${message}"
    echo "${message}" >> "${LOG_FILE}" 2>&1
}

genCA() {
    openssl genrsa -des3 -out ${CA_PATH}/rootCA.key -passout pass:cloud.sigsus.cn 4096
    openssl req -x509 -new -nodes -key ${CA_PATH}/rootCA.key -sha256 -days 3650 \
    -subj ${SUBJECT} -passin pass:cloud.sigsus.cn -out ${CA_PATH}/rootCA.crt
}

ensureCA() {
    if [ ! -e ${CA_PATH}/rootCA.key ] || [ ! -e ${CA_PATH}/rootCA.crt ]; then
        genCA
    fi
}

genCsr() {
    local name=$1
    openssl genrsa -out ${CERT_PATH}/${name}.key 2048
    openssl req -new -key ${CERT_PATH}/${name}.key -subj ${SUBJECT} -out ${CERT_PATH}/${name}.csr
}

genCert() {
    local name=$1
    openssl x509 -req -in ${CERT_PATH}/${name}.csr -CA ${CA_PATH}/rootCA.crt -CAkey ${CA_PATH}/rootCA.key \
    -CAcreateserial -passin pass:cloud.sigsus.cn -out ${CERT_PATH}/${name}.crt -days 365 -sha256
}

genCertAndKey() {
    local name="server"
    echo ${CA_PATH}
    echo ${CERT_PATH}
    ls -l ${CA_PATH}
    ls -l ${CERT_PATH}
    if [[ -e ${CA_PATH}/rootCA.key && -e ${CA_PATH}/rootCA.crt && -e ${CA_PATH}/rootCA.srl && -e ${CERT_PATH}/${name}.crt && -e ${CERT_PATH}/${name}.csr && -e ${CERT_PATH}/${name}.key ]]; then
        LOG_INFO "CA and Certs has been generated, use generated certificates now."
        return 0
    fi
    LOG_INFO "CA and Certs has not been generated, generating..."
    ensureCA
    genCsr $name
    genCert $name
}

envCheck()
{
    # neccessary software check
    if [[ ! `command -v openssl` ]]; then
        LOG_ERROR "Openssl is not found, please try again."
    fi
    if [[ ! `command -v wget` ]]; then
        LOG_ERROR "Wget is not found, please try again."
    fi
    if [[ ! `command -v md5sum` ]]; then
        LOG_ERROR "Md5sum is not found, please try again."
    fi

    # cloud domain check
    LOG_INFO "CLOUD DOMAIN has been set to: ${CLOUD_DOMAIN}"
    # arch type check
    LOG_INFO "ARCH type has been set to: ${ARCH}"
}

envInit()
{
    # init some environment variables
    TAR_PACKAGE_NAME="apulisedge_${ARCH}.tar.gz"
    LOG_INFO "Package will be save as: ${TAR_PACKAGE_NAME}"
    CA_PATH=${CA_PATH:-${PACKAGE_PATH}/ca}
    CA_SUBJECT=${CA_SUBJECT:-/C=CN/ST=Guangdong/L=Shenzhen/O=Apulis/CN=${CLOUD_DOMAIN}}
    CERT_PATH=${CERT_PATH:-${PACKAGE_PATH}/certs}
    SUBJECT=${SUBJECT:-/C=CN/ST=Guangdong/L=Shenzhen/O=Apulis/CN=${CLOUD_DOMAIN}}

    # create directory
    mkdir -p ${PACKAGE_PATH}
    mkdir -p ${CA_PATH}
    mkdir -p ${CERT_PATH}
    mkdir -p ${INSTALL_SCRIPT_PATH}

}

downloadScripts()
{
    cd ${PACKAGE_PATH}/..
    git clone https://apulis-gitlab.apulis.cn/apulis/apulisedge.git -b develop
    cp apulisedge/deployment/installation/install_edge.sh ${INSTALL_SCRIPT_PATH}
    cp apulisedge/deployment/installation/uninstall_edge.sh ${INSTALL_SCRIPT_PATH}
}

compressPackage()
{
    cd ${PACKAGE_PATH}

    find ./ -type f -print0 | xargs -0 md5sum > ./checksum.md5

    cd -
    cd ${PACKAGE_PATH}/..
    tar -cvzf ./${TAR_PACKAGE_NAME} ${PACKAGE_NAME}

    cd -
}

signPackage()
{
    cd ${PACKAGE_PATH}/..

    mkdir -p private
    openssl genrsa -out private/rsa_private_${ARCH}.key
    openssl rsa -in private/rsa_private_${ARCH}.key -pubout -out apulisedge_${ARCH}.key
    openssl dgst -sign private/rsa_private_${ARCH}.key -sha256 -out apulisedge_${ARCH}.sig ${TAR_PACKAGE_NAME}

    cd -
}

envClean()
{
    cd ${PACKAGE_PATH}/..

    # delete git project
    rm -rf apulisedge

    cd -
}

main()
{
    if which getopt > /dev/null 2>&1; then
        OPTS=$(getopt d:i: "$*" 2>/dev/null)
        if [ ! $? ]; then
            printf "%s\\n" "$USAGE"
            exit 2
        fi

        eval set -- "$OPTS"
        while true; do
            case "$1" in
                -d)
                CLOUD_DOMAIN="$2"
                shift;
                shift;
                ;;
                --)
                shift
                break
                ;;
                *)
                printf "ERROR: did not recognize option '%s'\\n" "$1"
                exit 1
                ;;
            esac
        done
    fi

    process=(
        envCheck
        envInit
        genCertAndKey
        downloadScripts
        compressPackage
        signPackage
        envClean
    )

    # === init log
    mkdir -p ${LOG_DIR}
    if [ ! -e "${LOG_FILE}}" ]; then
        touch "${LOG_FILE}"
    fi

    LOG_INFO "=== package generate begin"
    for i in "${!process[@]}";do
        LOG_INFO "process ${process[${i}]} begin"
        ${process[${i}]}
        if [ $? -ne 0 ]; then
            LOG_ERROR "process-${process[${i}]} failed"
            exit 1
        fi
    done
    LOG_INFO "=== package generate completed"

}


# ===
# === main code start here
# ===
# main "$@"

# now amd and arm packages are the same
ARCH="amd64"
main "$@"
ARCH="arm64"
main "$@"
