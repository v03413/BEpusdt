<template>
  <component :is="tag" ref="nodes" />
</template>

<script setup lang="ts">
import JsBarcode from "jsbarcode";
const props = defineProps({
  // 节点类型，可接收svg、img、canvas
  tag: {
    type: String,
    default: "svg"
  },
  // 条形码具体内容
  text: {
    type: String,
    default: ""
  },
  // 条形码配置
  // https://github.com/lindell/JsBarcode/wiki/Options
  options: {
    type: Object,
    default() {
      return {};
    }
  }
});

const nodes = ref();
onMounted(() => {
  JsBarcode(nodes.value, props.text, props.options);
});
</script>
