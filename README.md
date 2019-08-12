# kube-config

kube-config is a simple command line tool to generate a kubernetes configuration for Kubernetes clusters.

It will do the following upon execution:

  * Open up a web browser at localhost, and prompt the user to login
  * Pass all the login information through to [Dex](https://github.com/coreos/dex) which will generate a Kubernetes token
  * Generate a kubernetes bearer token which can be used to authenticate to the cluster
  * Create a kubeconfig file, either to stdout (default) or to a specified file

# Installing

### Brew

```bash
$ brew tap apptio/tap
$ brew update
$ brew install kube-config
```

### Manually

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

# Usage

### Config

There are a lot of options for kube-config, but the defaults should be fine for 99% of use cases. All you need is a configuration file at `$HOME/.kube-config/kube-config.yaml` which lists the clusters at Apptio - one has been provided for you, but it looks like this:

```yaml
clusters:
- name: cluster1
  address: cluster1.example.com
  tier: dev
- name: cluster2
  address: cluster2.example.com
  tier: dev
```

### Terminal

##### Simple Usage

Once you have your configuration file, simply run the tool:

```
$ kube-config
```

If you want to write your `KUBECONFIG` to a file, specify a path to an output file:

```
./kube-config --output /home/lbriggs/kubeconfig
```

##### Advanced Usage

`kube-config` flags from the `--help` output:

```
Flags:
  -c, --cluster string       Cluster to use as current-context (default is first in config)
      --debug                Print all request and responses from the OpenID Connect issuer
  -h, --help                 help for kube-config
      --list-tiers           If specified, the program will list the available tiers and then exit
  -n, --namespace string     Set your default namespace across ALL contexts in the tier
  -k, --no-verify-ssl        If specified, disable SSL cert checking (WARNING: UNSAFE)
  -o, --output string        Path to write Kubeconfig file (default "/home/lbriggs/.kube/config.d")
  -t, --tier string          Tier to authenticate for (default "dev")
      --caservername string  Override Servername used for verifying cluster CA server certificate
  -u, --username string      Username for login
```

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

> Note: Without `-c / --cluster` `kube-config` will select the first cluster in the config for your `current-context`.


# Building

