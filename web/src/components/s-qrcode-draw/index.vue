<template>
  <component :is="tag" ref="nodes" />
</template>

<script setup lang="ts">
import QRCode from "qrcode";

const props = defineProps({
  // 节点类型，可接收img、canvas
  tag: {
    type: String,
    default: "img"
  },
  // 条形码具体内容
  text: {
    type: String,
    default: ""
  },
  // 二维码配置
  // https://www.npmjs.com/package/qrcode
  options: {
    type: Object,
    default() {
      return {};
    }
  }
});

const nodes = ref();
const createQRCode = () => {
  if (props.tag === "canvas") {
    QRCode.toCanvas(nodes.value, props.text, props.options, function (error: any) {
      if (error) return console.error(error);
    });
  } else {
    QRCode.toDataURL(props.text, props.options, function (error: any, url: any) {
      if (error) return console.error(error);
      nodes.value.src = url; // 将二维码图像显示在 img 元素中
    });
  }
};

onMounted(() => {
  createQRCode();
});
</script>

<style lang="scss" scoped></style>
