# RubyLogWatcher

This Log watcher is developed to report "Fatal" type error found in the Log file. 
It only scan the logs for event that have occured in the past one hour (which is 
configuraled via attribute "LogInterval" found in the config.json file )

This Log watcher developed using GoLang is installed into any Linux OS
using "Systemctl" type daemon setting
