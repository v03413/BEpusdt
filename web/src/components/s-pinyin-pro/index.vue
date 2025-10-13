<template>
  <div>
    <div v-html="content"></div>
  </div>
</template>

<script setup lang="ts">
import { html } from "pinyin-pro";
// https://github.com/zh-lx/pinyin-pro
interface Props {
  text: string;
  options?: {
    // 汉字和拼音最外层 <span> 标签的类名
    resultClass?: string;
    // 拼音外层 <rt> 标签的类名
    pinyinClass?: string;
    // 汉字外层 <span> 标签的类名
    chineseClass?: string;
    // 是否用 <span> 标签包裹非汉字字符。
    wrapNonChinese?: boolean;
    // 非汉字字符外层 <span> 标签的类名 （仅在 wrapNonChinese 值为 true 时有效）
    nonChineseClass?: string;
    // 拼音中音调的显示形式
    tone?: boolean;
    // v3.24.2+  是否输出 <rp>(</rp> 和 <rp>)</rp> 标签
    rp?: boolean;
    // 指定文字类名 key：类名 value：文字数组
    customClassMap?: {
      [key: string]: string[];
    };
  };
}

// withDefaults为defineProps标注类型
// https://cn.vuejs.org/guide/typescript/composition-api.html#typing-component-props
const props = withDefaults(defineProps<Props>(), {
  text: "", // 默认值
  options: () => ({})
});

const content = ref<string>();
const createPinyin = () => {
  content.value = html(props.text, props.options);
};

watch(
  () => props.text,
  () => {
    createPinyin();
  }
);

createPinyin();
</script>

<style lang="scss" scoped></style>
