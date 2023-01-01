#!/bin/bash

# Package
sudo yum update -y
sudo yum install -y curl git

# Timezone
sudo timedatectl set-timezone Asia/Tokyo

# Locale
sudo localectl set-locale LANG=ja_JP.utf8
sudo localectl set-locale LC_CTYPE=ja_JP.utf8
sudo localectl set-keymap jp106

# Nat Instance (https://docs.aws.amazon.com/ja_jp/vpc/latest/userguide/VPC_NAT_Instance.html)
sudo sysctl -w net.ipv4.ip_forward=1 | sudo tee -a /etc/sysctl.conf
sudo /sbin/iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
sudo yum install -y iptables-services
sudo service iptables save
sudo systemctl enable --now iptables

# Package Manager (asdf - https://asdf-vm.com/guide/getting-started.html#_1-install-dependencies)
sudo git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.10.2
sudo echo '. $HOME/.asdf/asdf.sh' >> ~/.bashrc
