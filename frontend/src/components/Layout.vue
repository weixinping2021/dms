<template>
    <n-layout style="height: 100vh;">
      <!-- 顶部导航 -->
      <n-layout-header style="background: #fff; box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1); padding: 0 16px;">
        <Header />
      </n-layout-header>
  
      <n-layout has-sider>
        <!-- 左侧导航栏 -->
        <n-layout-sider bordered style="background: #f5f5f5; width: 200px;">
          <Sidebar />
        </n-layout-sider>
  
        <!-- 内容区域 -->
        <n-layout-content style="padding: 16px; background: #f9f9f9;">
          <router-view /> <!-- 路由视图 -->
        </n-layout-content>
      </n-layout>
    </n-layout>
  </template>
  
  <script>
  import Header from './Header.vue';
  import Sidebar from './Sidebar.vue';
  
  export default {
    components: {
      Header,
      Sidebar,
    },
  };
  </script>
  