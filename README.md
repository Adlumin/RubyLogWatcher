# RubyLogWatcher

This Log watcher is developed to track and report "Fatal" type error's found in the Log file (for this purpose, generated out of Ruby code). 
It only scans the log files for event(entires) that have occurred in the past ONE hour (which is configurable via attribute "LogInterval" found in the config.json file along with the file name and multiple paths )

These found events are scrapped out of the log file and posted on to AWS Elastic Search (where developers with less or no privilage to the EC2 instance can review it) 

This Log watcher can be installed on any Linux OS
using "systemctl" type daemon setting,  which can be found in the folder: InstallUninstallScript

