# 部署ethnode
部署官方的PPA(geth)版本，作为接入eth网络的入口

## 启动新实例

因需要部署以太坊全节点，相对比较耗费资源，所以推荐最低配置为8核_32g内存_300GB SSD以上的实例

建议部署两台以上eth节点，避免因单台故障、节点延时或运维不当等造成的服务不可用

申请2台EC2实例，参考启动aws [EC2实例](new_ec2_cn.md)，并且关联`ethnode-SecurityGroup`安全组

> 如果还没创建，请参考[aws安全组](security_group_cn.md)关于`ethnode-SecurityGroup`部分的说明，创建后再关联

## 部署
```
sudo apt-get -y install software-properties-common
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get -y install ethereum
sudo mkdir -p /data/ethereum
sudo chown -R ubuntu:ubuntu /data/ethereum
```

## 启停

* ### 启动
修改x.x.x.x为本实例内网ip

```
sudo nohup geth --datadir /data/ethereum --fast --cache=1024 --rpc --rpcaddr x.x.x.x --rpcport 8545 --rpccorsdomain * &
```

* ### 终止
`sudo kill -9 进程id`

## 日志

当前目录下 `nohup.out`
