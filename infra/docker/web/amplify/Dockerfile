FROM --platform=linux/amd64 amazonlinux:2023

RUN yum update -y
RUN yum install -y \
  openssh \
  git \
  nodejs

RUN npm install -g yarn
