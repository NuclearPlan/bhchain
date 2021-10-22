
<div align="">
  <h1> Bhchain </h1>
</div>

**This is the repo of the Bhex Chain's main program. Everything starts from this repo.  You are in the right place if you want to run the validator/settlement, learn how the chain works or contribute codes.** 


Bhex Chain is powered by bluehelix, the next-generation decentralized custody and clearing technology.

<img width="749" height="528" src="https://static.hcdncn.com/hbtcchain/static/img/FM.5af890d6.png"/>

![GitHub](https://img.shields.io/badge/License-Apache2.0-brightgreen)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bluehelix-chain/bhchain)
![LoC](https://tokei.rs/b1/github/bluehelix-chain/bhchain)

## Introduction of Code's Dependencies
Although the code has been design to one command to run, but it do rely others repo:

 - [bhchain](https://github.com/bluehelix-chain/bhchain) (This Repo) : where compile the bhcd and bhcli, which are the main program of this chain and the cli to interactive with its RPC.  This Repo forks from Cosmos SDK. 
 - [chainnode](https://github.com/bluehelix-chain/chainnode) : it is a library that imported by both Bhchain and Bhsettle. This library contains the logic the interaction with others chain(btc,doge,heco and etc.)
 - [tendermint](https://github.com/bluehelix-chain/tendermint):  tendermint forks from offical repo. Just use as it.  
 - [settle](https://github.com/bluehelix-chain/settle):  this will compile the bhsettle program, which deployed in the core node., that play a role in the multi parity signature process. It utilised the P2P network underlaid by tendermint and depends the dsign library on cryptographic algorithm. 
 - [dsign](https://github.com/bluehelix-chain/dsign): it is a library that provide multi parties signature generation function both for ed25519 curve and secp256k1 curve. 
 - [ed25519](https://github.com/bluehelix-chain/ed25519): it is a library forked from the golang internal, exposed some internal functions for design. 


## How-to run your own BHEX Chain Full node

### 1. Building
    git clone https://github.com/bluehelix-chain/bhchain
    cd bhchain
    
    #build bhcd for running node
    go build ./cmd/bhcd
    
    #build bhcli for cli tool
    go build ./cmd./bhcli

### 2. Configure Genesis.json

Hit following command to generate config and your tendermint private key set. 

    ./bhcd init

Fetch config from https://github.com/bluehelix-chain/build/blob/master/mainnet/chainconfig/config/

Replace three files app.toml, config.toml and genesis.json of your local home directory at ~/.bhexchain/config/


### 3. Run the full node

Use following command to start fetch block from other blockchain node.

    ./bhcd start


## How-to protect your node not be hacked

1.  Backup your config and data regularly
2.  Do not expose any port to the public network except p2p port (default 26656)
3. If your are running settle node, please make sure you backed up *All of your private data*  regularly. (suggest backup per hrs, and keep at least 5 days' backup)
4. Make sure all of your ssh login will be audited, and the management privilege should be distributed very carefully. 


## Others Useful Links

 #### 1. Home Page
 
https://www.bhexchain.com


#### 2. Quick Start 

https://docs.bhexchain.com/guide/quick-start.html

#### 3. RUN&CONFIG
https://docs.bhexchain.com/guide/node.html#id7

#### 4. Documentation
ENGLISH https://docs.bhexchain.com/en
中文 https://docs.bhexchain.com

#### 5. Dex
https://hdex.bhexchain.com/swap

#### 6. Wallet

android & iOS https://wallet.bhexchain.com


#### 7.Chain Explorer
https://explorer.bhexchain.com/index


