#!/sbin/openrc-run
  
name=$RC_SVCNAME
description="Notebook Agent"
supervisor="supervise-daemon"
command="/usr/local/bin/agent"
command_user="notebook:notebook"
pidfile="/run/agent.pid"
command_background="yes"

depend() {
  after net.eth0
}
