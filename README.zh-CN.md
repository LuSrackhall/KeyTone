<p align="center" style="text-align: center">
  <img src="./frontend/src/assets/img/KeyTone.png" width="35%"><br/>
</p>
<p align="center">
  让按键唤醒令人舒适的声音。
  <br/>
  <br/>
  <a href="https://github.com/LuSrackhall/KeyTone/blob/main/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/LuSrackhall/KeyTone"/>
  </a>
  <a href="https://github.com/LuSrackhall/KeyTone/tags" rel="nofollow">
    <img alt="GitHub tag (latest SemVer pre-release)" src="https://img.shields.io/github/v/tag/LuSrackhall/KeyTone?include_prereleases&label=version"/>
  </a>
</p>

<div align="center">
<strong>
<samp>

[English](README.md) · [简体中文](README.zh-CN.md)

</samp>
</strong>
</div>

# KeyTone

## 介绍

KeyTone是一款实用的软件，能够在需要保持安静的环境下实时模拟按键的声音，从而提升您使用电脑的舒适度。

尽管市面上已经有很多类似的软件，但KeyTone依然致力于为用户提供最佳的体验。

此外，KeyTone是一款遵循GPL协议的开源软件。

## 特性

* 功能简单、开箱即用。
* 按键的按下和抬起, 都可独立绑定声效。
* 基于真实按键的触发情况做设计, 仅按下与抬起的瞬间发出声效, 保持按压状态不会重复播放声效。
* 高度可定制, 用户可以根据自身需求来快速定制按键的声音, 本软件围绕键音专辑, 设计开发了丰富且便捷的一系列功能。

## 安装

您可以通过以下方式获取 **KeyTone**

1. **GitHub Releases**

   作为开源项目，您可以从我们的 [GitHub Releases](https://github.com/LuSrackhall/KeyTone/releases) 页面免费下载最新版本。

2. **官方网站**

   你也可以通过访问我们的[官方网站](https://keytone.xuanhall.com)免费获取最新版本和安装说明。

3. **应用商店**

   通过各平台的官方渠道商店获取 KeyTone 的最新版本。
   * Windows&nbsp; —— &nbsp;&nbsp;&nbsp;&nbsp;从[Microsoft Store(微软商店)](https://apps.microsoft.com/store/detail/9NGKDXHPGJXD?cid=DevShareMCLPCS)获取

      &nbsp;&nbsp;<a href="https://apps.microsoft.com/detail/9ngkdxhpgjxd?referrer=appbadge&mode=direct"><img src="https://get.microsoft.com/images/zh-cn%20dark.svg" width="200"/></a>

   * MacOS&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;——&nbsp;&nbsp;&nbsp;&nbsp;  暂无上架计划

   通过应用商店下载的优势是, 可获得应用商店对所上架软件的自动更新支持。 此外, 应用商店的内置沙盒机制、以及严格的应用审核管理机制，都可为用户提供更进一步的安全保障。
  
**KeyTone软件是开源且免费的,** 但也欢迎大家通过应用商店进行购买。
<blockquote style="border-left:none; padding-left:0;">
<p>与大多数开源产品一样，KeyTone 无法独自发展。因此, 购买行为或是其它形式的赞助, 能够为开发者提供资金支持, 这有利于项目的持续维护和迭代。</p>
</blockquote>

## 系统要求

* Windows 10 或更高版本。 macOS 11 或更高版本。
* 音频输出设备。

## 星路历程

[![Stargazers over time](https://starchart.cc/LuSrackhall/KeyTone.svg?variant=adaptive)](https://starchart.cc/LuSrackhall/KeyTone)

## 开发环境与本地调试

本项目开发环境依赖如下：

* Go 1.22.0
* Node.js 18.x
* Quasar CLI

### Windows 下环境搭建步骤

1. 安装 [Go 1.22.0](https://go.dev/dl/) 并配置环境变量。
2. 安装 [Node.js 18.x](https://nodejs.org/en/download/)。
3. 全局安装 Quasar CLIo

   ```shell
   npm install -g @quasar/cli
   ```

4. 安装前端依赖：

   ```shell
   cd frontend
   npm install
   ```

5. 启动开发环境（Electron 桌面端）：

   ```shell
   quasar dev -m electron
   ```

## 路线图

项目的路线图可以在[这里](https://github.com/LuSrackhall/KeyTone/milestones)找到。在创建增强请求之前，请先查阅它。😊

## 问题

欢迎任何想法!

## 贡献

欢迎任何形式的贡献！
<!-- 请查看我们的贡献指南以了解详情。 -->

<!-- ## 星星增长趋势

<a href="https://www.star-history.com/#LuSrackhall/KeyTone&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=LuSrackhall/KeyTone&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=LuSrackhall/KeyTone&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=LuSrackhall/KeyTone&type=Date" />
 </picture>
</a> -->



## 许可证

本项目遵循[GPL许可证](https://github.com/LuSrackhall/KeyTone/blob/main/LICENSE)。

## 著作权

版权所有 (C) 2024-现在 LuSrackhall

<!-- * 部分代码（如与 Steam API 相关的代码）因包含敏感信息而不在 GPL 许可范围内。这些代码将按照符合Steam平台要求的专有许可发布，用于商业用途。 -->