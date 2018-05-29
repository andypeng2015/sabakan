Getting started
===============

* [Setup sabakan](#setupsabakan)
  * [Prepare etcd](#prepareetcd)
  * [Install sabakan and sabactl](#installsabakan)
  * [Run sabakan](#runsabakan)
  * [Configure sabakan](#configuresabakan)
* [Netboot](#netboot)
* [Test](#test)
* [What's next](#whatsnext)

## <a name="setupsabakan" /> Setup sabakan

### <a name="prepareetcd" /> Prepare etcd

Sabakan requires [etcd][].  Install and run it at somewhere.

### <a name="installsabakan" /> Install sabakan and sabactl

Install `sabakan` and `sabactl`:

```console
$ go get -u github.com/cybozu-go/sabakan/cmd/sabakan
$ go get -u github.com/cybozu-go/sabakan/cmd/sabactl
```

`sabakan` Docker image is so available at [quay.io/cybozu/sabakan](https://quay.io/cybozu/sabakan)

### <a name="runsabakan" /> Run sabakan

```console
$ sabakan -etcd-servers http://etcd-host:2379
```

### <a name="configuresabakan" /> Configure sabakan

First of all, prepare JSON files

- ipam.json
```json
{
   "max-nodes-in-rack": 28,
   "node-ipv4-pool": "10.69.0.0/20",
   "node-ipv4-range-size": 6,
   "node-ipv4-range-mask": 26,
   "node-index-offset": 3,
   "node-ip-per-node": 3,
   "bmc-ipv4-pool": "10.72.16.0/20",
   "bmc-ipv4-range-size": 5,
   "bmc-ipv4-range-mask": 20
}
```

Read [ipam](ipam.md) if you want to know meaning of each parameter.

- dhcp.json
```json
{
   "gateway-offset": 100,
   "lease-minutes": 120
}
```

Read [dhcp](dhcp.md) if you want to know meaning of each parameter.

Use `sabactl` to configure `sabakan`. 

```console
$ sabactl ipam set -f ipam.json
$ sabactl dhcp set -f dhcp.json
```

Make sure current configuration.

```console
$ sabactl ipam get
$ sabactl dhcp get
```

Each output will be the same as [above JSON](#configuresabakan)

## <a name="netboot" /> Netboot

**ToDo**

## <a name="test" /> Test

**ToDo**

## <a name="whatsnext" /> What's next

Learn sabakan [concepts](concepts.md), then read other specifications.