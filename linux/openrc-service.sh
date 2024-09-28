#!/sbin/openrc-run
  
name=$RC_SVCNAME
command="/usr/bin/agent"
command_user="notebook:notebook"
pidfile="/var/run/agent.pid"
command_background="yes"

depend() {
  need net
}
