# Helmet

Helmet is a simple go script that aims to make managing multiple helm cli installations easier.

## Install

```
curl -o /usr/local/bin/helmet https://s3.amazonaws.com/darth-helmet/latest/macos/helmet && chmod +x /usr/local/bin/helmet
```

## Usage

```
Usage: helmet {ls,install} [VERSION]

Run with no sub commands, helmet symlinks the specified helm version to /usr/local/bin/helm.

SubCommands:
  - ls                      list installed versions
  - install VERSION         install the specified version

```

### Example

```
➜  ~ helmet ls
helm-2.12.3
helm-2.14.1
helm-2.14.3
helm-2.16.3
➜  ~ helmet install 3.1.2
Downloaded to /var/folders/m0/z4cbfm553_3c2hw1d8hzjjhc0000gn/T/helmenv185832043/helm-v3.1.2-darwin-amd64.tar.gz
Successfully copied helm excutable to /usr/local/bin/helm-3.1.2
Successfully installed helm 3.1.2
➜  ~ helmet ls
helm-2.12.3
helm-2.14.1
helm-2.14.3
helm-2.16.3
helm-3.1.2
➜  ~ helmet 3.1.2
Setting helm version: 3.1.2
➜  ~ helm version
version.BuildInfo{Version:"v3.1.2", GitCommit:"d878d4d45863e42fd5cff6743294a11d28a9abce", GitTreeState:"clean", GoVersion:"go1.13.8"}

```
