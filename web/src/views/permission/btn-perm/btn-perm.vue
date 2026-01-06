<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-alert>
        当前页面
        <a-tag color="green"> 超级管理员 </a-tag>
        和
        <a-tag color="green"> 普通角色 </a-tag>
        都可查看
      </a-alert>
      <br />
      <a-card hoverable title="根据用户权限判断">
        <div>
          <a-typography-paragraph> <a-typography-text bold> 当前用户权限 </a-typography-text></a-typography-paragraph>
          <s-code-view :code-json="permissions" style="width: 100%" />
          <br />
          <a-row :gutter="[24, 12]">
            <a-col :span="24">
              <s-code-view type="vue" :code-json="sysBtn" style="width: 100%" />
            </a-col>
            <a-col :span="24">
              <a-space>
                <span>权限按钮：</span>
                <a-button v-hasPerm="['sys:btn:add']" type="primary">新增</a-button>
                <a-button v-hasPerm="['sys:btn:edit']" type="primary" status="success">编辑</a-button>
                <a-button v-hasPerm="['sys:btn:delete']" type="primary" status="danger">删除</a-button>
              </a-space>
            </a-col>
          </a-row>
          <br />
          <a-row :gutter="[24, 12]">
            <a-col :span="24">
              <s-code-view type="vue" :code-json="commonBtn" style="width: 100%" />
            </a-col>
            <a-col :span="24">
              <a-space>
                <span>权限按钮：</span>
                <a-button v-hasPerm="['common:btn:add']" type="primary">新增</a-button>
                <a-button v-hasPerm="['common:btn:edit']" type="primary" status="success">编辑</a-button>
                <a-button v-hasPerm="['common:btn:delete']" type="primary" status="danger">删除</a-button>
              </a-space>
            </a-col>
          </a-row>
        </div>
      </a-card>
      <br />
      <a-card hoverable title="根据角色权限判断">
        <div>
          <a-typography-paragraph> <a-typography-text bold> 当前角色权限 </a-typography-text></a-typography-paragraph>
          <s-code-view :code-json="roles" style="width: 100%" />
          <br />
          <a-row :gutter="[24, 12]">
            <a-col :span="24">
              <s-code-view type="vue" :code-json="roleBtn" style="width: 100%" />
            </a-col>
            <a-col :span="24">
              <a-space>
                <span>权限按钮：</span>
                <a-button v-hasRole="['admin']" type="primary">超级管理员</a-button>
                <a-button v-hasRole="['common']" type="primary" status="success">普通角色</a-button>
              </a-space>
            </a-col>
          </a-row>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useUserInfoStore } from "@/store/modules/user-info";
let userInfoStore = useUserInfoStore();
let { account } = storeToRefs(userInfoStore);
const permissions = computed(() => JSON.stringify(account.value.permissions, null));
const roles = computed(() => JSON.stringify(account.value.roles, null));

let sysBtn = `<a-space>
    <a-button v-hasPerm="['sys:btn:add']" type="primary">新增</a-button>
    <a-button v-hasPerm="['sys:btn:edit']" type="primary" status="success">编辑</a-button>
    <a-button v-hasPerm="['sys:btn:delete']" type="primary" status="danger">删除</a-button>
</a-space>`;

let commonBtn = `<a-space>
    <a-button v-hasPerm="['common:btn:add']" type="primary">新增</a-button>
    <a-button v-hasPerm="['common:btn:edit']" type="primary" status="success">编辑</a-button>
    <a-button v-hasPerm="['common:btn:delete']" type="primary" status="danger">删除</a-button>
</a-space>`;

let roleBtn = `<a-space>
    <a-button v-hasRole="['admin']" type="primary">超级管理员</a-button>
    <a-button v-hasRole="['common']" type="primary" status="success">普通角色</a-button>
</a-space>`;
</script>

<style lang="scss" scoped></style>
