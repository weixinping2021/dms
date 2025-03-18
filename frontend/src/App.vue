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
                    <RouterLink to="/">数据库列表</RouterLink>
                </span>
            </a-menu-item>
            <a-divider style="margin: 0;" />
            <a-menu-item key="2">
                <span>
                    <RouterLink to="/about">about</RouterLink>
                </span>
            </a-menu-item>
            <a-divider style="margin: 0;" />
            <a-sub-menu key="sub3" title="mysql">
                <a-menu-item key="3">
                    <span>
                        <RouterLink to="/offlinebinlog">离线binlog解析</RouterLink>
                    </span>
                </a-menu-item>
                <a-menu-item key="4">
                    <span>数据库管理</span>
                </a-menu-item>
            </a-sub-menu>
            <a-divider style="margin: 0;" />
            <a-sub-menu key="sub2" title="redis">
                <a-menu-item key="5">
                    <span>
                        <RouterLink to="/redis">rdb文件分析</RouterLink>
                    </span>
                </a-menu-item>
            </a-sub-menu>
        </a-menu>
    </a-layout-sider>

    <a-layout>
        <!-- 内容区 -->
        <a-layout-content>
            <a-modal v-model:open="open" title="设置工作目录" @ok="handleOk">
                <a-input v-model:value="workdir" placeholder="输入工作目录" @focus="openChoseDirDlg" />
            </a-modal>
            <div class="content">
                <RouterView />
            </div>
        </a-layout-content>
    </a-layout>
</a-layout>
</template>

<script setup>
import {ref} from 'vue';
import {message} from 'ant-design-vue';
import {GetWorkDir,SetWorkDir,OpenDir} from "../wailsjs/go/main/App";
const selectedKeys = ref(['1']); // 选中的菜单项
const workdir = ref("")
const open = ref(false);

GetWorkDir().then(result => {
   // console.log(result)
    if (result === "") {
        open.value = true;
    }
}).catch(err => {
    message.error(err)
})

function openChoseDirDlg() {
    if (workdir.value === "") {
        OpenDir().then(res => {
            if (res != "") {
                workdir.value = res
            }
        })
    }
}

function handleOk() {
    SetWorkDir(workdir.value).then(result => {
        console.log(result)
        if (result === "success") {
            open.value = false;
        } else {
            message.error(result, 10);
        }
    })
}
</script>

<style scoped>
::v-deep .ant-layout-content {
  background: #ffffff; /* 设置为白色或其他颜色 */
}
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

.sider {
    border-right: 1px solid #ffffff;
}
</style>
