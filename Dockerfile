FROM resin/amd64-debian:jessie

RUN apt-get update && apt-get install -y golang curl python-pip

# Install AWS CLI
RUN pip install awscli

COPY . /
