# Introduction 
This repo contains example code how to build a compute engine based on Colonies and Kubernetes. We are going to build a simple compute engine that calculates Fibonacci numbers. 

* Using the Colonies CLI tool we are going to submit process specs to the Colonies Server. The process spec contains a Fibonnaci number that should be calculate by a Colony App/Worker running inside a Kubernetes pod. 
* The Colonies Server maintains a queue of incoming process specs (jobs). 
* Each Fibonacci worker connects to the Colonies Server and request a process spec, which it then executes. That means that all workers compete on being assigned a process spec. The Colonies Server ensure that only one worker get a certain process spec. If a worker doesn't complete a task in certain time, the Colonies Server then move the process spec back to the queue so that other workers can execute it.
* Using Kubernetes, it possible to dynamically scale number of workers (pods) up and down. If the compute engine is down-scaled and a certian worker is destroyed before it finish, the Colonies Server will then move the process spec back to the queue as mentioned before. 

![Compute Engine](docs/images/compute_engine.png?raw=true "Compute Engine")

# Colony App/Worker

# Kubernetes 
The Yaml below contains deployment code for Kubernetes. Note the **preStop** lifecycle hook. It is used to unregister a worker when the pods is destroyed. This is done by making sure the Golang code above stores its runtime id in a file, which can then be fetched to unregister the worker. The **preStop** is always called when the pod is undeployed.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fibonacci-deployment
  labels:
    app: fibonacci
spec:
  replicas: 1
  selector:
    matchLabels:
      app: fibonacci
  template:
    metadata:
      labels:
        app: fibonacci
    spec:
      containers:
      - name: fibonacci
        image: johan/fibonacci_k8s
        command:
            - "go"
            - "run"
            - "solver.go"
        lifecycle:
          preStop:
            exec:
              command: ["go","run","solver.go", "unregister"]
        resources:
            requests:
              memory: "1000Mi"
              cpu: "1000m"
            limits:
              memory: "1000Mi"
              cpu: "1000m"
        env:
        - name: COLONYID
          value: "d03b4b236a479622ee5542e4ead7d254315b557bba74391511b5942e3a05bffd"
        - name: COLONYPRVKEY
          value: "ccf0cdd308b62add43fb2555410838948dfbcc8b0f186550f6098470f20e6108"
        - name: COLONIES_SERVER_HOST
          value: "10.0.0.240"
        - name: COLONIES_SERVER_PORT
          value: "8080"
        - name: CORES
          value: "1"
        - name: MEM
          value: "1000"
```

```console
kubectl apply -f deployment.yaml -n test
```

Note that a namespace named **test** must be created before calling the command above. The scale-up, just increase the replicas field in the Yaml and re-apply the deployment.yaml file using the command above. Also note that all Colony ids and keys need to be created, and the Yaml file need to updated with correct keys before testing the code in this repo.

# Submitting a process spec
```json
{
    "conditions": {
        "runtimetype": "fibonacci"
    },
    "env": {
        "fibonacciNum": "10"
    }
}
```

```console
colonies process subitt --spec process.json
```

# Checking the queue
```console
colonies process psw --count 5 
```
Output:
```
+------------------------------------------------------------------+---------------------+
|                                ID                                |   SUBMISSION TIME   |
+------------------------------------------------------------------+---------------------+
| 2c42fbe3d729b3d145fc3288ccc785f0accd8c07aea65178854d6dcdb18a080f | 2022-02-02 14:37:36 |
| 63eefd3325f64c5866debc71171a754ef3de38f44db553eaf1875080014ee300 | 2022-02-02 14:37:36 |
| b4ab1fb17873fac30b86b5ad7f284e977e9e6a55719247cec619b92753de21cc | 2022-02-02 14:37:36 |
| 79d0f68d42d85df307ad0774dc01ae1a722fe85a996b28009356d9f7936e59c2 | 2022-02-02 14:37:35 |
+------------------------------------------------------------------+---------------------+
```
