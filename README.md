# TheOnlyMirror

TheOnlyMirror 是一个用于单域名实现多镜像源的项目。通过对来源 UA 与 URL 的判断实现了自动分流到不同的 source 进行反向代理。

访问根路径 `/` 会展示一个首页，列出所有已启用的镜像及配置说明。

## 可用镜像

- Debian 系
  - Ubuntu (含 Ports)
  - Debian
  - Kali
- Alpine
- PyPI
- Docker Hub
- npm
- GitHub Clone
- Go Modules (通过 GitHub 镜像)
- AOSP (Android 开源项目)

## 配置介绍

复制 `config.example.json` 为 `config.json` 后按需修改：

```jsonc
{
  "domain": "mirror.example.com",    // 对外域名
  "mirrorList": [                    // 首页展示的镜像列表，为空则展示全部
    "dockerhub",
    "github",
    "pypi",
    "ubuntu",
    "aosp"
  ],
  "sources": {                       // 上游源 URL
    "pypi_index": "https://pypi.org",
    "pypi_files": "https://files.pythonhosted.org",
    "dockerhub": "https://registry-1.docker.io",
    "dockerblob": "https://production.cloudflare.docker.com",
    "ubuntu": "http://archive.ubuntu.com",
    "ubuntu_ports": "http://ports.ubuntu.com",
    "debian": "http://deb.debian.org",
    "kali": "http://http.kali.org",
    "npm": "https://registry.npmjs.org",
    "alpine": "https://dl-cdn.alpinelinux.org",
    "github": "https://github.com",
    "aosp": "https://android.googlesource.com"
  },
  "port": 8080,                      // HTTP 监听端口
  "tlsport": 443,                    // HTTPS 监听端口
  "tls": false,                      // 是否启用 TLS
  "certFile": "example.crt",         // TLS 证书文件路径
  "keyFile": "example.key",          // TLS 私钥文件路径
  "hostControll": false,             // 是否启用 Host 白名单控制
  "hostList": [],                    // 允许访问的 Host 列表
  "proxy": "",                       // 全局上层代理 (支持 http/https/socks5)
  "sourceProxy": {}                  // 按源单独配置代理，优先级高于全局
}
```

## 食用方式

对于大部分应用，只需将原始源地址替换为你的镜像域名即可。UA 自动检测会将请求路由到正确的上游。

对于 APT 系包管理器，需要在域名后添加对应路径：

| 发行版 | 路径 |
|--------|------|
| Ubuntu | `/ubuntu` |
| Ubuntu Ports | `/ubuntu-ports` |
| Debian | `/debian` |
| Kali | `/kali` |

Alpine 使用路径 `/alpine`。

## 部署方式

```shell
git clone https://github.com/Jlan45/TheOnlyMirror
cd TheOnlyMirror
cp config.example.json config.json
# 编辑 config.json 修改 domain 等配置
go build .
./TheOnlyMirror
```

## TODO

- 可能会通过添加内部代理的方式方便校内部署
- 可能会有的 yum 系代理
- 可能会打包一个 Docker 镜像
- 想要什么站点可以提 issue