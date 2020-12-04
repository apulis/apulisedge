########################## NODE TEST REST API ##########################
#  list all nodes
curl http://127.0.0.1:32767/apulisEdge/api/node/listNode -d @listNode.json|jq .

# create node
curl http://127.0.0.1:32767/apulisEdge/api/node/createNode -d @createNode.json|jq .

# describe node
curl http://127.0.0.1:32767/apulisEdge/api/node/desNode -d @desNode.json|jq .

# delete node
curl http://127.0.0.1:32767/apulisEdge/api/node/deleteNode -d @deleteNode.json|jq .


########################## APPLICATION TEST REST API ##########################
# list application
curl http://127.0.0.1:32767/apulisEdge/api/application/listApplication -d @listApplication.json|jq .

# create application
curl http://127.0.0.1:32767/apulisEdge/api/application/createApplication -d @createApplication.json|jq .

# delete application
curl http://127.0.0.1:32767/apulisEdge/api/application/deleteApplication -d @deleteApplication.json|jq .

# list deploy
curl http://127.0.0.1:32767/apulisEdge/api/application/listApplicationDeploy -d @listApplicationDeploy.json|jq .

# deploy application
curl http://127.0.0.1:32767/apulisEdge/api/application/deployApplication -d @deployApplication.json|jq .

# undeploy application
curl http://127.0.0.1:32767/apulisEdge/api/application/undeployApplication -d @undeployApplication.json|jq .

