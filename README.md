# RESTful Ethereum validator
RESTful API application written in Go. App interacts with Ethereum data providing information about block rewards and validator sync committee duties per slot.

## Install dependencies
```bash
brew install docker
```

## Building and Running the Application
```bash
git clone https://github.com/zdarovich/restful-ethereum-validator.git
cd restful-ethereum-validator
docker build --build-arg RPC_DIAL_URL=https://palpable-twilight-violet.quiknode.pro/78799c2f342c3cf5f6e726b4fab3e64bd3888cee -t ethereum-validator .
docker run -p 8080:8080 ethereum-validator
```

## Libraries
- Gin for quick REST API implementation
- BeaconClient generated from Swagger API spec https://ethereum.github.io/beacon-APIs
- Logrus for structured logging
- Envconfig for environment variables parsing
- go-ethereum client to distinguish contract vs wallet addresses

## Design
- Clients, services implement interfaces to allow easy testing and mocking which is SOLID
- main.go is the dependency injection entry point
- requests are lru cached to avoid hitting the same endpoint multiple times

## Endpoints
### GET /blockreward/{slot}
## Example
```bash
curl http://localhost:8080/blockreward/8920269
```

### GET /syncduties/{slot}
## Example
```bash 
curl http://localhost:8080/syncduties/8920269
```