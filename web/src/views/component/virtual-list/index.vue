<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space align="start">
        <a-card :style="{ width: '330px' }" title="垂直滚动-单列"> <vertical /> </a-card>
        <a-card :style="{ width: '330px' }" title="垂直滚动-网格">
          <vertical-grid />
        </a-card>
        <a-card :style="{ width: '330px' }" title="横向滚动">
          <horizontal />
        </a-card>
      </a-space>
      <a-divider />
      <a-typography>
        <a-typography-paragraph>
          <ol>
            <li>垂直虚拟列表需要设置父级高度，高度可以是height固定值，也可以是max-height约束值</li>
            <li>
              垂直网格布局的虚拟列表，组件不直接支持，需要手动分组，每行给一个id，一行显示两个
              <a-typography-text code> [{id: uid, list: [1,2]},{id: uid, list: [3,4]}] </a-typography-text>
            </li>
            <li>
              横向滚动，需要设置父元素高度height，高度必须是固定值，无法使用max-height，此外组件还需设置属性<a-typography-text
                code
                >direction="horizontal"</a-typography-text
              >
            </li>
          </ol>
        </a-typography-paragraph>
        <a-typography-paragraph blockquote>
          <div>
            基于：<a-link href="https://github.com/Akryum/vue-virtual-scroller" target="_blank">vue-virtual-scroller</a-link>
          </div>
          <div>
            页面使用需引入：
            <a-typography-text code>
              import { DynamicScroller, DynamicScrollerItem } from "vue-virtual-scroller";
            </a-typography-text>
          </div>
          <div>
            以及css文件：
            <a-typography-text code> import "vue-virtual-scroller/dist/vue-virtual-scroller.css"; </a-typography-text>
          </div>
        </a-typography-paragraph>
      </a-typography>
    </div>
  </div>
</template>
<script setup lang="ts">
import "vue-virtual-scroller/dist/vue-virtual-scroller.css";
import vertical from "./vertical.vue";
import verticalGrid from "./vertical-grid.vue";
import horizontal from "./horizontal.vue";
defineOptions({ name: "virtual-list" });
</script>
<style lang="scss" scoped></style>
