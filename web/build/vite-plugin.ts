import vue from "@vitejs/plugin-vue";
import { resolve } from "path";
import { PluginOption } from "vite";
import { vitePluginForArco } from "@arco-plugins/vite-vue";
import { createSvgIconsPlugin } from "vite-plugin-svg-icons";
import AutoImport from "unplugin-auto-import/vite";
import { ArcoResolver } from "unplugin-vue-components/resolvers";
import Components from "unplugin-vue-components/vite";
import { viteMockServe } from "vite-plugin-mock";
import eslintPlugin from "vite-plugin-eslint";
/**
 * 创建 vite 插件
 * @param viteEnv
 */
export const createVitePlugins = (viteEnv: ViteEnv): (PluginOption | PluginOption[])[] => {
  const env = viteEnv;
  return [
    vue(),
    // esLint 报错信息显示在浏览器界面上
    eslintPlugin(),
    vitePluginForArco({
      style: "css"
    }),
    createSvgIconsPlugin({
      // 配置src下存放svg的路径，这里表示在src/assets/svgs文件夹下
      iconDirs: [resolve(process.cwd(), "src/assets/svgs")],
      symbolId: "icon-[dir]-[name]"
    }),
    AutoImport({
      // 自动导入 Vue 相关函数，如：ref, reactive, toRef 等
      imports: ["vue", "vue-router"],
      // 自动导入的目录-自定义全局函数
      dirs: ["src/globals"],
      // arco组件的按需加载
      resolvers: [ArcoResolver()],
      // 解决eslint报错问题
      eslintrc: {
        // 这里先设置成true然后npm run dev 运行之后会生成 .eslintrc-auto-import.json 文件之后，在改为false
        enabled: false,
        filepath: "./.eslintrc-auto-import.json", // 生成的文件路径
        globalsPropValue: true
      },
      // 配置文件生成位置
      dts: "src/auto-import.d.ts"
    }),
    Components({
      resolvers: [
        // arco组件的按需加载
        ArcoResolver({
          sideEffect: true
        })
      ],
      // 自动加载组件的目录配置,默认的为 'src/components'
      dirs: ["src/components"],
      // 组件的有效文件扩展名
      extensions: ["vue"],
      // 配置文件生成位置
      dts: "src/components.d.ts"
    }),
    viteMockServe({
      mockPath: "./src/mock/", // 目录位置
      logger: true, //  是否在控制台显示请求日志
      supportTs: true, // 是否读取ts文件模块
      localEnabled: env.VITE_APP_OPEN_MOCK === "true", // 设置是否启用本地mock文件
      prodEnabled: env.VITE_APP_OPEN_MOCK === "true", // 设置打包是否启用mock功能
      // 这样可以控制关闭mock的时候不让mock打包到最终代码内
      injectCode: `
          import { setupProdMockServer } from '../src/mock/index';
          setupProdMockServer();
        `
    })
  ];
};
