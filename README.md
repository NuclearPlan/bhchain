<div align="center">
  <h1> Bhchain </h1>
</div>

![GitHub](https://img.shields.io/badge/License-Apache2.0-brightgreen)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bluehelix-chain/bhchain)

Bhex Chain is powered by bluehelix, the next-generation decentralized custody and clearing technology.

<img width="749" height="528" src="https://static.hcdncn.com/hbtcchain/static/img/FM.5af890d6.png"/>

## Document
* EN https://docs.bhexchain.com/en/
* CN https://docs.bhexchain.com/


## Explorer
* https://explorer.bhexchain.com/index


## Dex
* https://hdex.bhexchain.com/swap


## Build
    go build ./cmd/bhcd/
    
## RUN&CONFIG
* https://docs.bhexchain.com/guide/node.html#id7


## Security Suggestion
* Backup your config and data regularly
* Do not expose any port to the public network except p2p port (default 26656)


## to add support to chain X
Add package X inside package chainadaptor and implement the interface chainadaptor.ChainAdaptor (embed the fallback.ChainAdaptor and override specific methods).

Provide NewXChainAdaptor factory method and register it in the dispatcher
