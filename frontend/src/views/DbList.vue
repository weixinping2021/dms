<template>
  <a-table bordered :data-source="data" :columns="columns" size="small">
    <template #bodyCell="{ column, record }">
      <template v-if="column.dataIndex === 'name'">
        <a @click="onClick(record)">{{ record.name }}</a>
      </template>
      <template v-if="column.dataIndex === 'password'">
        <a-input-password :value="record.password" disabled :bordered="false"/>
      </template>
    </template>
  </a-table>
</template>

<script setup>
import { GetFullCons } from "../../wailsjs/go/main/App";
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter()
const data = ref([])

const columns = [
  {
    title: '连接名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '主机',
    dataIndex: 'host',
    key: 'host',
  },
  {
    title: '端口',
    dataIndex: 'port',
    key: 'port',
  },
  {
    title: '用户名', 
    dataIndex: 'user',
    key: 'user',
  },
  {
    title: '密码',
    dataIndex: 'password',
    key: 'password',
  },
];

const onClick = record => {
  console.log(record)
  router.push({ name: 'Mysql', params: { id: record.name } });
};

GetFullCons().then(result => {
  console.log(result)
  data.value = result;
})

</script>