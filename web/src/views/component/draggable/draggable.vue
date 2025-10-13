<template>
  <div class="snow-page">
    <div class="snow-inner draggable-page">
      <a-row :gutter="16">
        <a-col :xs="24" :sm="24" :md="24" :lg="14" :xl="12" :xxl="12">
          <a-space direction="vertical">
            <div class="draggable-container">
              <div class="title">Favorites</div>
              <draggable
                :list="list"
                class="draggable-box"
                chosen-class="chosenClass"
                animation="300"
                item-key="value"
                @start="onStart"
                @end="onEnd"
              >
                <template #item="{ element }">
                  <div class="draggable-list" :class="!animated && 'animated-fade-up-' + element.value">
                    <div class="draggable-inner">
                      <s-svg-icon :name="element.svg" :size="30" class="svg-box" />
                      <div>{{ element.label }}</div>
                    </div>
                  </div>
                </template>
              </draggable>
            </div>
            <div>
              采用开源拖拽库：
              <a-link href="https://github.com/SortableJS/vue.draggable.next" target="_blank"> vue.draggable </a-link>
            </div>
          </a-space>
        </a-col>
        <a-col :xs="24" :sm="24" :md="24" :lg="10" :xl="12" :xxl="12">
          <a-scrollbar style="height: 396px; overflow: auto" outer-class="scrollbar">
            <s-code-view :code-json="codeJson" style="width: 500px" />
          </a-scrollbar>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import draggable from "vuedraggable";
const list = ref([
  { value: 0, label: "bilibili", svg: "bilibili" },
  { value: 1, label: "微信", svg: "wechat" },
  { value: 2, label: "QQ音乐", svg: "QQ音乐" },
  { value: 3, label: "百度", svg: "百度" },
  { value: 4, label: "优酷", svg: "优酷" },
  { value: 5, label: "网易云音乐", svg: "网易云音乐" },
  { value: 6, label: "抖音", svg: "抖音" },
  { value: 7, label: "西瓜视频", svg: "西瓜视频" },
  { value: 8, label: "微博", svg: "微博" },
  { value: 9, label: "豆瓣", svg: "豆瓣" },
  { value: 10, label: "github", svg: "github" },
  { value: 11, label: "gitee", svg: "gitee" }
]);
const codeJson = computed(() => JSON.stringify({ list: list.value }, null, "\t"));

const animated = ref(false);

// 拖拽开始
const onStart = () => {
  animated.value = true;
};

// 拖拽结束
const onEnd = () => {
  console.log("拖拽结束");
};
</script>

<style lang="scss" scoped>
.draggable-page {
  display: flex;
  column-gap: $padding;
}
.draggable-container {
  padding: $padding;
  background: $color-bg-1;
  border-radius: $radius-box-5;
  box-shadow: inset 0 0 0 1px $color-border-2;
  .title {
    margin-bottom: $padding;
    font-weight: bold;
  }
  .draggable-box {
    display: flex;
    flex-wrap: wrap;
    gap: $padding $padding;
    width: 448px;
    .draggable-list {
      box-sizing: border-box;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 100px;
      height: 100px;
      background: $color-fill-1;
      border: 1px solid $color-border-2;
      border-radius: $radius-box-5;
      .draggable-inner {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        .svg-box {
          margin-bottom: $margin;
        }
      }
    }
  }
  .chosenClass {
    color: $color-primary;
  }
}
.scrollbar {
  height: 100%;
}
</style>
