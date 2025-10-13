<template>
  <code-mirror
    v-model="codeValue"
    :extensions="extensions"
    :basic="config.basic"
    :tab-size="config.tabSize"
    :dark="config.dark"
    :readonly="config.readonly"
  />
</template>

<script setup lang="ts">
import CodeMirror from "vue-codemirror6";
import { oneDark } from "@codemirror/theme-one-dark";
import { json } from "@codemirror/lang-json";
import { javascript } from "@codemirror/lang-javascript";
import { vue } from "@codemirror/lang-vue";
// Props接口类型
interface Props {
  type?: "javascript" | "vue" | "json";
  codeJson: string;
}
// withDefaults为defineProps标注类型
// https://cn.vuejs.org/guide/typescript/composition-api.html#typing-component-props
const props = withDefaults(defineProps<Props>(), {
  type: "javascript", // 默认值
  codeJson: "" // 默认值
});

// 通过计算属性实现数据实时更新
// 将序列化的字符串美观输出，\t 换行
const codeValue = computed(() => props.codeJson);

// 基础配置
const config = ref({
  tabSize: 2,
  basic: true,
  dark: true,
  readonly: true
});

// 扩展
// 基础扩展：黑暗模式
// 根据默认值来判断添加什么扩展
const extensions = computed(() => {
  const arr = [oneDark];
  if (props.type === "javascript") {
    arr.push(javascript());
  }
  if (props.type === "vue") {
    arr.push(vue());
  }
  if (props.type === "json") {
    arr.push(json());
  }
  return arr;
});
</script>

<style lang="scss" scoped>
// 修改code中的字体样式
:deep(.ͼ1 .cm-scroller) {
  font-family: source-code-pro, Menlo, Monaco, Consolas, "Courier New", monospace;
}
</style>
