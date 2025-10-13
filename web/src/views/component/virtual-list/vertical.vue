<template>
  <div class="dynamic-scroller">
    <DynamicScroller :items="list" :min-item-size="54" class="scroller">
      <template #default="{ item, index, active }">
        <DynamicScrollerItem
          :item="item"
          :active="active"
          :size-dependencies="[item.id]"
          :data-index="index"
          :data-active="active"
          class="message"
        >
          <div class="list-item">
            <a-image :src="myImage" width="50" />
            <div class="item-right">
              <div>月之暗面-{{ item.id }}</div>
              <div>The Dark Side of The Moon</div>
            </div>
          </div>
        </DynamicScrollerItem>
      </template>
    </DynamicScroller>
  </div>
</template>
<script setup lang="ts">
import myImage from "@/assets/img/The Dark Side of The Moon.jpg";
import { DynamicScroller, DynamicScrollerItem } from "vue-virtual-scroller";
// 还需引入vue-virtual-scroller的css文件，此处父组件已引入
// import "vue-virtual-scroller/dist/vue-virtual-scroller.css";

const list = ref<any>([]);

for (let i = 0; i < 500; i++) {
  list.value.push({
    id: i
  });
}
</script>
<style lang="scss" scoped>
.dynamic-scroller {
  display: flex;
  flex-direction: column;
  width: 300px;
  max-height: 400px;
  overflow: hidden;
  border: 1px solid $color-primary;
  border-radius: 4px;
  .scroller {
    flex: 1;
    .message {
      padding: 12px;
      margin-bottom: 4px;
      border: 1px solid $color-border-2;
      .list-item {
        display: flex;
        column-gap: 4px;
        .item-right {
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          padding: 4px 0;
        }
      }
    }
  }
}
</style>
