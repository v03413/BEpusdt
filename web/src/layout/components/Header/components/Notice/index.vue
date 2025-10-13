<template>
  <a-tabs :default-active-key="current" :active-key="current" @tab-click="onTab">
    <a-tab-pane v-for="item in noticeData" :key="item.id">
      <template #title>{{ `${$t(`system.${item.title}`)}(${item.data.length})` }}</template>
      <div class="notice" v-for="content in item.data" :key="content.id">
        <a-image width="36" height="36" fit="cover" :src="myImage" class="notice_img" />
        <div class="content margin-left-text">
          <div>
            <span class="nickname">{{ content.nickname }}</span>
            <span class="time margin-left-text">{{ content.time }}</span>
          </div>
          <div class="context">{{ content.content }}</div>
        </div>
      </div>
      <a-empty v-if="item.data.length === 0" />
    </a-tab-pane>
  </a-tabs>
</template>

<script setup lang="ts">
import myImage from "@/assets/img/my-image.jpg";
const data = ref([
  { id: 100, img: "", time: "1分钟前", nickname: "兔子先森", content: "一键三连" },
  { id: 120, img: "", time: "1小时前", nickname: "捷克大力士", content: "与你握手问好" },
  { id: 130, img: "", time: "2小时前", nickname: "forever", content: "you forever" }
]);
const noticeData = ref([
  {
    id: 1,
    title: "notice",
    data: data.value
  },
  { id: 2, title: "message", data: [] },
  { id: 3, title: "backlog", data: [] }
]);
const current = ref<number>(1);
const onTab = (key: number) => {
  current.value = key;
};
</script>

<style lang="scss" scoped>
.notice {
  display: flex;
  align-items: center;
  margin-bottom: $margin;
  .notice_img {
    width: 36px;
    height: 36px;
    border-radius: 50%;
  }
  .content {
    width: 200px;
    .nickname {
      font-size: $font-size-body-3;
      color: $color-text-2;
    }
    .time {
      font-size: $font-size-body-1;
      color: $color-text-3;
    }
    .context {
      font-size: $font-size-body-3;
    }
  }
}
.margin-left-text {
  margin-left: $margin-text;
}

// 解决tabs宽度异常的bug
:deep(.arco-tabs-content .arco-tabs-content-list) {
  display: unset;
}
</style>
