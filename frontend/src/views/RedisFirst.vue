<template>
    <a-collapse v-model:activeKey="activeKey">
        <a-collapse-panel key="1" header="总内存分布">
            <a-descriptions>
                <a-descriptions-item label="永久">{{ expireData[0].forever }}({{ expireData[1].forever
                }}keys)</a-descriptions-item>
                <a-descriptions-item label="会过期">{{ expireData[0].expire }}({{ expireData[1].expire
                }}keys)</a-descriptions-item>
                <a-descriptions-item label="总和">{{ expireData[0].total }}({{ expireData[1].total
                }}keys)</a-descriptions-item>
            </a-descriptions>
        </a-collapse-panel>
        <a-collapse-panel key="2" header="永久key内存分布">
            <v-chart class="chart" :option="forver" style="width: 100%; height: 500px;" @click="forverClick" />
        </a-collapse-panel>
        <a-collapse-panel key="3" header="会过期key内存分布">
            <v-chart class="chart" :option="expire" style="width: 100%; height: 500px;" @click="onChartClick" />
        </a-collapse-panel>
        <a-collapse-panel key="4" header="prefix详情">
            <a-tabs v-model:activeKey="tabactiveKey" hide-add type="editable-card" @edit="removeTab">
                <a-tab-pane v-for="tab in tabs" :key="tab.key" :tab="tab.title">
                    <!-- 动态内容 -->
                    <component :is="tab.component" :data="tab.data" />
                </a-tab-pane>
            </a-tabs>
        </a-collapse-panel>
    </a-collapse>
</template>

<script setup>
import {ref,computed,onMounted,watch,shallowRef} from 'vue';
import RedisPrefix from './RedisPrefix.vue';
import RedisPrefixs from './RedisPrefixs.vue';
import { GetRedisMemory, GetExpireMemoryPic,GetPrefixDetail, GetForverMemoryPic} from "../../wailsjs/go/redis/Redis";
import { use} from 'echarts/core';
import { BarChart,LineChart,ScatterChart} from 'echarts/charts';
import { GridComponent,TooltipComponent,TitleComponent, DataZoomComponent} from 'echarts/components';
import {CanvasRenderer} from 'echarts/renderers';
import VChart from 'vue-echarts';

const props = defineProps({
    data: Object
})

const expireDataPic = ref([])
const activeKey = ref([])
const forverDataPic = ref([])
const expireData = ref([])
const tabactiveKey = ref()
// 存放动态创建的 Tabs
const tabs = ref([]);
watch(activeKey, val => {
    console.log(val);
});
use([
    CanvasRenderer,
    BarChart,
    GridComponent,
    TooltipComponent,
    TitleComponent,
    LineChart,
    ScatterChart,
    DataZoomComponent
]);

const forver = computed(() => {
    return {
        title: {
            left: 'center',
            text: '永久key内存分布(按前缀)'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        xAxis: {
            type: 'category',
            data: forverDataPic.value[0]
        },
        yAxis: [{
            type: 'value', // 第一个 y 轴，适用于柱状图的销售量
            name: '数量',
        },
        {
            type: 'value', // 第一个 y 轴，适用于柱状图的销售量
            name: '内存使用量(MB)',
        }
        ],
        series: [{
            name: '数量',
            data: forverDataPic.value[1],
            type: "bar",
            yAxisIndex: 0,
        },
        {
            name: '内存使用量',
            type: 'bar',
            data: forverDataPic.value[2],
            yAxisIndex: 1, // 使用第一个 y 轴
        }
        ]
    }
});

const expire = computed(() => {
    return {
        title: {
            left: 'center',
            text: '过期key内存分布'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        dataZoom: {
            type: 'inside',
            filterMode: 'none', // 或 'weak'/'empty'
            zoomOnMouseWheel: true, // 确保启用滚轮缩放
            moveOnMouseMove: true,
        },
        xAxis: {
            type: 'category',
            data: expireDataPic.value[0]
        },
        yAxis: [{
            type: 'value', // 第一个 y 轴，适用于柱状图的销售量
            name: '数量',
        },
        {
            type: 'value', // 第一个 y 轴，适用于柱状图的销售量
            name: '内存使用量(MB)',
        }
        ],
        series: [{
            name: '数量',
            data: expireDataPic.value[1],
            type: "bar",
            yAxisIndex: 0,
        },
        {
            name: '内存使用量',
            type: 'bar',
            data: expireDataPic.value[2],
            yAxisIndex: 1, // 使用第一个 y 轴
        }
        ]
    }
});

const removeTab = (targetKey) => {
    tabs.value = tabs.value.filter(tab => tab.key !== targetKey);
    // 处理删除后的 Tab
    if (tabs.value.length > 0) {
        tabactiveKey.value = tabs.value[tabs.value.length - 1].key;
    } else {
        tabactiveKey.value = null;
    }
};

function onChartClick(params) {
    if (params.seriesType === 'bar') {
        console.log('点击的柱子:', params.name);
        let newTab = {
            key: params.name,
            title:params.name,
            component: null,  // 默认没有组件
            data: {}  // 用来存放要传递的数据
        };
        newTab.component = shallowRef(RedisPrefixs);
        newTab.data = { message: props.data.message, keydate: params.name };
        tabs.value.push(newTab);
        tabactiveKey.value = params.name; // 切换到新 Tab
        activeKey.value = "4";
    }
}

function forverClick(params) {
    if (params.seriesType === 'bar') {
        console.log('点击的柱子:', params.name);
        let newTab = {
            key: params.name,
            title:'永久'+ params.name,
            component: null,  // 默认没有组件
            data: {}  // 用来存放要传递的数据
        };
        newTab.component = shallowRef(RedisPrefix);
        newTab.data = { message: props.data.message, prefix: params.name };
        tabs.value.push(newTab);
        tabactiveKey.value = params.name; // 切换到新 Tab
        activeKey.value = "4";
    }
}

onMounted(() => {
    console.log(props.data.message);
    GetRedisMemory(props.data.message).then((result) => {
        expireData.value = result;
    });
    GetForverMemoryPic(props.data.message).then((result) => {
        console.log(result);
        forverDataPic.value = result;
    });
    GetExpireMemoryPic(props.data.message).then((result) => {
        console.log(result);
        expireDataPic.value = result;
    });
})
</script>
