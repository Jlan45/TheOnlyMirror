# TheOnlyMirror

TheOnlyMirror 是一个用于单域名实现多镜像源的项目。通过对来源UA与URL的判断实现了自动分流到不同的source进行反向代理</br>

测试站点 mirror.nothinglikethis.asia

## 另一种实现
如果可以的话请看看隔壁用config实现的版本 https://github.com/huangzheng2016/TheOnlyMirror
## 可用镜像

- Debian系
  - Ubuntu
  - Kali
  - Debian
- Pypi
- Docker
- npm
- Github Clone


## 配置介绍

```
{
  "domain": "mirror.nothinglikethis.asia", //域名
  "mirrorList": [ //目前没用
    "dockerhub",
    "github",
    "pypi",
    "ubuntu"
  ],
  "sources": { //源url
    "pypi_index": "https://pypi.org",
    "pypi_files": "https://files.pythonhosted.org",
    "dockerhub":"https://docker.mirror.nothinglikethis.asia",
    "ubuntu": "http://archive.ubuntu.com",
    "ubuntu_ports": "http://ports.ubuntu.com",
    "debian": "http://deb.debian.org",
    "kali": "http://http.kali.org",
    "npm":"https://registry.npmjs.org"
  },
  "port": 8080, //监听端口
  "tlsport": 443, //https监听端口,如果对应端口被占用不会终止程序
  "tls": false, //根据该配置判断对部分有修改的源返回http还是https
  "hostControll": false, //可通过host实现用户分配
  "hostList": [] //host列表
}
```

## 食用方式

对于大部分的应用来说，只需要将对应镜像位置的url进行替换即可，少部分使用APT作为包管理器的，请在对应的源后添加如下Path

- Ubuntu /ubuntu
- Ubuntu-ports /ubuntu-ports
- Debian /debian
- Kali /kali

## 部署方式
```shell
git clone https://github.com/Jlan45/TheOnlyMirror
cd TheOnlyMirror
go build .
./TheOnlyMirror
```
## TODO

- 可能会通过添加内部代理的方式方便校内部署
- 可能会有的yum系代理
- 可能会添加的go代理
- 可能会打包一个docker
- 想要什么站点可以提issue