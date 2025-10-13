import { computed } from "vue";
import { useWindowSize } from "@vueuse/core";

interface DeviceSize {
  isPc: ComputedRef<boolean>;
  isPad: ComputedRef<boolean>;
  isMobile: ComputedRef<boolean>;
}
/**
 * 小于 768px：适配移动设备（如手机）。
 * 768px - 1024px：适配平板设备。
 * 大于 1024px：适配桌面设备。
 * @returns PC窗口isPc, 平板窗口isPad, 移动端窗口isMobile
 */
export const useDevicesSize = (): DeviceSize => {
  const { width } = useWindowSize();
  const isPc = computed(() => width.value > 1024);
  const isPad = computed(() => width.value > 768 && width.value <= 1024);
  const isMobile = computed(() => width.value <= 768);
  return { isPc, isPad, isMobile };
};
