first, let's make sure all services are built and deployed
command: kubectl get pods -n student-cafe

Can't place order it is showing 'Food catalog service not available" in postman and Failed to place order. please try again

check the logs
kubectl get pods -n student-cafe | grep order
kubectl logs order-deployment-65dd6db867-fkrjx -n student-cafe --tail=20



The issue was: The order service couldn't find the food-catalog-service through Consul's health checks, even though both services were registered.

The fix: I modified the order service to continue processing orders even if it can't find the catalog service through Consul, using the Kubernetes service name as a fallback.

