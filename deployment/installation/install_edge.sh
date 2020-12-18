#!/bin/bash
# ===
# === ENV variables
# ===
# common
ARCH= # will be init latter
NODENAME=$(hostname)
LOG_DIR=/var/log/apulisedge
INSTALL_LOG_FILE=${LOG_DIR}/installer.log
SCRIPT_DIR=/opt/apulisedge
DESIRE_DOCKER_VERSION=17.06

# kubeedge
KUBEEDGE_LOG_DIR=/var/log/kubeedge
KUBEEDGE_DATABASES_DIR=/var/lib/kubeedge
KUBEEDGE_TAR_FILE= # need arch, will be init later
DEFAULT_KUBEEDGE_EDGE_IMAGE=apulis/kubeedge-edge:1.0
KUBEEDGE_EDGE_IMAGE=${DEFAULT_KUBEEDGE_EDGE_IMAGE}
KUBEEDGE_HOME_PATH=/etc/kubeedge
EDGECORE_LOG_FILE=${LOG_DIR}/edgecore.log

# apulisedge
APULISEDGE_PACKAGE_DOWNLOAD_PATH=/tmp/apulisedge
APULISEDGE_IMAGE=apulisedge-agent:1.0


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
    # === check no kubeedge
    containerID=`docker ps | grep ${KUBEEDGE_EDGE_IMAGE} | awk '{print $1}'`
    if [[ "$containerID" != "" ]]; then
        LOG_ERROR "There is already a kubeedge client running. Please stop it and retry."
        return 1
    fi

    # === check hostname case
    HOSTNAME=`hostname`
    REGEX_OUTPUT=`echo ${HOSTNAME} | grep -P "[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*" -o`
    if [[ ! "${HOSTNAME}" == "${REGEX_OUTPUT}" ]]; then
        LOG_ERROR "Subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character (e.g. 'example.com', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*') "
        return 1
    fi

    # === check docker install status and version
    if [[ ! `command -v docker` ]]; then
        LOG_ERROR "Docker is not found but is required on node."
        LOG_ERROR "Please install docker and then try again."
    fi
    result=$(docker info 2>&1 | sed -n '1p' | grep Cannot | grep connect | grep Docker)
    if [ -n "${result}" ]; then
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
        else
            LOG_ERROR "docker version is too low, mini support version is ${DESIRE_DOCKER_VERSION}."
            return 1
        fi
    fi

    # === check download package
    LOG_INFO "check download package:${APULISEDGE_PACKAGE_DOWNLOAD_PATH}/${KUBEEDGE_TAR_FILE}"
    if [[ ! -e ${APULISEDGE_PACKAGE_DOWNLOAD_PATH}/${KUBEEDGE_TAR_FILE} ]];then
        LOG_ERROR "Can't find kubeedge.tar.gz"
        LOG_ERROR "Please fix and then try again."
        return 1
    fi

    # === check input params
    LOG_INFO "check input params:SERVER_DOMAIN"
    if [[ "${SERVER_DOMAIN}" = "" ]]; then
        LOG_ERROR "Cloud server domain is not specified."
        LOG_ERROR "Please fix and then try again."
        return 1
    fi
}

envInit()
{
    # === init edgecore env
    LOG_INFO "Initializing environment......"
    LOG_INFO "create directory..."
    mkdir -p ${KUBEEDGE_HOME_PATH}
    mkdir -p ${KUBEEDGE_LOG_DIR}
    mkdir -p ${KUBEEDGE_DATABASES_DIR}
    mkdir -p ${SCRIPT_DIR}
    LOG_INFO "directory ready."
    cd ${KUBEEDGE_HOME_PATH}
    LOG_INFO "decompress file..."
    cp ${APULISEDGE_PACKAGE_DOWNLOAD_PATH}/${KUBEEDGE_TAR_FILE} ${KUBEEDGE_HOME_PATH}
    tar -zxvf ${KUBEEDGE_HOME_PATH}/${KUBEEDGE_TAR_FILE}
    cp -r ${KUBEEDGE_HOME_PATH}/package/ca ${KUBEEDGE_HOME_PATH}/
    cp -r ${KUBEEDGE_HOME_PATH}/package/certs ${KUBEEDGE_HOME_PATH}/
    LOG_INFO "file decompressed."
    LOG_INFO "Initializing completed."

}

runEdgecore()
{
    LOG_INFO "pulling images..."
    docker pull ${KUBEEDGE_EDGE_IMAGE}
    docker tag ${KUBEEDGE_EDGE_IMAGE} ${DEFAULT_KUBEEDGE_EDGE_IMAGE}
    LOG_INFO "images ready."
    cd ${KUBEEDGE_HOME_PATH}
    # generate edgecore runtime config
    mkdir -p config
    mkdir -p /var/lib/edged
    LOG_INFO "create edgecore config file..."
    docker run ${KUBEEDGE_EDGE_IMAGE} /bin/bash -c "edgecore --minconfig" | tee config/edgecore.yaml
    sed -i "s#httpServer:\ .*10002#httpServer: https://${SERVER_DOMAIN}:10002#g" config/edgecore.yaml
    sed -i "s#server:\ .*10001#server: ${SERVER_DOMAIN}:10001#g" config/edgecore.yaml
    sed -i "s#server:\ .*10000#server: ${SERVER_DOMAIN}:10000#g" config/edgecore.yaml
    sed -i "s#hostnameOverride:\ .*#hostnameOverride:\ ${NODENAME}#g" config/edgecore.yaml
    sed -i "/.*token:\ .*/d" config/edgecore.yaml
    # run edgecore image
    LOG_INFO "config file generated."
    systemctl enable docker.service
    systemctl start docker
    LOG_INFO "run edgecore container..."
    docker run -d -P \
    --restart=always \
    --privileged=true \
    --network=host \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/lib/edged:/var/lib/edged \
    -v /var/lib/docker:/var/lib/docker \
    -v ${KUBEEDGE_DATABASES_DIR}:${KUBEEDGE_DATABASES_DIR} \
    -v ${KUBEEDGE_LOG_DIR}:${KUBEEDGE_LOG_DIR} \
    -v ${KUBEEDGE_HOME_PATH}:${KUBEEDGE_HOME_PATH} \
    ${KUBEEDGE_EDGE_IMAGE}
    LOG_INFO "container started"
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
                -l)
                KUBEEDGE_EDGE_IMAGE="$2"
                shift;
                shift;
                ;;
                -h)
                NODENAME="$2"
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

    # === init some variables
    ARCH="$(uname -m)"
    KUBEEDGE_TAR_FILE=apulisedge_${ARCH}.tar.gz

    # === init log
    mkdir -p ${LOG_DIR}
    if [ ! -e "${INSTALL_LOG_FILE}}" ]; then
        touch "${INSTALL_LOG_FILE}"
    fi

    process=(
        envCheck
        envInit
        runEdgecore
    )

    LOG_INFO "=== edge node install begin"
    for i in "${!process[@]}";do
        LOG_INFO "process ${process[${i}]} begin"
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
if [ "${1}" != "--source-only" ]; then
    main "$@"
fi