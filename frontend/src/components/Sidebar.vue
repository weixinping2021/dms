<template>
    <n-menu :options="menuOptions" v-model:value="selectedKey" />
  </template>
  
  <script>
  import { ref } from 'vue';
  import { useRouter } from 'vue-router';
  
  export default {
    setup() {
      const router = useRouter();
      const selectedKey = ref('/');
  
      const menuOptions = [
        { label: '首页', key: '/', path: '/' },
        { label: '关于我们', key: '/about', path: '/about' },
        { label: '联系我们', key: '/contact', path: '/contact' },
      ];
  
      const handleMenuClick = (key) => {
        selectedKey.value = key;
        router.push(key);
      };
  
      return {
        menuOptions,
        selectedKey,
        handleMenuClick,
      };
    },
  };
  </script>
  