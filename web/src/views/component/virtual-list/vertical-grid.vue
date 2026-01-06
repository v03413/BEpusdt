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
            <div v-for="curr in item.list" :key="curr.id" class="grid-item">
              <a-image :src="myImage" width="50" />
              <div class="item-right">
                <div>甲壳虫-{{ curr.id }}</div>
                <div class="desc-text">Abbey Road</div>
              </div>
            </div>
          </div>
        </DynamicScrollerItem>
      </template>
    </DynamicScroller>
  </div>
</template>
<script setup lang="ts">
import myImage from "@/assets/img/Abbey Road.jpg";
import { v4 as uuidv4 } from "uuid";
import { DynamicScroller, DynamicScrollerItem } from "vue-virtual-scroller";
// 还需引入vue-virtual-scroller的css文件，此处父组件已引入
// import "vue-virtual-scroller/dist/vue-virtual-scroller.css";

const list = ref<any>([]);

const getList = () => {
  let array = [];
  for (let i = 0; i < 500; i++) {
    array.push({
      id: i
    });
  }
  let group = [];
  for (let i = 0; i < array.length; i += 2) {
    group.push({
      id: uuidv4(),
      list: array.slice(i, i + 2)
    });
  }
  list.value = group;
};

getList();
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
      padding: 4px;
      margin-bottom: 4px;
      border: 1px solid $color-border-2;
      .list-item {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        grid-gap: 4px;
        .grid-item {
          display: flex;
          column-gap: 4px;
          border: 1px solid $color-success;
          border-radius: 4px;
          .item-right {
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            padding: 4px 0;
            .desc-text {
              font-size: 12px;
            }
          }
        }
      }
    }
  }
}
</style>
