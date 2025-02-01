#!/bin/bash

docker stop forum-container
docker rm forum-container
docker rmi forum-img
docker builder prune