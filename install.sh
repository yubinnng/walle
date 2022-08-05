echo "* Creating 'walle' namespace"
kubectl apply -f deployments/walle-namespace.yml

echo "* Deploying PostgreSQL"
kubectl apply -f deployments/postgres.yml

echo "* Deploying NATS"
kubectl apply -f deployments/nats.yml

echo "* Deploying worfklow engine"
faas-cli deploy -f walle-engine.yml -g $OPENFAAS_GATEWAY

echo "* Deploying api-server"
kubectl apply -f deployments/api-server.yml

echo "* Deploying ui dashboard"
kubectl apply -f deployments/ui.yml