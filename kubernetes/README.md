# Deployment

Consul is deployed via the `stable/consul` Helm chart.

```bash
helm install --name my-release stable/consul --set ui.enabled=true,uiService.enabled=true
```

A unified Helm chart is on the way, but for now deploy the `ingress.yml` and `app.yml`
manifests.
Replace the Load Balancer Service's`spec.loadBalancerIP` with your own static IP address.
Replace the Ingress's  `spec.rules[].backend.serviceName` for your Helm release
