import { App, Directive } from "vue";
import antiShake from "@/directives/modules/global/anti-shake";
import throttle from "@/directives/modules/global/throttle";
import custom from "@/directives/modules/global/custom";
import hasPerm from "@/directives/modules/permission/has-perm";
import hasRole from "@/directives/modules/permission/has-role";

// 定义安装函数
// install 函数是一个对象中的方法，其作用是将一系列指令对象安装到 Vue 应用实例中，它自带两个参数：app 和 options
// app就是vue实例，options则是安装函数的参数(可选)
// install 方法的定义就插件对象，install 方法会在引入插件时自动被 vue 调用，并传参 vue 实例和 options

const group: { [key: string]: Directive } = {
  antiShake,
  throttle,
  custom,
  hasPerm,
  hasRole
};

const directives = {
  install(app: App<Element>) {
    // 将一系列自定义指令对象安装到 Vue 应用实例中
    for (const key in group) {
      app.directive(key, group[key]);
    }
  }
};

// 导出安装函数
export default directives;
