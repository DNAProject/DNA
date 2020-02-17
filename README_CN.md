[![Build Status](https://travis-ci.org/DNAProject/DNA.svg?branch=master)](https://travis-ci.org/DNAProject/DNA)

[English](README.md) | 中文

# DNA (Distributed Networks Architecture)

DNA是go语言实现的基于区块链技术的去中心化的分布式网络协议。可以用来数字化资产和金融相关业务包括资产注册，发行，转账等。

## 特性

* 可扩展的轻量级通用智能合约
* 跨链交互协议（进行中）
* 抗量子密码算法 (可选择模块)
* 中国商用密码算法 (可选择模块)
* 高度优化的交易处理速度
* 节点访问权限控制
* P2P连接链路加密
* 多种共识算法支持 (DBFT/VBFT)
* 可配置区块生成时间
* 可配置电子货币模型
* 可配置的分区共识(进行中)

# 编译
成功编译DNA需要以下准备：

* Go版本在1.12.5及以上
* 正确的Go语言开发环境

克隆DNA仓库到$GOPATH/src/DNAProject目录


```shell
$ git clone https://github.com/DNAProject/DNA.git
```

用make编译源码

```shell
$ make
```

成功编译后会生成可执行程序

* `dnaNode`: 节点程序

# 部署

成功运行DNA需要至少4个节点，可以通过两种方式进行部署

* 多机部署
* 单机部署

## 多机部署配置

我们可以通过修改默认的配置文件`config.json`进行快速部署。

1. 将相关文件复制到目标主机，包括：
    - 默认配置文件`config.json`
    - 节点程序`dnaNode`

3. 种子节点配置
    - 将4个主机作为种子节点地址分别填写到每个配置文件的`SeelList`中，格式为`种子节点IP地址 + 种子节点NodePort`

4. 创建钱包文件
    - 通过命令行程序，在每个主机上分别创建节点运行所需的钱包文件wallet.dat 
      
        `$ ./dnaNode account add -d` 

5. 记账人配置
    - 为每个节点创建钱包时会显示钱包的公钥信息，将所有节点的公钥信息分别填写到每个节点的配置文件的`peers`项中
    
        注：每个节点的钱包公钥信息也可以通过命令行程序查看：
    
        `$ ./dnaNode account list -v` 


多机部署配置完成，每个节点目录结构如下

```shell
$ ls
config.json dnaNode wallet.dat
```

一个配置文件片段如下, 其中10.0.1.100、10.0.1.101等都是种子节点地址:
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
## 单机部署配置

如果您希望以测试模式运行DNA区块链，那么不需要做任何以上多机配置。直接通过以下命令启动测试模式的DNA区块链。

```shell
$ ./dnaNode --testmode
$ - input your wallet password
```

## 运行
以任意顺序运行每个节点node程序，并在出现`Password:`提示后输入节点的钱包密码

```shell
$ ./dnaNode
$ - 输入你的钱包口令
```

了解更多请运行 `./dnaNode --help`.


# 贡献代码

请您以签过名的commit发送pull request请求，我们期待您的加入！
您也可以通过邮件的方式发送你的代码到开发者邮件列表，欢迎加入DNA邮件列表和开发者论坛。

另外，在您想为本项目贡献代码时请提供详细的提交信息，格式参考如下：

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

# 开源社区

## 邮件列表

我们为开发者提供了一下邮件列表

- OnchainDNA@googlegroups.com

可以通过两种方式订阅并参与讨论

- 发送任何内容到邮箱地址 OnchainDNA+subscribe@googlegroups.com

- 登录 https://groups.google.com/forum/#!forum/OnchainDNA 


# 许可证

DNA遵守LGPL License, 版本3.0。 详细信息请查看项目根目录下的LICENSE文件。
