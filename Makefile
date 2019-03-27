INSTANCE_ID := $(shell aws ec2 describe-instances --filters "Name=tag:Name,Values=labdeviceserver" --max-items 1 --region eu-north-1 | jq -r '.Reservations[0].Instances[0].InstanceId')

.PHONY: cross-compile-udpserver
cross-compile-udpserver:
	GOOS=linux GOARCH=amd64 go build -o udpserver-linux-amd64 ./cmd/udpserver/udpserver.go

.PHONY: deploy-udpserver
deploy-udpserver: cross-compile-udpserver
	scp udpserver-linux-amd64 labdeviceserver:
	ssh labdeviceserver "./udpserver-linux-amd64"

.PHONY: start-udpserver
start-udpserver:
	aws ec2 start-instances --instance-ids $(INSTANCE_ID) --region eu-north-1
	aws ec2 wait instance-running --instance-id $(INSTANCE_ID) --region eu-north-1

stop-udpserver:
	aws ec2 stop-instances --instance-ids $(INSTANCE_ID) --region eu-north-1
	aws ec2 wait instance-stopped --instance-id $(INSTANCE_ID) --region eu-north-1

