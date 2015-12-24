FROM resin/armv7hf-debian:jessie

ENV QEMU_EXECVE 1
COPY . /usr/bin
