<template>
    <div class="components-page-header-demo-responsive" style="border: 1px solid rgb(235, 237, 240)">
        <a-page-header title="Mysql" :sub-title="route.params.id" @back="() => $router.go(-1)">
            <template #footer>
                <a-tabs v-model:activeKey="activeKey" type="card">
                    <a-tab-pane key="1" tab="会话管理">
                        <a-divider style="height: 2px; background-color: black" />
                        <a-row>
                            <a-statistic title="总连接数" :value=total style="margin-right: 30px" />
                            <a-statistic title="活跃连接" :value=active style="margin-right: 30px" />
                            <a-statistic title="连接使用率" :precision="2" suffix="%" :value=active/total style="margin-right: 30px" />
                        </a-row>
                        <a-divider />
                        <a-flex gap="small">
                            <a-radio-group v-model:value="processType" button-style="solid">
                                <a-radio-button value="all">全部会话</a-radio-button>
                                <a-radio-button value="alive">活跃会话</a-radio-button>
                            </a-radio-group>
                            <a-button type="primary" @click="GetMysqlProcesslistFresh">刷新</a-button>
                            <a-button @click="Kill">杀死选中会话</a-button>
                        </a-flex>
                        <a-table :columns="columns" :data-source="data" :row-selection="rowSelection" :scroll="{ x: 1500, y: 1000 }" bordered size="small">
                            <template #bodyCell="{ column, record }">
                                <template v-if="record.children != null && column.key === 'id'">
                                    <a-tag color="pink">{{ Object.keys(record.children).length }}</a-tag>
                                </template>
                            </template>
                        </a-table>
                        <a-flex gap="small">
                            <a-card title="按用户">
                                <a-table :columns="columnsU" :data-source="dataU" bordered size="small"></a-table>
                            </a-card>
                            <a-card title="按访问来源">
                                <a-table :columns="columnsI" :data-source="dataI" bordered size="small"></a-table>
                            </a-card>
                            <a-card title="按数据库">
                                <a-table :columns="columnD" :data-source="dataD" bordered size="small"></a-table>
                            </a-card>
                        </a-flex>
                    </a-tab-pane>
                    <a-tab-pane key="2" tab="锁分析">
                        <a-divider style="height: 2px; background-color: black" />
                        <div>
                            <a-flex gap="small">
                            <a-button type="primary" @click="GetMysqlLocks">刷新</a-button>
                            <a-button @click="LockKill">杀死选中会话</a-button>
                        </a-flex>
                            <a-table :columns="columnsL" :data-source="dataL" :row-selection="rowSelection" :scroll="{ x: 1500, y: 1000 }" bordered size="small"></a-table>
                        </div>
                    </a-tab-pane>
                </a-tabs>
            </template>
        </a-page-header>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import { GetMysqlProcesslist, KillMysqlProcesses,GetConsStatus ,GetConspercent,GetMysqlLock} from "../../wailsjs/go/main/App";
import { message } from 'ant-design-vue';
const total = ref("")
const active = ref("")
const dbId = ref("")
const route = useRoute()
const activeKey = ref('1');
const processType = ref("alive");
const data = ref([])
const dataL = ref([])
const selectedRowsToKill = ref([])
dbId.value = route.params.id
const columnsU = [
    {
        title: '用户',
        dataIndex: 'Name',
    },
    {
        title: '活跃数',
        dataIndex: 'Actives',
    },
    {
        title: '总数',
        dataIndex: 'Total',
    },
];
const columnsI = [
    {
        title: 'ip',
        dataIndex: 'Name',
    },
    {
        title: '活跃数',
        dataIndex: 'Actives',
    },
    {
        title: '总数',
        dataIndex: 'Total',
    },
];
const columnD = [
    {
        title: '数据库名',
        dataIndex: 'Name',
    },
    {
        title: '活跃数',
        dataIndex: 'Actives',
    },
    {
        title: '总数',
        dataIndex: 'Total',
    },
];
const dataU = ref([]);
const dataI = ref([]);
const dataD = ref([]);

