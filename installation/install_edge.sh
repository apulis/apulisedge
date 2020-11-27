# ===
# === ENV variables
# ===

if [[ ! command -v docker ]]; then
    echo "ERROR !!!"
    echo "Docker is not found but is required on node."
    echo "Please install docker and then try again."
fi
echo "=== edge node install begin"

edgecore --minconfig > edgecore.yaml
edgecore --config edgecore.yaml

echo "=== edge node install completed"