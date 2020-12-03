# create node
curl http://127.0.0.1:32767/apulisEdge/api/node/createNode -d @createNode.json|jq .

#  list all nodes
curl http://127.0.0.1:32767/apulisEdge/api/node/listNodes -d @listNodes.json|jq .

# describe node
curl http://127.0.0.1:32767/apulisEdge/api/node/desNode -d @desNode.json|jq .

# create application
curl http://127.0.0.1:32767/apulisEdge/api/application/createApplication -d @createApplication.json|jq .

# deploy application
curl http://127.0.0.1:32767/apulisEdge/api/application/deployApplication -d @deployApplication.json|jq .