const columnsL = [
    {
        title: 'ID',
        dataIndex: 'id',
        key: 'id',
        width: '70',
        fixed: 'left',
    },
    {
        title: '用户',
        dataIndex: 'user',
        key: 'user',
    },
    {
        title: '主机',
        dataIndex: 'host',
        key: 'host',
        ellipsis: true,
    },
    {
        title: '数据库名',
        dataIndex: 'dbname',
        key: 'dbname',
    },
    {
        title: '命令',
        dataIndex: 'command',
        key: 'command',
    },
    {
        title: '执行时间',
        dataIndex: 'time',
        key: 'time',
    },
    {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        ellipsis: true,
    },
    {
        title: 'sql',
        dataIndex: 'sql',
        key: 'sql',
        width: '30%',
        ellipsis: true,
    }
];

const columns = [
    {
        title: 'ID',
        dataIndex: 'id',
        key: 'id',
        width: '70',
        fixed: 'left',
    },
    {
        title: '用户',
        dataIndex: 'user',
        key: 'user',
    },
    {
        title: '主机',
        dataIndex: 'host',
        key: 'host',
        ellipsis: true,
    },
    {
        title: '数据库名',
        dataIndex: 'dbname',
        key: 'dbname',
    },
    {
        title: '命令',
        dataIndex: 'command',
        key: 'command',
    },
    {
        title: '执行时间',
        dataIndex: 'time',
        key: 'time',
    },
    {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
        ellipsis: true,
    },
    {
        title: 'sql',
        dataIndex: 'sql',
        key: 'sql',
        width: '30%',
        ellipsis: true,
    }
];

const rowSelection = ref({
    checkStrictly: false,
    onChange: (selectedRowKeys, selectedRows) => {
        console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
    },
    onSelect: (record, selected, selectedRows) => {
        console.log(record, selected, selectedRows);
        selectedRowsToKill.value = selectedRows
    },
    onSelectAll: (selected, selectedRows, changeRows) => {
        console.log(selected, selectedRows, changeRows);
    },
});

function GetMysqlProcesslistFresh() {
    GetMysqlProcesslist(dbId.value, processType.value).then(result => {
        console.log(result)
        data.value = result;
    })
}

function GetMysqlLocks(){
    GetMysqlLock(dbId.value).then(result => {
        console.log(result)
        dataL.value = result;
    })
}

function Kill() {
    //console.log(selectedRowsToKill.value)
    KillMysqlProcesses(dbId.value, selectedRowsToKill.value).then(result => {
        console.log(result)
        if (result === "success") {
            message.success('kill done');
        }
        GetMysqlProcesslist(dbId.value, processType.value).then(result => {
            console.log(result)
            data.value = result;
        })
    })
}

GetMysqlProcesslist(dbId.value, processType.value).then(result => {
    console.log(result)
    data.value = result;
    GetConsStatus(dbId.value).then(result =>{
        console.log(result)
        dataU.value = result.User
        dataI.value = result["Ip"]
        dataD.value = result["Db"]
    })
})

GetConspercent(dbId.value).then(result => {
    console.log(result)
    total.value = result.total
    active.value = result.active
})


</script>

<style scoped>
.components-page-header-demo-responsive {
    padding-bottom: 20px;
}

.components-page-header-demo-responsive tr:last-child td {
    padding-bottom: 0;
}

#components-page-header-demo-responsive .content {
    display: flex;
}

#components-page-header-demo-responsive .ant-statistic-content {
    font-size: 20px;
    line-height: 28px;
}

@media (max-width: 576px) {
    #components-page-header-demo-responsive .content {
        display: block;
    }

    #components-page-header-demo-responsive .main {
        width: 100%;
        margin-bottom: 12px;
    }

    #components-page-header-demo-responsive .extra {
        width: 100%;
        margin-left: 0;
        text-align: left;
    }
}
</style>