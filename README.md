<p align="center" style="text-align: center">
  <img src="./frontend/src/assets/img/KeyTone.png" width="35%"><br/>
</p>
<p align="center">
  Let the keys awaken a comfortable sound.
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

<!-- [English](README.md) · [简体中文](README.zh-CN.md) -->
[English](README.md) · [简体中文](README.zh-CN.md)

</samp>
</strong>
</div>

# KeyTone

## Introduction

KeyTone is a practical software that can simulate the sound of keystrokes in real-time in environments where silence is required, thereby enhancing your comfort when using a computer.

Although there are already many similar software on the market, KeyTone is still committed to providing users with the best experience.

KeyTone's **core is implemented in high-performance Go**, which lays a solid foundation for the application's performance and allows for greater optimization potential.
> At this stage, our hybrid architecture uses Electron as the UI layer because it enables **rapid and stable implementation of complex cross-platform features**. Although Electron applications tend to be large in size, this is **not a long-term bottleneck for the KeyTone project architecture**. From the early stages of the project, we have included the evaluation and integration of lighter Go desktop frameworks (such as Wails) in our technical roadmap, and **will smoothly migrate or refactor when their ecosystems are mature enough to completely solve the size issue.**

In addition, KeyTone is an open-source software that follows the GPL license.

## Features

* Simple functionality, ready to use out of the box.
* Key press and release can be independently bound to sound effects.
* Designed based on real key trigger conditions, sound effects are only played at the moment of pressing and releasing, and will not be repeated while holding down.
* Highly customizable, users can quickly customize the sound of keys according to their own needs. This software has designed and developed a rich and convenient series of functions around key sound albums.
* Provides rich [documentation support](https://keytone.xuanhall.com/zh/guide/getting-started/quick-start/) to help users quickly understand the software's features as much as possible.

## Highlights

Unleash your creativity with keystroke sounds to craft a unique and personalized auditory experience.

KeyTone does not provide any audio files—this is by design!
> We encourage you to tap into your creativity by uploading your own or collected audio to create your personalized keystroke sound collection.

With KeyTone, you can:
* Combine multiple sounds into a single ***advanced sound*** and bind it to a single key, setting these sounds to play randomly or in sequence each time the key is pressed (playing one sound at a time).
* Allow these ***advanced sounds*** to be combined, inherited, and nested with each other, helping you create rich, varied sound effects and unlocking endless possibilities for keystroke sounds.
* **Participate in the production and sharing of Keytone albums. The main sharing community is the [Keytone community on itch](https://lusrackhall.itch.io/keytone/community). Welcome to contribute topics or join discussions on existing topics!**
* Of course, as a prerequisite, the sound files (or audio files) are not provided in this project and will not be provided in the future.
   <blockquote>
   <details>
   <summary>However, the open-source community and some free audio sharing websites provide rich resources:</summary>

   >
   > `You can freely use these audio resources locally; but if you need to share them further, be sure to check their specific license agreements.`
   >
   > * [Nigh/OpenKeySound](https://github.com/Nigh/OpenKeySound) — This repository is provided by [Nigh](https://github.com/Nigh), containing high-quality switch sounds that they have personally recorded and edited, along with related usage instructions.
   >   > Thanks to **Nigh** for their dedication and effort. Also, thanks to the [Appinn Community](https://meta.appinn.net/) for providing a sharing and communication channel.
   > * [Pixabay](https://pixabay.com/) — Pixabay is a vibrant community of authors, artists and creators sharing royalty-free images, video, audio and other media.
   > * [Freesound](https://freesound.org/search/?q=typing) — This website offers some sounds under the CC license.
   >   > Thanks to **Appinn Community** user [feeshy](https://meta.appinn.net/t/topic/72445/4#:~:text=12%20%E5%A4%A9-,%E5%8F%AF%E4%BB%A5%E5%86%85%E7%BD%AE%E4%B8%80%E4%BA%9BCC%E5%8D%8F%E8%AE%AE%E7%9A%84%E5%BD%95%E9%9F%B3,-freesound.org) for recommending this website to everyone.
   >
   > Of course, everyone can also record their own sounds to permanently capture the sound of each key on their cherished keyboard and use KeyTone to create a key sound album for preservation.
   </details>
   </blockquote>

KeyTone offers you a stage for free creation, turning every keystroke into a unique piece of sound art.

## Installation

You can obtain **KeyTone** through the following methods:

1. **GitHub Releases**

   As an open-source project, you can download the latest version for free from our [GitHub Releases](https://github.com/LuSrackhall/KeyTone/releases) page.

2. **Official Website**

   You can also visit our [official website](https://keytone.xuanhall.com) to get the latest version and installation instructions for free.

3. **itch.io**

   Download via the itch website, [click here to enter the release page](https://lusrackhall.itch.io/keytone).

4. **App Stores**

   Get the latest version of KeyTone through official app stores on various platforms.
   * Windows &nbsp;—— &nbsp;&nbsp;&nbsp; Get from [Microsoft Store](https://apps.microsoft.com/store/detail/9NGKDXHPGJXD?cid=DevShareMCLPCS)

      &nbsp;&nbsp;<a href="https://apps.microsoft.com/detail/9ngkdxhpgjxd?referrer=appbadge&mode=direct"><img src="https://get.microsoft.com/images/en-us%20dark.svg" width="200"/></a>

   * MacOS &nbsp;&nbsp;&nbsp; —— &nbsp;&nbsp;&nbsp; No plans to list yet

   The advantage of downloading through app stores is that you can get automatic update support for the listed software. In addition, the built-in sandbox mechanism of app stores and the strict application review management mechanism can provide users with further security guarantees.
  
**KeyTone software is open-source and free,** but we also welcome everyone to purchase through app stores.
<blockquote style="border-left:none; padding-left:0;">
<p>Like most open-source products, KeyTone cannot develop on its own. Therefore, purchasing or other forms of sponsorship can provide financial support to developers, which is beneficial for the continuous maintenance and iteration of the project.</p>
</blockquote>

## System Requirements

* Windows 10 or higher. macOS 11 or higher.
* Audio output device.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/LuSrackhall/KeyTone.svg?variant=adaptive)](https://starchart.cc/LuSrackhall/KeyTone)

## Development Environment & Local Debugging

The project requires the following development environment:

* Go 1.22.0
* Node.js 18.x
* Quasar CLI

### Setup Steps on Windows

1. Install [Go 1.22.0](https://go.dev/dl/) and configure your environment variables.
2. Install [Node.js 18.x](https://nodejs.org/en/download/).
3. Install Quasar CLI globally:

   ```shell
   npm install -g @quasar/cli
   ```

4. Install frontend dependencies:

   ```shell
   cd frontend
   npm install
   ```

5. Start the development environment (Electron desktop):

   ```shell
   quasar dev -m electron
   ```

## Roadmap

The project's roadmap can be found [here](https://github.com/LuSrackhall/KeyTone/milestones). Please check it before creating enhancement requests. 😊

## Issues

Welcome any ideas!

## Contributions

Welcome contributions in any form!
<!-- Please check our contribution guidelines for details. -->

<!-- ## Star History

<a href="https://www.star-history.com/#LuSrackhall/KeyTone&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=LuSrackhall/KeyTone&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=LuSrackhall/KeyTone&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=LuSrackhall/KeyTone&type=Date" />
 </picture>
</a> -->

## License

This project follows the [GPL License](https://github.com/LuSrackhall/KeyTone/blob/main/LICENSE).

## Copyright

Copyright (C) 2024-present LuSrackhall