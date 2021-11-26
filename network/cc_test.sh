#!/bin/bash

set -x

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addUser", "u1"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "u1"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addWork", "a1", "title1", "m1", "100", "100"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addWork", "a2", "title2", "m2", "100", "100"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "a2"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addWork", "a3", "title3", "m3", "100", "100"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addWork", "a4", "titl4", "m4", "100", "100"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addWork", "a5", "title5", "m5", "100", "100"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addUser", "u2"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "u2"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "a1"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getInfos", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["tradeProps", "myCompany", "u1", "a1", "30"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["tradeProps", "u1", "u2", "a1", "30"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["endTradeProps", "a1", "1000"]}' --peerAddresses peer0.org1.artwork.com:7051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getHistory", "u1"]}'
docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getHistory", "a1"]}'