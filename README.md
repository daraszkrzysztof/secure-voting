# secure-voting
Application for secure voting. Provides security, anonymity, ballot casting assurance, verifiability with functionality of multiple trustee.

## installation

Building image:
```
sudo docker build -t kdarasz/secure-voting .
```

Starting container:
```
sudo docker run -d -p 8080:8080  kdarasz/secure-voting .
```

## testing

```
curl -X PUT --verbose -H "Content-Type: application/json" -H "organizer-password: zaq12wsx" http://localhost:8080/secure-voting/new-organizer/org1
```