/**
 * https://cn.vitejs.dev/config/dep-optimization-options#optimizedeps-exclude
 * 在 Vite 的配置文件 vite.config.ts 中，optimizeDeps 选项用于配置依赖优化，特别是预构建步骤，以提高开发和构建性能。
 * 这种情况尤其需要注意，当你禁止浏览器缓存时就必须将对应模块加载到依赖预构建，否则跳转时会频繁出现页面刷新的问题，因为它无法在浏览器缓存也没有在本地node_modules/.vite内缓存。
 * 如果使用的第三方库是全局引入，也就是src/main.ts文件里，那么它就会自动被vite缓存到node_modules/.vite，此时就不用将该第三方库做依赖预构建(添加到include内)。
 */
const include = [
  "vue-codemirror6",
  "qrcode",
  "jsbarcode",
  "vuedraggable",
  "@wangeditor/editor",
  "fingerprintjs2",
  "xgplayer",
  "print-js",
  "@codemirror/theme-one-dark",
  "@codemirror/lang-json",
  "@codemirror/lang-javascript",
  "@codemirror/lang-vue",
  "pinyin-pro"
];

/**
 * 在预构建中强制排除的依赖项
 * 如果你有某些依赖包不希望被预构建，可以使用 exclude 属性来排除它们。
 */
// const exclude = [];

export { include };
