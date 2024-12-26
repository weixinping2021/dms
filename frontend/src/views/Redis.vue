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
                <a-select
                ref="select"
                v-model:value="formattedTime"
                style="width: 400px"
                :options="rdbresult"
                @focus="focus"
                @change="handleChange"
                ></a-select>
            </a-flex>
        </a-spin>
    </a-row>
    <a-divider />
    <a-tabs v-model:activeKey="activeKey" type="card">
        <a-tab-pane key="1" tab="内存占比图">
            <a-card>
                <a-card-grid style="width: 50%; text-align: center"><v-chart class="chart" :option="option"
                        autoresize /></a-card-grid>
                <a-card-grid style="width: 50%; text-align: center"><v-chart class="chart" :option="option1"
                        autoresize /></a-card-grid>
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
            <a-input-search
            v-model:value="prefix"
            placeholder="input search text"
            style="width: 200px"
            @search="onSearch"
            />
            <a-table :columns="columns" :data-source="dataPrefixSearch"></a-table>
        </a-tab-pane>
    </a-tabs>
</template>

<script setup>
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { PieChart } from 'echarts/charts';
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components';
import VChart, { THEME_KEY } from 'vue-echarts';
import { ref, provide } from 'vue';
import { FileAddOutlined } from '@ant-design/icons-vue';
import { GetRedisMemory, OpenDialog, GetRedisKeys, AnalyseRdb,GetRedisTop500Prefix,GetPrefixkeys,GetRdbResultTitle } from "../../wailsjs/go/main/App";
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
use([
    CanvasRenderer,
    PieChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
]);

provide(THEME_KEY, 'light');

const option = ref({
    title: {
    text: '内存占用',
    left: 'center',
  },
    tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b} : {c} ({d}%)',
    },
    series: [
        {
            name: '内存占比',
            type: 'pie',
            data: [],
            emphasis: {
                itemStyle: {
                    shadowBlur: 10,
                    shadowOffsetX: 0,
                    shadowColor: 'rgba(0, 0, 0, 0.5)',
                },
            },
        },
    ],
})

const option1 = ref({
    title: {
    text: 'key个数占用',
    left: 'center',
  },
    tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b} : {c} ({d}%)',
    },
    series: [
        {
            name: 'key个数占比',
            type: 'pie',
            data: [],
            emphasis: {
                itemStyle: {
                    shadowBlur: 10,
                    shadowOffsetX: 0,
                    shadowColor: 'rgba(0, 0, 0, 0.5)',
                },
            },
        },
    ],
})


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
            console.log(result);
            totalMomery.value = result;
            option.value.series[0].data = [
                { value: result.momeryForever, name: '不过期' },
                { value: result.momerys3, name: '3天内过期' },
                { value: result.momeryb3s7, name: '3-7天内过期' },
                { value: result.momeryb7, name: '>7天过期' },
                // 添加或更新数据
            ];
            option1.value.series[0].data = [
                { value: result.countForever, name: '不过期' },
                { value: result.counts3, name: '3天内过期' },
                { value: result.countb3s7, name: '3-7天内过期' },
                { value: result.countb7, name: '>7天过期' },
                // 添加或更新数据
            ];
        });
        GetRedisKeys("size",formattedTime.value).then((result) => {
            console.log(result);
            dataMemory.value = result;
        });
        GetRedisKeys("days",formattedTime.value).then((result) => {
            console.log(result);
            dataExpire.value = result;
        });
        GetRedisTop500Prefix(formattedTime.value).then((result) => {
            console.log(result);
            dataPrefix.value = result;
        });
    });

}
function onSearch(){
    GetPrefixkeys(prefix.value,formattedTime.value).then((result) => {
            console.log(result);
            dataPrefixSearch.value = result;
        });
}
function focus(){
    GetRdbResultTitle().then((result) => {
            console.log(result);
            rdbresult.value = result;
        });
}

function handleChange(){
    GetRedisMemory(formattedTime.value).then((result) => {
            console.log(result);
            totalMomery.value = result;
            option.value.series[0].data = [
            { value: result.momeryForever, name: '不过期' },
                { value: result.momerys3, name: '3天内过期' },
                { value: result.momeryb3s7, name: '3-7天内过期' },
                { value: result.momeryb7, name: '>7天过期' },
                // 添加或更新数据
            ];
            option1.value.series[0].data = [
            { value: result.countForever, name: '不过期' },
                { value: result.counts3, name: '3天内过期' },
                { value: result.countb3s7, name: '3-7天内过期' },
                { value: result.countb7, name: '>7天过期' },
                // 添加或更新数据
            ];
        });
        GetRedisKeys("size",formattedTime.value).then((result) => {
            console.log(result);
            dataMemory.value = result;
        });
        GetRedisKeys("days",formattedTime.value).then((result) => {
            console.log(result);
            dataExpire.value = result;
        });
        GetRedisTop500Prefix(formattedTime.value).then((result) => {
            console.log(result);
            dataPrefix.value = result;
        });
}

</script>

<style scoped>
.chart {
    height: 50vh;
}
</style>