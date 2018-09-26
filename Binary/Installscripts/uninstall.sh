#!/bin/bash
sudo kill -9 $(ps aux | grep './logwatcher' | awk '{print $2}')
sudo systemctl stop adluminslogwatcher.service
sudo systemctl disable adluminslogwatcher.service
sudo rm /etc/systemd/system/adluminslogwatcher.service
sudo rm -rf /home/ubuntu/productionlogwatcher
sudo systemctl daemon-reload