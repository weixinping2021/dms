<template>
    
    <a-table :columns="columnPrefixs" :data-source="prefixs" size="small">
        <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'key'">
                <a-button type="link" @click="buttonClick(record)">{{record.key}}</a-button>
            </template>
        </template>
    </a-table>
    <a-table :columns="columnPrefix" :data-source="prefix" size="small">
    </a-table>
</template>


<script setup>
import { ref, onMounted, } from 'vue';
import { GetExpirePrefixs, GetPrefixDetail } from "../../wailsjs/go/redis/Redis";

const props = defineProps({
    data: Object
})

const prefixs = ref([])
const prefix = ref([])

const columnPrefixs = [{
    title: 'prefixs',
    dataIndex: 'key',
    width: 300,
    ellipsis: true,
},
{
    title: 'Size',
    dataIndex: 'size',
},
{
    title: 'SizeReadable',
    dataIndex: 'sizereadable'
},
{
    title: '元素个数',
    dataIndex: 'elementcount',
},
{
    title: 'DB',
    dataIndex: 'db',
},
{
    title: '过期时间',
    dataIndex: 'expire',
}
];

const columnPrefix = [{
    title: 'key',
    dataIndex: 'key',
    width: 300,
    ellipsis: true,
},
{
    title: 'Size',
    dataIndex: 'size',
},
{
    title: 'SizeReadable',
    dataIndex: 'sizereadable'
},
{
    title: 'Type',
    dataIndex: 'type',
},
{
    title: '元素个数',
    dataIndex: 'elementcount',
},
{
    title: 'DB',
    dataIndex: 'db',
},
{
    title: '过期时间',
    dataIndex: 'expire',
}
];

const buttonClick = (record) => {
    console.log(record);
    GetPrefixDetail(props.data.message, record.key, record.expire).then((result) => {
            //console.log(result);
            prefix.value = result;
        });
};


onMounted(() => {
    console.log(props.data.message);
    GetExpirePrefixs(props.data.message, props.data.keydate).then((result) => {
        //console.log(result);
        prefixs.value = result;
    });
})
</script>