import { defineStore } from 'pinia';
import { ConfigGet, GetAudioPackageList, GetAudioPackageName, LoadConfig } from 'src/boot/query/keytonePkg-query';
import { ref, watch } from 'vue';
import { useSettingStore } from './setting-store';
import { useQuasar } from 'quasar';
import { StoreGet, StoreSet } from 'src/boot/query/store-query';
import { useKeytoneAlbumStore } from 'src/stores/keytoneAlbum-store';

export const useMainStore = defineStore('main', () => {
  const keytoneAlbum_store = useKeytoneAlbumStore();
  // 添加音量百分比状态
  const volumePercentage = ref(1); // 默认100%

  // 添加初始化函数
  async function initVolumePercentage() {
    const savedPercentage = await StoreGet('volume_percentage');
    if (savedPercentage !== undefined && savedPercentage !== null) {
      volumePercentage.value = savedPercentage;
    }
  }

  // 监听音量百分比变化并保存
  watch(volumePercentage, (newValue) => {
    // setTimeout(() => { // 因为本身就是异步的, 有一定延后性, 所以不需要setTimeout。我们延后的目的是:
    //////////////////////// * 对某个同时涉及多个相关变量的写入持久化的操作, 所可能造成的sse回馈结果的混乱的问题。(即短时间内几乎同时的多次触发, 可能破坏其本身的延后性)
    //////////////////////// * 而这个问题在做一些某些仅涉及单一数据变更的操作时会被异步本身的滞后性所天然解决。
    //////////////////////// 只有当无法天然解决时, 才需要手动延后。
    //////////////////////// * 某些逻辑可能会依靠这些天然性质, 所有在此处做手动延后会违背异步本身的天然延后性质, 从而影响到依靠这些性质的逻辑。(当然这个字段目前不涉及, 但我只是举例子, 本项目中还是有些字段符合此介绍的, 毕竟太长的延后会影响相关部分的性能。)
    //////////////////////// * 因此, 最好在引发多个相关变量的写入持久化的操作时, 具体问题具体分析, 通过在相关源头处对次要变量进行手动延后来解决这类问题。

    StoreSet('volume_percentage', newValue);
    // }, 30);
  });

  // 记录每次组件卸载时, 音量百分比的变化, 防止调整音量降幅后, 回到主页面时, 音量百分比也恢复到调整前的状态。(此时我不希望恢复到调整前的状态)
  const volumeNormalReduceScope = ref(5);

  //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  /////////////////////////////////////////////////下方是键音包相关的配置////////////////////////////////////////////////
  //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
  const keyTonePkgOptions = ref([]);
  const keyTonePkgOptionsName = ref(new Map());
  const keyTonePkgOptionsCopyright = ref(new Map()); // Store copyright info for each album

  // Helper function to format copyright information for display
  const formatCopyrightInfo = (copyrightData: any): string => {
    if (!copyrightData || typeof copyrightData !== 'object') {
      return '';
    }

    const authors: Array<{name: string, lastExportTime: number}> = [];
    
    // Process each copyright entry
    Object.values(copyrightData).forEach((entry: any) => {
      if (entry && entry.Author && entry.ExportTime && Array.isArray(entry.ExportTime)) {
        const lastExportTime = Math.max(...entry.ExportTime);
        authors.push({
          name: entry.Author,
          lastExportTime: lastExportTime
        });
      }
    });

    if (authors.length === 0) {
      return '';
    }

    // Sort by last export time (newest first)
    authors.sort((a, b) => b.lastExportTime - a.lastExportTime);

    // Get first author (most recent) and optionally show additional authors count
    const firstAuthor = authors[0].name;
    if (authors.length === 1) {
      return ` - ${firstAuthor}`;
    } else {
      return ` - ${firstAuthor} +${authors.length - 1}`;
    }
  };

  // Get copyright information for all albums
  const loadAlbumCopyrightInfo = async () => {
    keyTonePkgOptionsCopyright.value.clear();
    
    for (const albumPath of keyTonePkgOptions.value) {
      try {
        // Load the album config to get copyright data
        await LoadConfig(albumPath, false);
        const copyrightData = await ConfigGet('copyright');
        const formattedInfo = formatCopyrightInfo(copyrightData);
        keyTonePkgOptionsCopyright.value.set(albumPath, formattedInfo);
      } catch (error) {
        console.error(`Failed to load copyright info for ${albumPath}:`, error);
        keyTonePkgOptionsCopyright.value.set(albumPath, '');
      }
    }
  };

  /**
   * 获取键音包列表
   * * 可用于在应用启动时初始化键音包列表, 或是在创建/编辑键音包后, 更新键音包列表。
   */
  function GetKeyToneAlbumList() {
    // 获取键音包列表的初始化逻辑, 没必要在main_store中进行, 而是应该在其所相关的对应逻辑中进行(比如在App.vue或是boot中)。
    GetAudioPackageList().then(async (res) => {
      if (res.list) {
        keyTonePkgOptions.value = res.list;
        console.log('keyTonePkgOptions', keyTonePkgOptions.value);
        keyTonePkgOptionsName.value.clear();
        
        // Load album names
        for (const item of keyTonePkgOptions.value) {
          try {
            const nameRes = await GetAudioPackageName(item);
            keyTonePkgOptionsName.value.set(item, nameRes.name);
          } catch (error) {
            console.error(`Failed to get name for ${item}:`, error);
            keyTonePkgOptionsName.value.set(item, 'Unknown Album');
          }
        }
        
        // Load copyright information
        await loadAlbumCopyrightInfo();
      } else {
        keyTonePkgOptions.value = [];
      }
    });
  }

  // Function to generate the display label for album selector
  const getAlbumDisplayLabel = (albumPath: string): string => {
    const albumName = keyTonePkgOptionsName.value.get(albumPath) || 'Unknown Album';
    const copyrightInfo = keyTonePkgOptionsCopyright.value.get(albumPath) || '';
    return albumName + copyrightInfo;
  };

  const setting_store = useSettingStore();
  const q = useQuasar();

  /**
   * 加载用户所选的键音包
   * * 这个函数能够保证, 不会重复加载与目前正在使用键音包相同uuid的键音包。
   * * -- 比如, 从设置页面或是其它页面 返回主页面时, 不会重复加载。
   */
  function LoadSelectedKeyTonePkg() {
    ConfigGet('audio_pkg_uuid').then((res) => {
      // 若当前正在新建键音包, 则不需要进行反复确认的加载逻辑
      if (keytoneAlbum_store.isCreateNewKeytoneAlbum) {
        return;
      }
      console.log('res= ', res);
      console.log('setting_store.mainHome.selectedKeyTonePkg= ', setting_store.mainHome.selectedKeyTonePkg);
      // 这里的路径处理, 是为了兼容不同操作系统。(我们简单使用了quasar的platform.is.win来判断)
      // * TIPS: 若后续发现兼容性问题, 可考虑替换为node的path.basename:
      // *       * 利用vite的vite - plugin - node - polyfills插件引入(推荐)
      // *       * 或是第三方的path - browserify这个前端简单的路径处理库引入
      // *       * 以上两者均可(不过据目前所知, 只有win系统是使用'\\'作为路径分隔符的, 其他系统都是使用'/'作为路径分隔符)
      // *       * * Windows:          '\'  (反斜杠)
      // *       * * macOS:            '/'  (正斜杠)
      // *       * * Linux:            '/'  (正斜杠)
      // *       * * Unix:             '/'  (正斜杠)
      // *       * * Android:          '/'  (正斜杠)
      // *       * * iOS:              '/'  (正斜杠)
      // *       * * HarmonyOS:        '/'  (正斜杠)
      // *       * * HarmonyOS Next:   '/'  (正斜杠)
      // *       * 也就是说, 我们此处甚至可以无需借助quasar的platform.is.win来判断, 直接对字符串使用.replace(/\\/g, '/')后, 统一转换为正斜杠来处理即可。
      const UUID = setting_store.mainHome.selectedKeyTonePkg.split(q.platform.is.win ? '\\' : '/').pop();
      console.log('UUID= ', UUID);

      if (res !== UUID) {
        // 若当前的配置文件中的uuid 与 实际使用的键音包uuid不一致, 以配置文件中用户选择的键音包uuid为准, 重新加载对应键音包。
        // * setting_store.mainHome.selectedKeyTonePkg 由SSE保证,始终与配置文件的相关配置一致。
        // * ConfigGet('audio_pkg_uuid')读取到的uuid, 可能受 新建/编辑 键音包操作的影响, 导致与配置文件中的uuid不一致。
        // * 因此, 此处需要以配置文件中用户选择的键音包uuid为准, 重新加载对应键音包。否则无需重新加载。(比如从设置页面返回主页面时。)
        //
        LoadConfig(setting_store.mainHome.selectedKeyTonePkg, false).then((res) => {
          if (!res) {
            // 如果LoadConfig加载失败, 说明用户所选的键音包在当前环境下(可能已被外力删除), 因此我们将其置空。
            setting_store.mainHome.selectedKeyTonePkg = '';
            console.log(
              '重新加载持久化中用户所选的键音包失败, 此键音包可能已被破坏, 已清空所选键音包以供用户重新选择。'
            );
            return;
          }
          console.log('重新加载用户所选的键音包成功');
        });
      }
    });
  }

  return {
    volumePercentage,
    initVolumePercentage,
    volumeNormalReduceScope,
    keyTonePkgOptions,
    keyTonePkgOptionsName,
    keyTonePkgOptionsCopyright,
    GetKeyToneAlbumList,
    LoadSelectedKeyTonePkg,
    getAlbumDisplayLabel,
    formatCopyrightInfo,
    loadAlbumCopyrightInfo,
  };
});
