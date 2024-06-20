#  (2024-06-20)


### Bug Fixes

* **frontend | electron:** 成功修复阴影不显示问题, 顺便解决了右侧圆角似乎不全的问题。<采用方案二,单透明原生窗口+纯css的解决方案> ([4694a61](https://github.com/LuSrackhall/KeyTone/commit/4694a613f61f81b417562a3b806e05a7b44aa5b4))
* **frontend | electron:** 修复阴影不显示问题失败。 本次修复失败的原因是。 对于electron窗口机制认识较浅<方案一无法实施>。对css的认识不够深刻<方案二无法实施>。 ([d40fb49](https://github.com/LuSrackhall/KeyTone/commit/d40fb499c45c47a752313f5d5fa7b59bc581b53d))
* **frontend:** 成功修复了展开抽屉后左下角圆角被直角覆盖的问题。<不过打开瞬间的过程中, 还是会出现短暂的直角残留, 不过这个问题不大且考虑为quasar框架的自身问题, 故暂不解决> ([91f2f9f](https://github.com/LuSrackhall/KeyTone/commit/91f2f9fe17b4eb87193e487693c5fc8811755d0f))


### Features

* **frontend | electron:** 尝试通过electron透明窗口+css的方式增加毛玻璃效果失败<因此方案的毛玻璃无法直接透到操作系统桌面>。因此本次仅将主界面背景改为渐变色。 ([3b58b4f](https://github.com/LuSrackhall/KeyTone/commit/3b58b4f0ca77d00d67e4374aa22dbe321b6311e8))
* **frontend | electron:** 初始化electron的入口文件, 在quasar配置文件中配置electron的相关项, 本次主要为win的相关配置。 ([50f22db](https://github.com/LuSrackhall/KeyTone/commit/50f22db02d61ff68728fd9416212ba83ba4e03b3))
* **frontend | electron:** 实现了托盘图标后台功能, 并初始化了托盘菜单。 至此, 正常关闭窗口不再退出应用, 而是进入托盘后台继续运行。 ([a547184](https://github.com/LuSrackhall/KeyTone/commit/a547184df52a0ff7caf2beecad28483807b2623a))
* **frontend | electron:** 使用主界面logo作为应用的临时logo ([069689f](https://github.com/LuSrackhall/KeyTone/commit/069689f9baf49a7296769b0292d6f672d1925881))
* **frontend | electron:** 新增窗口圆角, 不过由于内容占满<已查明确认是此原因>, 阴影无法正常展示。且右侧圆角似乎不全。 ([467ce32](https://github.com/LuSrackhall/KeyTone/commit/467ce329a44a5f587a3d7f580cd3fdfc41af98df))
* **frontend | electron:** 新增制作人展示, 并通过electron的预加载脚本调用原生api, 打开系统默认浏览器展示制作人的github主页。 ([d126b71](https://github.com/LuSrackhall/KeyTone/commit/d126b712193a331c83bcfa3110687f34ea5a8143))
* **frontend:** 在应用名称后, 增加版本号的显示 ([843aea4](https://github.com/LuSrackhall/KeyTone/commit/843aea47c91b804a0fd64539adc8a547bbb314c4))
* **frontend:** 主界面暂时仅使用图标填充即可, 最初套壳版本无需任何多余功能。 ([1aa79f5](https://github.com/LuSrackhall/KeyTone/commit/1aa79f51f2c519e84f27aa9b7eef22a5180816f1))
* **mvp:** 1.实现了键盘的全局监听功能。2.实现了音频播放功能。3.实现了mvp键盘音的核心功能逻辑。 ([#1](https://github.com/LuSrackhall/KeyTone/issues/1)) ([e30ea76](https://github.com/LuSrackhall/KeyTone/commit/e30ea7656d72a7b40afc129e52bc5fa961e9c98a))



