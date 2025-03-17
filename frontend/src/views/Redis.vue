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
    <a-divider />
    <a-tabs v-model:activeKey="activeKey" type="card">
        <a-tab-pane key="1" tab="过期时间统计数据">
            <a-card>
                <a-card-grid style="width: 50%; text-align: center"> <a-table :dataSource="dataMemorys"
                        :columns="memoryCols">
                        <template #summary>
                            <a-table-summary-row>
                                <a-table-summary-cell>总和</a-table-summary-cell>
                                <a-table-summary-cell>
                                    {{ totalMomery.momeryForever + totalMomery.momerys3 + totalMomery.momeryb3s7 + totalMomery.momeryb7 }}
                                </a-table-summary-cell>
                            </a-table-summary-row>
                        </template>
                    </a-table>
                </a-card-grid>
                <a-card-grid style="width: 50%; text-align: center">
                    <a-table :dataSource="dataCounts"
                        :columns="countsCols">
                        <template #summary>
                            <a-table-summary-row>
                                <a-table-summary-cell>总和</a-table-summary-cell>
                                <a-table-summary-cell>
                                    {{ totalMomery.countForever + totalMomery.counts3 + totalMomery.countb3s7 + totalMomery.countb7 }}
                                </a-table-summary-cell>
                            </a-table-summary-row>
                        </template>
                    </a-table>
                </a-card-grid>
            </a-card>


        </a-tab-pane>
        <a-tab-pane key="2" tab="Top 500 BigKey(按内存)">
            <a-table :columns="columns" :data-source="dataMemory"></a-table>
        </a-tab-pane>
        <a-tab-pane key="3" tab="Top 500 BigKey(按过期时间)">
            <a-table :columns="columns" :data-source="dataExpire"></a-table>
        </a-tab-pane>
        <a-tab-pane key="4" tab="Top 500 key前缀(按内存)">
            <a-table :columns="columnPrefix" :data-source="dataPrefix"></a-table>
        </a-tab-pane>
        <a-tab-pane key="5" tab="Top 500 key前缀匹配查询(按内存)">
            <a-input-search v-model:value="prefix" placeholder="input search text" style="width: 200px"
                @search="onSearch" />
            <a-table :columns="columns" :data-source="dataPrefixSearch"></a-table>
        </a-tab-pane>
    </a-tabs>
</template>

<script setup>
import { ref } from 'vue';
import { FileAddOutlined } from '@ant-design/icons-vue';
import { OpenDialog } from "../../wailsjs/go/main/App";
import { GetRedisMemory, GetRedisKeys, AnalyseRdb, GetRedisTop500Prefix, GetPrefixkeys, GetRdbResultTitle } from "../../wailsjs/go/redis/Redis";
const spinning = ref(false);
const filename = ref("")
const totalMomery = ref([])
const dataMemory = ref([])
const dataExpire = ref([])
const dataPrefix = ref([])
const dataPrefixSearch = ref([])
const prefix = ref("")
const formattedTime = ref("")
const rdbresult = ref([])
const dataMemorys = ref([])
const dataCounts = ref([])


const columns = [
    {
        title: 'Name',
        dataIndex: 'name',
        ellipsis: true,
    },
    {
        title: 'Size',
        dataIndex: 'size',
    },
    {
        title: 'Expire',
        dataIndex: 'expire',
    },
    {
        title: 'Type',
        dataIndex: 'type',
    },
];

const memoryCols = [
    {
        title: 'key分类',
        dataIndex: 'keystype',
        key: 'keystype',
    },
    {
        title: '内存占用(bytes)',
        dataIndex: 'memory',
        key: 'memory',
        sorter: (a, b) => a.memory - b.memory,
    }
]

const countsCols = [
    {
        title: 'key分类',
        dataIndex: 'keystype',
        key: 'keystype',
    },
    {
        title: '个数',
        dataIndex: 'count',
        key: 'count',
        sorter: (a, b) => a.count - b.count,
    }
]
const columnPrefix = [
    {
        title: 'Prefix',
        dataIndex: 'name',
        ellipsis: true,
    },
    {
        title: 'Size',
        dataIndex: 'size',
    },
    {
        title: 'key个数',
        dataIndex: 'type',
    }
];

function openChoseFileDlg() {
    OpenDialog().then(res => {
        console.log("chose file :", res)
        if (res != "") {
            filename.value = res
            console.log(res)
        }
    })
}

const activeKey = ref('1');

function Analyse() {
    spinning.value = true;
    AnalyseRdb(filename.value).then((result) => {
        spinning.value = false;
        formattedTime.value = result
        GetRedisMemory(formattedTime.value).then((result) => {
            totalMomery.value = result;
            dataMemorys.value = [
                { key: '1', keystype: '不过期', memory: result.momeryForever },
                { key: '2', keystype: '3天内过期', memory: result.momerys3 },
                { key: '3', keystype: '3-7天内过期', memory: result.momeryb3s7 },
                { key: '4', keystype: '>7天过期', memory: result.momeryb7 }
            ];
            dataCounts.value = [
                { key: '1', keystype: '不过期', count: result.countForever },
                { key: '2', keystype: '3天内过期', count: result.counts3 },
                { key: '3', keystype: '3-7天内过期', count: result.countb3s7 },
                { key: '4', keystype: '>7天过期', count: result.countb7 }
            ];
        });
        GetRedisKeys("size", formattedTime.value).then((result) => {
            console.log(result);
            dataMemory.value = result;
        });
        GetRedisKeys("days", formattedTime.value).then((result) => {
            console.log(result);
            dataExpire.value = result;
        });
        GetRedisTop500Prefix(formattedTime.value).then((result) => {
            console.log(result);
            dataPrefix.value = result;
        });
    });

}
function onSearch() {
    GetPrefixkeys(prefix.value, formattedTime.value).then((result) => {
        console.log(result);
        dataPrefixSearch.value = result;
    });
}
function focus() {
    GetRdbResultTitle().then((result) => {
        console.log(result);
        rdbresult.value = result;
    });
}

function handleChange() {
    GetRedisMemory(formattedTime.value).then((result) => {
        totalMomery.value = result;
        dataMemorys.value = [
            { key: '1', keystype: '不过期', memory: result.momeryForever },
            { key: '2', keystype: '3天内过期', memory: result.momerys3 },
            { key: '3', keystype: '3-7天内过期', memory: result.momeryb3s7 },
            { key: '4', keystype: '>7天过期', memory: result.momeryb7 }
        ];
        dataCounts.value = [
                { key: '1', keystype: '不过期', count: result.countForever },
                { key: '2', keystype: '3天内过期', count: result.counts3 },
                { key: '3', keystype: '3-7天内过期', count: result.countb3s7 },
                { key: '4', keystype: '>7天过期', count: result.countb7 }
            ];
    });
    GetRedisKeys("size", formattedTime.value).then((result) => {
        console.log(result);
        dataMemory.value = result;
    });
    GetRedisKeys("days", formattedTime.value).then((result) => {
        console.log(result);
        dataExpire.value = result;
    });
    GetRedisTop500Prefix(formattedTime.value).then((result) => {
        console.log(result);
        dataPrefix.value = result;
    });
}

</script>