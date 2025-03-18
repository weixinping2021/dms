<template>
<a-table :columns="columnPrefix" :data-source="prefixDetail" size="small" >
</a-table>
</template>
<script setup>
import {ref,onMounted,} from 'vue';
import {GetPrefixDetail} from "../../wailsjs/go/redis/Redis";

const props = defineProps({
    data: Object
})

const prefixDetail = ref([])

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
        dataIndex: 'Expire',
    }
];

onMounted(() => {
    console.log(props.data.message);
    GetPrefixDetail(props.data.message, props.data.prefix, "").then((result) => {
            //console.log(result);
            prefixDetail.value = result;
        });
})
</script>