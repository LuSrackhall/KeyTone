/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

/* eslint-env node */

/*
 * This file runs in a Node context (it's NOT transpiled by Babel), so use only
 * the ES6 features that are supported by your Node version. https://node.green/
 */

// Configuration for your app
// https://v2.quasar.dev/quasar-cli-vite/quasar-config-js

const { configure } = require('quasar/wrappers');
const path = require('path');
const { version } = require('./package.json');

module.exports = configure(function (/* ctx */) {
  return {
    eslint: {
      // fix: true,
      // include: [],
      // exclude: [],
      // rawOptions: {},
      warnings: true,
      errors: true,
    },

    // https://v2.quasar.dev/quasar-cli-vite/prefetch-feature
    // preFetch: true,

    // app boot file (/src/boot)
    // --> boot files are part of "main.js"
    // https://v2.quasar.dev/quasar-cli-vite/boot-files
    boot: ['i18n', 'axios', 'unocss', 'quasar-lang-pack'],

    // https://v2.quasar.dev/quasar-cli-vite/quasar-config-js#css
    css: ['app.scss'],

    // https://github.com/quasarframework/quasar/tree/dev/extras
    extras: [
      // 'ionicons-v4',
      // 'mdi-v5',
      // 'fontawesome-v6',
      // 'eva-icons',
      // 'themify',
      // 'line-awesome',
      // 'roboto-font-latin-ext', // this or either 'roboto-font', NEVER both!

      'roboto-font', // optional, you are not bound to it
      'material-icons', // optional, you are not bound to it
    ],

    // Full list of options: https://v2.quasar.dev/quasar-cli-vite/quasar-config-js#build
    build: {
      target: {
        browser: ['es2019', 'edge88', 'firefox78', 'chrome87', 'safari13.1'],
        node: 'node16',
      },

      // TIPS: 可以通过以下代码, 让build生成的生产环境的代码, 删除console.log、debugger等。
      minify: 'terser',
      extendViteConf(viteConf) {
        viteConf.build.terserOptions = viteConf.build.terserOptions || {};
        viteConf.build.terserOptions = {
          ...viteConf.build.terserOptions,
          compress: {
            ...viteConf.build.terserOptions.compress,
            drop_console: true,
            drop_debugger: true,
          },
        };
      },

      vueRouterMode: 'hash', // available values: 'hash', 'history'
      // vueRouterBase,
      // vueDevtools,
      // vueOptionsAPI: false,

      // rebuildCache: true, // rebuilds Vite/linter/etc cache on startup

      // publicPath: '/',
      // analyze: true,
      env: {
        APP_VERSION: version,
      },
      // rawDefine: {}
      // ignorePublicFolder: true,
      // minify: false,
      // polyfillModulePreload: true,
      // distDir

      // extendViteConf (viteConf) {},
      // viteVuePluginOptions: {},

      vitePlugins: [
        [
          '@intlify/vite-plugin-vue-i18n',
          {
            // if you want to use Vue I18n Legacy API, you need to set `compositionOnly: false`
            // compositionOnly: false,

            // if you want to use named tokens in your Vue I18n messages, such as 'Hello {name}',
            // you need to set `runtimeOnly: false`
            // runtimeOnly: false,

            // you need to set i18n resource including paths !
            include: path.resolve(__dirname, './src/i18n/**'),
          },
        ],
        ['unocss/vite', {}],
      ],
    },

    // Full list of options: https://v2.quasar.dev/quasar-cli-vite/quasar-config-js#devServer
    devServer: {
      // https: true
      open: true, // opens browser window automatically
    },

    // https://v2.quasar.dev/quasar-cli-vite/quasar-config-js#framework
    framework: {
      config: {},

      // iconSet: 'material-icons', // Quasar icon set
      // lang: 'en-US', // Quasar language pack

      // For special cases outside of where the auto-import strategy can have an impact
      // (like functional components as one of the examples),
      // you can manually specify Quasar components/directives to be available everywhere:
      //
      // components: [],
      // directives: [],

      // Quasar plugins
      plugins: ['Notify', 'Dialog'],
    },

    // animations: 'all', // --- includes all animations
    // https://v2.quasar.dev/options/animations
    animations: [],

    // https://v2.quasar.dev/quasar-cli-vite/quasar-config-js#sourcefiles
    // sourceFiles: {
    //   rootComponent: 'src/App.vue',
    //   router: 'src/router/index',
    //   store: 'src/store/index',
    //   registerServiceWorker: 'src-pwa/register-service-worker',
    //   serviceWorker: 'src-pwa/custom-service-worker',
    //   pwaManifestFile: 'src-pwa/manifest.json',
    //   electronMain: 'src-electron/electron-main',
    //   electronPreload: 'src-electron/electron-preload'
    // },
    sourceFiles: {
      electronPreload: 'src-electron/main-process/electron-preload',
    },

    // https://v2.quasar.dev/quasar-cli-vite/developing-ssr/configuring-ssr
    ssr: {
      // ssrPwaHtmlFilename: 'offline.html', // do NOT use index.html as name!
      // will mess up SSR

      // extendSSRWebserverConf (esbuildConf) {},
      // extendPackageJson (json) {},

      pwa: false,

      // manualStoreHydration: true,
      // manualPostHydrationTrigger: true,

      prodPort: 3000, // The default port that the production server should use
      // (gets superseded if process.env.PORT is specified at runtime)

      middlewares: [
        'render', // keep this as last one
      ],
    },

    // https://v2.quasar.dev/quasar-cli-vite/developing-pwa/configuring-pwa
    pwa: {
      workboxMode: 'generateSW', // or 'injectManifest'
      injectPwaMetaTags: true,
      swFilename: 'sw.js',
      manifestFilename: 'manifest.json',
      useCredentialsForManifestTag: false,
      // useFilenameHashes: true,
      // extendGenerateSWOptions (cfg) {}
      // extendInjectManifestOptions (cfg) {},
      // extendManifestJson (json) {}
      // extendPWACustomSWConf (esbuildConf) {}
    },

    // Full list of options: https://v2.quasar.dev/quasar-cli-vite/developing-cordova-apps/configuring-cordova
    cordova: {
      // noIosLegacyBuildFlag: true, // uncomment only if you know what you are doing
    },

    // Full list of options: https://v2.quasar.dev/quasar-cli-vite/developing-capacitor-apps/configuring-capacitor
    capacitor: {
      hideSplashscreen: true,
    },

    // Full list of options: https://v2.quasar.dev/quasar-cli-vite/developing-electron-apps/configuring-electron
    electron: {
      // extendElectronMainConf (esbuildConf)
      // extendElectronPreloadConf (esbuildConf)

      inspectPort: 5858,

      bundler: 'builder', // 'packager' or 'builder'

      packager: {
        // https://github.com/electron-userland/electron-packager/blob/master/docs/api.md#options
        // OS X / Mac App Store
        // appBundleId: '',
        // appCategoryType: '',
        // osxSign: '',
        // protocol: 'myapp://path',
        // Windows only
        // win32metadata: { ... }
      },

      builder: {
        // https://www.electron.build/configuration/configuration
        productName: 'KeyTone', // 这是您的应用程序的名称，用户在安装时看到的名称。
        appId: 'top.srackhall.keytone',
        // asar: false, // (推荐使用默认的true,以简单地保护下前端源代码<虽然实际作用不大>)这表明您的应用代码不会被打包到 ASAR 存档中。ASAR 是 Electron 用来打包应用源代码的格式，设置为 false 表示源代码将以普通文件夹的形式包含在应用中。
        compression: 'maximum', // 压缩级别, 但这个主要与构建速度有关, 如指定为maximum将会使用更多的时间来构建。<与此同时, 最终生成的包的大小, 却并没有带来实质性的体积优化; 因此一般大家喜欢指定为store, 以或获得显著的构建时间的提升>
        copyright: 'Copyright (C) 2024 LuSrackhall',
        artifactName: 'KeyTone-${version}-${os}-${arch}.${ext}', // 这个配置是用于定义生成的文件名模板。artifactName 可以包含变量，比如 ${version} 会用应用程序的版本号替换，${os} 会用目标操作系统替换，${ext} 会用文件的扩展名替换。这样，生成的安装文件名能清楚地反映出版本信息和适用的操作系统。
        extraFiles: {
          //指定额外要包括在应用程序包中的文件或目录。可以配置 from 和 to 字段来定义源文件或目录的路径以及它们在打包应用中的目标路径。
          from: '../LICENSE',
          to: 'LICENSE',
        },
        win: {
          // icon: 'src/assets/icon.ico', // 错误的论点: ~~某些应用模板或框架（如 Quasar），已经在项目中包含了默认的图标文件，并在构建过程中自动配置好了这些图标, 因此无需在此处配置。~~
          icon: 'icons/icon.ico', // 不指定的话, 安装后是没有图标的。(win下, 使用ico格式结尾的图标名) 这里之前未指定也拥有图标的原因我也不知道, 但鉴于Linux上的经历, 最好还是指定下。
          extraResources: [
            {
              from: 'dist/key_tone_sdk',
              to: 'key_tone_sdk',
              // filter:[]
            },
          ], // TIPS: 核心, 有了它就不需要再依赖makefile了
          signingHashAlgorithms: ['sha256'], // 这个配置项意味着, 如果要启用代码签名，将使用SHA-256哈希算法进行签名。代码签名对于确认应用程序的完整性和来源是非常重要的。
          rfc3161TimeStampServer: 'http://timestamp.entrust.net/TSS/RFC3161sha2TS', //这个配置项用于指定一个RFC 3161兼容的时间戳服务器的URL。当你的应用程序进行代码签名时，时间戳能确保签名在证书过期之后仍然有效。
          //这个网址不是私人的，它是一个公开的RFC 3161兼容时间戳服务的URL。时间戳服务由Entrust公司提供，是用来为数字签名生成时间戳的。这对于确保在软件开发和分发过程中代码签名的完整性和长期有效性非常重要。当您对软件进行数字签名时，时间戳服务会记录和证明签名是在证书有效期内完成的。这样，即使证书在未来某个时候过期，带有时间戳的签名也仍然会被操作系统和用户视为有效和可信的。
          // 通常情况下，使用公开的RFC 3161时间戳服务，如Entrust提供的服务，是不需要付费的。这些服务一般是免费提供给公众使用的，以便于签名过程中实现时间戳的功能。但是，有一些服务提供商可能会对频繁或大量的使用提出收费，所以具体是否需要付费可能取决于你的使用情况以及服务提供商的政策。如果你是商业用户或对时间戳服务有高要求，建议直接咨询服务提供商，以了解其服务详情和是否需要付费。
          // target: 'nsis', // 这个配置是特定于 Windows 平台的。target 指定生成安装程序的类型。这里设置为 "NSIS"（Nullsoft Scriptable Install System），它是一个流行的 Windows 安装程序制作工具，允许开发者创建带有自定义逻辑和界面的安装程序。
          target: [
            'nsis', // 普通安装包
            {
              target: 'appx',
              arch: ['x64'],
            },
          ],
        },
        // 定制安装程序流程和相关操作。
        nsis: {
          oneClick: false, // 这个选项设置为false意味着安装时不会直接一键安装，用户将会看到更多的安装选项。
          perMachine: false, // 表示安装默认是为当前用户安装，而不是为所有用户安装。
          allowToChangeInstallationDirectory: true, // 允许用户在安装期间更改安装目录。
          allowElevation: true, // 允许安装程序请求提升权限，如果需要的话。
          deleteAppDataOnUninstall: true, // 表示当应用被卸载时将删除应用数据目录。
          createDesktopShortcut: true, // 表示安装程序将创建桌面快捷方式。
          createStartMenuShortcut: true, //表示将创建开始菜单快捷方式。
          shortcutName: 'KeyTone', // 指定了快捷方式的名称。
          license: '../LICENSE', // 指定了安装过程中会显示的许可协议文件。
          // include: 'nsis/installer.nsh', // 包含额外的NSIS脚本文件，用于自定义安装过程。 // 比如检测系统版本是否是win10及以上, 若不满足触发退出窗口, 用户无法继续安装只能退出。
          warningsAsErrors: false, // 这意味着在NSIS脚本编译时，警告不会作为错误处理
          // installerSidebar: 'nsis/installerSidebar.bmp', // 安装侧边栏使用的图片。
          // uninstallerSidebar: 'nsis/uninstallerSidebar.bmp', // 卸载侧边栏使用的图片。
        },
        appx: {
          applicationId: 'LuSrackhall.KeyTone', // 应用程序在 Windows Store 中的唯一标识符
          displayName: 'KeyTone', // Windows Store 中显示的应用名称
          publisher: 'CN=54940991-8DDB-4049-B299-2C3D65A6BBCC', // 发布者信息
          publisherDisplayName: 'LuSrackhall', // Windows Store 中显示的发布者名称
          identityName: 'LuSrackhall.KeyTone', // 应用的身份名称
          languages: [
            'zh-CN',
            'en-US',
            'zh-TW',
            'ja',
            'ko-KR',
            'de',
            'ru',
            'fr',
            'it',
            'es',
            'pt',
            'pt-BR',
            'pl',
            'tr',
            'vi',
            'id',
            'ar',
          ], // 支持的语言
          // backgroundColor: 'transparent', // 设置为透明
          backgroundColor: '#ffffff', // 应用启动时的背景色
          // showNameOnTiles: true, // 在磁贴上显示应用名称
          addAutoLaunchExtension: true, // 是否添加自动启动扩展
          // identityType: 'Store', // 指定应用程序的认证类型。'Store' 表示应用程序已经通过微软的认证，可以在Windows Store中发布。
          // identityCertificate: '../../KeyTone.pfx', // 指定应用程序的证书文件路径。这个文件通常用于对应用程序进行数字签名，以确保应用程序的完整性和安全性。
          // identityCertificatePassword: '12345678', // 指定应用程序的证书密码。这个密码用于解密和验证应用程序的数字签名。
        },
        linux: {
          // icon: 'src/assets/icon.icns', // 错误的论点: ~~某些应用模板或框架（如 Quasar），已经在项目中包含了默认的图标文件，并在构建过程中自动配置好了这些图标, 因此无需在此处配置。~~
          icon: 'icons/icon.icns', // 不指定的话, 安装后是没有图标的。(Linux下, 使用icns格式结尾的图标名)// 某些应用模板或框架（如 Quasar）, 只会帮助你将图标处理到指定路径和指定尺寸, 但用于何处还是要自己主动指定的。
          icon: 'src/assets/img/KeyTone.png', // 由于在Ubuntu上安装后, 没有图标, 因此我尝试指定它来解决这个问题。
          category: 'Utility', // 指定应用在 Linux 应用菜单中的分类。Utility 是一个常见的类别，表示工具类应用。此配置有助于将应用放置在合适的类别中，使用户更容易找到。
          executableName: 'keytone', // 定义构建后可执行文件的名称。在你的例子中，构建后的可执行文件将命名为 Keytone。在 Linux 中，用户通常通过这个名称来启动应用, 因此常用纯小写来命名。

          // 定义 Linux 桌面文件（.desktop 文件）的属性。这些文件通常放置在 /usr/share/applications/ 或 ~/.local/share/applications/ 中，用于在应用菜单中显示应用。
          desktop: {
            Name: 'KeyTone', // Name: 应用程序的显示名称。它将在菜单或启动器中显示。
            Comment: 'let keyboard bring forth pleasing sounds', // 应用的简短描述。当用户将鼠标悬停在应用图标上时会显示这个描述。
            Type: 'Application', //  指定桌面文件的类型，通常为 Application，表示这是一个应用程序。
            Categories: 'Utility;', // 指定应用的类别，用于将应用归类到系统的应用菜单中。多个类别可以用分号分隔。
          },

          // 指定需要包含在应用包中的额外资源文件或目录。
          extraResources: {
            from: 'dist/key_tone_sdk',
            to: 'key_tone_sdk',
          },
          // 指定构建的目标包格式，可以是单一的包格式或多个格式。
          target: [
            'tar.gz', // 生成 tar.gz 格式 // 一个使用 tar 和 gzip 组合压缩的归档格式，适用于发布和分发 Linux 应用。
            'AppImage', // 生成 AppImage 格式 // 一种通用的 Linux 可执行文件格式，允许用户在大多数 Linux 发行版上运行应用，而无需安装。
            'deb', // 生成 deb 格式 // Debian 软件包格式，广泛用于基于 Debian 的 Linux 发行版（如 Ubuntu）中。此格式方便用户通过软件包管理器（如 dpkg 或 apt）安装应用。
          ],
        },

        mac: {
          // WARN: 我绝对不再折腾这个问题了。 -> 即 为 应用名称.app 授予的辅助功能权限, 无法作用于 应用名称.app/Contents/MacOS/应用名称 这个二进制可执行文件。 (为 应用名称.app 授予的权限实际上可能是给到了 builderId/appid 或是整个包的签名证书上, 很有可能是没有使用签名证书的缘故造成的包内部可执行程序无法访问权限。) 总之, 在我不打算买苹果开发者账号之前, 不需要解决此问题。
          identity: 'Copyright (C) 2024 LuSrackhall',
          provisioningProfile: '../../KeyTone.provisionprofile', // 这个文件通常由 Apple 提供，包含应用程序的分发许可和其他信息，确保应用程序可以在特定设备或环境中运行。(这应该就是苹果平台的签名文件了吧)
          icon: 'icons/icon.icns', // 不多解释了, 它将在应用程序包中用于显示应用程序图标
          gatekeeperAssess: false, //  控制是否执行 Gatekeeper 评估。 Gatekeeper 是 macOS 的一项安全功能，用于确保应用程序来自受信任的开发者。将此设置为 false 将跳过 Gatekeeper 的评估。一般情况下，不建议关闭。
          hardenedRuntime: true, // 启用或禁用强化运行时。强化运行时是一项 macOS 的安全功能，增强应用程序的安全性和保护措施。启用它可以帮助防止应用程序被恶意修改。
          // `plutil -lint entitlements.plist`命令可用于检测.plist文件的正确性, 如此命令的执行结果应为`entitlements.plist: OK`
          entitlements: '../../entitlements.mas.plist', // 指定权限文件的路径。权限文件定义了应用程序请求的系统资源（如文件系统、网络访问、硬件等）。这个文件通常需要进行签名，并且在打包和分发过程中与应用一起提供。
          // 用于向 Info.plist 文件添加额外的键值对。( Info.plist 是 macOS 应用程序的核心配置文件，包含了应用的基本信息。这个字段可以用来添加一些应用程序运行所需的特定键值。)
          // extendInfo: {
          // NSMicrophoneUsageDescription: 'This app requires microphone access to record audio.', // 用于说明为什么应用需要访问麦克风。
          // },
          entitlementsInherit: '../../entitlements.mas.plist', // 为子进程或辅助进程指定继承的权限文件路径。(当应用程序生成子进程时，这些进程会继承指定的权限文件中的权限。对于沙盒化的应用程序，这一点尤为重要。)
          minimumSystemVersion: '10.15', // 指定应用程序支持的最低 macOS 版本。 (确保应用程序不会在低于指定版本的 macOS 上安装或运行。这里的设置是 macOS 10.15。)

          // 指定需要包含在应用包中的额外资源文件或目录。
          extraResources: {
            from: 'dist/key_tone_sdk',
            to: 'key_tone_sdk',
          },
          // 指定打包目标格式和cpu架构。
          target: [
            {
              target: 'dmg',
              // arch: 'x64',
              // arch: 'arm64',
            },
            // {
            //   target: 'dmg',
            //   // arch: 'x64',
            //   // arch: 'arm64',
            // },
          ],
        },
      },
    },

    // Full list of options: https://v2.quasar.dev/quasar-cli-vite/developing-browser-extensions/configuring-bex
    bex: {
      contentScripts: ['my-content-script'],

      // extendBexScriptsConf (esbuildConf) {}
      // extendBexManifestJson (json) {}
    },
  };
});
