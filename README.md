<div align="">
  <h1> Bhchain </h1>
</div>




Bhex Chain is powered by bluehelix, the next-generation decentralized custody and clearing technology.

<img width="749" height="528" src="https://static.hcdncn.com/hbtcchain/static/img/FM.5af890d6.png"/>

![GitHub](https://img.shields.io/badge/License-Apache2.0-brightgreen)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bluehelix-chain/bhchain)
![LoC](https://tokei.rs/b1/github/bluehelix-chain/bhchain)


## Home Page

https://www.bhexchain.com


## Building

    git clone https://github.com/bluehelix-chain/bhchain
    cd bhchain
    
    #build bhcd for running node
    go build ./cmd/bhcd
    
    #build bhcli for cli tool
    go build ./cmd./bhcli
## Quick Start

https://docs.bhexchain.com/guide/quick-start.html

## RUN&CONFIG
* https://docs.bhexchain.com/guide/node.html#id7

## Documentation

* ENGLISH https://docs.bhexchain.com/en
* 中文 https://docs.bhexchain.com

## Dex
https://hdex.bhexchain.com/swap

## Wallet

android & ios https://wallet.bhexchain.com

web https://chrome.google.com/webstore/detail/metamask/nkbihfbeogaeaoehlefnkodbefgpgknn 

## Chain Explorer
* https://explorer.bhexchain.com/index



## Security Suggestion
* Backup your config and data regularly
* Do not expose any port to the public network except p2p port (default 26656)


## to add support to chain X
Add package X inside package chainadaptor and implement the interface chainadaptor.ChainAdaptor (embed the fallback.ChainAdaptor and override specific methods).

Provide NewXChainAdaptor factory method and register it in the dispatcher
