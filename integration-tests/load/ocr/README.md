# OCR Load tests

## Setup

These tests can connect to any cluster create with [chainlink-cluster](../../../charts/chainlink-cluster/README.md)

<<<<<<< HEAD
Create your cluster

```sh
kubectl create ns my-cluster
devspace use namespace my-cluster
=======
Create your cluster, if you already have one just use `kubefwd`
```
kubectl create ns cl-cluster
devspace use namespace cl-cluster
>>>>>>> 06656fac80999d1539e16951a54b87c6df13a9c7
devspace deploy
sudo kubefwd svc -n cl-cluster
```

Change environment connection configuration [here](../../../charts/chainlink-cluster/connect.toml)

If you haven't changed anything in [devspace.yaml](../../../charts/chainlink-cluster/devspace.yaml) then default connection configuration will work

## Usage

```sh
export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestOCRLoad
go test -v -run TestOCRVolume
```

Check test configuration [here](config.toml)