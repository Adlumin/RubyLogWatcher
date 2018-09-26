#!/bin/bash
######################################################################
# Install Adlumin Ruby Log Watcher.
#
# Copyright (c) 2018 Adlumin Inc
# All rights reserved.
# https://www.adlumin.com/
######################################################################

function makefolder() {
	sudo mkdir -p /opt/productionlogwatcher
	sleep 1
}

function transfer_binary() {
	sudo cp logwatcher /opt/productionlogwatcher
	sudo chown ubuntu:ubuntu /opt/productionlogwatcher/logwatcher
	sudo chmod 0775 /opt/productionlogwatcher/logwatcher
	sleep 1
}

function transfer_watchstart() {
	sudo cp watchmgr /opt/productionlogwatcher
	sudo chown ubuntu:ubuntu /opt/productionlogwatcher/watchmgr
	sudo chmod 0775 /opt/productionlogwatcher/watchmgr
	sleep 1
}


function transfer_systemddsetup() {
	sudo cp adluminslogwatcher.service /etc/systemd/system/
	sudo chown ubuntu:ubuntu /etc/systemd/system/adluminslogwatcher.service
	sudo chmod 0664 /etc/systemd/system/adluminslogwatcher.service
	sleep 2
}

function reload_systemd_and_start_adluminslogwatcher() {
	sudo systemctl daemon-reload
	sudo systemctl enable adluminslogwatcher.service
	sleep 2
	sudo systemctl start adluminslogwatcher.service
}

## MainScript
makefolder
transfer_binary
transfer_watchstart
transfer_systemddsetup
reload_systemd_and_start_adluminslogwatcher

########################################################
# Installer to crosscheck Installation was successful  #
########################################################
sudo systemctl status adluminslogwatcher
