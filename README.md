# Sticker - SymbolTicker

The best(maybe) app to fetch a history of specific stock.

## Build the image

Build the docker image:

```console
docker build -t sticker .
```

Then start the container:

```console
docker run -e SYMBOL=<stock's symbol> \
           -e NDAYS=<number of days> \
           -e APIKEY=<YOUR API KEY> \
           -p 8080:8080 \
           sticker
```

The service is available at http://localhost:8080

OR just use github image:

```console
docker run -e SYMBOL=<stock's symbol> \
           -e NDAYS=<number of days> \
           -e APIKEY=<YOUR API KEY> \
           -p 8080:8080 \
           ghcr.io/ggrafu/sticker:latest
```

## Setup Minikube (macOS)

Install minikube per instructions: https://minikube.sigs.k8s.io/docs/start/

Install `kubectl`: https://kubernetes.io/docs/tasks/tools/

Start minikube:

```console
minikube start
```

Create configmap:

```console
kubectl apply -f configmap.yaml
```

Update `secret.yaml` with base64 encoded value of APIKEY. Then apply manifest:

```console
kubectl apply -f secret.yaml
```

Create deployment:

```console
kubectl apply -f deployment.yaml
```

Create service to access the app:

```console
kubectl apply -f service.yaml
```

Verify that pod is running:

```console
kubectl get pods

NAME                      READY   STATUS    RESTARTS      AGE
sticker-6f49ccf8f-bdb9f   1/1     Running   1 (1h ago)   1h
```

To access the service you need to expose the port using `minikube service` command:

```console
minikube service sticker --url

http://127.0.0.1:59268
```

Keep that terminal window open and use URL from the output to access the service(your port number may differ):

```console
curl http://127.0.0.1:59268/v1/data

{"values":[236.96,238.73,238.19,244.43,241.8,240.45,244.69,249.01,257.22,256.92,252.51,245.42,247.4,244.37,245.12],"average":245.54799}
```

## Setup ingress on minikube (macOS)

In order to use ingress controller the minikube addon needs to be enabled:

```console
minikube addons enable ingress
```

On macOS addon requres the usage of the tunnel. Run the following command and keep terminal window open:

```console
minikube tunnel
```

Create ingress:

```console
kubectl apply -f ingress.yaml
```

In the terminal window of minikube tunnel you may be asked to provide admin password to get access to ports 80/443. After allowing the access you can test ingress:

```console
curl http://localhost/sticker/v1/data

{"values":[236.96,238.73,238.19,244.43,241.8,240.45,244.69,249.01,257.22,256.92,252.51,245.42,247.4,244.37,245.12],"average":245.54799}
```

## Scalability notes

By using local cache, application significantly reduces the load on source API, however after scaling up to more then 5 replicas some of the pods will inevitably hit the source API's rate limit. To solve this problem the usage of shared cache is suggested. The following PR contains required changes: https://github.com/ggrafu/sticker/pull/1

After adding shared redis cache the application can scale up reasonably more. Keep in mind that the changes don't include HA mode for redis. In real application the redis needs to be deployed in cluster mode.
