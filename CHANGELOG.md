#  (2025-03-14)


### Bug Fixes

* **keytonepkg:** 修复了主页面中主动清空所选键音包后, 整个应用的键音包使用状况没有恢复默认内嵌音的bug。 ([afb7371](https://github.com/LuSrackhall/KeyTone/commit/afb73713649423516b13b3bfb2ada8852bbcb6a0)), closes [#65](https://github.com/LuSrackhall/KeyTone/issues/65)
* **volume:** 对于预览音量的操作, 应该基于音频文件原本的声音去预览, 不应该受到全局音量设置的影响。本次提交修复了这个问题。 ([1d0dfef](https://github.com/LuSrackhall/KeyTone/commit/1d0dfeffb8291fb6cd0931815bd7a9272b4c9b0c)), closes [#64](https://github.com/LuSrackhall/KeyTone/issues/64)



#  (2025-02-27)


### Bug Fixes

* 1.让所有弹出的对话框, 都向左偏移。2.补上最近几次提交中的i18n遗漏。 ([5f51e3f](https://github.com/LuSrackhall/KeyTone/commit/5f51e3ff66646b58d8b30644b7a802712520dc03))
* 修复了键音专辑页面中, 至臻键音菜单中对播放模式选项相关的国际化的不兼容问题。 ([35b4149](https://github.com/LuSrackhall/KeyTone/commit/35b414924ab8d433828d704d8a2b53438ff172ef))
* **i18n:** 补充国际化配置未覆盖的地方。 ([dff2413](https://github.com/LuSrackhall/KeyTone/commit/dff2413813acdca907aea923f760bcafd84a96cb))
* **i18n:** 补充最近提交中关于i18n翻译的遗漏。 ([9e69272](https://github.com/LuSrackhall/KeyTone/commit/9e69272dc4defd51edf6f5a4198ea8892eaa82dc))
* **i18n:** 更新国际化文本中, 隐私政策和用户协议文本，并添加空状态提示信息文本 ([d8095f7](https://github.com/LuSrackhall/KeyTone/commit/d8095f78edfc43b9cccf9fd1fd4eb5699783f078))
* **i18n:** 更新键音专辑页面的英文国际化文本，简化提示信息, 防止与中文的样式差得太多。 ([51b2795](https://github.com/LuSrackhall/KeyTone/commit/51b27955c97b971d4bed4c4c9152eab521601833))
* **i18n:** 国际化文本中, 添加键音专辑页面的国际化文本，包括新建、导入、导出和删除专辑的提示信息 ([bbaeae5](https://github.com/LuSrackhall/KeyTone/commit/bbaeae5ba41736b453a0af9fd9ee1a1f073f3846))
* **i18n:** 为键音专辑组件配置国际化, 本次提交配置了 按键联动声效 这一步骤中的 单键局部配置中单键增设对话框的国际化文本。(其中涉及的notify除外)(修改某个单键的国际化文本未配置) ([1ad9afd](https://github.com/LuSrackhall/KeyTone/commit/1ad9afd10db850d53db99652ab11cd5719746380))
* **i18n:** 为键音专辑组件配置国际化, 本次提交配置了 按键联动声效 这一步骤中的 单键局部配置中修改某个单键的国际化文本。(其中涉及的notify除外) ([ff852ef](https://github.com/LuSrackhall/KeyTone/commit/ff852ef4c2a2d5ef335041f8b3a59133521be982))
* **i18n:** 为键音专辑组件配置国际化, 本次提交配置了 按键联动声效 这一步骤中的 全局配置的 国际化文本。(单键局部配置的国际化后续整) ([48a59fc](https://github.com/LuSrackhall/KeyTone/commit/48a59fcab3f31daeda7ba11b253be502e31ef036))
* **i18n:** 为键音专辑组件配置国际化, 本次提交主要是为了补全相关通知内容的国际化配置。 ([06b0532](https://github.com/LuSrackhall/KeyTone/commit/06b05328972ae273defe661886da7d4bc768ac9d))
* **i18n:** 为键音专辑组件配置国际化。本次提交配置了 至臻键音 这一步骤中, 编辑已有至臻键音按钮及相关对话框内部国际化。(内部二重对话框内播放模式选项的相关国际化,因报错问题暂未适配) ([3728586](https://github.com/LuSrackhall/KeyTone/commit/3728586c8d4dffe3f775d9632a7346e971faa34a))
* **i18n:** 为键音专辑组件配置国际化。本次提交配置了 至臻键音 这一步骤中, 制作新的至臻键音按钮及相关对话框内部的国际化。(tips: 编辑已有至臻键音按钮及相关对话框内部国际化未适配。) ([874ad52](https://github.com/LuSrackhall/KeyTone/commit/874ad5243de4adbeeebd44c909845915b091850d))
* **i18n:** 为键音专辑组件配置国际化。本次提交配置了 至臻键音 这一步骤中,辑已有至臻键音按钮及相关对话框内部二重对话框内播放模式选项的相关国际化。(但仍有问题, 此设计似乎需要重构) ([b6a709a](https://github.com/LuSrackhall/KeyTone/commit/b6a709a7b269c6a40bf707efbe3cf51e23e5c410))
* **i18n:** 为键音专辑组件配置国际化。本次提交配置了裁剪定义声音这一步骤的相关国际化文本。 ([a6e217c](https://github.com/LuSrackhall/KeyTone/commit/a6e217cdfe639d7f310286169adb3ad1e8e6d89c))
* **i18n:** 为键音专辑组件配置国际化。本次提交配置了载入音频源文件这一步骤的相关国际化文本。 ([f01c7fe](https://github.com/LuSrackhall/KeyTone/commit/f01c7fecf1dde1f0da39a0dbb10abcb5bfaef549))
* **i18n:** 再次复查键音专辑页面的国际化文本, 做到一个页面一个结构体，新增了标签、空状态和通知信息。 ([4b4cda6](https://github.com/LuSrackhall/KeyTone/commit/4b4cda68118c88de4923ad3db8bc81186d75eb95))
* **ui:** 将打开对话框后背景蒙版的透明度调整为完全透明, 这样既能够保证蒙版的功能, 又不会因蒙版的存在带来视觉体验的下降(特别是仅存在于本软件中的蒙版溢出问题) ([6eb732f](https://github.com/LuSrackhall/KeyTone/commit/6eb732f1135740932db43956d1e43071a1ee857d))
* **ui:** 修复国际化没覆盖到的地方, 并顺便重构此部分代码的沉余情况。 ([35258a3](https://github.com/LuSrackhall/KeyTone/commit/35258a3fc5fd497d445434bb2c813b83780cf568))
* **ui:** 修复设置页面样式异常的问题。 ([ceda084](https://github.com/LuSrackhall/KeyTone/commit/ceda084d403a3f07e6887a07ca126c22bc1779fa))



#  (2025-02-17)


### Bug Fixes

* 对导出专辑 按钮的具体功能逻辑, 进行了重构。直接依赖前端的下载功能。 不再使用上传某个文件来获取路径的方式(需要依赖electron的api)交由sdk来保存。 ([225aeef](https://github.com/LuSrackhall/KeyTone/commit/225aeef6f1a017e199cc23220a4ebf3581d8c01e))
* 修复 导出专辑 时, 压缩的文件夹改变了原始文件夹根目录名称的bug ([968654c](https://github.com/LuSrackhall/KeyTone/commit/968654c03f835f47676ccbf99177945426cad1cd))
* 修复导出专辑时, 用户还未选择保存的路径, 就弹窗导出成功通知的bug。(本次修复选择使用 File System Access API 替代当前的 download 属性方案) ([f751aa4](https://github.com/LuSrackhall/KeyTone/commit/f751aa4871a5c694c52ad669a7142871f7feb474)), closes [#61](https://github.com/LuSrackhall/KeyTone/issues/61)
* 修复了带芯片的选择器组件, 受h-5.8影响, 无法显示完全, 以及显示纵向滚动条的bug。 ([d6906d5](https://github.com/LuSrackhall/KeyTone/commit/d6906d5b65bd43d18505faf50aeaa956a64eb926))
* 修复了上传专辑覆盖现有专辑时, 若被覆盖是已选择专辑时, ui没有更新的bug, 且优化了更新时机。另外,  添加获取专辑文件元数据的功能，支持上传并验证专辑文件格式，返回相关元数据信息。 ([32b81f8](https://github.com/LuSrackhall/KeyTone/commit/32b81f8abaf0a391837e8dcae675aec84545c2ac))
* 修复删除专辑时, 若删除的是最后一个专辑, 响应的专辑列表仍有不可用残留的bug。 ([be7ebdb](https://github.com/LuSrackhall/KeyTone/commit/be7ebdb135048db032de418cefefba7d612fbb61))
* **frontend:** 修复了潜在的bug。即当用户所选的键音包受到外部破坏(比如被恶意删除)后, 在重新加载页面时仍会看到一个名称为空白的假性不可用键音包的bug。 ([372f8b7](https://github.com/LuSrackhall/KeyTone/commit/372f8b762541f4bea15c15bfc6387152b43ef05b))
* **sdk:** 调整临界区, 以继续根除sdk可能产生的死锁或其它意外情况。 ([c5b48e9](https://github.com/LuSrackhall/KeyTone/commit/c5b48e9619f8f6d0eb0ab21f7c960abcc132e633))
* **sdk:** 对于所加载键音包路径不存在的情况, 应该回收Viper这个全局变量, 以免发生未知的错误。 ([3b40784](https://github.com/LuSrackhall/KeyTone/commit/3b407840275532050cffa68bb34d34eac0f2ad76))
* **sdk:** 修复了因上次修复引发的新bug。即当Viper==nil时造成panic引发sdk崩溃的问题。 ([21908e6](https://github.com/LuSrackhall/KeyTone/commit/21908e6f71e1cd8683a79517f334caf7a1b5856f))
* **sdk:** 修复了sdk的panic问题。重新梳理了键音包模块的生命周期, 以及并发的临界区, 删除了不需要的判断, 以及对必要的临界区进行上锁操作, 以避免panic的风险。 ([226384b](https://github.com/LuSrackhall/KeyTone/commit/226384b6eb888eca0276e2488be41e3b33d5fcae))
* **ui:** 1.修复了编译已有声音的选择框, 在没有任何选择时, 样式过窄的问题;2 修复了 选择后, 内部关于源文件的选择框 名字过长时的溢出问题。 ([8f48a5a](https://github.com/LuSrackhall/KeyTone/commit/8f48a5a44571e30d445fab695e19f1709d3079dc))
* **ui:** 进一步修复了上个提交中的问题。本次修复位置迁移值sse回调是由于历史问题的积累造成的->函数initData()中初次主动获取并初始化这一操作的内部生命周期,因长期忽视维护, 使得目前变得难以管理。(若后续有相关bug报告, 则仍需迁回处理, 且需要处理的不止此一项) ([06961ce](https://github.com/LuSrackhall/KeyTone/commit/06961ce059ceaa4a732b28f7937d411345bd6f4e)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)
* **ui:** 进一步修复前两个提交中的问题。本次修复是因为当全键已选择声音存在源文件(soundFileList)时, 会由于源文件的生命周期更新延迟造成错误的问题(上上个提交的注释中介绍过)。(解决方案:直接在initData中完善此字段的初始化以补全其生命周期。) ([ce1d101](https://github.com/LuSrackhall/KeyTone/commit/ce1d101d6762f25806c6c2a5e3d6dfcec6bd1ca1)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)
* **ui:** 临时修复了 编辑键音专辑时, 按键联动声效的全键声效设置对话框中, 相关配置初始化缺失不符合真实配置的问题。 ([55be04c](https://github.com/LuSrackhall/KeyTone/commit/55be04cf7575aec449d17377814e6ba6f2f29924)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)
* **ui:** 明确键音专辑初始化的必要内容,明确其生命周期,保证连续创建键音包时,不会错误的中止正在初始化的键音专辑。同时对于初始化的键音专辑,在初始化完成之前不予展示(此时展示加载动画)。 ([18b59d5](https://github.com/LuSrackhall/KeyTone/commit/18b59d53d9e19fc1f1c37e517007de0d62c348fc))
* **ui:** 通过deep,进一步修复了主页面中, 选择键音名称过长时引发界面异常的问题。(本次修复主要解决了label溢出过了选择框内关闭按钮后, 且滚动时会带着选择框整体一起滚动的问题) ([ab27938](https://github.com/LuSrackhall/KeyTone/commit/ab2793877d570cf871d48a7b21b1adf73165bd10))
* **ui:** 修复键音专辑页面上唯一的选择组件, 因所选键音专辑名称过长而造成的横向页面溢出问题。(这是与刚刚修复的主页面的选择组件相同的问题) ([ab9a6a4](https://github.com/LuSrackhall/KeyTone/commit/ab9a6a4b61f23d2a2012b6f3a7ceae6ca96b3768))
* **ui:** 修复键音专辑页面中,切换键音包后,发生的错误的初始化问题。(本次错误是因为: 漏掉了持久化中 按下测试音/抬起测试音 相关字段的初始读取, 并错误的按照默认数据给重新进行了初始化) ([3d705a9](https://github.com/LuSrackhall/KeyTone/commit/3d705a996128c1e9a8399d8cf2b5637630f9ff01)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)
* **ui:** 修复了 编辑对话框(如音频文件、声音、按键音)中, 因右上角名称过长造成的溢出问题。(对声音和按键音的问题通过移除组件处理, 对音频文件采取仅显示后缀名称来处理。) ([ba67f28](https://github.com/LuSrackhall/KeyTone/commit/ba67f28619d9ef52bb763fc948d2da0129b04f64))
* **ui:** 修复了按键联动声效步骤条在关闭后, 对当前步骤是否被操作过会作出错误判断的问题。 ([385608d](https://github.com/LuSrackhall/KeyTone/commit/385608dada66c5a1d69adc18f563a0b988573ae9)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)
* **ui:** 修复了创建键音包时, 初始化过程被打断而造成创建的键音包不可用的bug。 ([3aeb848](https://github.com/LuSrackhall/KeyTone/commit/3aeb848564f2c80f4cf00bb6e3e8dd52db49ed6b))
* **ui:** 修复了从选择框删除所选键音专辑后, 键音专辑组件未随之消失的bug。 本次提交中还包含了预防键音专辑组件在创建新专辑过程中的重复渲染bug。 ([53213af](https://github.com/LuSrackhall/KeyTone/commit/53213afc6d618e1a4ac63b5eb3415123e1c0d4ae))
* **ui:** 修复了键音专辑页面, 在切换键音包后失去滚动监听的bug。这个bug造成了预期的向下滚动触发收起按钮 以及 在顶部向上滚动触发展开按钮的触发 的逻辑 变得失效。 ([f2749e7](https://github.com/LuSrackhall/KeyTone/commit/f2749e784f94fd5b149bfb8d4306d70e853762af)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)
* **ui:** 修复了键音专辑页面顶部选择器在某些情况下被下方键音专辑组件覆盖的问题。复现方式->对滚动到底后, 点击展开按钮。 ([53533a0](https://github.com/LuSrackhall/KeyTone/commit/53533a0dd8f744ee046e2547297c9acfdb121a90))
* **ui:** 修复了可能存在的潜在的内存泄漏问题。(对键音专辑组件中, 所创建的侦听器, 在组件卸载前进行删除操作) ([b4ecfbd](https://github.com/LuSrackhall/KeyTone/commit/b4ecfbdcd47e4d7f017aa71c6565ef4291a41af0))
* **ui:** 修复了连续创建键音包时的bug。这个bug会使得从创建第二个键音包开始, 将新内容错误的赋给上一个键音包。 ([665b14c](https://github.com/LuSrackhall/KeyTone/commit/665b14c3327795be3ea7acf7e70e8bab23290163))
* **ui:** 修复了刷新页面后造成的之前 为 quasar选择器菜单中的选项, 添加的细微滚动条消失, 退回换行策略的bug。 ([6743096](https://github.com/LuSrackhall/KeyTone/commit/6743096f049807a3bef3015d1f598ee7b844cfb8))
* **ui:** 修复了刷新页面后造成的之前 为 quasar选择器添加的滚动条 消失且样式重新溢出的bug。 ([1f612b6](https://github.com/LuSrackhall/KeyTone/commit/1f612b61dd066ae8b450dbcd27d6caab517adb0e))
* **ui:** 修复了音频专辑组件中, 步骤框内部, 普通选择框的文字溢出问题。 对于带芯片的多选框, 将原有省略号策略改为溢出滚动策略。 ([35d84b8](https://github.com/LuSrackhall/KeyTone/commit/35d84b85f0bd65486fd478726274a5958a8284a8))
* **ui:** 修复了音频专辑组件中, 步骤框上方的专辑名称显示会因名称过长而溢出的问题。(采用不换行+溢出滚动的策略)(注意, 对于quasar的输入器来说, 不换行+溢出滚动的策略是自带的) ([3e3fc50](https://github.com/LuSrackhall/KeyTone/commit/3e3fc502b18164b492629b91988c819aa7b81575))
* **ui:** 修复了q-select选择器的弹出菜单, 普遍存在的当选项名称过长时的溢出问题。另外, 本次提交还统一了q-select选择其弹出菜单的滚动条样式, 弃用了默认的经典滚动条样式。 ([b63cdde](https://github.com/LuSrackhall/KeyTone/commit/b63cdde8d2befe8123ccee9aa7f81da8645dff7f))
* **ui:** 修复了q-select选择器的弹出菜单, 普遍存在的英文内容过长时超出部分无法正常显示的bug。(中文会自动换行) 本次修复方案是, 使得所有字体的内容都遵循超出自动换行的策略。 ([36ecac8](https://github.com/LuSrackhall/KeyTone/commit/36ecac8c9c8535738eb05a3d8f542690e861b62a))
* **ui:** 修复主页面q-select选择器的弹出菜单(目前使用的对话框模式)中, 最上方已选择部分因所选项的名称过长的超出部分溢出的问题。 ([d106b98](https://github.com/LuSrackhall/KeyTone/commit/d106b98204ed9c375b635bb0dfb39ab357c82fbb))
* **ui:** 修复主页中, 选择的键音专辑名称过长时引发的页面样式异常的bug。(采取不换行+溢出滚动的策略) ([2c81c63](https://github.com/LuSrackhall/KeyTone/commit/2c81c6340b957624096c9640786409676e1369cd))
* **ui:** 在创建键音包期间, 禁用选择器, 以避免引发类似连续创建时的问题。 ([88bb2b2](https://github.com/LuSrackhall/KeyTone/commit/88bb2b235ce49751503d3af3d583a63f40791da2))
* **ui:** 在创建键音包期间, 应禁止可能造成创建中止问题的一些按钮。 (不过做的可能有些过度了, 不过先加上, 后续再考虑是否删除吧) ([6ac55cb](https://github.com/LuSrackhall/KeyTone/commit/6ac55cb5cddcf8fea311e299b1c1bce910f35381))
* **ui:** 在键音专辑页面中, 删除按键联动音效设置中单键声效设置对话框的不符合设计的"确定"和"取消"按钮, 改为更合适单独的"Close"按钮。 ([8a31d4e](https://github.com/LuSrackhall/KeyTone/commit/8a31d4e05b5baae7a6a35504621ff4f6b0cffe09)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)


### Features

* 实现 .ktalbum 文件格式，包含加密和验证 ([575351a](https://github.com/LuSrackhall/KeyTone/commit/575351a64022a8f193d1c7d2f3b8d75eea0f0bd2))
* 实现了 导出专辑 按钮的具体功能。 ([485b8df](https://github.com/LuSrackhall/KeyTone/commit/485b8df361b2224e376bd786fb14d91696a61df3))
* 实现了 导入专辑 按钮的实际功能。 ([6fe8936](https://github.com/LuSrackhall/KeyTone/commit/6fe8936c290db5335171f6ba7b5c6a7b16183e8a))
* 添加专辑覆盖导入功能，处理已存在专辑的情况。现在可以在专辑已存在时, 选择覆盖导入或是取消导入, 而不是直接报错。 ([ef50487](https://github.com/LuSrackhall/KeyTone/commit/ef504875ab655f8885497d81711492bb3a6277e1))
* 添加专辑元数据支持，优化导入专辑验证逻辑 ([1d58f48](https://github.com/LuSrackhall/KeyTone/commit/1d58f48ff08d3e9d46257bc638e9291ad563208a))
* 完善键音专辑页面, 删除当前专辑的按钮的功能。 ([caefec4](https://github.com/LuSrackhall/KeyTone/commit/caefec45442675767c2b02560a1b35fb6acad1cc))
* 为导入专辑遇到已存在的情况, 添加导入专辑为新专辑的功能，允许用户在导入时选择覆盖、保存为新专辑或取消导入。 ([865eaa3](https://github.com/LuSrackhall/KeyTone/commit/865eaa3092c85bb304e7364d8b566d0236ded847))
* 为键音专辑页面, 添加键音专辑选择器的空状态提示和相关样式，优化用户体验。 ([2db2e67](https://github.com/LuSrackhall/KeyTone/commit/2db2e671a351fe9674f6c55d3c6e4199c9c62490))
* 为主页面, 添加键音专辑选择器的空状态提示和导航功能，优化用户体验 ([09620b3](https://github.com/LuSrackhall/KeyTone/commit/09620b31d8a8d48e9dbbd75fdc874055fc50bdb4))
* 重新实现新设计的键音专辑页面中, 创建键音包按钮的功能。 ([c391399](https://github.com/LuSrackhall/KeyTone/commit/c3913999a29e54d769093be3b9300efbb3f6d865))
* **tools:** 初始化 ktalbum-tools CLI 工具，实现文件解包功能 ([640a0de](https://github.com/LuSrackhall/KeyTone/commit/640a0dece640f5ecc63f279aa7214d838a4e375f))
* **tools:** 添加 Web 服务和前端界面到 ktalbum-tools ([3cb5bb2](https://github.com/LuSrackhall/KeyTone/commit/3cb5bb229b6893db5aecbbed18c4dbe8f46c61cc))
* **tools:** 增强文件信息检索并添加文件完整性检查 ([a22b2b8](https://github.com/LuSrackhall/KeyTone/commit/a22b2b83830c4ca61018f763e0e4849c30d13178))
* **ui:** 对quasar选择器菜单中的选项, 不再采取换行方案, 而是改为滚动方案, 并新增了对应样式的细微滚动条, 以方便用户查看。 ([44dc336](https://github.com/LuSrackhall/KeyTone/commit/44dc336687f45e4c1f9cc4e2fd4012864deafdda))
* **ui:** 放弃旧的键音专辑页面的设计, 改用新设计的键音专辑页面。 ([a6b4d14](https://github.com/LuSrackhall/KeyTone/commit/a6b4d14349aee00424e6fb3a7661829316372229)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)
* **ui:** 为键音专辑组件中, 编辑对话框(如音频文件、声音、按键音)的管理选择器, 添加清除所选项的功能。 ([5e9f601](https://github.com/LuSrackhall/KeyTone/commit/5e9f601a42e4e718ebce89b17eb90c74fdcbb1da))
* **ui:** 为所有页面中的 quasar 选择器 的 输入框 添加 符合 选择器样式的 滚动条。 ([5f14469](https://github.com/LuSrackhall/KeyTone/commit/5f144692092e12398b7b3807a0efa833cd0a838e))
* **ui:** 新增 键音专辑的  创建、导入、导出、删除 等四个功能按钮。 ([dcc3a1a](https://github.com/LuSrackhall/KeyTone/commit/dcc3a1a6c07dd292847e474f18e27663e3b4b92a))


### Reverts

* Revert "chore(ui): 因上方组件收起会引发下方组件向上靠拢, 而这个靠拢是没有动画的, 本次提交就是为其添加动画。" ([9ea8478](https://github.com/LuSrackhall/KeyTone/commit/9ea8478a576fb20d399b1b5bdf4ad336026c3837)), closes [#55](https://github.com/LuSrackhall/KeyTone/issues/55)



#  (2025-01-29)


### Bug Fixes

* **electron | ui:** 关于应用商店版本 自动启动功能 的进一步修复。 改为在打包时引入并默认启用, 并在设置界面引导用户前往系统 设置>应用>启动 中, 关闭/打开自启动功能。 ([e03c9a9](https://github.com/LuSrackhall/KeyTone/commit/e03c9a91768fb516025846fcef8d79cadf444f0f))
* **ui:** 当音量为0%时, 静音按钮应该保持关闭状态。 ([30c04fb](https://github.com/LuSrackhall/KeyTone/commit/30c04fbb43d20adcb75237dfc12559ba050d3309))
* **ui:** 修复了主界面音量滑块通过单击来调整后, 会被自动退回更改前的位置的bug。(此bug发生的原因是: 对需要实时持久化的变量数据, 仅单变量变更的操作会因通信延迟而天然可靠, 但多个变量变更的操作破坏了这个特性)(解决方式, 在次要变量的变化源头手动制造延迟来使其重新变得可靠) ([4aa721f](https://github.com/LuSrackhall/KeyTone/commit/4aa721f9a4f2dc916c877f47c77a6af02bc122db))



#  (2025-01-28)


### Bug Fixes

* **electron:** 修复了当开启自动启动时隐藏窗口的设置后,会引发sdkIsRun属性无法更新而造成的托盘i18n等功能无效的bug。同时本次提交中将'启动时是否隐藏窗口'的依赖改为nodejs。 ([23f99c6](https://github.com/LuSrackhall/KeyTone/commit/23f99c60f923f671d320706e3d4cc95358deface))
* **electron:** 修复了当启动时隐藏窗口设置项主动 做开关操作并最终设置为开启后, 二次启动仍会展示窗口的bug(预期是仅聚焦)。(问题原因仍然是viper的老毛病,) ([5cc191a](https://github.com/LuSrackhall/KeyTone/commit/5cc191aa509529685880574aeb842ed7e3ce6c69))
* **sse:** 修复潜在的因端口占用问题造成的sse的链接不起作用的问题。 ([c2041ac](https://github.com/LuSrackhall/KeyTone/commit/c2041ac10ff3a84a8d31702d1d8626af6db86919))
* **ui | electron | sdk:** 修复潜在的因端口占用问题造成的electron中ui及node主进程的restful请求无法访问sdk实际端口的问题(sse暂未涉及)(spa除外) ([9b2e705](https://github.com/LuSrackhall/KeyTone/commit/9b2e70533c43ba412047c972e39878b6f347067a))
* **ui | electron:** 修复了 启动与自动启动设置项 在appx格式下 完全不显示的bug。 ([72ea0f0](https://github.com/LuSrackhall/KeyTone/commit/72ea0f085be0b5cc99dd95d8cf2185dcb932bf5e)), closes [#58](https://github.com/LuSrackhall/KeyTone/issues/58)
* **ui:** 解决潜在的 过分监听端口变化引起的 资源消耗的问题。 ([b7f8c6a](https://github.com/LuSrackhall/KeyTone/commit/b7f8c6a0db86365ee103605a12fdc0193213b289))
* **ui:** 修复调整 音量增减调节 的音量增减幅度后,回到主页面时 主页面的音量选择的百分比变化的bug,修复后会保持不变。(本次修复包括实时调整时的百分比保持不变--不过有极小机率会发生变化) ([1100566](https://github.com/LuSrackhall/KeyTone/commit/1100566dfcb2a43c67e2bbbc313898b4b7e48d93))
* **ui:** 修复了和之前两次修复相关的, 当进入设置页面 更改主页相关设置的音量降幅后, 回到主页面时 主页面音量正常变化。紧接着进入设置页面 更改音量增减调节的音量增减幅度后, 回到主页面时 主页面百分比异常的按照由降幅引起的变化之前的旧百分比保持的bug。 ([5f4a4fd](https://github.com/LuSrackhall/KeyTone/commit/5f4a4fd503b1e6c09def8f1ac0f803f4b1e79870))
* **ui:** 修复了主页面调整音量的进度条在使用时, 调整至0%后声音不会完全消失的bug。(增加了音量为0%时自动开启静音功能的逻辑) ([df97c32](https://github.com/LuSrackhall/KeyTone/commit/df97c324bd9e8f20149635eac16d664cf6aa05a3))
* **ui:** 修复因上个修复引起的 修改主页相关设置的音量降幅后,回到主页面时 主页面的音量选择百分比不变的bug,修复后会正常变化。(本次修复无需包括实时调整时的情况,因实时调整不存在此bug) ([9aa1b7e](https://github.com/LuSrackhall/KeyTone/commit/9aa1b7efd53e214c981e25775da3023a2cfa45dd))
* **ui:** 修复主页 volume 音量进度条, 在min绝对值由大变小时, 若进度条的所选值(即我定义的Normal值)比min更靠左就出现负数百分比的bug。(min受 设置页面中 主页面相关音量降幅 及 原始音量增减调节 的影响) ([f852da1](https://github.com/LuSrackhall/KeyTone/commit/f852da1c0e2f1643e0ec82e31e76dcf1a2d353cb))



#  (2025-01-20)


### Bug Fixes

* **ui | electron | appx:** 将win商店和桌面快捷方式下的图标背景, 改回白色。 ([0bfb32e](https://github.com/LuSrackhall/KeyTone/commit/0bfb32ebb727785b26bf3b5c74b0cb46632a091c)), closes [#57](https://github.com/LuSrackhall/KeyTone/issues/57)
* **ui | electron:** 修复了应用页面内部的url路由在通过`ctrl+鼠标左键`点击后, 触发新的electron窗口展示的问题。(通过监听mainWindow下的新窗口的产生事件, 在新窗口产生前判断是否是预期的产生行为, 若不是则纠正) ([aad9f3d](https://github.com/LuSrackhall/KeyTone/commit/aad9f3d1482e490f5206545b13941df2b6e5da0e)), closes [#40](https://github.com/LuSrackhall/KeyTone/issues/40)
* **ui | electron:** 由于appx格式的自启动功能暂时修复失败, 故在设置页面暂时移除相关设置项。(仅暂时移除appx打包格式下自启动设置) ([8bed9d4](https://github.com/LuSrackhall/KeyTone/commit/8bed9d4dc8e5a98b8c96299ce27934dadfc7f291)), closes [#58](https://github.com/LuSrackhall/KeyTone/issues/58)



#  (2025-01-17)


### Bug Fixes

* **electron:** 再次修复微软商店版本自启动相关设置无效的问题。 ([d100718](https://github.com/LuSrackhall/KeyTone/commit/d100718868938bb29fdff42dfb6ff66e7f33bab2))
* **ui | electron | appx:** 修复图标背景在win商店和桌面快捷方式时, 为白色的问题。(预期是按照图片中正常的透明色来展示) ([99477f8](https://github.com/LuSrackhall/KeyTone/commit/99477f8e44c7155c9f08b9c6705d022477f166ee)), closes [#57](https://github.com/LuSrackhall/KeyTone/issues/57)
* **ui:** 更改设置界面原有的点标题进入整页的默认行为, 改为双击进入, 原有的单击标题后的行为改为与单击箭头一致的展开设置。 ([9f00684](https://github.com/LuSrackhall/KeyTone/commit/9f00684dff011742b4bca32a42c833de27a256d5))
* **ui:** 修复设置界面改为双击进入整页行为后, 容易被单击误触进入的问题。(放弃了默认的双击事件, 改用单击事件加自定义双击函数处理) ([09d7428](https://github.com/LuSrackhall/KeyTone/commit/09d7428c9da8a0eb9cdf8091688722c1b4f5f3f7))
* **ui:** 修复设置页面中, 主页相关设置的图标错用音量图标的问题。(虽然目前里面的设置确实仅与音量相关) ([87be226](https://github.com/LuSrackhall/KeyTone/commit/87be2264f955fa1bc79a688629ecd04ee472827c))



#  (2025-01-16)


### Bug Fixes

* **build:** 解决构建appx包时, 图标缺失的问题。本次提交中利用quasar的icongenie图标处理工具的自定义配置, 在正确的路径下生成了所需的图标资源。 ([033eacb](https://github.com/LuSrackhall/KeyTone/commit/033eacbf9d0a4d66f62a9b8f144efebe5c478338)), closes [#57](https://github.com/LuSrackhall/KeyTone/issues/57)
* **electron:** 修复了微软商店版本自启动相关设置无效的问题 ([ac10af0](https://github.com/LuSrackhall/KeyTone/commit/ac10af0962e970ee0789a01f566ba684bafad247))
* **ui | electron | appx:** 尝试修复图标尺寸即背景与exe版本不一致的问题 ([6c6e5f5](https://github.com/LuSrackhall/KeyTone/commit/6c6e5f54e3f2f865339ff96554de51e3f2ea2f2c)), closes [#57](https://github.com/LuSrackhall/KeyTone/issues/57)
* **ui | electron | refactor:** 重构了开机自启动相关设置项的逻辑。修复了必须重启应用才能声效的问题(目前设置后就会立即生效)。 ([b79430c](https://github.com/LuSrackhall/KeyTone/commit/b79430c3b86d00c7a5920f520d64103e5d04396c))
* **ui:** 修复了顶部导航栏在路由名称过长时因换行而造成的样式异常的bug。 ([c7eae76](https://github.com/LuSrackhall/KeyTone/commit/c7eae76b208de46b9836cfddacd4e5629257dbb3))



#  (2025-01-07)


### Bug Fixes

* **sdk:** 经验证, 即使键音包不存在, 也不会引发报错之类的问题, 而初始化时所选键音包为空字符串的情况实际上也可以归为此类, 因此sdk中无需对空字符串的restful返回错误。 ([75d8dce](https://github.com/LuSrackhall/KeyTone/commit/75d8dce448abf175bf124b82ddd865602359a30a)), closes [#56](https://github.com/LuSrackhall/KeyTone/issues/56)
* **ui:** 解决在主界面中当用户选择某个键音包后, 焦点仍保持在选择组件而造成的容易引发误触的问题。 ([c45c4bd](https://github.com/LuSrackhall/KeyTone/commit/c45c4bdb7350b0e3812f984fca259f4c22574228)), closes [#56](https://github.com/LuSrackhall/KeyTone/issues/56)
* **ui:** 修复了 键音包列表 为空时 引发的界面卡死问题。 ([e2dbf2d](https://github.com/LuSrackhall/KeyTone/commit/e2dbf2d8b2b55ae3493321cafc33968081c6c818)), closes [#56](https://github.com/LuSrackhall/KeyTone/issues/56)
* **ui:** 修复了 用户新建键音包并返回主页面后,  键音包列表未能更新至最新的问题。 并且本次提交中, 将键音包列表的初始化逻辑移动至其本该在的App.vue文件中。 ([6f94635](https://github.com/LuSrackhall/KeyTone/commit/6f9463557b77f56d249eb2a2797d7bdd60a78959)), closes [#56](https://github.com/LuSrackhall/KeyTone/issues/56)
* **ui:** 在主界面上提供相应的 快速的清楚所选键音包 的按钮, 以解决当前选择键音包后, 无法清空所选键音包 至 软件初始状态的问题。 ([8db1a47](https://github.com/LuSrackhall/KeyTone/commit/8db1a47085019e0df3feb6ebb0b21544c19d418a)), closes [#56](https://github.com/LuSrackhall/KeyTone/issues/56)



#  (2025-01-07)


### Bug Fixes

* 解决 全键声效设置 时, 无法对按下或抬起的全键声效 设置空值, 以至于无法单独地仅定义 按下声效 或 抬起声效, 甚至是全都取消定义。解决方式是, 使用空字符串。 ([a9d56b8](https://github.com/LuSrackhall/KeyTone/commit/a9d56b80b17bf8e6e0d442dc4d5ab2b736bec4ca)), closes [#45](https://github.com/LuSrackhall/KeyTone/issues/45)
* 修复了 单键Dik码 与 name 实时映射是,  多个按键被按下 时的通知 在应用未获取焦点时被意外触发的bug。 ([1c51e42](https://github.com/LuSrackhall/KeyTone/commit/1c51e4271c1e781ed1cb8c62b08d70ed94d3c941)), closes [#47](https://github.com/LuSrackhall/KeyTone/issues/47)
* **sdk:** 当音频包内某个 声音 被删除被, 在播放过程中不应该panic, 仅通过log记录对应的错误事件即可。 ([36502ba](https://github.com/LuSrackhall/KeyTone/commit/36502ba55f83fd69617f98ef52b6fb43cc97ab0d)), closes [#52](https://github.com/LuSrackhall/KeyTone/issues/52)
* **sdk:** 当音频包内某个源文件被删除后, 在播放过程中不应该panic, 仅通过log记录对应的事件即可。(后续需要做进一步处理,如提示前端或是继续判断之后优先级的声效直到嵌入测试音的播放) ([0caa183](https://github.com/LuSrackhall/KeyTone/commit/0caa183ca6bfc06d6b75cae7a79c241c6df00b72)), closes [#51](https://github.com/LuSrackhall/KeyTone/issues/51)
* **sdk:** 还是应该正视内存泄漏问题--其中一部分好像是goroutine太多引起的<可用协程池解决>, 另一部分更大泄漏的可能和beep有关。取巧的方式来清理内存, 会造成播放后的声音无法正常走向结束,从而造成更大的内存泄漏。采取手动结束的方式, 又会造成应该发出的声音中断的不佳体验。 ([cefcbae](https://github.com/LuSrackhall/KeyTone/commit/cefcbae0a25f66454242c3ead34045856c797656))
* **sdk:** 解决删除源文件时,相关sha256下无name时,对应的sha256字段与真实音频源文件没有被删除的问题。但目前在真实源文件路径名称处存在小问题,没能获得audioPkgUUID。 ([efd4d64](https://github.com/LuSrackhall/KeyTone/commit/efd4d645d949e8012d9efa8363a74417bd66d127)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **sdk:** 通过补足参数信息的方式, 彻底解决了删除文件功能无法删除实际文件和对应sha256配置字段的问题。 ([6b81822](https://github.com/LuSrackhall/KeyTone/commit/6b8182275b7d6d530371e5d600beace1bbab84cf)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **sdk:** 新版本的beep在一定程度上修复了内存泄漏问题, 因此升级beep依赖版本, 并着手解决本项目中的相关问题。 ([ed5c59e](https://github.com/LuSrackhall/KeyTone/commit/ed5c59e14352f01ab9b31149be2b198a67ea97b5))
* **sdk:** 修复了内存泄漏问题! ([4df2433](https://github.com/LuSrackhall/KeyTone/commit/4df2433bba9262c820f5bc65c031b209ef6ae8d9))
* **sdk:** 修复了首次启动[特指无配置文件时的启动]时, 无法正确操纵配置项的bug。并且修复了操纵配置项时, 潜在的可能造成同级其它配置项被删除的bug[特指同级中高被赋nil的项]。 ([3f541c7](https://github.com/LuSrackhall/KeyTone/commit/3f541c755e90ce81a9b8250115996010e1b308c9))
* **sdk:** 修复了因至臻键音中 type为 random 或 loop 时的值为空时 被选择为 对应的按键联动声效 所引发 sdk 发生的 panic 崩溃问题。 ([f36cb52](https://github.com/LuSrackhall/KeyTone/commit/f36cb52c95777d662333455fc91ff4dfe6131bda)), closes [#44](https://github.com/LuSrackhall/KeyTone/issues/44) [#45](https://github.com/LuSrackhall/KeyTone/issues/45)
* **sdk:** 修复总被防火墙告知KeyTone正在监听某某端口的警告的bug。 ([dbe52ec](https://github.com/LuSrackhall/KeyTone/commit/dbe52ec9723d388f408bf62ed0d5fbc2e6712f57)), closes [#34](https://github.com/LuSrackhall/KeyTone/issues/34)
* **sdk:** 修复beep的错误用法, 遵循全生命周期仅调用一次speaker.Init()的预期规则。以防止潜在的内存泄漏风险, 并尽量将碎片内存的缓冲区抑制在可接受的范围内。 ([95a9671](https://github.com/LuSrackhall/KeyTone/commit/95a96719ac3f02eb6e3d03a65c81b2c9548f6d93))
* **ui | sdk:** 重新审视 按键联动声效 步骤中 是否启用内嵌测试音 的选项。不应该只提供单一的 启用/禁用 内嵌测试音选项,而是应该将其分为 按下/抬起 内嵌测试音的 启用/禁用 选项。 ([6f77c72](https://github.com/LuSrackhall/KeyTone/commit/6f77c72c5c7af3a00149b6f23664f5dab25449d9)), closes [#43](https://github.com/LuSrackhall/KeyTone/issues/43)
* **ui:** 30ms的防抖延时,不止为何在键音包界面的键音包名称输入框输入时不起作用[主设置则完全可用]。此单独对名称输入框增加防抖延时。注意[由于vue中onBeforeMount等生命周期钩子回调中定义函数或变量常量时,如果此钩子是单独触发的,则无伤大雅[但最好还是定义在外面]。] ([21fc36d](https://github.com/LuSrackhall/KeyTone/commit/21fc36da680ed054e72bad5e394cd34f1a47ec95))
* **ui:** 编辑已有声音的代码中, 相关model的类型定义名称不合适, 甚至有遗漏的字段。-- 源头是声音制作的相关代码在报错配置时, 就遗漏了相关valume字段。 ([ab4fdc0](https://github.com/LuSrackhall/KeyTone/commit/ab4fdc0e2fdc27d42e16ff07add802cf61a5ef2b)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31) [#33](https://github.com/LuSrackhall/KeyTone/issues/33)
* **ui:** 彻底修复 制作新的按键音对话框 的按下抬起配置  中  最近几次提交提到的这类bug->即所谓的对象引用变更引发的bug ([e6aa9d7](https://github.com/LuSrackhall/KeyTone/commit/e6aa9d7f8b38a266b39291a55115136b44b738af)), closes [#36](https://github.com/LuSrackhall/KeyTone/issues/36) [#41](https://github.com/LuSrackhall/KeyTone/issues/41)
* **ui:** 从根源上修复了在某些特殊情景后(如 win+j win+k 或 左键点击选择框后的esc键), 即使实际上只按下了一个按键, 但也仍会不断弹出'同时按下多个按键的提示消息'的bug。 ([09db473](https://github.com/LuSrackhall/KeyTone/commit/09db4732062a277b2448b0202798c6c996e2fff0)), closes [#44](https://github.com/LuSrackhall/KeyTone/issues/44) [#47](https://github.com/LuSrackhall/KeyTone/issues/47)
* **ui:** 对 制作新的按键音 对话框中, 按下和抬起 配置中的 选择输入框, 尽可能适配响应式的label。 但目前仍有些许问题->更改声音名称后, 1.对象引用变更;2.没有触发响应式,这个主要在声音name有空字符串转非空字符串时 ([7745e0a](https://github.com/LuSrackhall/KeyTone/commit/7745e0afa159901808d1246e544923b95f12b122)), closes [#36](https://github.com/LuSrackhall/KeyTone/issues/36) [#41](https://github.com/LuSrackhall/KeyTone/issues/41)
* **ui:** 对应StoreSet的逻辑, 虽然使用中极少关心是否传输成功, 但可以不用却不能没有。 因此,返回成功或失败的boolean标志。 ([9232941](https://github.com/LuSrackhall/KeyTone/commit/9232941fe1cb3fba10426116cb32d37fddd89955))
* **ui:** 对于 单键录制的逻辑 改变之前错误的代码结构。 修改为正确的将结构, 即相关逻辑中的api应该分开为独立的 记录相关 与 录制相关 的两组。 ([5ddb7e6](https://github.com/LuSrackhall/KeyTone/commit/5ddb7e61c0816db124789ce1172bf70b095455ce)), closes [#47](https://github.com/LuSrackhall/KeyTone/issues/47) [#50](https://github.com/LuSrackhall/KeyTone/issues/50)
* **ui:** 放弃使用对象整体作为 选择输入框组件的 uuid, 而是改用自己指定uuid。 彻底避免最近几次提交中所解决的这类bug->即所谓的对象引用变更引发的bug, 之后想改类型结构就改, js/ts嘛, 图个方便而已。 ([4d5a336](https://github.com/LuSrackhall/KeyTone/commit/4d5a336f18d325a0d7a590dbcfb4cc048405cf47)), closes [#33](https://github.com/LuSrackhall/KeyTone/issues/33) [#41](https://github.com/LuSrackhall/KeyTone/issues/41)
* **ui:** 给单键声效定义中, 将 已设置单键声效按键 的识别逻辑修改为 只有 down/up 至少一个被正确设置且value不为空字符串时, 才算作 已设置单键声效的按键。 ([8d4c9b6](https://github.com/LuSrackhall/KeyTone/commit/8d4c9b6721e9ad5f51085ba33b4b07405e3fc433)), closes [#44](https://github.com/LuSrackhall/KeyTone/issues/44)
* **ui:** 给单键声效定义中, 修复了 所点击单键 的声效的原数据初始化时, 遇到 有 但为空字符串值时 的转换处理方式。   由返回本身 改为 返回null。 ([b64511d](https://github.com/LuSrackhall/KeyTone/commit/b64511d54209824fcff97b4f4de444577932756a)), closes [#44](https://github.com/LuSrackhall/KeyTone/issues/44)
* **ui:** 将对应逻辑放到其该放的为准, 使得代码逻辑更清晰。1.启动时 加载持久化键音包的逻辑 移动值 App.vue中或boot中; 2.选择键音包后分别在watch和sse回调中加载调用, 双重保险。 ([df8cd2f](https://github.com/LuSrackhall/KeyTone/commit/df8cd2fece23ea2274b79721b9eb7dfb1945aeba)), closes [#54](https://github.com/LuSrackhall/KeyTone/issues/54)
* **ui:** 解决了 制作新的按键音 对话框内, 多选声音时, 不会在每次选择后自动关闭列表 的bug。 ([9fbddac](https://github.com/LuSrackhall/KeyTone/commit/9fbddac3fd4da0ef792bba5959a84ac4fbf06665)), closes [#36](https://github.com/LuSrackhall/KeyTone/issues/36)
* **ui:** 适配新架构下, 编辑已有声音 对话框 中的业务逻辑, 目前适配进度-> 1.利用绑定uuid解决潜在的引用变更问题;2.解决了展示label时, 因类型不匹配造成的报错问题; ([8662828](https://github.com/LuSrackhall/KeyTone/commit/8662828fdae65c017aec15f7ca24d1763c77d947)), closes [#37](https://github.com/LuSrackhall/KeyTone/issues/37)
* **ui:** 适配新架构下,编辑已有声音 对话框 中的业务逻辑, 目前适配进度->已完成->解决了最终遗留的列表中label显示异常问题[原有是列表项元素类型 与 已选项元素类型 不一致]。 ([f1e937d](https://github.com/LuSrackhall/KeyTone/commit/f1e937de63ffcc8d15db1894658f5d033785e3dd)), closes [#37](https://github.com/LuSrackhall/KeyTone/issues/37) [#42](https://github.com/LuSrackhall/KeyTone/issues/42) [#39](https://github.com/LuSrackhall/KeyTone/issues/39)
* **ui:** 完善了新声音添加到配置文件时的校验逻辑和通知栏信息。 ([7830f49](https://github.com/LuSrackhall/KeyTone/commit/7830f4996cd1d94f92691043d5a040e416458480)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31)
* **ui:** 修复 编辑已有按键音 前, 对是否有按键音的判断异常问题。 ([34f6bca](https://github.com/LuSrackhall/KeyTone/commit/34f6bca3c2cf6cad891d2d880f13f87641c638b8)), closes [#37](https://github.com/LuSrackhall/KeyTone/issues/37)
* **ui:** 修复 新架构下, 制作新的按键音 时, 若选择类型为源文件时无法反向定位的bug。 原因是存储时漏掉了关键字段, 本次提交中已修复。 ([b4e407d](https://github.com/LuSrackhall/KeyTone/commit/b4e407df354c28ca215d820b36942dba5fe4a34b)), closes [#39](https://github.com/LuSrackhall/KeyTone/issues/39) [#36](https://github.com/LuSrackhall/KeyTone/issues/36)
* **ui:** 修复 在键音包创建界面退出 后, 无法重新加载 用户选择的 持久化至 设置文件的  键音包的问题。 现在在退出此界面的生命周期钩子中, 会执行相应逻辑。 ([19b601b](https://github.com/LuSrackhall/KeyTone/commit/19b601b9818eea4ce5b100c6bdf7cf409470c0c2)), closes [#54](https://github.com/LuSrackhall/KeyTone/issues/54)
* **ui:** 修复 制作新的按键音 和 编辑已有按键音 对话框内, 保存配置时的 通知条幅t提示不准确的bug。 修复完成后, 如果是制作, 则为创建成功, 如果是编辑则为修改成功。 ([e18c569](https://github.com/LuSrackhall/KeyTone/commit/e18c5695c772800970c29a40858b08c953441a3f)), closes [#36](https://github.com/LuSrackhall/KeyTone/issues/36) [#37](https://github.com/LuSrackhall/KeyTone/issues/37)
* **ui:** 修复, 编辑已有声音 对话框的一个bug->对象引用变更。以及 制作新的按键音 对话框的按下抬起配置对话框中的两个bug->1.对象引用变更;2.没有触发响应式[主要指在声音name有空字符串转非空字符串时] ([203b4a1](https://github.com/LuSrackhall/KeyTone/commit/203b4a18827b34c7167934fe234966c834caafad)), closes [#33](https://github.com/LuSrackhall/KeyTone/issues/33) [#36](https://github.com/LuSrackhall/KeyTone/issues/36) [#41](https://github.com/LuSrackhall/KeyTone/issues/41)
* **ui:** 修复了  单键录制时, 有一些具有特定功能的按键, 会影响录制行为 的bug。 ([2f9306e](https://github.com/LuSrackhall/KeyTone/commit/2f9306ea9f2fb2a5948235928211ec214a19ce4a)), closes [#50](https://github.com/LuSrackhall/KeyTone/issues/50)
* **ui:** 修复了 按键联动音效 步骤中, 内嵌测试音是否启动组件与sse的循环依赖bug。 同时, 顺手修改了部分提示框内容, 以及ui中按钮的位置。 ([c2b97c8](https://github.com/LuSrackhall/KeyTone/commit/c2b97c8f26d04cad73e44f82e1d45b4c834df248)), closes [#43](https://github.com/LuSrackhall/KeyTone/issues/43)
* **ui:** 修复了 单键录制时, 数字键盘中'enter'按键会造成之后再次按下任何按键都 报多个按键同时按下的消息的bug。 ([ceb5dd4](https://github.com/LuSrackhall/KeyTone/commit/ceb5dd4e2fe01df4d565949b8f0f9eb8dae6fcd1)), closes [#50](https://github.com/LuSrackhall/KeyTone/issues/50)
* **ui:** 修复了 全键声效设置 对话框中, 锚定至臻键音时非至臻键音也被意外锚定的bug。 本次修复后, 选择或完全删除时仅会对至臻键音其锚定作用, 非至臻键音将完全不受锚定效果影响。 ([c3dae27](https://github.com/LuSrackhall/KeyTone/commit/c3dae27808ed4533dd22648eb19cd17f1b50330d)), closes [#45](https://github.com/LuSrackhall/KeyTone/issues/45)
* **ui:** 修复了 全局步进器 中 自由 展开/关闭 每个步骤的事件, 在步骤展开后的内容区域也会触发的bug。 并且本次提交, 为unocss引入了用于兼容tailwindcss中 [@apply](https://github.com/apply) 指令语法 的指令转换器。 ([dbdc024](https://github.com/LuSrackhall/KeyTone/commit/dbdc024973efcecc392b9fe7470725399052f908)), closes [#21](https://github.com/LuSrackhall/KeyTone/issues/21)
* **ui:** 修复了成功删除文件后, 主动给selectedSoundFile.value赋空值而触发SoundFileRename调用而引起的skd判断错误参数并返回前端的问题。 ([35b8dfd](https://github.com/LuSrackhall/KeyTone/commit/35b8dfdf752733344bc4787f42f1288005e3a521)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **ui:** 修复了创建键音包界面输入框的错误提示不会响应式中英文切换的bug。这个bug在此应用中可能并不是真实用户需求能够接触到的, 但这的确算是一个bug。 ([780877f](https://github.com/LuSrackhall/KeyTone/commit/780877f37621a8f26ed1ac8626015198e2890d3e)), closes [#28](https://github.com/LuSrackhall/KeyTone/issues/28)
* **ui:** 修复了给单键的声效定义中, 按下 及 抬起 声效设置相关的组件中, 锚定的逻辑无效问题。   问题原因是在解决全局声效选择的代码时, 遗漏了此部分的适配性改动。 ([5238682](https://github.com/LuSrackhall/KeyTone/commit/523868272cad3264106de22f0126927384e7e8b5)), closes [#44](https://github.com/LuSrackhall/KeyTone/issues/44)
* **ui:** 修复了删除最后一个声音源文件时, ui界面仍然存在此文件选项的bug。--这是一个纯前端的问题, 因为sdk已经删除了音频源文件并正确更新了配置文件。 ([26166f7](https://github.com/LuSrackhall/KeyTone/commit/26166f757ac1b165a4fb26b00fd36b326894cc53)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **ui:** 修复了一些代码中之前未发现的逻辑问题 ([18a23d5](https://github.com/LuSrackhall/KeyTone/commit/18a23d5b0fba36e7c73780b80cbff4a096e5e943)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **ui:** 修复了因删除某个源文件后, 造成配置name缺失而引起的sse回调后的逻辑无法正确执行, 进而影响驱动ui的数据结构无法正确变更的bug。 ([9e8f2d6](https://github.com/LuSrackhall/KeyTone/commit/9e8f2d62000c37ccae70a12ff6a55e9393f60175)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **ui:** 修复了重新启动 或是 从创建/编辑键音包界面 返回后, 用户所选择的持久化的 音频包, 未被正确加载到sdk中的问题。以及潜在的重复加载的问题。 ([fdf4c48](https://github.com/LuSrackhall/KeyTone/commit/fdf4c48e4b237bd406c2ca12ef83b6d70da3013d)), closes [#54](https://github.com/LuSrackhall/KeyTone/issues/54)
* **ui:** 修复上个提交的第二个问题,其并不是因为多余的name字段赋值失败造成的。真实原因是声音对象的name_id与声音文件对象的nameID无法对应造成赋值报错而引发的undefine。 ([caf9df0](https://github.com/LuSrackhall/KeyTone/commit/caf9df0b6bd821562ecef6a0ee4fd790b9c85583)), closes [#33](https://github.com/LuSrackhall/KeyTone/issues/33)
* **ui:** 修复上上个提交的第一个问题。 解决方案-对name名称的获取,由原本的直接获取,转为通过数组的find方法和 sha256+name_id字段 来查找自身,以重新获取name的方式。 ([41a5631](https://github.com/LuSrackhall/KeyTone/commit/41a5631c022091807087a84bbab89ac1191fb9de)), closes [#33](https://github.com/LuSrackhall/KeyTone/issues/33)
* **ui:** 修复在编辑音频源文件名称时, 因空名称而造成的对应声音源文件选择后无法识别编辑的bug。 ([0046e7d](https://github.com/LuSrackhall/KeyTone/commit/0046e7d6976c94d9b48b5d56bfd22cba05109af7)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **ui:** 修复制作新的声音时, 点击确认添加并添加成功后, 未能清理volume输入框的bug。 ([7684c1c](https://github.com/LuSrackhall/KeyTone/commit/7684c1c4300c54d969f45a29f9ea44ee849d6069)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31)
* **ui:** 修复主界面键音包选择器,选择声效后 无法在重新启动 或是 切换界面并返回主界面后 持久化的问题,或者说本次提交避免了每次启动应用或是切回主界面都要重新选择键音包的问题。因为本次提交中持久化了用户的最终选择。 ([cb03897](https://github.com/LuSrackhall/KeyTone/commit/cb03897198bcc9aed04109c69a02bd4359e5c3d7)), closes [#54](https://github.com/LuSrackhall/KeyTone/issues/54)
* **ui:** 之前的重构, 因js/ts传递参数时的引用传递问题,造成了不小心引入了不该存入配置文件的字段, 本次commit修复了此bug。 ([49c9055](https://github.com/LuSrackhall/KeyTone/commit/49c90551466df3f2ca65ceff6143a45f083cf879)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31) [#33](https://github.com/LuSrackhall/KeyTone/issues/33)
* **ui:** vue3中,当一个组件的setup函数返回一个Promise时，它需要被包裹在一个<Suspense>组件中才能正确渲染。简单说不能在setup函数中使用await,否则无法正确渲染页面。使用生命周期钩子,在其内部使用await可解决。本次提交还额外明确了些本页面的生命周期。 ([0e0dd47](https://github.com/LuSrackhall/KeyTone/commit/0e0dd4700c9dbdfef652cf8d3e44535a1becc7fa))


### Features

* 实现了 可控制的 局部启动的 单键Dik码 与 前端事件名称的 映射。(会在多个按键时自动拒绝映射, 并给出警告; 会根据空字符的判断逻辑, 来防止 错误映射的产生(原因见注释)。) ([f67dcb6](https://github.com/LuSrackhall/KeyTone/commit/f67dcb6cabbc2820e435570486086e75653a7ee5)), closes [#47](https://github.com/LuSrackhall/KeyTone/issues/47)
* **sdk | api:** 对kaytonePkg的api, 在通用的get和set基础上, 新增一个delete。 使得后续实现删除配置功能时, 可以更加便捷。 ([d331046](https://github.com/LuSrackhall/KeyTone/commit/d33104623a83ad952ab83322fcf3788e870fe2b6))
* **sdk | ui:** 通过在键音包模块中引入其专用的sse, 使得可以在ui界面获取以加载的音频文件列表, 从而使得进一步的编辑成为可能[主要指重命名和删除]。 ([22fd6bc](https://github.com/LuSrackhall/KeyTone/commit/22fd6bc9d073230a2bd0007595408ad4baf39c84)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **sdk:** 升级了键音播放器, 以兼容性优先的策略, 适配了音频裁切和分级音频指定。 ([9e98802](https://github.com/LuSrackhall/KeyTone/commit/9e988029a26f53c82eeb10ec10c07c349c065f61))
* **sdk:** 完善了 按联动音效 步骤中, 内嵌测试音是否启动的控制。 -> 之前仅实现了前端ui代码, 本次commit在sdk中为此功能做了真正的适配。 ([a998cae](https://github.com/LuSrackhall/KeyTone/commit/a998caedd456eefd0fcc5ae1c2aa9d606d829dad)), closes [#43](https://github.com/LuSrackhall/KeyTone/issues/43)
* **sdk:** 在后端进行最终的文件保存之前, 对文件进行sha256哈希, 并使用哈希值作为最终文件名称以保存文件。 ([58265d0](https://github.com/LuSrackhall/KeyTone/commit/58265d0451003fd5776e364b923b7c6f9d1bda2d)), closes [#23](https://github.com/LuSrackhall/KeyTone/issues/23)
* **sdk:** 支持更多格式的音频文件, 如mp3、ogg等 ([6ffce2b](https://github.com/LuSrackhall/KeyTone/commit/6ffce2b65f2decc40a9f6a7680b0211d4128c36a))
* **sse-keyevent:** 将sdk中的keyevent映射至前端, 利用map存储keyCode与其对映的State值。 ([5b68e29](https://github.com/LuSrackhall/KeyTone/commit/5b68e29df76cc2c23e16cb2025bd23920f07d8b9))
* **ui | sdk:** 引入对键音包配置文件进行读写的restfulAPI桥梁。-> 此桥梁理论上通用与新建和编辑, 桥梁的初始化在"新建"或"编辑"之初由后端进行。 ([9a482f6](https://github.com/LuSrackhall/KeyTone/commit/9a482f6697636599578ccb3c2d2c0335620ca9ba)), closes [#27](https://github.com/LuSrackhall/KeyTone/issues/27)
* **ui | sdk:** 在制作新声音的界面提供预览声音的按钮, 从后端提供预览播放声音的接口, 供前端实现预览声音的功能。 ([125ec42](https://github.com/LuSrackhall/KeyTone/commit/125ec4263b93fcc965f49fd8e9f1ddb108f7720a)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31)
* **ui | sdk:** 在sdk中, 对声音播放器的第二个参数Cut的结构体增加volume字段, 使得可以初始化所传入声音的音量。从而使得ui中, 预览声音时, 可以使用音量参数。 ([b7cb036](https://github.com/LuSrackhall/KeyTone/commit/b7cb0364e062883563d588662047e6d44b138506)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31) [#32](https://github.com/LuSrackhall/KeyTone/issues/32)
* **ui:** 编辑已有声音界面初始化, 新增了选择输入框, 可供选择已有的声音进行进一步的编辑。 写代码时注意使用此输入框时, 需要对原model的类型做准确的定义, 慎用鸭子类型。 ([18b429c](https://github.com/LuSrackhall/KeyTone/commit/18b429cd4192283aad062651b393796a451ffd8e)), closes [#33](https://github.com/LuSrackhall/KeyTone/issues/33)
* **ui:** 编辑已有声音源文件/管理已载入的源文件 | 实现重命名功能 ([0c3ab18](https://github.com/LuSrackhall/KeyTone/commit/0c3ab184740b18d3638188fe9908c245b4f71268)), closes [#29](https://github.com/LuSrackhall/KeyTone/issues/29)
* **ui:** 初始化了制作新的声音相关的ui界面。 ([bc76369](https://github.com/LuSrackhall/KeyTone/commit/bc76369bd50308ab4291fa15e9cf78f7ce665841)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31)
* **ui:** 初始化声音编辑时选择声音后的编辑卡片。注意,卡片中源文件的选择框目前存在1.无法跟随name的变化实时更新;2.选择后仅会赋值源类型字段, 多余的name字段仍是undefine。 ([e1ea1b0](https://github.com/LuSrackhall/KeyTone/commit/e1ea1b0caf790f9983c23b6705696c8d9c75a774)), closes [#33](https://github.com/LuSrackhall/KeyTone/issues/33)
* **ui:** 键音包架构变更, 使得可以在键音中选择 源文件、声音、按键音 三种类型。 本次提交主要针对 制作新的按键音 对话框。  因破坏性变更 编辑已有按键音 在适配之前将无法正常工作。 ([9e5ffc6](https://github.com/LuSrackhall/KeyTone/commit/9e5ffc65224cbb8efb41758a82e63120d6a92764)), closes [#39](https://github.com/LuSrackhall/KeyTone/issues/39) [#36](https://github.com/LuSrackhall/KeyTone/issues/36)
* **ui:** 键音包需要一个名字, 且是必填的。名字需要国际化的默认值, 且会自动写入对应配置文件中去。 ([e59f87d](https://github.com/LuSrackhall/KeyTone/commit/e59f87da4d2bc1cf9008807de643f18c417f3297)), closes [#28](https://github.com/LuSrackhall/KeyTone/issues/28)
* **ui:** 声音编辑时, 在选择某个声音后的编辑卡片上, 增加预览声音、确认修改、删除声音等三个按钮。 ([cef5d53](https://github.com/LuSrackhall/KeyTone/commit/cef5d5304dcd9170dddb2ae8f1d27222b86f78ea)), closes [#33](https://github.com/LuSrackhall/KeyTone/issues/33)
* **ui:** 实现了新声音添加到配置文件的功能, 并增加了一定量的校验逻辑, 以及相关的通知栏。 ([d5e835f](https://github.com/LuSrackhall/KeyTone/commit/d5e835f4faf79c81d4cf991c9fda6e247dc9a8ff)), closes [#31](https://github.com/LuSrackhall/KeyTone/issues/31)
* **ui:** 适配新架构下,编辑已有声音 对话框 中的业务逻辑, 目前适配进度->1.适配新架构下按下抬起配置的声音选择列表;2.遇到问题,列表中label显示异常[推测是关键字段类型路径不一致引起] ([ad1a3da](https://github.com/LuSrackhall/KeyTone/commit/ad1a3da6b275ece24ba7bc70cc073a1f14ff83ff)), closes [#37](https://github.com/LuSrackhall/KeyTone/issues/37) [#39](https://github.com/LuSrackhall/KeyTone/issues/39)
* **ui:** 完善了 编辑已有按键音 对话框内 的基本可用的业务逻辑, 不过目前对话框内暂未支持预览按钮 以及 必要的删除按钮  , 后续会陆续支持。 ([8b8c81a](https://github.com/LuSrackhall/KeyTone/commit/8b8c81a10292961e28834bc114d9cc61906dc645)), closes [#37](https://github.com/LuSrackhall/KeyTone/issues/37)
* **ui:** 完善了 制作新的按键音 对话框内 的基本可用的业务逻辑, 不过目前对话框内暂未支持预览按钮, 后续会陆续支持。  还有个遗留问题是, 多选声音时, 不会在每次选择后自动关闭列表。 ([388fe05](https://github.com/LuSrackhall/KeyTone/commit/388fe05934f747637fa81e6556891ca1aef4f0ac)), closes [#36](https://github.com/LuSrackhall/KeyTone/issues/36)
* **ui:** 在主界面提供键音包选择器, 以方便用户快速选择不同的键音包来使用。 ([29f9eb2](https://github.com/LuSrackhall/KeyTone/commit/29f9eb28830dee486b9d7a18eb1257d6b3b71632)), closes [#54](https://github.com/LuSrackhall/KeyTone/issues/54)
* **ui:** 增加快速增设单键声效按钮及其对话框页面。(具体功能逐步完善->本次仅在页面内初始化了用于 搜索/录制 单键 的 select 组件->主要适配了 录制单键的初始功能) ([78502db](https://github.com/LuSrackhall/KeyTone/commit/78502dbfcca63d4d41f3a3459e46663c74c9b547)), closes [#44](https://github.com/LuSrackhall/KeyTone/issues/44)
* **ui:** 增加了 单独Dik码 与 name 实时映射功能, 如何被其它组件调用使用的api, 可使代码逻辑更清晰。 ([51fc510](https://github.com/LuSrackhall/KeyTone/commit/51fc510a688ec464f4b467309bfd76879cd88a3f)), closes [#47](https://github.com/LuSrackhall/KeyTone/issues/47)
* **ui:** 支持 编辑已有按键音 对话框内 必要的删除按钮。 ([2a06d22](https://github.com/LuSrackhall/KeyTone/commit/2a06d2264fc5053ad4817c593520d0ce65ec3dbd)), closes [#37](https://github.com/LuSrackhall/KeyTone/issues/37)



#  (2024-08-29)


### Bug Fixes

* 改善了从配置文件更改被监听到前端收到并反馈至ui的实时性。<调整了防抖延后的毫秒数值(缩小)> ([941e175](https://github.com/LuSrackhall/KeyTone/commit/941e175c84bda2012f336b390817ea19452c635a))
* 修复了当音量提升/缩减滑块选择为小于-5或者说-g时, 主界面音量条件滑块异常的bug。<甚至仅小于0就有滑块范围缩小隐患, 原因是主页min音量的设计问题> ([1a0a88f](https://github.com/LuSrackhall/KeyTone/commit/1a0a88f96e19993256dbc1b75a0b19e0bd8c9bdb)), closes [#15](https://github.com/LuSrackhall/KeyTone/issues/15)
* **frontend | electron:** 修复了托盘后台情况下<即隐藏窗口情况>, 重复开启单例, 不会弹出窗口的bug。 ([1657cb9](https://github.com/LuSrackhall/KeyTone/commit/1657cb99317e8a87c20890988c6fbff0ca8cabf7))
* **frontend | ui:** 修复了页面可选中的bug。<这是由于我希望这个工具应用的界面更偏向与界面而不是页面, 因此禁止页面的文本选择以及图片的拖动。> ([09e688a](https://github.com/LuSrackhall/KeyTone/commit/09e688a0be5a2116e59bc6946874851e5f49ed0e))
* **mute:** 对于主页面的静音按钮,对于我们键盘音软件来说,太容易被键盘误触给重新打开了。<造成不好的体验,误以为无法静音,因此禁用click中的键盘事件触发机制,仅保留鼠标单击事件。> ([82103c0](https://github.com/LuSrackhall/KeyTone/commit/82103c0f5b2998915f9f6a4f54bab7725b445e00))
* **sdk:** 继续修复sdk的小概率崩溃bug。<虽然文档中说viper是并发读写不安全的, 但我简单测试了下, 纯并发读的情况很少崩溃, 反而读写同时进行或纯并发写时容易奔溃,因此换读写锁。> ([ecda87c](https://github.com/LuSrackhall/KeyTone/commit/ecda87c698bffe907cb7aaa0eb5b238a3e8d1144))
* **sdk:** 解决了viper.Set()的高优先级覆盖问题, 使其使用方式更符合本应用的场景需求。<虽不知这样会不会增加损耗, 不过性能瓶颈在纯客户端应用中几乎可以忽略不计> ([f97fac5](https://github.com/LuSrackhall/KeyTone/commit/f97fac592066718973ccc72996422bb7f633441f)), closes [#18](https://github.com/LuSrackhall/KeyTone/issues/18) [#11](https://github.com/LuSrackhall/KeyTone/issues/11)
* **sdk:** 修复了sdk的小概率崩溃bug。 <虽然机率很小, 但毕竟viper是读写并发不安全的, 因此暂且对齐set和get行为粗暴的上个互斥锁好了> ([1c884a6](https://github.com/LuSrackhall/KeyTone/commit/1c884a6cb2ded11d544c7f6bbf0ec5400a7288ed))
* **tray:** 更改托盘的'关闭<close>'选项名称为'退出<quit>' ([adc9439](https://github.com/LuSrackhall/KeyTone/commit/adc9439c4de7e7c73ebfa0569578de37ca2cfefa))
* **ui交互方式变更:** 对于设置界面的各组配置, 由默认展开, 改为默认收起。 ([42c84e4](https://github.com/LuSrackhall/KeyTone/commit/42c84e4b0c492cc6c027134afd9bc29b576782c6)), closes [#13](https://github.com/LuSrackhall/KeyTone/issues/13)
* **ui:** 修复了界面溢出问题, 现在可能超出的内容和可能产生的滚动条, 不再会溢出界面。 ([d060ad6](https://github.com/LuSrackhall/KeyTone/commit/d060ad67e6415a8adc09d93b1ee3e75a9c2af33a)), closes [#12](https://github.com/LuSrackhall/KeyTone/issues/12)


### Features

* 对导航栏做了功能提升, 使得点击应用名称可以直接到主页面; 在侧边导航中, 新增了关闭侧边导航的按钮; 优化了侧边导航的关闭逻辑, 使得可在点击选项但未发生路由时也可关闭。 ([2215d01](https://github.com/LuSrackhall/KeyTone/commit/2215d018d6c7e96f03170d1b5fcb182b1b3161dc))
* 在主页面的增加音量调整滑块, 以及静音图标。即新增了主页面的音量调整功能, 以及静音功能。 ([537f349](https://github.com/LuSrackhall/KeyTone/commit/537f34949093490586b86578fee3fe36fbbcaab0)), closes [#15](https://github.com/LuSrackhall/KeyTone/issues/15)
* **主页面配置项:** 新增了主页面设置1.音量降低幅度的输入框;2.是否打开主页面音量调试滑块的开关;3.默认隐藏的主页面音量调试滑块;4.相关项的i18n。 ([1c29d04](https://github.com/LuSrackhall/KeyTone/commit/1c29d046129896b878c8109c639714a30c4c427a)), closes [#17](https://github.com/LuSrackhall/KeyTone/issues/17)
* **fix:** 完善了国际化多语言设置的功能。<此提交为小阶段总结提交--懒得用pr处理这部分了: 之前的一小段提交, 我们初始化了设置页面, 并在ui中加入了简单的页面导航功能。> ([5b3eaa7](https://github.com/LuSrackhall/KeyTone/commit/5b3eaa7d506841ae69e099ac2d1d67c54e2474b8))
* **sse:** 引入sse, 并通过sse和viper的文件监听回调,为前端ui同步配置文件的实时配置变更。<比如在electron侧完成的静音设置,或是直接手动修改配置文件完成的静音设置。> ([62bc498](https://github.com/LuSrackhall/KeyTone/commit/62bc4989264cacd4c5175c21539e4340edb1592f)), closes [#19](https://github.com/LuSrackhall/KeyTone/issues/19) [#11](https://github.com/LuSrackhall/KeyTone/issues/11) [#18](https://github.com/LuSrackhall/KeyTone/issues/18)
* **tray:** 在系统托盘<tray>引入静音/取消静音选项。 ([a55a20e](https://github.com/LuSrackhall/KeyTone/commit/a55a20eb7c49205e7b14c684d2729ab808131565)), closes [#18](https://github.com/LuSrackhall/KeyTone/issues/18)
* **ui交互方式变更:** 对于设置页面的各组配置的展开或收起的状态, 新增运行期间管理保留的被动功能。 ([26303d6](https://github.com/LuSrackhall/KeyTone/commit/26303d63a453dbf81ac647f15e38a36471a69d08)), closes [#13](https://github.com/LuSrackhall/KeyTone/issues/13)
* **ui:** 实现了设置界面提升/缩减原始音频包音量的功能,以及功能相关的国际化内容配置; 引入了lodash包防抖; 此外,完美解决了此功能可能存在的一些已知bug。 ([73d1605](https://github.com/LuSrackhall/KeyTone/commit/73d1605de48a57657f3e0690ce44354f80bf981f)), closes [#6](https://github.com/LuSrackhall/KeyTone/issues/6) [#14](https://github.com/LuSrackhall/KeyTone/issues/14)
* **ui:** 在ui设置页面, 新增"启动与自动启动"的系列设置 ([c05c646](https://github.com/LuSrackhall/KeyTone/commit/c05c6465553fe07e480a77771ec0b11c14744eaa)), closes [#4](https://github.com/LuSrackhall/KeyTone/issues/4) [#5](https://github.com/LuSrackhall/KeyTone/issues/5)



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



