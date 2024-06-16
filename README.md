> 🚀 本项目使用 Hyperledger Fabric 构建底层区块链网络, go 编写智能合约，应用层使用 gin+fabric-sdk-go ，前端使用
> vue3+element-ui

## 手动部署

环境要求： 安装了 Docker 和 Docker Compose 的 Linux 或 Mac OS 环境

附 Linux Docker 安装教程：[点此跳转](Install.md)

> 🤔 Docker 和 Docker Compose 需要先自行学习。本项目的区块链网络搭建、链码部署、前后端编译/部署都是使用 Docker 和 Docker
> Compose 完成的。

1. 下载本项目放在任意目录下，例：`/root/fabric-realty`

2. 给予项目权限，执行 `sudo chmod -R +x /root/fabric-realty/`

3. 进入 `network` 目录，执行 `./start.sh` 部署区块链网络和智能合约

4. 进入 `application` 目录，执行 `./start.sh`
   启动前后端应用，然后就可使用浏览器访问前端页面 [http://localhost:8000](http://localhost:8000)
   ，其中后端接口地址为 [http://localhost:8888](http://localhost:8888)

5. （可选）进入 `network/explorer` 目录，执行 `./start.sh`
   启动区块链浏览器后，访问 [http://localhost:8080](http://localhost:8080)，用户名 admin，密码
   123456

## 完全清理环境

注意，该操作会将所有数据清空。按照该先后顺序：

1. （如果启动了区块链浏览器）进入 `network/explorer` 目录，执行 `./stop.sh` 关闭区块链浏览器

2. 进入 `application` 目录，执行 `./stop.sh` 关闭区块链应用

3. 最后进入 `network` 目录，执行 `./stop.sh` 关闭区块链网络并清理链码容器

## 目录结构

- `application/server` : `fabric-sdk-go` 调用链码（即智能合约），`gin` 提供外部访问接口（RESTful API）


- `application/web` : `vue` + `element-ui` 提供前端展示页面


- `chaincode` : go 编写的链码（即智能合约）


- `network` : Hyperledger Fabric 区块链网络配置

## 功能介绍
前端使用vue3所有登录注册功能都支持以及验证码，7天免登录，必须同意某些协议等基础前端功能：
包含 1.登录注册2.道路信息管理（包含高德API请求可视化）3.音乐信息（使用MINIO）4.其他业务功能信息，可扩展性极强，可以随时改变5.Fabric区块链后端使用智能合约实现交易和捐赠功能
其中所有中间件Mysql和MINIO都是用docker-compose部署
包含所有的区块链服务DockerFile实现
## 截图展示
![image](https://github.com/gebilxs/fabric-realty/assets/87237189/86225240-9fa2-431c-831a-570a8d1c72d5)
![image](https://github.com/gebilxs/fabric-realty/assets/87237189/f76dcb34-6fec-429b-852c-a29f3481d13c)
![image](https://github.com/gebilxs/fabric-realty/assets/87237189/3199f3f5-e41b-4f92-a52c-fcec93268545)
![image](https://github.com/gebilxs/fabric-realty/assets/87237189/5708e0df-6f68-4472-812e-b3cb7c92f119)
![image](https://github.com/gebilxs/fabric-realty/assets/87237189/06ecb39a-59f8-4f5f-90cd-f01efd55ec0f)
![image](https://github.com/gebilxs/fabric-realty/assets/87237189/aa9b7479-7904-4d94-9636-b9237c24d255)
![image](https://github.com/gebilxs/fabric-realty/assets/87237189/34247f08-2a1f-4095-88b7-dcd9d8c8d4fd)

如需前端可以联系gebilxs@gmail.com
