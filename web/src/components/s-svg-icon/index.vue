<template>
  <svg
    aria-hidden="true"
    :class="svgClass"
    :style="{
      color: color,
      fill: color,
      width: iconSize,
      height: iconSize
    }"
    v-bind="$attrs"
  >
    <use :xlink:href="iconName" :fill="color" />
  </svg>
</template>

<script lang="ts" setup>
import { computed } from "vue";
defineOptions({ name: "s-svg-icon" });
const props = defineProps({
  name: {
    type: String,
    default: ""
  },
  color: {
    type: String,
    default: ""
  },
  size: {
    type: [Number, String],
    default: 15
  }
});

// 判断传入的值，是否带有单位，如果没有，就默认用px单位
const getUnitValue = (value: string | number): string | number => {
  return /(px|em|rem|%)$/.test(value.toString()) ? value : value + "px";
};

// svg大小
const iconSize = computed<string | number>(() => {
  return getUnitValue(props.size);
});

// svg名称-对应资源文件夹的svg名称
const iconName = computed<string>(() => `#icon-${props.name}`);

// svg动态类名
const svgClass = computed<string>(() => {
  if (props.name) return `svg-icon icon-${props.name}`;
  return "svg-icon";
});
</script>

<style lang="scss" scope>
.svg-icon {
  flex-shrink: 0;
  width: auto;
  height: auto;
  vertical-align: middle;
}
</style>
