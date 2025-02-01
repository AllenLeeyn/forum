#!/bin/bash

docker build -t forum-img .
docker create --name forum-container -p 8080:8080 forum-img
docker images
docker ps -a