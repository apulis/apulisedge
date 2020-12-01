#!/bin/bash
# ===
# === ENV variables
# ===
ARCH= # will be init latter
LOG_DIR=/var/log/apulisedge
KUBEEDGE_LOG_DIR=/var/log/kubeedge
INSTALL_LOG_FILE=${LOG_DIR}/installer.log
EDGECORE_LOG_FILE=${LOG_DIR}/edgecore.log
APULISEDGE_PACKAGE_DOWNLOAD_PATH=/tmp/apulisedge
KUBEEDGE_TAR_FILE= # need arch, will be init later
KUBEEDGE_HOME_PATH=/etc/kubeedge
DESIRE_DOCKER_VERSION=17.06
APULISEDGE_IMAGE=apulisedge-agent:1.0
KUBEEDGE_EDGE_IMAGE=apulis/kubeedge-edge:1.0


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
    # === check docker install status and version
    if [[ ! `command -v docker` ]]; then
        LOG_ERROR "ERROR !!!"
        LOG_ERROR "Docker is not found but is required on node."
        LOG_ERROR "Please install docker and then try again."
    fi
    result=$(docker info 2>&1 | sed -n '1p' | grep Cannot | grep connect | grep Docker)
    if [ -n "${result}" ]; then
        LOG_ERROR "ERROR !!!"
        LOG_ERROR "docker is not start, please start first."
        return 1
    fi
    result=$(docker info 2>&1 | sed -n "/Version/p" | grep Server)
    if [ -n "${result}" ]; then
        LOG_INFO "got docker ${result}"
        version=$(echo ${result} | awk '{print $3}')
        big_version=$(echo ${version} | awk -F '[.]' '{print $1}')
        small_version=$(echo ${version} | awk -F '[.]' '{print $2}')
        docker_verison="${big_version}.${small_version}"
        if [ $(expr ${docker_verison} \>= ${DESIRE_DOCKER_VERSION}) -eq 1 ]; then
            LOG_INFO "check docker version success"
            return 0
        else
            LOG_ERROR "docker version is too low, mini support version is ${DESIRE_DOCKER_VERSION}."
            LOG_ERROR "docker version is too low, mini support version is ${DESIRE_DOCKER_VERSION}."
            return 1
        fi
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
    mkdir -p ${KUBEEDGE_LOG_DIR}
    cd ${KUBEEDGE_HOME_PATH}
    cp ${APULISEDGE_PACKAGE_DOWNLOAD_PATH}/${KUBEEDGE_TAR_FILE} ${KUBEEDGE_HOME_PATH}
    tar -zxf ${KUBEEDGE_HOME_PATH}/${KUBEEDGE_TAR_FILE}

}

runEdgecore()
{
    docker pull ${KUBEEDGE_EDGE_IMAGE}
    cd ${KUBEEDGE_HOME_PATH}
    # generate edgecore runtime config
    mkdir -p config
    docker run ${KUBEEDGE_EDGE_IMAGE} /bin/bash -c "edgecore --minconfig" | tee config/edgecore.yaml
    sed -i "s#httpServer:\ .*10002#httpServer: https://${SERVER_DOMAIN}:10002#g" config/edgecore.yaml
    sed -i "s#server:\ .*10001#server: ${SERVER_DOMAIN}:10001#g" config/edgecore.yaml
    sed -i "s#server:\ .*10000#server: ${SERVER_DOMAIN}:10000#g" config/edgecore.yaml
    # run edgecore image
    systemctl enable docker.service
    systemctl start docker
    docker run -d -P --restart=always --privileged=true --network=host -v ${KUBEEDGE_LOG_DIR}:${KUBEEDGE_LOG_DIR} -v /var/run/docker.sock:/var/run/docker.sock -v ${KUBEEDGE_HOME_PATH}:${KUBEEDGE_HOME_PATH} ${APULISEDGE_IMAGE}
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
                SERVER_DOMAIN="$2"
                shift;
                shift;
                ;;
                -i)
                APULISEDGE_IMAGE="$2"
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