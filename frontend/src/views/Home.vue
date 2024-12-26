<template>
  <a-table bordered :data-source="data" :columns="columns">
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.dataIndex === 'name'">
        <a @click="onClick(record)">{{ text }}</a>
      </template>
      <template v-if="['host', 'user', 'password', 'port'].includes(column.dataIndex)">
        <div>
          <a-input v-if="editableData[record.key]" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
          <template v-else>
            {{ text }}
          </template>
        </div>
      </template>
      <template v-else-if="column.dataIndex === 'operation'">
        <div class="editable-row-operations">
          <span v-if="editableData[record.key]">
            <a-typography-link @click="save(record.key)">Save</a-typography-link>
            <a-popconfirm title="Sure to cancel?" @confirm="cancel(record.key)">
              <a>Cancel</a>
            </a-popconfirm>
          </span>
          <span v-else>
            <a @click="edit(record.key)">Edit</a>
          </span>
        </div>
      </template>
    </template>
  </a-table>
</template>
<script setup>

import { GetFullCons } from "../../wailsjs/go/main/App";
import { ref ,reactive} from 'vue';
import { useRouter } from 'vue-router';
const router = useRouter()

const columns = [
  {
    title: '连接名',
    dataIndex: 'name',
    key: 'name',
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
    title: 'operation',
    dataIndex: 'operation',
  },
];
const data = ref([])

const onClick = record => {
  console.log(record)
  router.push({ name: 'MysqlUser', params: { id: record.name } });
};

GetFullCons().then(result => {
  console.log(result)
  data.value = result;
})

const editableData = reactive({});
const edit = key => {
  editableData[key] = cloneDeep(dataSource.value.filter(item => key === item.key)[0]);
};
const save = key => {
  Object.assign(dataSource.value.filter(item => key === item.key)[0], editableData[key]);
  delete editableData[key];
};
const cancel = key => {
  delete editableData[key];
};

</script>