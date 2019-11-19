#!/usr/bin/env bash

export KEY_FILE=$HOME/certificati/ticketag/ovh_cloud/key.pem
export GOPATH=$HOME/go

#scp -i $KEY_FILE $GOPATH/src/github.com/ticketag/docker-stopper/docker-stopper ubuntu@51.68.115.215:/home/ubuntu/
scp -i $KEY_FILE $GOPATH/src/github.com/ticketag/docker-stopper/restart_server/restart-server ubuntu@51.68.115.215:/home/ubuntu/

#curl localhost:30001