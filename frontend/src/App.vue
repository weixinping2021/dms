<template>
  <a-layout style="min-height: 100vh">
    <!-- 左侧菜单 -->
    <a-layout-sider :trigger="null" theme="light" class="sider">

      <!-- 左侧菜单的Logo，点击切换 -->
      <div class="logo">
        <a-avatar src="logo.jpeg" shape="square" :size="64" />
      </div>
      <!-- 左侧菜单 -->
      <a-menu v-model:selectedKeys="selectedKeys" mode="inline" theme="light">
        <a-menu-item key="1">
          <span>
            <RouterLink to="/">Go to Home</RouterLink>
          </span>
        </a-menu-item>
        <a-divider style="margin: 0;" />
        <a-menu-item key="2">
          <span>
            <RouterLink to="/about">about</RouterLink>
          </span>
        </a-menu-item>
        <a-divider style="margin: 0;" />
        <a-sub-menu key="sub3" title="mysql" >
          <a-menu-item key="3">
            <span><RouterLink to="/offlinebinlog">离线binlog解析</RouterLink></span>
          </a-menu-item>
          <a-menu-item key="4">
            <span>数据库管理</span>
          </a-menu-item>
        </a-sub-menu>
        <a-divider style="margin: 0;" />
        <a-sub-menu key="sub2" title="redis" >
          <a-menu-item key="5">
            <span><RouterLink to="/redis">rdb文件分析</RouterLink></span>
          </a-menu-item>
          <a-menu-item key="6">
            <span>日志管理</span>
          </a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <!-- 内容区 -->
      <a-layout-content>
        <div class="content">
          <RouterView />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref, computed } from 'vue';

const open = ref(false);

const dataSource = ref([]);

const count = computed(() => dataSource.value.length + 1);



const selectedKeys = ref(['1']);  // 选中的菜单项



</script>

<style scoped>
.sider .logo {
  padding: 16px;
  font-size: 18px;
  color: #1890ff;
  font-weight: bold;
  text-align: center;
}

.trigger {
  font-size: 20px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
}

/* 内容区域样式 */
.content {
  padding: 20px;
  background: #fff;
}
</style>
