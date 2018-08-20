# Avengers, simulation tool for multi node environment

블록체인 코어 엔진 개발 시 consenesus 혹은 p2p network 를 구현함에 있어 다중 노드 환경을 다루는 많은 기능에 대한 구현이 필요합니다. Avengers 는 이러한 다중 노드 환경에서 노드 사이의 통신을 고루틴 사이의 채널링을 통해 가상으로 구현하여 다중 노드를 다루는 많은 기능들을 간편하게 테스트 할 수 있게 해줍니다.

## Getting started

```go
//initialize network manager
networkManager := mock.NewNetworkManager()

//make process
process := mock.NewProcess()
//initialize process
process.Init(processId)

//generate client and server and inject network manager function
client := mock.NewClient(processId, networkManager.GrpcCall)
server := mock.NewServer(processId, networkManager.GrpcConsume)

//register command handlers
server.Register("message.receive", func(){})

//add process to network manager
networkManager.AddProcess(process)

```



## Logical design

Avengers는 크게 network manager, process, rpc client, rpc server, components 의 5개의 파트로 구성됩니다.



**network manager** 는 각 프로세스 사이의 통신에 필요한 인프라를 구축해 주며, 각 프로세스가 서로의 존재를 모르는 상태에서 원활하게 원하는 프로세스에게 프로세스 아이디만을 통해 요청을 전달할 수 있게 해줍니다.



**process** 는 블록체인에서 하나의 노드를 나타내며 개발자가 필요한 함수 혹은 구현체를 등록하여 가상으로 노드내에서의 동작을 실행시켜 줍니다.



**rpc client** 는 it-chain 의 rabbitmq client 의 역할을 수행합니다. Production 환경에서는 각 노드별로 rabbitmq server가 구동하여 컴포넌트사이의 통신을 담당하고 다른 노드로 grpc 요청을 전달하는 경우에도 rabbitmq를 거치게 됩니다. 하지만 본 테스트 환경에서 각 노드를 담당하는 가상의 프로세스들은 자신만의 rabbit mq server 를 가질 수 없기때문에 rabbitmq server와 통신하는 부분을 mock을 이용하여 가려줄 필요가 있습니다. 여기서 rabbitmq client 부분을 **rpc client** 가 담당하고 server 부분을 **rpc server** 가 담당하게 됩니다. 



**components** 는 사용자가 process에 등록하는 많은 구현체이며 현재 버전에서는 grpc receive command를 처리하기 위한 handler로 제한됩니다.