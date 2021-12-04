# ArtWork - Hyperledger Fabric Project

[PPT](https://docs.google.com/presentation/d/1R_YuVRk6d4nVjM2tsvUI4r1VCIgzRnZ2/edit?rtpof=true&sd=true) / [Youtube](https://www.youtube.com/watch?v=3GxkCB55qPY)

## pre-condition

- curl, docker, docker-compose, go, nodejs
- hyperledger fabric-docker images are installed
- GOPATH are configured
- hyperledger bineries are installed (cryptogen, configtxgen ... etcs)

## -network

1. generating crypto-config directory, genesis.block, channel and anchor peer transactions

   - cd network
   - ./generate.sh

2. starting the network, create channel and join

   - ./start.sh

3. chaincode install, instsantiate and test(invoke, query, invoke)

   - ./cc.sh
   - If you want to test => ./cc_test.sh

## -prototype (turn on prototype)

- cd ../prototype

1. nodejs module install

   - npm install

2. certification works

   - node enrollAdmin.js
   - node registerUser.js

3. server start

   - node server.js

4. open web browser and connect to localhost:8080

## -application (turn on application)

- cd ../application
- same with prototype method
