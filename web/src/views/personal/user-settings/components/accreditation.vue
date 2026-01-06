<template>
  <div>
    <a-row align="center" :gutter="[0, 16]">
      <a-col :span="24">
        <a-card title="实名认证">
          <template #extra>
            <a-link>修改认证信息</a-link>
          </template>
          <a-space direction="vertical" size="large">
            <a-descriptions :column="3" :align="{ label: 'right' }" table-layout="fixed">
              <a-descriptions-item v-for="item of data" :key="item.label" :label="item.label">
                <a-tag v-if="item.label == '实名认证：'" :color="item.value == 1 ? 'green' : 'red'">{{
                  item.value == 1 ? "已认证" : "未认证"
                }}</a-tag>
                <span v-else>{{ item.value }}</span>
              </a-descriptions-item>
            </a-descriptions>
          </a-space>
        </a-card>
      </a-col>
      <a-col :span="24">
        <a-card title="认证记录">
          <a-table :columns="tableColumns" :data="tableData" size="small">
            <template #status="{ record }">
              <a-space>
                <a-space>
                  <a-badge status="success" text="已通过" v-if="record.status == 1" />
                  <a-badge status="processing" text="审核中" v-else />
                </a-space>
              </a-space>
            </template>
            <template #optional>
              <a-space>
                <a-button size="mini" type="text">查看</a-button>
                <a-button size="mini" type="text">撤回</a-button>
              </a-space>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
const data = ref([
  {
    label: "用户ID：",
    value: "1000001"
  },
  {
    label: "用户名：",
    value: "admin"
  },
  {
    label: "用户昵称：",
    value: "兔子先森"
  },
  {
    label: "角色：",
    value: "系统管理员"
  },
  {
    label: "权限：",
    value: "全部"
  },
  {
    label: "实名认证：",
    value: 1
  },
  {
    label: "邮箱：",
    value: "2547096351@qq.com"
  },
  {
    label: "手机号：",
    value: "15888888888"
  },
  {
    label: "注册时间：",
    value: "2024-06-20"
  }
]);

const tableColumns = ref([
  {
    title: "认证类型",
    dataIndex: "authenticationType"
  },
  {
    title: "认证内容",
    dataIndex: "authenticationContent"
  },
  {
    title: "当前状态",
    dataIndex: "status",
    align: "center",
    slotName: "status"
  },
  {
    title: "创建时间",
    dataIndex: "createTime"
  },
  {
    title: "操作",
    dataIndex: "operation",
    slotName: "optional",
    fixed: "right",
    width: 200
  }
]);
const tableData = reactive([
  {
    key: "1",
    authenticationType: "证件认证",
    authenticationContent: "实名认证，鉴定信息",
    status: "0",
    createTime: "2024-06-20 17:35:25"
  },
  {
    key: "2",
    authenticationType: "证件认证",
    authenticationContent: "实名认证，鉴定信息",
    status: "1",
    createTime: "2024-06-08 15:25:37"
  },
  {
    key: "3",
    authenticationType: "证件认证",
    authenticationContent: "实名认证，鉴定信息",
    status: "1",
    createTime: "2024-06-01 10:55:21"
  }
]);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
.margin-bottom {
  margin-bottom: $padding;
}
</style>
