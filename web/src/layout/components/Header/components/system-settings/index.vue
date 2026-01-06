<template>
  <a-drawer :width="340" :visible="props.systemOpen" @ok="handleCancel" @cancel="handleCancel" unmount-on-close>
    <template #title> 系统设置 </template>
    <div>
      <div>
        <div class="title">界面设置</div>
        <div class="flex-row">
          <div>菜单折叠</div>
          <a-switch v-model="collapsed" />
        </div>
        <div class="flex-row">
          <div>菜单手风琴</div>
          <a-switch v-model="isAccordion" />
        </div>
        <div class="flex-row">
          <div>面包屑</div>
          <a-switch v-model="isBreadcrumb" />
        </div>
        <div class="flex-row">
          <div>标签栏</div>
          <a-switch v-model="isTabs" @change="tabsChange" />
        </div>
        <div class="flex-row">
          <div>页脚</div>
          <a-switch v-model="isFooter" />
        </div>
      </div>
      <div>
        <div class="title">水印设置</div>
        <div class="flex-row">
          <div>水印颜色</div>
          <pick-colors
            :theme="darkMode ? 'dark' : 'light'"
            show-alpha
            format="rgb"
            v-model:value="watermarkStyle.color"
            :colors="['rgba(0, 0, 0, 0.15)']"
            :z-index="2000"
          />
        </div>
        <div class="flex-row">
          <div>水印文案</div>
          <a-input :style="{ width: '100px' }" v-model="watermark" placeholder="请输入" allow-clear />
        </div>
        <div class="flex-row">
          <div>水印大小</div>
          <a-slider v-model="watermarkStyle.fontSize" :min="10" :max="50" :style="{ width: '100px' }" />
        </div>
        <div class="flex-row">
          <div>水印角度</div>
          <a-slider v-model="watermarkRotate" :min="0" :max="360" :style="{ width: '100px' }" />
        </div>
        <div class="flex-row">
          <div>水印间隙</div>
          <a-slider :default-value="gapInfo[0]" :min="0" :max="300" :style="{ width: '100px' }" @change="onWatermarkGap" />
        </div>
      </div>
      <div>
        <div class="title">系统设置</div>
        <div class="flex-row">
          <div>防调试</div>
          <a-switch v-model="debugPrevention" />
        </div>
      </div>
    </div>
  </a-drawer>
</template>

<script setup lang="ts">
import PickColors from "vue-pick-colors";
import { storeToRefs } from "pinia";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { useThemeConfig } from "@/store/modules/theme-config";
import { currentlyRoute } from "@/router/route-output";
import { DebugControl } from "@/utils/debug-prevention";
const themeStore = useThemeConfig();
const routerStore = useRouteConfigStore();
const {
  collapsed,
  isAccordion,
  isBreadcrumb,
  isTabs,
  isFooter,
  watermark,
  watermarkStyle,
  watermarkRotate,
  watermarkGap,
  darkMode,
  debugPrevention
} = storeToRefs(themeStore);
const { tabsList, cacheRoutes } = storeToRefs(routerStore);
const route = useRoute();
const props = defineProps({
  systemOpen: {
    type: Boolean,
    default: false
  }
});

const gapInfo = ref(watermarkGap.value);
const onWatermarkGap = (e: number) => {
  watermarkGap.value = watermarkGap.value.map(() => e);
};

/*
  是否关闭tabs栏
  如果关闭，那么所有tabs全部取消、所有页面缓存全部取消
  如果开启，那么添加当前路由到tabs
*/
const tabsChange = (e: boolean) => {
  if (!e) {
    tabsList.value = [];
    cacheRoutes.value = [];
  } else {
    currentlyRoute(route);
  }
};

// 监听debug开关;
const debugControl = new DebugControl();
watch(
  () => debugPrevention.value,
  newValue => {
    if (newValue) {
      debugControl.start();
    } else {
      debugControl.stop();
    }
  },
  {
    immediate: true
  }
);

const emits = defineEmits(["systemCancel"]);
const handleCancel = () => {
  emits("systemCancel");
};
</script>

<style lang="scss" scoped>
.title {
  margin-bottom: $margin;
  font-size: $font-size-title-1;
  font-weight: bold;
  color: $color-text-1;
}
.flex-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $margin;
}
</style>
