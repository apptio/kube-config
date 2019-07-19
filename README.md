# kube-config

kube-config is a simple command line tool to generate a kubernetes configuration for Kubernetes clusters.

It will do the following upon execution:

  * Open up a web browser at localhost, and prompt the user to login
  * Pass all the login information through to [Dex](https://github.com/coreos/dex) which will generate a Kubernetes token
  * Generate a kubernetes bearer token which can be used to authenticate to the cluster
  * Create a kubeconfig file, either to stdout (default) or to a specified file


# Usage

There are a lot of options for kube-config, but the defaults should be fine for 99% of use cases. All you need is a configuration file at `$HOME/.kube-config.yaml` which lists the clusters at Apptio - one has been provided for you, but it looks like this:

```yaml
clusters:
- name: cluster1
  address: cluster1.example.com
  tier: dev
- name: cluster2
  address: cluster2.example.com
  tier: dev
  certificate: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFM01URXhOakExTURVeE5Wb1hEVEkzTVRFeE5EQTFNRFV4TlZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTUxSClJTYmNwT1EzVWFFS0V6VEpBTHp2SXhOVkFURE9LRENYNFVpZExoT2xpUGhsbG9ORDZMbWVWK0tBTDZ4TFJOd2kKRWM3VjlJMmNrWHdCNjVQYXlPNVY3MzNEWDZ2RytYVFlUaVNybi9Bcm5aVzlsUkFjM3FWWmt5VU9aRnpsVDFndQo3eGx5cWwvY1Jrem1sRFQvM3Y0NXZ0eXdya0FZam54ek54cCtsUkxGNUIwWHVDdnpnbXZlbVlzbUNNWFZva2RjCmw0cEhKNzk0UWJEUGlLMTRPTkpTU0NvTTFrZTNBNFppU29sMkFzbzgxdGdXaFhzVTBiV050NFQ0Zk43ZjhRNDgKKzJGcTFiY2xza215T29nc2c2cDJ3aXVwYjVmMEw5Z01ReTRwei9QNFR1dk5rb0lVOVJCbFVXcWRKa2FCblVuawpWdHdvekM0bkM3RE5aZlJpeStrQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFMV2FRU2dSb3FXeFN2TnZONGtYaCszY2xLVlQKTFNGdkR0N0kyb1Y0UjVjeTBLYldIU0ZDZWI4SUgxbE5NeWs1WnVtOEFESXVQTkVBb3hNU0pPVnpqLzd0THI4Two5eXZlYkt1RTNlUWRRbTVua0hjRmRKZHlGNUZLZXNyNkpPUFdIYWxpMlpKVHFKREZzZ1lvMXdXb1F4K1NWeDBsCm5MbW1RQkpHUlRsTmY0dXdDd1BEdmppK1NIWlZrUHBPeVJsMzN5OXJiZ0pTdFlzT29MRTVDc1VYdDZMUUhNdkwKSkRQNktPTGZxQUhDUXhVNmU2N2kzMkx6MXpXM1M5R0xpdHAwMXRyVEt0TmRnSStaWHpjcXFWQ1gwN0JvM0szaQpvdTloZHh2ZzE4OXQ2WW1yRjdjamt3Smc0YWJ1c1diRXh2YmdBTFBQa2lsd3A0SnZlMWNkY2RBQ2tIZz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
```

Note: certificate is an optional value - if that's not specified, kube-config will go and grab the certificate for you.

Once you have your configuration file, simply run the tool:

```
./kube-config
```

If you want to write your `KUBECONFIG` to a file, specify a path to an output file:

```
./kube-config --output /home/lbriggs/kubeconfig
```

# Building

