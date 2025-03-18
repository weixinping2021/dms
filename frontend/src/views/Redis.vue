<template>
    <a-row>
        <a-spin :spinning="spinning">
            <a-flex gap="middle">
                <a-input v-model:value=filename>
                    <template #addonAfter>
                        <FileAddOutlined @click="openChoseFileDlg" />
                    </template>
                </a-input>
                <a-button @click="Analyse">开始分析</a-button>
                <a-select ref="select" v-model:value="formattedTime" style="width: 400px" :options="rdbresult"
                    @focus="focus" @change="handleChange"></a-select>
            </a-flex>
        </a-spin>
    </a-row>
    <a-tabs v-model:activeKey="activeKey" hide-add type="editable-card" @edit="removeTab">
        <a-tab-pane 
        v-for="tab in tabs" 
        :key="tab.key" 
        :tab="tab.title">
        <!-- 动态内容 -->
        <component :is="tab.component" :data="tab.data" />
        </a-tab-pane>
    </a-tabs>
</template>

<script setup>
import { ref ,shallowRef} from 'vue';
import { FileAddOutlined } from '@ant-design/icons-vue';
import { OpenDialog } from "../../wailsjs/go/main/App";
import { GetRedisMemory, AnalyseRdb, GetRdbResultTitle } from "../../wailsjs/go/redis/Redis";
import RedisFirst from './RedisFirst.vue';
const spinning = ref(false);
const filename = ref("")
const formattedTime = ref("")
const rdbresult = ref([])
const activeKey = ref();

// 存放动态创建的 Tabs
const tabs = ref([]);


const removeTab = (targetKey) => {
  tabs.value = tabs.value.filter(tab => tab.key !== targetKey);
  // 处理删除后的 Tab
  if (tabs.value.length > 0) {
    activeKey.value = tabs.value[tabs.value.length - 1].key;
  } else {
    activeKey.value = null;
  }
};

function openChoseFileDlg() {
    OpenDialog().then(res => {
        console.log("chose file :", res)
        if (res != "") {
            filename.value = res
            console.log(res)
        }
    })
}

function Analyse() {
    spinning.value = true;
    AnalyseRdb(filename.value).then((result) => {
        spinning.value = false;
        formattedTime.value = result
    });
}


function focus() {

    GetRdbResultTitle().then((result) => {
        //console.log(result);
        rdbresult.value = result;
    });
}

function handleChange() {
    let newTab = {
      key: formattedTime.value,
      title: formattedTime.value,
      component: null,  // 默认没有组件
      data: {}  // 用来存放要传递的数据
    };
    newTab.component = shallowRef(RedisFirst);
    newTab.data = { message:  formattedTime.value};
    tabs.value.push(newTab);
    activeKey.value = formattedTime.value; // 切换到新 Tab
}
</script>