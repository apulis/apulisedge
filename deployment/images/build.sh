# ===
# === function statement
# ===
buildAgentDocer()
{
    
}

buildCloudDocer()
{

}

main()
{
    case "$1" in
        cloud)
        SERVER_DOMAIN="$2"
        shift;
        shift;
        ;;
        agent)
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
}



# ===
# === main code start here
# ===
main "$@"