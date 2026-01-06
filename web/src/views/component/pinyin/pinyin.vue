<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-card title="拼音转换" hoverable>
        <a-space direction="vertical" size="large" fill>
          <a-textarea v-model="myText" placeholder="请输入内容" allow-clear />
          <s-pinyin-pro
            v-if="myText"
            :text="myText"
            :options="{
              resultClass: 'my-text-item'
            }"
          />
          <a-empty v-else />
        </a-space>
      </a-card>
      <br />
      <a-card title="默认字符拼音-自定义字体大小" hoverable>
        <s-pinyin-pro
          :text="text"
          :options="{
            resultClass: 'my-text-item'
          }"
        />
      </a-card>
      <br />
      <a-card title="字符拼音-自定义颜色" hoverable>
        <s-pinyin-pro
          :text="text"
          :options="{
            resultClass: 'my-text-item',
            chineseClass: 'my-chinese-item',
            pinyinClass: 'my-pinyin-item'
          }"
        />
      </a-card>
      <br />
      <a-card title="字符拼音-指定文字颜色" hoverable>
        <s-pinyin-pro
          :text="text"
          :options="{
            resultClass: 'my-text-item',
            customClassMap: {
              'red-item': ['节', '气'],
              'blue-item': ['科', '学']
            }
          }"
        />
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
let text =
  "二十四节气是我国独创的传统历法，也是我国历史长河中不可多得的瑰宝，上至风雨雷电，下至芸芸众生，包罗万象。在长期的生产实践中，我国劳动人民通过对太阳、天象的不断观察，开创出了节气这种独特的历法。经过不断的探索、分析和总结，节气的划分逐渐变得科学和丰富，到距今两千多年的秦汉时期，二十四节气已经形成了完整的体系，并一直沿用至今。";
const myText = ref("");
</script>

<style lang="scss" scoped>
:deep(.my-text-item) {
  font-size: 16px;
  line-height: 30px;
}
:deep(.my-chinese-item) {
  color: $color-primary;
}
:deep(.my-pinyin-item) {
  color: $color-danger;
}
:deep(.red-item) {
  color: $color-danger;
}
:deep(.blue-item) {
  color: $color-primary;
}
</style>
