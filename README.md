[![Build Status](https://travis-ci.org/DNAProject/DNA.svg?branch=master)](https://travis-ci.org/DNAProject/DNA)

English | [中文](README_CN.md)

# DNA (Distributed Networks Architecture)

 DNA is a decentralized distributed network protocol based on blockchain technology and is implemented in Golang.
 Through peer-to-peer network, DNA can be used to digitize assets and provide financial service, including asset
 registration, issuance, transfer, etc.

## Highlight Features

 *	Scalable Lightweight Universal Smart Contract
 *	Crosschain Interactive Protocol
 *	Quantum-Resistant Cryptography (optional module)
 *	China National Crypto Standard (optional module)
 *	High Optimization of TPS
 *	P2P Link Layer Encryption
 *	Node Access Control
 *	Multiple Consensus Algorithm Support (DBFT/VBFT)
 *	Configurable Block Generation Time
 *	Configurable Digital Currency Incentive
 *	Configable Sharding Consensus (in progress)


# Building
The requirements to build DNA are:
 *	Go version 1.12.5 or later
 *	Properly configured Go environment
 
Clone the DNA repository into the appropriate `$GOPATH/src/DNAProject` directory.


```shell
$ git clone https://github.com/DNAProject/DNA.git

```

Build the source code with make.

```shell
$ make
```

After building the source code, you should see two executable programs:

* `dnaNode`: the node program

Follow the procedures in Deployment section to give them a shot!


# Deployment
 
To run DNA successfully, at least 4 nodes are required. The four nodes can be deployed in the following two way:

* multi-hosts deployment
* testmode deployment

## Configurations for multi-hosts deployment

 We can do a quick multi-host deployment by modifying the default configuration file `config.json`. Change the IP
 address in `SeedList` section to the seed node's IP address, and then copy the changed file to the hosts that you
 will run on.
 On each host, put the executable program `dnaNode` and the configuration file `config.json` into the same directory.
 Like :
 
```shell
$ ls
config.json dnaNode

```
 Each node also needs a `wallet.dat` to run. The quickest way to generate wallets is to run `./dnaNode account add -d`
 on each host.Then, change the `peerPubkey` and `address` field to the 4 nodes' wallet public keys, which you can get
 from the last command's echo. The public key sequence does not matter.
 Now all configurations are completed.
 
 Here's an snippet for configuration, note that `10.0.1.100` and `10.0.1.101` are public seed node's addresses:

```shell
$ cat config.json
{
  "SeedList": [
    "10.0.1.100:20338",
    "10.0.1.101:20338",
    "10.0.1.102:20338",
    "10.0.1.103:20338"
  ],
  "ConsensusType":"vbft",
  "VBFT":{
    "n":40,
    "c":1,
    "k":4,
    "l":64,
    "block_msg_delay":10000,
    "hash_msg_delay":10000,
    "peer_handshake_timeout":10,
    "max_block_change_view":3000,
    "admin_ont_id":"did:dna:AMAx993nE6NEqZjwBssUfopxnnvTdob9ij",
    "min_init_stake":10000,
    "vrf_value":"1c9810aa9822e511d5804a9c4db9dd08497c31087b0daafa34d768a3253441fa20515e2f30f81741102af0ca3cefc4818fef16adb825fbaa8cad78647f3afb590e",
    "vrf_proof":"c57741f934042cb8d8b087b44b161db56fc3ffd4ffb675d36cd09f83935be853d8729f3f5298d12d6fd28d45dde515a4b9d7f67682d182ba5118abf451ff1988",
    "peers":[
      {
        "index":1,
        "peerPubkey":"0289ebcf708798cd4c2570385e1371ba10bdc91e4800fa5b98a9b276eab9300f10",
        "address":"ANT97HNwurK2LE2LEiU72MsSD684nPyJMX",
        "initPos":10000
      },
      {
        "index":2,
        "peerPubkey":"039dc5f67a4e1b3e4fc907ed430fd3958d8b6690f4f298b5e041697bd5be77f3e8",
        "address":"AMLU5evr9EeW8G1WaZT1n1HDBxaq5GczeC",
        "initPos":10000
      },
      {
        "index":3,
        "peerPubkey":"0369f4005b006166e988af436860b8a06c15f3eb272ccbabff175e067e6bba88d7",
        "address":"AbSAwqHQmNMoUT8ps8N16HciYtgprbNozF",
        "initPos":10000
      },
      {
        "index":4,
        "peerPubkey":"035998e70d829eea58998ec743113cf778f66932a063efc1a0a0496717c4a0d93d",
        "address":"AemhQtcPTGegSk1UAsiLnePVcut1MLXSPg",
        "initPos":10000
      }
    ]
  }
}
```

## Configurations for testmode deployment

If you like to run in test mode, there's no configuration needed.
With the following command, you can start DNA in test mode.

```shell
$ ./dnaNode --testmode
$ - input your wallet password
```

## Getting Started

Start the seed node program first and then other nodes. Just run:

```shell
$ ./dnaNode
$ - input your wallet password
```

Run `./dnaNode --help` for more details.


# Contributing

Can I contribute patches to DNA project?

Yes! Please open a pull request with signed-off commits. We appreciate your help!

You can also send your patches as emails to the developer mailing list.
Please join the DNA mailing list or forum and talk to us about it.

Either way, if you don't sign off your patches, we will not accept them.
This means adding a line that says "Signed-off-by: Name <email>" at the
end of each commit, indicating that you wrote the code and have the right
to pass it on as an open source patch.

Also, please write good git commit messages.  A good commit message
looks like this:

	Header line: explain the commit in one line (use the imperative)

	Body of commit message is a few lines of text, explaining things
	in more detail, possibly giving some background about the issue
	being fixed, etc etc.

	The body of the commit message can be several paragraphs, and
	please do proper word-wrap and keep columns shorter than about
	74 characters or so. That way "git log" will show things
	nicely even when it's indented.

	Make sure you explain your solution and why you're doing what you're
	doing, as opposed to describing what you're doing. Reviewers and your
	future self can read the patch, but might not understand why a
	particular solution was implemented.

	Reported-by: whoever-reported-it
	Signed-off-by: Your Name <youremail@yourhost.com>


# License

DNA blockchain is licensed under the LGPL License, Version 3.0. See LICENSE for the full license text.
