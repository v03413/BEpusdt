<template>
  <div>
    <a-input ref="inputRef" :style="{ width: '100%' }" placeholder="请选择图标" v-model="modelValue" @focus="onFocus">
      <template #suffix v-if="modelValue">
        <s-svg-icon v-if="type == 'svg'" :name="modelValue" :size="size" />
        <component v-else :is="modelValue" :size="size"></component>
      </template>
      <template #append>
        <span class="icon-reset" @click="reset">重置</span>
      </template>
    </a-input>
    <a-modal v-model:visible="visible" :unmount-on-close="true" width="70vw" :top="'50px'" :align-center="false" :footer="false">
      <template #title> 请选择图标 </template>
      <a-input
        :style="{
          width: '100%',
          marginBottom: '14px'
        }"
        placeholder="请输入关键词搜索"
        allow-clear
        v-model="searchName"
      >
        <template #prefix>
          <icon-search />
        </template>
      </a-input>
      <div class="over-scroll">
        <div class="icon-grid">
          <div v-for="item in iconList" :key="item" class="grid-item" @click="onIcon(item)">
            <s-svg-icon v-if="type == 'svg'" :name="item" :size="25" />
            <component v-else :is="item" :size="25"></component>
            <span>{{ item }}</span>
          </div>
        </div>
        <div v-if="iconList.length === 0" class="empty-row">
          <a-empty>
            <template #image>
              <icon-exclamation-circle-fill />
            </template>
            <span>未查询到您要找的图标</span>
          </a-empty>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import * as ArcoIcons from "@arco-design/web-vue/es/icon";
interface Props {
  modelValue: string;
  type?: string;
  size?: number;
}
const props = withDefaults(defineProps<Props>(), {
  modelValue: "", // 图标名称-默认空字符串
  type: "arco", // 图标库类型，arco 或 svg
  size: 25 // 图标大小-默认25px
});

const { modelValue, type, size } = toRefs(props);

const emit = defineEmits(["update:modelValue"]);

const visible = ref<boolean>(false);

const reset = () => {
  emit("update:modelValue", "");
};

const onFocus = () => {
  searchName.value = "";
  visible.value = true;
};

// svg模块-assets/svgs文件夹下的所有svg图标
const SvgIconModules = import.meta.glob("@assets/svgs/*.svg");
// 搜索关键字
const searchName = ref<string>("");
// icon列表
const iconList = computed(() => {
  let icons: string[] = [];
  if (type.value === "arco") {
    if (!ArcoIcons) return [];
    for (let key in ArcoIcons) {
      if (key != "default") icons.push(key);
    }
  } else {
    if (!SvgIconModules) return [];
    for (let path in SvgIconModules) {
      icons.push(path.replace("/src/assets/svgs/", "").split(".svg")[0]);
    }
  }
  // 若按照关键字搜索，则返回搜索结果
  if (searchName.value) {
    return icons.filter(item => item.toLowerCase().includes(searchName.value.toLowerCase()));
  }
  // 若关键字为空则返回全部列表
  return icons;
});

// 选择图标
const onIcon = (name: string) => {
  visible.value = false;
  emit("update:modelValue", name);
};
</script>

<style lang="scss" scoped>
.over-scroll {
  max-height: 70vh;
  overflow: auto;
  .icon-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, 150px); /* 自适应填充列，每列宽为 150px */
    justify-content: space-evenly; /* 水平均匀分配列之间的空隙 */
    overflow: hidden;
    .grid-item {
      display: flex;
      flex-direction: column;
      row-gap: $margin;
      align-items: center;
      justify-content: center;
      height: 100px;
      transition: transform 0.2s ease-in-out; /* 添加过渡效果 */
      &:hover {
        transform: scale(1.3); /* 鼠标移入时放大到 1.3 倍 */
      }
    }
  }
}
.empty-row {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 300px;
}
.icon-reset {
  cursor: pointer;
}
</style>
