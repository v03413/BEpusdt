import { defineStore } from "pinia";
import persistedstateConfig from "@/store/config/index";
/**
 * 全局配置
 * @methods setCollapsed 设置菜单折叠
 * @methods setRefreshPage 刷新页面
 * @methods setLanguage 设置语言
 */
const themeConfig = () => {
  // 是否折叠菜单
  const collapsed = ref<boolean>(false);
  // 刷新页面
  const refreshPage = ref<boolean>(true);
  // 系统语言
  const language = ref<string>("zh-CN");
  // 黑暗模式
  const darkMode = ref<boolean>(false);
  // 菜单手风琴
  const isAccordion = ref<boolean>(true);
  // 面包屑渲染
  const isBreadcrumb = ref<boolean>(true);
  // 标签栏渲染
  const isTabs = ref<boolean>(true);
  // 页脚渲染
  const isFooter = ref<boolean>(true);
  // 水印
  const watermark = ref<string>("");
  // 水印风格
  const watermarkStyle = ref<{ fontSize: number; color: string }>({
    fontSize: 12,
    color: "rgba(0, 0, 0, 0.15)"
  });
  // 水印角度
  const watermarkRotate = ref<number>(330);
  // 水印间隙
  const watermarkGap = ref<[number, number]>([100, 100]);
  // 防止调试
  const debugPrevention = ref<boolean>(false);
  // 布局模式：layoutDefaults、layoutHead、layoutMixing
  const layoutType = ref<string>("layoutDefaults");
  // 色弱模式
  const colorWeakMode = ref<boolean>(false);
  // 灰色模式
  const grayMode = ref<boolean>(false);
  // 侧边栏深色
  const asideDark = ref<boolean>(false);
  // 页面过渡方式
  const transitionPage = ref<string>("fadeInOut");
  // 主题色
  const themeColor = ref<string>("#165DFF");
  // 调色盘
  const presetColors = ref<string[]>([
    "#165DFF",
    "#F53F3F",
    "#F77234",
    "#FF7D00",
    "#F7BA1E",
    "#FADC19",
    "#9FDB1D",
    "#00B42A",
    "#3491FA",
    "#c34d9c",
    "#722ED1",
    "#E6A23C",
    "#F56C6C",
    "#409EFF",
    "#42b883",
    "#b1d391"
  ]);

  // 折叠菜单
  function setCollapsed(data: boolean) {
    collapsed.value = data;
  }
  // 刷新页面
  function setRefreshPage(data: boolean) {
    refreshPage.value = data;
  }
  // 设置语言
  function setLanguage(data: string) {
    language.value = data;
  }

  return {
    collapsed,
    refreshPage,
    language,
    darkMode,
    isAccordion,
    isBreadcrumb,
    isTabs,
    isFooter,
    watermark,
    watermarkStyle,
    watermarkRotate,
    watermarkGap,
    debugPrevention,
    layoutType,
    colorWeakMode,
    grayMode,
    asideDark,
    transitionPage,
    themeColor,
    presetColors,
    setCollapsed,
    setRefreshPage,
    setLanguage
  };
};

export const useThemeConfig = defineStore("theme-config", themeConfig, {
  persist: persistedstateConfig("theme-config")
});
