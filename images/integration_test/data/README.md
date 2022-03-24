# Testing Data for Sigma Engine

## Description

Each folder contains the raw sample data of the log type as provided by the elastic/beats repository on GitHub.

## Notes

- Netflow module test data was in binary format so it was omitted
- Using example data from another project for winlog.sysmon data see if we can find example sysmon json from another source to test with.
- Using my own generated CDNS data, hard to find examples for it in the wild
- The topic name can be generated from traversing the directory and using the folder names

## Raw Data Ingest topics

Raw data will be expected to be produced on the following topics:

- cdns.cdns
- checkpoint.firewall
- cisco-custom.asa
- cisco-custom.ftd
- juniper.junos
- juniper.netscreen
- netflow.log
- panw.panos
- system.syslog
- winlog.sysmon
- zeek.conn
- zeek.dce_rpc
- zeek.dhcp
- zeek.dnp3
- zeek.dns
- zeek.dpd
- zeek.files
- zeek.http
- zeek.intel
- zeek.irc
- zeek.kerberos
- zeek.modbus
- zeek.mysql
- zeek.notice
- zeek.ntlm
- zeek.pe
- zeek.radius
- zeek.rdp
- zeek.rfb
- zeek.sip
- zeek.smb_cmd
- zeek.smb_files
- zeek.smb_mapping
- zeek.smtp
- zeek.snmp
- zeek.socks
- zeek.ssh
- zeek.ssl
- zeek.traceroute
- zeek.tunnel
- zeek.weird
