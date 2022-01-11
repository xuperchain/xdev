# XDEV 文档
XDEV 是 [XuperChain](https://github.com/xuperchain/xuperchain) 合约构建测试工具，支持
- WASM 合约构建
- 不同语言和类型的合约测试

### XDEV 安装
1. 安装依赖项

    XDEV 使用 [Docker](https://docs.docker.com/engine/install/) 进行 WASM 合约构建，如果你使用XDEV 构建 WASM 合约，在使用前你需要 [安装Docker](https://docs.docker.com/engine/install/) 

2. 构建 xdev

   xdev 需要从源码开始构建 

``` bash
    git clone https://github.com/xuperchain/xdev.git 
    cd xdev 
    git checkout v1.0.0
    cd xdev 
    make 
```
构建产出在当前目录下的 bin 目录下

3. 设置环境变量（可选）

   可以将<XDEV_ROOT>/bin 目录加入到环境变量,以便可以在任意路径使用xdev，其中 XDEV_ROOT 是 xdev 源码的根目录。

### 合约构建

1. 单文件构建

以 C++ 合约[counter](https://github.com/xuperchain/contract-sdk-cpp/blob/main/example/counter.cc) 合约为例

``` bash
    xdev build -o counter.wasm example/counter.cc 
```

1. 多文件构建
以 [xuper_relay](https://github.com/xuperchain/contract-sdk-cpp/tree/main/example/xuper_relayer) 为例，执行
``` bash
    xdev build 
``` 

### 合约测试
1. 单文件测试
合约测试需要编写测试文件,以 [xuper_relay](https://github.com/xuperchain/contract-sdk-cpp/blob/main/test/xuper_relay.test.js) 合约的测试为例
``` bash
    xdev test xuper_relay.test.js
```


## 参与贡献

如果你遇到问题或需要新功能，欢迎创建issue。

如果你可以解决某个issue, 欢迎发送PR。

如项目对您有帮助，欢迎star支持。


