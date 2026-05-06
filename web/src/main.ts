import { createApp } from "vue";
import "@/style.css";
import App from "@/App.vue";
// arco-design
import ArcoVue from "@arco-design/web-vue";
// vue-router
import router from "@/router/index";
// pinia
import pinia from "@/store/index";
// arco-css
import "@arco-design/web-vue/dist/arco.css";
// 额外引入图标库
import ArcoVueIcon from "@arco-design/web-vue/es/icon";
// 注册全局svg
import "virtual:svg-icons-register";
// 引入i18n
import i18n from "@/lang/index";
// 引入字体
import "@/assets/fonts/fonts.scss";
// 引入自定义指令
import directives from "@/directives/index";

const app = createApp(App);

// app.use(plugin, options)
// 其中 plugin 表示要传递的插件对象， options 参数是可选的，表示选项配置
// https://cn.vuejs.org/api/application.html#app-use

app.use(ArcoVue, {
  componentPrefix: "arco"
});
app.use(pinia);
app.use(ArcoVueIcon);
app.use(router);
app.use(i18n);
app.use(directives);
app.mount("#app");
