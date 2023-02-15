# Gopanel

A control panel that is written in Golang and is able to manage [Caddy 2](https://caddyserver.com/) and Guard web server.

Committed to becoming a server-type, router-type, NAS-type all-round management panel.

[![AUR package](https://repology.org/badge/version-for-repo/aur/gopanel.svg)](https://repology.org/project/gopanel/versions)
[![pipeline status](https://gitlab.com/xiayesuifeng/gopanel/badges/master/pipeline.svg)](https://gitlab.com/xiayesuifeng/gopanel/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/xiayesuifeng/gopanel)](https://goreportcard.com/report/gitlab.com/xiayesuifeng/gopanel)
[![GoDoc](https://godoc.org/gitlab.com/xiayesuifeng/gopanel?status.svg)](https://godoc.org/gitlab.com/xiayesuifeng/gopanel)
[![Sourcegraph](https://sourcegraph.com/gitlab.com/xiayesuifeng/gopanel/-/badge.svg)](https://sourcegraph.com/gitlab.com/xiayesuifeng/gopanel)
![license](https://img.shields.io/badge/license-GPL3.0-green.svg)

# Gopanel 前端

[gopanel-web](https://gitlab.com/xiayesuifeng/gopanel-web.git)

## Features (功能)
> PS: 以下大部分仍在开发中

* 图形化配置反向代理，静态文件服务器等所有 `Caddy` 支持的操作
* 应用管理
  * 优雅管理 web 服务
  * 支持应用中心，一键安装部署
  * 一站式从路由配置到 web 服务进程管理
  * 自定义应用图形化管理
  * 支持与 systemd 服务集成
* DDNS (基于 Caddy 动态 DNS 插件实现)
* 容器化支持
  * 支持 Docker/Podman 后端
  * 与应用中心集成，一键部署容器服务
  * 图形化配置容器
* 网络 (从基础配置到像路由器那样丰富的功能)
  * 防火墙
  * DHCP
  * DNS
  * 交换机
  * VLAN
  * 端口转发
* 虚拟化支持
* 存储
  * 硬盘管理
  * Samba
  * NFS
* 系统管理 (依赖使用 systemd 的发行版)
  * 服务管理
  * 日志管理
* 插件
  * 支持插件中心
  * 通过插件扩展更多更强大的功能

## Installer (安装)

> ArchLinux (AUR)

假设 `AUR Helper` 为 `yay`
```bash
yay -S gopanel-bin
```

> Other GNU/Linux Distro (其他 `GNU/Linux` 发行版)
> 请确保已安装好 `caddy2`
```bash
wget https://gitlab.com/xiayesuifeng/gopanel/-/jobs/artifacts/master/download?job=build-gopanel -o gopanel.zip
unzip gopanel.zip
cd gopanel
sudo mkdir -p /etc/gopanel/app.conf.d
sudo mkdir -p /usr/share/gopanel
sudo install -D -m 0755 gopanel /usr/bin/gopanel
sudo install -D -m 0644 systemd/gopanel.service /usr/lib/systemd/system/gopanel.service
sudo install -D -m 0644 config.default.json /etc/gopanel/config.json
sudo cp -rf web /usr/share/gopanel/web
sudo chmod -R 0644 /usr/share/gopanel/web
sudo systemctl daemon-reload
```

## 开机自启与运行
```bash
sudo systemctl enable --now gopanel
```

## 配置文件详解
```json5
{
  // 日志配置
  "log": {
    // 日志级别，可选 debug, info, warn, error
    "level": "debug",
    // 日志输出，可选 stdout, stderr, 文件路径
    "output": "stderr",
    // 日志格式，可选  text, json
    "format": "text"
  },
  // 运行模式，如调试好请改 release
  "mode":"debug",
  // 登录密码，默认为 admin，如要修改请参考下方的 ‘加密密码生成’
  "password": "0925e15d0ae6af196e6295923d76af02b4a3420f",
  // app 配置文件存储路径
  "appConf": "app.conf.d",
  // jwt 加密密钥设置，如不设置则每次启动设置为 gopanel-secret-[随机数]
  "secret": "",
  // gopanel 访问配置
  "panel": {
    // 域名绑定，如没有可删除或者留空并且必须设置 port
    "domain": "example.com",
    // 端口设置，如没设置域名则必须设置否则无特殊要求不需要设置
    "port": 2020,
    // 禁用 SSL (取代原来的 automaticHttps)
    "disableSSL": true
  },
  // caddy 2 配置
  "caddy": {
    // caddy 2 API 
    "adminAddress": "http://localhost:2019",
  },
  // 预留配置，未来用于后端异常状态通知
  "smtp":{
    "username":"",
    "password": "",
    "host": ""
  },
  // 集成 netdata 设置
  "netdata":{
    // 是否启用该功能
    "enable": false,
    // 反代地址
    "host": "localhost:19999",
    "path": "",
    // 使用 SSL 访问
    "ssl": false
  }
}
```

## 加密密码生成
```
echo -n yourpassword | openssl dgst -md5 -binary | openssl dgst -sha1
```

## License

Gopanel is licensed under [GPLv3](LICENSE).