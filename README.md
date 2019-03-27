# Lab device tester

Automate lab testing of the devices to make the results consistent and repeatable.

## UDP server

The UDP server listens for UDP packages from the test devices and logs each message received with a timestamp.

### Preparations

- Add the private key to ~/.ssh/
- Add a host entry for the server in ~/.ssh/config

        Host labdeviceserver
        HostName 13.53.172.78 # this will change
        User ubuntu
        IdentityFile ~/.ssh/labdeviceserver.pem

### Starting/stopping server

```bash
# start the EC2 instance
make start-udpserver

# cross-compile, copy the binary to the instance and start it
make deploy-udpserver

# stop the EC2 instance
make stop-udpserver
```
