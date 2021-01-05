#! /bin/bash

declare -i CA_DAYS CERT_DAYS
CA_DAYS=3650
CERT_DAYS=365
CA_PATH=`pwd`
CERT_PATH=`pwd`
CERT_DOMAIN="apulis.cn"
SUBJECT="/C=CN/ST=Guangdong/L=Shenzhen/O=Apulis/CN=${CERT_DOMAIN}"

LOG_INFO()
{
    message="$(date +%Y-%m-%dT%H:%M:%S) | INFO | $*"
    echo "${message}"
}

LOG_ERROR()
{
    message="$(date +%Y-%m-%dT%H:%M:%S) | ERROR | $*"
    echo "${message}"
}

genCA() {
    openssl genrsa -des3 -out ${CA_PATH}/rootCA.key -passout pass:apulis.cn 4096
    openssl req -x509 -new -nodes -key ${CA_PATH}/rootCA.key -sha256 -days $CA_DAYS \
    -subj ${SUBJECT} -passin pass:apulis.cn -out ${CA_PATH}/rootCA.crt
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
    -CAcreateserial -passin pass:apulis.cn -out ${CERT_PATH}/${name}.crt -days $CERT_DAYS -sha256
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

main()
{
    OPTS=`getopt -q -o d:c:m: --long ca-days:,cert-days:,cert-domain: -n "$0" -- "$@"`
    eval set -- "$OPTS"

    while true
    do
      case "$1" in
        -d|--ca-days)
          CA_DAYS=$2
          shift 2
          ;;
        -c|--cert-days)
          CERT_DAYS=$2
          shift 2
          ;;
        -c|--cert-domain)
          CERT_DOMAIN=$2
          shift 2
          ;;
        --)
          shift
          break
          ;;
        *)
          echo "Internal error!";
          exit 1
          ;;
      esac
    done

    echo "ca days = $CA_DAYS";
    echo "cert days = $CERT_DAYS";
    echo ""

    process=(
        genCertAndKey
    )

    # === init log
    LOG_INFO "=== certificate generate begin"
    for i in "${!process[@]}";do
        LOG_INFO "process ${process[${i}]} begin"
        ${process[${i}]}
        if [ $? -ne 0 ]; then
            LOG_ERROR "process-${process[${i}]} failed"
            exit 1
        fi
    done
    LOG_INFO "=== certificate generate completed"
}


# ===
# === main code start here
# ===
# main "$@"
main "$@"
