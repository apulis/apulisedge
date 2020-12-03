#!/bin/bash
# ===
# === ENV variables
# ===
. /opt/apulisedge/install_edge.sh --source-only


# ===
# === function statement
# ===
doubleCheck()
{
    printf "!!! CAUTION !!!\\n"
    printf "Edge node will be removed from cluster!\\n"
    printf "Is it right to continue? [y/N(default)] >>>"

	read -r ans
	while [ "$ans" != "y" ] && [ "$ans" != "N" ] && [ "$ans" != "N" ]
	do
		printf "Please answer 'y' or 'N':'\\n"
		printf ">>> "
		read -r ans
	done

	if [ "$ans" == "y" ]; then
        echo "Uninstall proess continue..."
	else
        echo "Uninstall process stop."
        exit 0
	fi
}

stopEdgecore()
{
    containerID=`docker ps | grep ${KUBEEDGE_EDGE_IMAGE} | awk '{print $1}'`
    docker stop ${containerID}
}

deleteEdgecore()
{
    docker rmi ${KUBEEDGE_EDGE_IMAGE}
}

deleteFileAndDir()
{
    LOG_INFO "delete apulisedge downloading file"
    rm -rf ${APULISEDGE_PACKAGE_DOWNLOAD_PATH}
    rm -rf ${SCRIPT_DIR}
    LOG_INFO "delete kubeedge runtime"
    rm -rf ${KUBEEDGE_HOME_PATH}
}

main()
{
    process=(
        doubleCheck
        stopEdgecore
        deleteEdgecore
        deleteFileAndDir
    )
    LOG_INFO "=== uninstall begin"
    for i in "${!process[@]}";do
        LOG_INFO "process ${process[${i}]} begin"
        ${process[${i}]}
        if [ $? -ne 0 ]; then
            LOG_ERROR "process-${process[${i}]} failed"
            exit 1
        fi
    done
    LOG_INFO "=== uninstall completed"
}


# ===
# === main code start here
# ===
main "$@"