<template>
    <a-layout style="height: 100vh;">
      <!-- 顶部导航 -->
      <a-layout-header style="background: #fff; box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1); padding: 0 16px;">
        <Header />
      </a-layout-header>
  
      <a-layout has-sider>
        <!-- 左侧导航栏 -->
        <a-layout-sider bordered style="background: #f5f5f5; width: 200px;">
          <Sidebar />
        </a-layout-sider>
  
        <!-- 内容区域 -->
        <a-layout-content style="padding: 16px;">
          <router-view /> <!-- 路由视图 -->
        </a-layout-content>
      </a-layout>
    </a-layout>
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
  
  <style scoped>
  ::v-deep .ant-layout-content {
    background: #ffffff !important; /* 设置为白色或其他颜色 */
  }
  </style>