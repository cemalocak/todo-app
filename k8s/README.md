# Kubernetes Deployment Guide

## Prerequisites

1. **Kubernetes Cluster** (v1.20+)
2. **kubectl** configured to connect to your cluster
3. **Docker images** built and pushed to registry
4. **Ingress Controller** (nginx-ingress recommended)

## Quick Deployment

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get all -n todo-app

# Check pod logs
kubectl logs -n todo-app -l app=todo-backend
kubectl logs -n todo-app -l app=todo-frontend
```

## Step-by-Step Deployment

### 1. Create Namespace
```bash
kubectl apply -f k8s/namespace.yaml
```

### 2. Deploy Backend
```bash
kubectl apply -f k8s/backend.yaml
```

### 3. Deploy Frontend  
```bash
kubectl apply -f k8s/frontend.yaml
```

### 4. Setup Ingress
```bash
# Install nginx-ingress controller (if not installed)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Apply ingress
kubectl apply -f k8s/ingress.yaml
```

## Configuration

### Update Docker Images
Edit the image references in the deployment files:
- `backend.yaml`: Update `image: ghcr.io/username/todo-app-backend:latest`
- `frontend.yaml`: Update `image: ghcr.io/username/todo-app-frontend:latest`

### Update Domain
Edit `ingress.yaml`:
- Replace `todo-app.yourdomain.com` with your actual domain

## Monitoring

```bash
# Check pod status
kubectl get pods -n todo-app

# Check services
kubectl get svc -n todo-app

# Check ingress
kubectl get ingress -n todo-app

# View pod logs
kubectl logs -f deployment/todo-backend -n todo-app
kubectl logs -f deployment/todo-frontend -n todo-app
```

## Scaling

```bash
# Scale backend
kubectl scale deployment todo-backend --replicas=3 -n todo-app

# Scale frontend  
kubectl scale deployment todo-frontend --replicas=5 -n todo-app
```

## Troubleshooting

### Common Issues

1. **ImagePullBackOff**: Check if Docker images exist and are accessible
2. **CrashLoopBackOff**: Check pod logs for application errors
3. **Service Unavailable**: Verify service selectors match pod labels

### Debug Commands

```bash
# Describe failing pods
kubectl describe pod <pod-name> -n todo-app

# Get pod logs
kubectl logs <pod-name> -n todo-app

# Execute into pod
kubectl exec -it <pod-name> -n todo-app -- /bin/sh

# Port forward for testing
kubectl port-forward svc/todo-frontend-service 8080:80 -n todo-app
kubectl port-forward svc/todo-backend-service 8081:8080 -n todo-app
```

## Production Considerations

1. **Resource Limits**: Adjust CPU/memory limits based on your workload
2. **Storage**: Configure appropriate storage class for PVC
3. **Security**: 
   - Use non-root containers
   - Configure network policies
   - Enable Pod Security Standards
4. **Monitoring**: 
   - Install Prometheus/Grafana
   - Configure alerting
5. **Backup**: 
   - Backup persistent volumes
   - Export database regularly

## Cleanup

```bash
# Delete all resources
kubectl delete -f k8s/

# Or delete namespace (removes everything)
kubectl delete namespace todo-app
``` 