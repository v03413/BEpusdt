import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
import { generate, getRgbStr } from "@arco-design/color";
/**
 * 主题处理hooks，内置多种主题处理场景
 * @returns 主题方法
 */
export const useThemeMethods = () => {
  /**
   * @description: 初始化主题
   */
  const initTheme = () => {
    // 黑暗模式和主题色
    setDarkMode();
    // 色弱模式和灰色模式
    const themeStore = useThemeConfig();
    const { grayMode } = storeToRefs(themeStore);
    if (grayMode.value) {
      setGray();
    } else {
      setColorWeak();
    }
  };

  /**
   * @description: 暗黑模式
   */
  const setDarkMode = () => {
    const themeStore = useThemeConfig();
    const { darkMode, asideDark } = storeToRefs(themeStore);
    if (darkMode.value) {
      // 设置为暗黑主题
      document.body.setAttribute("arco-theme", "dark");
      asideDark.value = false; // 黑暗模式与侧边栏深色互斥
    } else {
      // 恢复亮色主题
      document.body.removeAttribute("arco-theme");
    }
    // 黑暗模式切换后需要更新主题色
    setThemeColor();
  };

  /**
   * @description: 主题色设置
   */
  const setThemeColor = () => {
    const themeStore = useThemeConfig();
    const { themeColor, darkMode } = storeToRefs(themeStore);
    // 对于给定的颜色，使用算法生成包含10种颜色的渐变样本。这适用于光模式和暗模式。
    let list = generate(themeColor.value, { list: true, dark: darkMode.value });
    list.forEach((color: string, index: number) => {
      // https://arco.design/palette/list
      // arco的颜色等级为1-10，这里需要index+1
      document.body.style.setProperty(`--primary-${index + 1}`, getRgbStr(color));
    });
  };

  /**
   * @description: 色弱模式
   */
  const setColorWeak = () => {
    // 获取html
    const htmlCase = document.querySelector("html") as HTMLHtmlElement | null;
    if (!htmlCase) return;
    const themeStore = useThemeConfig();
    const { colorWeakMode, grayMode } = storeToRefs(themeStore);
    if (colorWeakMode.value) {
      grayMode.value = false; // 色弱模式和灰色模式互斥
      htmlCase.style.filter = "invert(80%)"; // 反转色80%
    } else {
      htmlCase.style.filter = "";
    }
  };

  /**
   * @description: 灰色模式
   */
  const setGray = () => {
    // 获取html
    const htmlCase = document.querySelector("html") as HTMLHtmlElement | null;
    if (!htmlCase) return;
    const themeStore = useThemeConfig();
    const { colorWeakMode, grayMode } = storeToRefs(themeStore);
    if (grayMode.value) {
      colorWeakMode.value = false; // 色弱模式和灰色模式互斥
      htmlCase.style.filter = "grayscale(100%)"; // 灰度100%
    } else {
      htmlCase.style.filter = "";
    }
  };

  return {
    initTheme,
    setDarkMode,
    setThemeColor,
    setColorWeak,
    setGray
  };
};
