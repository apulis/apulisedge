# ===
# === ENV variables
# ===
ARCH= # will be init latter
LOG_DIR=/var/log/apulisedge
INSTALL_LOG_FILE=${LOG_DIR}/installer.log
EDGECORE_LOG_FILE=${LOG_DIR}/edgecore.log
APULISEDGE_PACKAGE_DOWNLOAD_PATH=/tmp
KUBEEDGE_TAR_FILE= # need arch, will be init later
KUBEEDGE_HOME_PATH=/etc/kubeedge


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
    echo "${message}" >> "${INSTALL_LOG_FILE}" 2>&1
}

LOG_ERROR()
{
    message="$(date +%Y-%m-%dT%H:%M:%S) | ERROR | $*"
    echo "${message}"
    echo "${message}" >> "${INSTALL_LOG_FILE}" 2>&1
}

envCheck()
{
    if [[ ! command -v docker ]]; then
        LOG_ERROR "ERROR !!!"
        LOG_ERROR "Docker is not found but is required on node."
        LOG_ERROR "Please install docker and then try again."
    fi

    if [[ ! -e ${APULISEDGE_PACKAGE_DOWNLOAD_PATH}/kubeedge.tar.gz ]];then
        LOG_ERROR "ERROR !!!"
        LOG_ERROR "Can't find kubeedge.tar.gz"
        LOG_ERROR "Please fix and then try again."
    fi

    if [[ "${SERVER_DOMAIN}" = "" ]]; then
        LOG_ERROR "ERROR !!!"
        LOG_ERROR "Cloud server domain is not specified."
        LOG_ERROR "Please fix and then try again."
    fi
}

envInit()
{
    # === init some variables
    ARCH="$(getArch $(uname -m))"
    KUBEEDGE_TAR_FILE=kubeedgeRuntime-${ARCH}.tar.gz

    # === init log
    mkdir -p ${LOG_DIR}
    if [ -e "${INSTALL_LOG_FILE}}" ]; then
        rm "${INSTALL_LOG_FILE}"
    fi
    touch "${INSTALL_LOG_FILE}"

    # === init edgecore env
    mkdir -p ${KUBEEDGE_HOME_PATH}
    cd ${KUBEEDGE_HOME_PATH}
    cp ${APULISEDGE_PACKAGE_DOWNLOAD_PATH}/${KUBEEDGE_TAR_FILE} ${KUBEEDGE_HOME_PATH}
    tar -zxf ${KUBEEDGE_HOME_PATH}/${KUBEEDGE_TAR_FILE}
    cp kubeedgeRuntime/edge/edgecore /usr/local/bin

}

runEdgecore()
{
    cd ${KUBEEDGE_HOME_PATH}
    mkdir -p config
    edgecore --minconfig > config/edgecore.yaml
    sed -i "s#httpServer:\ .*10002#httpServer: https://${SERVER_DOMAIN}:10002#g" config/edgecore.yaml
    sed -i "s#server:\ .*10001#server: ${SERVER_DOMAIN}:10001#g" config/edgecore.yaml
    sed -i "s#server:\ .*10000#server: ${SERVER_DOMAIN}:10000#g" config/edgecore.yaml
    nohup edgecore --config config/edgecore.yaml > ${EDGECORE_LOG_FILE} 2>&1 &
}

main()
{
    if which getopt > /dev/null 2>&1; then
        OPTS=$(getopt d:n:r:m:ulhez "$*" 2>/dev/null)
        if [ ! $? ]; then
            printf "%s\\n" "$USAGE"
            exit 2
        fi

        eval set -- "$OPTS"
        while true; do
            case "$1" in
                -d)
                SERVER_DOMAIN="$2"
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
        runEdgecore
    )

    LOG_INFO "=== edge node install begin"
    for i in "${!process[@]}";do
        LOG_INFO "process ${i} begin"
        ${process[${i}]}
        if [ $? -ne 0 ]; then
            LOG_ERROR "process-${process[${i}]} failed"
            exit 1
        fi
    done
    LOG_INFO "=== edge node install completed"
}

# ===
# === main code start here
# ===
main "$@"