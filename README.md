# Kubernetes 12-Factor Demo

## Build Yourself (optional)

If you just want to run the container, you can skip to "Run Container" part.

### Build go app for linux

Prerequisite: you have Golang installed on your local machine. For Ubuntu, you can follow this article: https://medium.com/better-programming/install-go-1-11-on-ubuntu-18-04-16-04-lts-8c098c503c5f. 

You are also assumed that you have installed all dependency needed for this proect.

To build the Golang application, navigate to root project folder, and run this command:

``` sh
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o web ./main.go
```

### Build/push container

``` sh
sudo docker build -t asatrya/k8s-12-factor-demo .
sudo docker push asatrya/k8s-12-factor-demo
```

## Run Container Locally

This is the docker image: https://hub.docker.com/r/asatrya/k8s-12-factor-demo

To run container on port 5005, execute this command:

``` sh
sudo docker run -p 5005:5005 --rm asatrya/k8s-12-factor-demo
```

Access http://localhost:5005/ on your browser, then it will show a simple screen displaying hostname and environment variable like this

![](screenshot.png)

## Deploy on Kubernetes Cluster

Prerequisite: you have succesfully connect your `kubectl` client with your cluster. Make sure this is working well by running 

``` sh
kubectl get nodes
```

it should display list of nodes in your cluster.

To deploy your container to Kubernetes cluster, run this command:

```sh
kubectl -f kubernetes/demo-configmap.yaml apply
kubectl -f kubernetes/demo-secret.yaml apply
kubectl -f kubernetes/demo-deployment.yaml apply
kubectl -f kubernetes/demo-service.yaml apply
```

## Resize Cluster (GCP)

``` sh
gcloud container clusters resize standard-cluster-1 --node-pool default-pool --num-nodes 0
```

## Where is The 12-Factor?

### I. One Codebase

> One codebase tracked in revision control, many deploys

The source code, Dockerfile, and Kubernetes declaration files (insiede `kubernetes` folder) are all text files and can be versioned in single Git repository.

From this repo you can deploy to many version and environment. It's recommended to utilize a CI/CD pipeline to automate the build and deployment process.

### II. Dependencies

> Explicitly declare and isolate dependencies.

In `Dockerfile`, everything we need is already defined there and later will be build as a Docker image and run as Docker container. That means all dependency is already stated explicitly and isolated in a container.

### III. Config

> Store config in the environment

In `kubernetes/demo-deployments.yaml` line 18 you can see all configurations and secrets are defined as container's environment variables.

```yaml
          env:
            - name: ENV
              valueFrom:
                  configMapKeyRef:
                    name: demo-configmap
                    key: ENV
            - name: DB_HOST
              valueFrom:
                  configMapKeyRef:
                    name: demo-configmap
                    key: DB_HOST
            - name: DB_PORT
              valueFrom:
                  configMapKeyRef:
                    name: demo-configmap
                    key: DB_PORT
            - name: DB_USER
              valueFrom:
                  secretKeyRef:
                    name: demo-secret
                    key: DB_USER_BASE64
            - name: DB_PASSWORD
              valueFrom:
                  secretKeyRef:
                    name: demo-secret
                    key: DB_PASSWORD_BASE64
```

Configuration and secret values themself defined in Kubernetes' ConfigMap and Secret objects in `kubernetes/demo-configmap.yaml` and `kubernetes/demo-secret.yaml` files.

### VI. Processes

> Execute the app as one or more stateless processes

By declaring a Deployment object, it means we expect our container runs in a stateless Pod. When you run an application in a stateless environment, you cannot store any persisent data in local machine (such as memory or filesystem). This demo doesn't cover this example.

### IV. Backing services

> Treat backing services as attached resources

Backing service, i.e database, is treated as an attched service by defining the connection credentials in ConfigMap and Secret object (see `kubernetes/demo-configmap.yaml` and `kubernetes/demo-secret.yaml` files).

The connection configurations the will be injected when we run the container (see "III. Config" part above). Once the database fails, we can simply provide a new configuration for our new database and then recreate the container.

### V. Build, release, run

> Strictly separate build and run stages

The build process is executed through this command

``` sh
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o web ./main.go
```

and

``` sh
sudo docker build -t asatrya/k8s-12-factor-demo .
```

The release process is executed with this command:

```sh
sudo docker push asatrya/k8s-12-factor-demo
```

The run process is executed with this command

for Docker:

```sh
sudo docker run -p 5005:5005 --rm asatrya/k8s-12-factor-demo
```

for Kubernetes:

```sh
kubectl -f kubernetes/demo-configmap.yaml apply
kubectl -f kubernetes/demo-secret.yaml apply
kubectl -f kubernetes/demo-deployment.yaml apply
kubectl -f kubernetes/demo-service.yaml apply
```

We simply strictly separate build, release, and run process because when doing a process, we just do the process without compromising one with another.

### X. Dev/prod parity

> Keep development, staging, and production as similar as possible

By using container, we ensure the similarity of the OS, platform, configurations, and dependency across environments (i.e.: prod, test, staging, production).

### VII. Port binding

> Export services via port binding

The application itself is listening on port 5005, as can be seen on `main.go` file line 65:

```go
http.ListenAndServe(":5005", handlers.CompressHandler(router))
```

Then, we map port 5005 to port 5005 (in this case the port is the same, but it can be different) exposed by the container. See `kubernetes/demo-deployments.yaml` line 48

```yaml
          ports:
            - containerPort: 5005
```

Then, the Service object exposes as port 80 to the internet. See `kubernetes/demo-service.yaml` line 9

```yaml
  ports:
  - name: demo-http-port
    port: 80
    targetPort: 5005
```

### VIII. Concurrency

> Scale out via the process model

Scaling up/down can be done by defining number of replicas in `kubernetes/demo-deployments.yaml` line 9

```yaml
  replicas: 2
```

### IX. Disposability

> Maximize robustness with fast startup and graceful shutdown

In Kubernetes, we can check the health status of a Pod. When a Pod go down into unhealthy condition, the Service will stop sending it traffic or eventually destroy and recreate the Pod.

In this case, we use `livenessProbe` property to check whether our Pod is still alive. This can be check on `kubernetes/demo-deployment.yaml` file 

```yaml
          livenessProbe:
            httpGet:
              path: /healthz
              port: 5005
              httpHeaders:
                - name: Custom-Header
                  value: Awesome
            initialDelaySeconds: 60
            periodSeconds: 10
```

### XI. Logs

> Treat logs as event streams

In `main.go`, we simply print this log to `stdout`

```go
fmt.Println("Listening on port 5005...")
```

If you use managed Kubernetes service like GKE, this log content can be analyzed using Google Stackdriver.

### XII. Admin processes

> Run admin tasks as one-off processes

This demo doesn't cover this example.
