/// <reference types="vite/client" />
declare module "*.vue" {
  import { DefineComponent } from "vue";
  const component: DefineComponent<object, object, any>;
  export default component;
}
declare module "vue-i18n";
declare module "@arco-design/color";
declare module "sortablejs";
declare module "vue-codemirror6";
declare module "@codemirror/theme-one-dark";
declare module "@codemirror/lang-json";
declare module "@codemirror/lang-javascript";
declare module "@codemirror/lang-vue";
declare module "nprogress";
declare module "@wangeditor/editor-for-vue";
declare module "@/directives/modules/custom";
declare module "mockjs";
declare module "@/store/modules/route-config";
declare module "fingerprintjs2";
declare module "@arco-design/web-vue";
declare module "vue-router";
declare module "pinia";
declare module "postcss-preset-env";
declare module "qrcode";
declare module "vue-virtual-scroller";
