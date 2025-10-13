<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-row :gutter="24">
        <a-col :span="12">
          <a-card title="普通路由">
            <a-alert style="width: 100%" type="success" class="content" :show-icon="false">
              <div>1. 无需在路由配置中声明，可以随意添加任意参数</div>
              <div>2. 参数以键值对的形式出现在 URL 的 ? 之后，用 & 分隔</div>
              <div>3. 例如：/user?id=123&name=John，其中 id=123 和 name=John 是 query 参数</div>
              <div>4. 参数是可选的非必填项，可以自由传递</div>
            </a-alert>
            <a-alert style="width: 100%" type="success" class="content" :show-icon="false">
              <div>路由：/functions/routing-operation/common-route</div>
            </a-alert>
            <a-space direction="vertical" fill>
              <a-input placeholder="请输入路由" allow-clear disabled v-model="connom.path">
                <template #prepend> 路由： </template>
              </a-input>
              <a-input placeholder="请输入参数" allow-clear v-model="connom.query.name">
                <template #prepend> 参数： </template>
              </a-input>
              <a-input placeholder="请输入参数" allow-clear v-model="connom.query.text">
                <template #prepend> 参数： </template>
              </a-input>
              <a-space>
                <a-button type="primary" @click="onQuery">跳转-query</a-button>
              </a-space>
            </a-space>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="动态匹配">
            <a-alert style="width: 100%" type="success" class="content" :show-icon="false">
              <div>1. 需要在路由配置中明确定义动态参数，否则参数会被忽略</div>
              <div>2. 参数会直接体现在 URL 路径中，通常用于动态路由匹配</div>
              <div>3. 例如：/user/123，其中 123 是 params 参数（对应路由 /user/:id）</div>
              <div>4. 如果在路由配置中声明了动态参数（如 /:id），则参数是必填的</div>
            </a-alert>
            <a-alert style="width: 100%" type="success" class="content" :show-icon="false">
              <div>路由：/functions/routing-operation/dynamic-route/:name/:text</div>
            </a-alert>
            <a-space direction="vertical" fill>
              <a-input placeholder="请输入路由" allow-clear disabled v-model="dynamic.name">
                <template #prepend> 路由： </template>
              </a-input>
              <a-input placeholder="请输入参数" allow-clear v-model="dynamic.params.name">
                <template #prepend> 参数： </template>
              </a-input>
              <a-input placeholder="请输入参数" allow-clear v-model="dynamic.params.text">
                <template #prepend> 参数： </template>
              </a-input>
              <a-space>
                <a-button type="primary" @click="onParams">跳转-params</a-button>
              </a-space>
            </a-space>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: "routing-operation" });
const router = useRouter();
// 普通跳转
const connom = ref({
  path: "/functions/routing-operation/common-route",
  query: {
    name: "SnowAdmin",
    text: "月之暗面"
  }
});
const onQuery = () => {
  let route = {
    path: connom.value.path, // 跳转路由
    query: connom.value.query // 路由参数
  };
  router.push(route);
};

// 动态匹配
const dynamic = ref({
  name: "dynamic-route",
  params: {
    name: "SnowAdmin",
    text: "月之暗面"
  }
});
const onParams = () => {
  try {
    let route = {
      name: dynamic.value.name, // 跳转路由
      params: dynamic.value.params // 路由参数
    };
    router.push(route);
  } catch {
    arcoMessage("warning", "请完整填写参数");
  }
};
</script>

<style lang="scss" scoped>
.content {
  width: 400px;
  margin: $margin 0;
}
</style>
