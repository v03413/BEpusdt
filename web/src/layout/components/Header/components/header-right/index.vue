<template>
  <div class="header_setting" :class="isMobile && 'head-absolute-fix'">
    <!-- 国际化 -->
    <!-- 暂时也是真没精力搞国际化 -->
    <!-- <a-dropdown trigger="hover" @select="onLange">
      <a-button size="mini" type="text" class="icon_btn" id="system-language">
        <template #icon>
          <icon-language :size="18" />
        </template>
      </a-button>
      <template #content>
        <a-doption :disabled="language === 'zh-CN'">{{ $t(`system.zh-CN`) }}</a-doption>
        <a-doption :disabled="language === 'en-US'">{{ $t(`system.en-US`) }}</a-doption>
      </template>
    </a-dropdown> -->
    <!-- 皮肤设置 -->
    <a-tooltip :content="$t(`system.system-settings`)">
      <a-button size="mini" type="text" class="icon_btn" id="system-settings" @click="onSystemSetting">
        <template #icon>
          <icon-skin :size="18" />
        </template>
      </a-button>
    </a-tooltip>
    <!-- 切换黑夜模式 -->
    <a-tooltip :content="$t(`system.${!darkMode ? 'switch-to-night-mode' : 'switch-to-daytime-mode'}`)">
      <a-button size="mini" type="text" class="icon_btn" id="system-dark" @click="onNightMode">
        <template #icon>
          <icon-sun-fill :size="18" v-if="!darkMode" />
          <icon-moon-fill :size="18" v-else />
        </template>
      </a-button>
    </a-tooltip>
    <!-- 我的 -->
    <a-dropdown trigger="hover">
      <div class="my_setting" id="system-my-setting">
        <!-- <a-image width="32" height="32" fit="cover" :src="account.user.avatar || myImage" class="my_image" /> -->
        <span class="user-nickname">Hi {{ admin_username || "BEpusdt" }}</span>
        <div class="icon_down">
          <icon-down style="stroke-width: 3" />
        </div>
      </div>
      <template #content>
        <!-- 项目地址 -->
        <a-doption @click="onProject">
          <template #default>
            <s-svg-icon :name="'github'" :size="18" />
            <span class="margin-left-text">{{ $t(`system.project-address`) }}</span>
          </template>
        </a-doption>
        <a-divider margin="0" />
        <!-- 退出登录 -->
        <a-doption @click="logOut">
          <template #default>
            <s-svg-icon :name="'exit'" :size="18" />
            <span class="margin-left-text">{{ $t(`system.logout`) }}</span>
          </template>
        </a-doption>
      </template>
    </a-dropdown>
  </div>
  <SystemSettings :system-open="systemOpen" @system-cancel="systemOpen = false" />
  <ThemeSettings :theme-open="themeOpen" @theme-cancel="themeOpen = false" />
</template>

<script setup lang="ts">
// import Notice from "@/layout/components/Header/components/Notice/index.vue";
import SystemSettings from "@/layout/components/Header/components/system-settings/index.vue";
import ThemeSettings from "@/layout/components/Header/components/theme-settings/index.vue";
// import { useI18n } from "vue-i18n";
import { Modal } from "@arco-design/web-vue";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useUserInfoStore } from "@/store/modules/user-info";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useThemeMethods } from "@/hooks/useThemeMethods";
import { useDevicesSize } from "@/hooks/useDevicesSize";
import { useRouteConfigStore } from "@/store/modules/route-config";

// const i18n = useI18n();
const router = useRouter();
const { isMobile } = useDevicesSize();
const userStore = useUserInfoStore();
const themeStore = useThemeConfig();
// const { language, darkMode } = storeToRefs(themeStore);
const { darkMode } = storeToRefs(themeStore);
const { admin_username } = storeToRefs(userStore); // 添加这行来获取 admin_username

// 系统设置
const systemOpen = ref(false);
const onSystemSetting = () => {
  systemOpen.value = true;
};

// 主题设置
const themeOpen = ref(false);
// const onThemeSetting = () => {
//   themeOpen.value = true;
// };
// 全屏
// const fullScreen = ref(true);
// const onFullScreen = () => {
//   if (!document.fullscreenElement) {
//     document.documentElement.requestFullscreen();
//     fullScreen.value = false;
//   } else {
//     if (document.exitFullscreen) {
//       document.exitFullscreen();
//       fullScreen.value = true;
//     }
//   }
// };

// 黑暗模式
const onNightMode = () => {
  darkMode.value = !darkMode.value;
  let { setDarkMode } = useThemeMethods();
  setDarkMode();
};

// 语言
// const onLange = (value: string) => {
//   if (value === "Chinese" || value === "中文") {
//     themeStore.setLanguage("zh-CN");
//   } else {
//     themeStore.setLanguage("en-US");
//   }
//   i18n.locale.value = language.value;
// };

// 项目地址
const onProject = () => {
  window.open("https://github.com/v03413/BEpusdt", "_blank");
};

// 退出登录
const logOut = () => {
  Modal.warning({
    title: "提示",
    content: "确定退出登录？",
    hideCancel: false,
    closable: true,
    onBeforeOk: async () => {
      try {
        // 用户退出
        await userStore.logOut();
        router.replace("/login");
        // 清除路由数据
        useRouteConfigStore().resetRoute();
        return true;
      } catch {
        return false;
      }
    }
  });
};
</script>

<style lang="scss" scoped>
.head-absolute-fix {
  position: absolute;
  top: 0;
  right: $padding;
}
.header_setting {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  background-color: $color-bg-2;
  > .icon_btn {
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: space-around;
    width: $icon-box;
    height: $icon-box;
    margin-left: $margin;
    color: $color-text-1;
    border-radius: $radius-box-1;
  }
  .my_setting {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 32px;
    margin-left: $margin;
    overflow: hidden;
    .my_image {
      margin-right: 8px;
      border-radius: 50%;
    }
    .user-nickname {
      white-space: nowrap;
    }
    .icon_down {
      margin: 0 0 0 5px;
      transform: rotate(0deg);
      transition: transform 0.2s;
    }
  }
}
.notice {
  position: relative;
  &::before {
    position: absolute;
    top: -4px;
    right: -2px;
    width: 6px;
    height: 6px;
    content: "";
    background: $color-danger;
    border: 2px solid #ffffff;
    border-radius: 50%;
  }
}
:deep(.arco-dropdown-open) {
  .icon_down {
    transform: rotate(180deg) !important;
  }
}
.margin-left-text {
  margin-left: $margin-text;
}
</style>
