import { defineStore } from 'pinia';
import { is } from 'quasar';
import { ref } from 'vue';

export const useKeytoneAlbumStore = defineStore('KeytoneAlbum', () => {
  const isCreateNewKeytoneAlbum = ref<boolean>(false);
  return { isCreateNewKeytoneAlbum };
});
