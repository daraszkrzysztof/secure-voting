[![Build Status](https://travis-ci.org/daraszkrzysztof/secure-voting.svg?branch=master)](https://travis-ci.org/daraszkrzysztof/secure-voting.svg?branch=master)

# secure-voting
Application for secure voting. Provides security, anonymity, ballot casting assurance, verifiability with functionality of multiple trustee.

## building
```
go install github.com/daraszkrzysztof/secure-voting
```

## testing
```
go test -v  github.com/daraszkrzysztof/secure-voting
``

## installation

Building volume:
```
docker volume create --name mongo-db-volume
```
Creating bridge network:
```
docker network create -d bridge secure-voting-net
```
Building image:
```
docker build -t kdarasz/secure-voting .
```

Starting container:
```
docker run -d -p 8080:8080  kdarasz/secure-voting .
```
- with changed port number :
```
docker run -d -e SECURE_VOTING_PORT=8999 -p 8999:8999 kdarasz/secure-voting .
```
- running application with database:
```
sudo docker-compose up -d
```

## testing

Testing functionality:
```
curl -X PUT --verbose -H "Content-Type: application/json" -H "organizer-password: abc123" http://localhost:8080/secure-voting/organizer/orgXY
```