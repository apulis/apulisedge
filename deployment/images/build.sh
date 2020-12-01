# ===
# === ENV variables
# ===
DOCKERFILE_DIR=deployment/images/
APULISEDGE_AGENT_DOCKERFLE=Dockerfile-apulisEdgeAgent
APULISEDGE_CLOUD_DOCKERFLE=Dockerfile-apulisEdgeCloud
KUBEEDGE_EDGE_DOCKERFLE=Dockerfile-kubeedgeEdge
KUBEEDGE_CLOUD_DOCKERFLE=Dockerfile-kubeedgeCloud
DOWNLOAD_KUBEEDGE_SCRIPT_PATH=deployment/images/download_kubeedge.sh


# ===
# === function statement
# ===
buildAgent()
{
    docker build . -f ${DOCKERFILE_DIR}/${APULISEDGE_AGENT_DOCKERFLE} -t apulis/apulisedge-agent:1.0
}

buildCloud()
{
    docker build . -f ${DOCKERFILE_DIR}/${APULISEDGE_CLOUD_DOCKERFLE} -t apulis/apulisedge-cloud:1.0
}

buildKubeedgeEdge()
{
    docker build . -f ${DOCKERFILE_DIR}/${KUBEEDGE_EDGE_DOCKERFLE} --build-arg downloadscript_path=${DOWNLOAD_KUBEEDGE_SCRIPT_PATH} -t apulis/kubeedge-edge:1.0
}

buildKubeedgeCloud()
{
    docker build . -f ${DOCKERFILE_DIR}/${KUBEEDGE_CLOUD_DOCKERFLE} --build-arg downloadscript_path=${DOWNLOAD_KUBEEDGE_SCRIPT_PATH} -t apulis/kubeedge-cloud:1.0
}

buildAll()
{
    buildAgent
    buildCloud
    buildKubeedgeEdge
    buildKubeedgeCloud
}

buildKubeedge()
{
    buildKubeedgeEdge
    buildKubeedgeCloud
}

buildApulisedge()
{
    buildAgent
    buildCloud
}

main()
{
    cd ../..
    $1 $2
}



# ===
# === main code start here
# ===
main "$@"