<template>
  <div>
    <a-config-provider :global="true" :locale="locale">
      <slot></slot>
    </a-config-provider>
  </div>
</template>

<script setup lang="ts">
import zhCN from "@arco-design/web-vue/es/locale/lang/zh-cn";
import enUS from "@arco-design/web-vue/es/locale/lang/en-us";
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
const themeStore = useThemeConfig();
const { language } = storeToRefs(themeStore);
interface Lang {
  [key: string]: any;
}
const locales = ref<Lang>({
  "zh-CN": zhCN,
  "en-US": enUS
});

const locale = computed(() => {
  return locales.value[language.value] || zhCN;
});
</script>

<style lang="scss" scoped></style>
