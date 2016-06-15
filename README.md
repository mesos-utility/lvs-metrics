# **lvs-metrics** <sup><sub>_lvs metrics collector for open-falcon_</sub></sup>
[![Build Status](https://travis-ci.org/mesos-utility/lvs-metrics.svg?branch=master)](https://travis-ci.org/mesos-utility/lvs-metrics)

lvs metrics collector for open-falcon.

## Dependencies

* [libnl][]
* [ipvs][]
* [netlink][]
* [toolkits][]
* [glog][]

Dependencies are handled by [godep][], simple install it and type `godep restore` to fetch them.

## Install

```console
$ git clone https://github.com/mesos-utility/lvs-metrics.git
$ cd lvs-metrics
$ make bin
```

```console
# sudo yum install -y libnl3.x86_64
or
# sudo apt-get install libnl-3-dev libnl-genl-3-dev
```

## Metrics
| Counters | Notes|
|-----|------|
|lvs.in.bytes|network in bytes per host|
|lvs.out.bytes|network out bytes per host|
|lvs.in.packets|network in packets per host|
|lvs.out.packets|network out packets per host|
|lvs.total.conns|lvs total connections per vip now|
|lvs.active.conn|lvs active connections per vip now|
|lvs.inact.conn|lvs inactive connections per vip now|
|lvs.realserver.num|lvs live realserver num per vip now|
|lvs.vip.conns|lvs conns counter from service start per vip|
|lvs.vip.inbytes|lvs inbytes counter from service start per vip|
|lvs.vip.outbytes|lvs outpkts counter from service start per vip|
|lvs.vip.inpkts|lvs inpkts counter from service start per vip|
|lvs.vip.outpkts|lvs outpkts counter from service start per vip|


[libnl]: https://www.infradead.org/~tgr/libnl
[ipvs]: https://github.com/google/seesaw/ipvs
[netlink]: https://github.com/google/seesaw/netlink
[toolkits]: https://github.com/toolkits
[glog]: https://github.com/golang/glog
[godep]: https://github.com/tools/godep
