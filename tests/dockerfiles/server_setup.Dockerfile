FROM ubuntu:22.04

RUN apt-get update && apt-get install -y curl git tar python3

COPY scripts/server_setup.py /root/server_setup.py

WORKDIR /root