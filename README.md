# kube-config [![Build Status](https://travis-ci.org/apptio/kube-config.svg?branch=master)](https://travis-ci.org/apptio/kube-config) <!-- omit in toc -->

`kube-config` is a simple command line tool to authenticate and generate kubernetes configurations for your clusters

At its core, kube-config is based on the [reference documentation of a dex client](https://github.com/dexidp/dex/blob/master/Documentation/using-dex.md), that heavily utilizes [go-oidc](https://godoc.org/github.com/coreos/go-oidc) and golang package [oauth2](https://godoc.org/golang.org/x/oauth2)

- [Installing](#installing)
    - [Brew](#brew)
    - [Manually](#manually)
- [Usage](#usage)
    - [Config](#config)
    - [Terminal](#terminal)
      - [Simple Usage](#simple-usage)
      - [Advanced Usage](#advanced-usage)
- [Building](#building)

#### Key features <!-- omit in toc -->

-   Tier-based approach to cluster configs and authentication (`-t`) to allow for varying environments (i.e. `dev`, `prod`)
-   Ability to directly select current-context cluster, namespace and username (`-c`, `-n`, `-u` resp.) 
-   Specify config directory location with (`-o`)

Upon execution, `kube-config` will:

-   Open up the user's default web browser to localhost, and prompt the user to login
-   Pass all the login information through to [Dex](https://github.com/coreos/dex) which will generate a Kubernetes token
-   Generate a kubernetes bearer token which can be used to authenticate to the cluster
-   Create a kubeconfig file at the default `~/.kube/config.d/` or specified location

## Installing

#### Brew

```bash
$ brew tap apptio/tap
$ brew update
$ brew install kube-config
```

#### Manually

```bash
$ go version # Ensure you have go at least
go version go1.12.6 linux/amd64
$ git clone https://github.com/apptio/kube-config
$ cd kube-config
$ GOBIN=$HOME/go/bin go install
$ mkdir -p ~/.kube-config
$ cp ./kube-config.yml ~/.kube-config/
```

> Note: Ensure `$HOME/go/bin/` is added to your path like so: `PATH="${PATH}:${HOME}/go/bin"` at the end of your `~/.bashrc` or `~/.zshrc`

## Usage

#### Config

There are a lot of options for kube-config, but the defaults should be fine for 99% of use cases. All you need is a YAML configuration file at `$HOME/.kube-config/kube-config.yml` that describes your various clusters and tiers:

```yaml
client-secret: <your_dex_secret_here>
clusters:
- name: cluster1
  address: cluster1.example.com
  tier: dev
- name: cluster2
  address: cluster2.example.com
  tier: dev
```

#### Terminal

##### Simple Usage

Once you have your configuration file, simply run the tool:

    # defaults to "--tier=dev"
    kube-config 

If you want to write your `KUBECONFIG` to a specific directory, use `--output|-o` and specify the path:

    # defaults to outputting "<tier>-config.conf" in the specified dir
    kube-config --output /home/lbriggs/.kube/config.d

##### Advanced Usage

`kube-config` flags from the `--help` output:

    Flags:
        --caservername string   Servername for CA service
    -c, --cluster string        Cluster to use as current-context (default is first in config)
        --debug                 Print all request and responses from the OpenID Connect issuer
    -h, --help                  help for kube-config
        --list-tiers            If specified, the program will list the available tiers and then exit
    -n, --namespace string      Set your default namespace across ALL contexts in the tier
    -k, --no-verify-ssl         If specified, disable SSL cert checking (WARNING: UNSAFE)
    -o, --output string         Path to write Kubeconfig file (default "/home/lbriggs/.kube/config.d")
    -g, --print-groups          If specified, will print out the AD groups you belong to
    -t, --tier string           Tier to authenticate for (default "dev")
    -u, --username string       Username for login


Run kube config with a tier and cluster target:

```bash
$ kube-config --tier=dev --cluster=cluster1   # OR
$ kube-config -t dev -c cluster1
```

Targeting a specific tier, cluster and namespace:

```bash
$ kube-config --tier=dev --cluster=cluster1 --namespace=kube-system   #OR
$ kube-config -t dev -c cluster1 -n kube-system
```

> Note: `kubectx` is a great tool to change contexts quickly, but it's covenient sometimes to just start where you want with `kube-config` if you haven't yet auth'd for the day.
>
> Note: Without `-c / --cluster` `kube-config` will select the first cluster in the config for your `current-context`.

## Building

We use go itself for dependencies, so simply download them by running:

```bash
cd kube-config
go mod download
```

and then build using the standard GO build tools:

```bash
go build -o kube-config .
```
