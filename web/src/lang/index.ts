import { createI18n } from "vue-i18n";
import zhCN from "@/lang/modules/zhCN";
import enUS from "@/lang/modules/enUS";

/* 获取语言 */
const getLang = () => {
  let store = localStorage.getItem("theme-config");
  if (store) {
    return JSON.parse(store)?.language || "zh-CN";
  } else {
    return "zh-CN";
  }
};
/* 这里必须是messages名称 */
const messages = {
  "zh-CN": zhCN,
  "en-US": enUS
};

const i18n = createI18n({
  legacy: false, // Composition API模式需要设为false
  globalInjection: true, // 全局生效: $
  locale: getLang(), // 默认语言
  messages, // 数据源
  missing: (_: string, key: string) => {
    return removeBeforeFirstDot(key);
  }
});

/**
 * @param { string } str 国际化key
 * @returns 去掉第一个点之前的字符串，例如："menu.home" => "home"，否则返回原字符串
 */
function removeBeforeFirstDot(str: string) {
  const dotIndex = str.indexOf(".");
  return dotIndex >= 0 ? str.slice(dotIndex + 1) : str;
}

export default i18n;
