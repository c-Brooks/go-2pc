# go-2pc

WIP - A simple implementation of Two-Phase Commit protocol.


## Architecture
```
  Each Physical Machine:

    |===============|
    | Consul Server |
    |===============|
      |
      |
      |
      |----  Consul Agent
      |    |=============|
      |    |    Inst 1   |
      |    |             |
      |    | gRPC Server |
      |    | gRPC Client |
      |    |=============|
      |
      |---------------|
      |               |
  Consul Agent      Consul Agent
|=============|    |=============|
|    Inst 2   |    |    Inst 3   |
|             |    |             |
| gRPC Server |    | gRPC Server |
| gRPC Client |    | gRPC Client |
|=============|    |=============|

```

The default config I will be using is:
- 3 Consul Servers
- 3 Application Instances for each server with 1 Consul agent each

The application instances communicate through peer-to-peer gRPC.

This is achieved by running a server on each instance and as many clients as there are other nodes.
The connections are stored in local memory as Consul K/V pairs as well as registering the addresses as Consul services.
That way, we don't have to open new connections if the k/v pair is there, but if not,
there is an address to dial in the service registry :)

We use Kubernetes to achieve this architecture. Each Pod contains an Application instance, and a Consul Agent.
Each physical Node also has a Consul server. This is achieved through
`spec.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution`

## Slides

https://slides.com/coreybrooks/consistency-in-distributed-systems
