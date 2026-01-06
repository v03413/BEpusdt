<template>
  <a-drawer :width="340" :visible="props.themeOpen" @ok="handleCancel" @cancel="handleCancel" unmount-on-close>
    <template #title> 主题设置 </template>
    <div>
      <div>
        <a-divider orientation="center">导航模式</a-divider>
        <div class="flex-center">
          <a-tooltip v-for="item in layoutList" :key="item.value" :content="item.label" position="top" mini>
            <div
              :class="layoutType === item.value ? `current-layout ${item.class}` : item.class"
              @click="layouetChange(item.value)"
            >
              <icon-check-circle-fill class="layout-icon" />
            </div>
          </a-tooltip>
        </div>
      </div>
      <div class="box-gap">
        <a-divider orientation="center">主题设置</a-divider>
        <div class="flex-center">
          <ColorPicker
            :theme="darkMode ? 'dark' : 'light'"
            sucker-hide
            :colors-default="presetColors"
            :color="themeColor"
            @change-color="themeColorChange"
          />
        </div>
      </div>
      <div class="box-gap">
        <div class="flex-row">
          <div>色弱模式</div>
          <a-switch v-model="colorWeakMode" @change="onColorWeak" />
        </div>
        <div class="flex-row">
          <div>灰色模式</div>
          <a-switch v-model="grayMode" @change="onGray" />
        </div>
        <div class="flex-row">
          <div>侧边栏深色</div>
          <a-switch :disabled="darkMode" v-model="asideDark" />
        </div>
        <div class="flex-row">
          <div>页面过渡</div>
          <a-select v-model="transitionPage" :style="{ width: '120px' }" placeholder="请选择">
            <a-option v-for="item in transitions" :key="item.value" :value="item.value" :label="item.label" />
          </a-select>
        </div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useThemeMethods } from "@/hooks/useThemeMethods";
import { useDevicesSize } from "@/hooks/useDevicesSize";
import { ColorPicker } from "vue-color-kit";
import "vue-color-kit/dist/vue-color-kit.css";

const themeStore = useThemeConfig();
const { layoutType, collapsed, colorWeakMode, grayMode, darkMode, asideDark, transitionPage, themeColor, presetColors } =
  storeToRefs(themeStore);

const layoutList = reactive({
  layoutDefaults: {
    value: "layoutDefaults",
    label: "默认布局",
    class: "layout-defaults"
  },
  layoutHead: {
    value: "layoutHead",
    label: "横向布局",
    class: "layout-head"
  },
  layoutMixing: {
    value: "layoutMixing",
    label: "混合布局",
    class: "layout-mixing"
  }
});

const transitions = ref([
  { value: "fadeInOut", label: "轻过渡" },
  { value: "cardInOut", label: "卡片" },
  { value: "fadeOut", label: "渐退" }
]);

interface IThemeColor {
  hex: string;
  hsv: {
    h: number;
    s: number;
    v: number;
  };
  rgba: {
    r: number;
    g: number;
    b: number;
    a: number;
  };
}
// 主题色设置
const themeColorChange = (value: IThemeColor) => {
  themeColor.value = value.hex;
  const { setThemeColor } = useThemeMethods();
  setThemeColor();
};
// 色弱模式
const onColorWeak = () => {
  const { setColorWeak } = useThemeMethods();
  setColorWeak();
};

// 灰色模式
const onGray = () => {
  const { setGray } = useThemeMethods();
  setGray();
};

// 布局变化
const layouetChange = (type: string) => {
  layoutType.value = type;
  const { isPc } = useDevicesSize();
  collapsed.value = isPc.value ? false : true;
};

const props = defineProps({
  themeOpen: {
    type: Boolean,
    default: false
  }
});

const emits = defineEmits(["themeCancel"]);
const handleCancel = () => {
  emits("themeCancel");
};
</script>

<style lang="scss" scoped>
.box-gap {
  margin-top: 30px;
}
.flex-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $margin;
}
.flex-center {
  display: flex;
  column-gap: $margin;
  align-items: center;
  justify-content: center;
  margin-bottom: $margin;
}
.layout-defaults,
.layout-head,
.layout-mixing {
  width: 70px;
  height: 50px;
  overflow: hidden;
  background: $color-fill-1;
  border-radius: $radius-box-1;
  box-shadow: $shadow-special;
  .layout-icon {
    display: none;
  }
}
.current-layout {
  position: relative;
  .layout-icon {
    position: absolute;
    right: 2px;
    bottom: 2px;
    display: block;
    font-size: $font-size-body-3;
    color: $color-primary;
  }
}
.layout-defaults {
  position: relative;
  &::before {
    position: absolute;
    left: 0;
    width: 10px;
    height: 100%;
    content: "";
    background: #232324;
  }
  &::after {
    position: absolute;
    top: 0;
    right: 0;
    width: 60px;
    height: 10px;
    content: "";
    background: #c2c4c8;
  }
}
.layout-head {
  position: relative;
  &::before {
    position: absolute;
    top: 0;
    width: 100%;
    height: 10px;
    content: "";
    background: #232324;
  }
}
.layout-mixing {
  position: relative;
  &::before {
    position: absolute;
    top: 0;
    width: 100%;
    height: 10px;
    content: "";
    background: #232324;
  }
  &::after {
    position: absolute;
    top: 10px;
    left: 0;
    width: 10px;
    height: calc(100% - 10px);
    content: "";
    background: #232324;
  }
}
:deep(.arco-divider-text) {
  background: $color-bg-3;
}
</style>
