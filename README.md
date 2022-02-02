# Introduction 
This repo contains example code how to build a compute engine based on Colonies and Kubernetes. We are going to build a simple compute engine that calculates Fibonacci numbers. 

* Using the Colonies CLI tool we are going to submit process specs to the Colonies Server. The process spec contains a Fibonnaci number which should be calculate. 
* The Colonies Server maintains a queue of incoming process specs (jobs). 
* Each Fibonacci Colony Apps/Workers connects to the Colonies Server and request a process spec, which it then executes. That means that all workers compete on being assigned a process spec. The Colonies Server ensure that only one worker get a certain process spec.    
* If a worker doesn't complete a task in certain time, the Colonies Server then move it back to the queue so that other workers can execute it.

![Compute Engine](docs/images/compute_engine.png?raw=true "Compute Engine")
