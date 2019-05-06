# iotex-payout

This is a tool that calculates the reward shares for voters at each epoch of
[IoTeX](https://github.com/iotexproject)

## Get started

### Minimum requirement

Please refer the minimum requirement for building IoTeX. Specifically,

| Components | Version | Description |
|----------|-------------|-------------|
| [Golang](https://golang.org) | &ge; 1.11.5 | Go programming language |
| [Dep](https://golang.github.io/dep/) | &ge; 0.5.0 | Dependency management tool, required only when you update dependencies |

### IoTeX compatibility

Iotex-payout currently supports iotex-core v0.5.0-rc8.

### Build from code

Check out code from
```
mkdir -p ~/go/src/github.com/
cd ~/go/src/github.com/
git clone https://github.com/Infinity-Stones/iotex_payout.git
cd iotex_payout
```

Build the tool by
```
dep ensure [--vendor-only]
go build
```

### Run iotex payout
```
iotex_payout [arguments]
```


### Run iotex payout
Build the container
```
docker build -f Dockerfile . -t iotexpayout
```

Run the container interactively
```
docker run -ti iotexpayout bash
```

Use ioctl to setup operator account
```
ioctl account import operator
```

Run iotex_payout 
```
iotex_payout delegate operator -e 100 -b 95 -p 90 -f 80
```

