#!/bin/sh
set -e

mkdir /var/hallucinet/
touch /var/hallucinet/hosts
/bin/coredns -conf /etc/hallucinet/Corefile &
/bin/monitor
