#!/bin/sh
set -e

mkdir /var/hallucinet/ -p
touch /var/hallucinet/hosts
/bin/coredns -conf /etc/hallucinet/Corefile &
/bin/monitor
