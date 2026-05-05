import { computed } from "vue";
import type { ComputedRef } from "vue";
import { useDevicesSize } from "@/hooks/useDevicesSize";

interface LayoutModel {
  formLayout: ComputedRef<"vertical" | "horizontal">;
  tableFixed: ComputedRef<"" | "right">;
  descriptionsLayout: ComputedRef<"inline-vertical" | "right">;
  descriptionsColumn: (min?: number, max?: number) => number;
  dialogWidth: (min?: string, max?: string) => string;
}

export const useLayoutModel = (): LayoutModel => {
  const { isMobile } = useDevicesSize();

  const formLayout = computed(() => (isMobile.value ? "vertical" : "horizontal"));
  const tableFixed = computed(() => (isMobile.value ? "" : "right"));
  const descriptionsLayout = computed(() => (isMobile.value ? "inline-vertical" : "right"));
  const descriptionsColumn = (min: number = 1, max: number = 2) => (isMobile.value ? min : max);
  const dialogWidth = (min: string = "40%", max: string = "95%") => (isMobile.value ? max : min);

  return { formLayout, tableFixed, descriptionsLayout, descriptionsColumn, dialogWidth };
};
