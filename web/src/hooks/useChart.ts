import { computed } from "vue";
import { storeToRefs } from "pinia";
import { EChartsOption } from "echarts";
import { useThemeConfig } from "@/store/modules/theme-config";

interface optionsFn {
  (isDark?: boolean, palette?: string[]): EChartsOption;
}

/**
 * ECharts图表hooks
 * @param sourceOption ECharts配置项函数
 * @returns { option: ComputedRef<EChartsOption>, theme: ComputedRef<string | undefined> } ECharts配置项和主题
 */
export function useChart(sourceOption: optionsFn): {
  option: ComputedRef<EChartsOption>;
  theme: ComputedRef<string | undefined>;
} {
  const themeStore = useThemeConfig();
  const { darkMode } = storeToRefs(themeStore);

  // 预设调色盘，适配arco.design主题：https://visactor.io/vchart/guide/tutorial_docs/Theme/Arco_Design
  const palette = ["#4080FF", "#55C5FD", "#FF7D00", "#4CD263", "#A871E3", "#F7BA1E", "#9FDB1D", "#F979B7", "#0FC6C2", "#E865DF"];

  // echarts support https://echarts.apache.org/zh/theme-builder.html
  const option = computed<EChartsOption>(() => {
    return sourceOption(darkMode.value, palette);
  });

  // 主题
  const theme = computed(() => (darkMode.value ? "dark" : undefined));

  return {
    // ECharts配置项
    option,
    // 黑暗模式
    theme
  };
}
