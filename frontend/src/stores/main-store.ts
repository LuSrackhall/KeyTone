import { defineStore } from 'pinia';
import { GetAudioPackageList, GetAudioPackageName } from 'src/boot/query/keytonePkg-query';
import { ref } from 'vue';
export const useMainStore = defineStore('main', () => {
  const keyTonePkgOptions = ref([]);
  const keyTonePkgOptionsName = ref(new Map());

  GetAudioPackageList().then((res) => {
    keyTonePkgOptions.value = res.list;
    console.log('keyTonePkgOptions', keyTonePkgOptions.value);
    keyTonePkgOptionsName.value.clear();
    keyTonePkgOptions.value.forEach((item: any) => {
      GetAudioPackageName(item).then((res) => {
        // console.log('res', res);
        keyTonePkgOptionsName.value.set(item, res.name);
      });
    });
  });

  return { keyTonePkgOptions, keyTonePkgOptionsName };
});
