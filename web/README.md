# SnowAdmin

<div>
	<p>
    <a href="https://github.com/WangFan-io/SnowAdmin" target="_blank">
		  <img src="https://img.shields.io/badge/License-MIT-%2343aefc" alt="License">
		</a>
		<a href="https://github.com/WangFan-io/SnowAdmin" target="_blank">
		  <img src="https://img.shields.io/badge/vue_.js-vue3_.x-%2300b42a" alt="Vue">
		</a>
    <a href="https://github.com/WangFan-io/SnowAdmin" target="_blank">
		  <img src="https://img.shields.io/badge/Vite-6x-%2311B337" alt="Vite">
		</a>
     <a href="https://github.com/WangFan-io/SnowAdmin" target="_blank">
		 <img alt="arco.design" src="https://img.shields.io/badge/arco.design-2.57.0-%23%2028%2C175%2C%2063">
		</a>
	</p>
</div>

### 介绍 📖

SnowAdmin 是一个清新优雅、高颜值且功能强大的后台管理模板，基于最新的前端技术栈，包括 Vue3, Vite7, TypeScript, Pinia。它内置了丰富的主题配置和组件，代码规范严谨，实现了自动化的文件路由系统。SnowAdmin 为您提供了一站式的后台管理解决方案，无需额外配置，开箱即用。同样是一个快速学习前沿技术的最佳实践。

### 版本 📦

**预览及文档**

项目预览地址：[SnowAdmin](http://115.190.79.132/#/login)

项目文档地址：[SnowAdmin-开发指南](http://115.190.79.132:81/)

**代码仓库**

| 平台   | 仓库地址                                               |
| ------ | ------------------------------------------------------ |
| GitHub | [Github 仓库](https://github.com/WangFan-io/SnowAdmin) |
| Gitee  | [Gitee 仓库](https://gitee.com/wang_fan_w/SnowAdmin)   |

> 完整版分支：main
>
> 极简版分支：main-thin
>
> ⚠️主要的维护在github上，提issue/pr请认准github主战场
>
> 大家也可以进wx群交流，二维码在介绍末尾

### 先行预览 🖼️

![登录页](https://ooo.0x0.ooo/2025/04/11/O0lhbG.png)

![首页](https://ooo.0x0.ooo/2025/04/11/O0lnVI.png)

![文件库](https://ooo.0x0.ooo/2025/04/11/O0l7qD.png)

![国际化](https://ooo.0x0.ooo/2025/04/11/O0lNAF.png)

![横向布局](https://ooo.0x0.ooo/2025/04/11/O0leQ6.png)

![混合布局](https://ooo.0x0.ooo/2025/04/11/O0lrCP.png)

![主题配置](https://ooo.0x0.ooo/2025/04/11/O0lKnb.png)

![黑暗模式-1](https://ooo.0x0.ooo/2025/04/11/O0lZ4l.png)

![黑暗模式-2](https://ooo.0x0.ooo/2025/04/11/O0lv6g.png)

![黑暗模式-3](https://ooo.0x0.ooo/2025/04/11/O0lgkB.png)

### 项目功能 🔨

- 清晰的项目架构：结构清晰，优雅易懂，代码全注释
- 前沿技术应用：采用 Vue3, Vite7, TypeScript等最新流行的技术栈
- 内置RBAC角色权限模型，实现路由、角色、按钮权限访问控制
- 使用 Pinia 替代 Vuex，轻量、简单、易用，集成 Pinia 持久化插件
- 支持 Arco-Design 组件大小切换、多主题布局、暗黑模式、i18n 国际化
- 使用 VueRouter 配置动态路由权限拦截、路由懒加载，支持页面按钮权限控制
- 使用 KeepAlive 对页面进行缓存，支持多级嵌套路由缓存
- 使用 Prettier 统一格式化代码，集成 ESLint、Stylelint 代码校验规范
- 使用 husky、lint-staged、commitlint等插件，规范提交信息
- 丰富的页面组件：内置多样页面和组件，包括 401、404、500 页面，以及布局组件、标签组件、主题配置组件等

### 环境准备 🏝️

确保你的环境满足以下要求（重要）：

- **git**: 你需要 git 来克隆和管理项目版本。
- **NodeJS**: >=18.12.0，推荐 20.12.0 或更高。
- **pnpm**: >= 8.7.0，推荐最新版本。

### 安装使用步骤 📔

### 从 GitHub 获取代码 🔗

```bash
# 克隆代码
git clone https://github.com/WangFan-io/SnowAdmin.git
```

### 从 Gitee 获取代码 🔗

```bash
# 克隆代码
git clone https://gitee.com/wang_fan_w/SnowAdmin.git
```

### 安装依赖 📌

安装项目依赖

```bash
pnpm install
```

### 插件配置 🛠️

安装 Vue - Official，禁用 Vetur

- [Vue - Official](https://marketplace.visualstudio.com/items?itemName=Vue.volar) - Vue 服务插件

### 项目启动 🚀

```json
{
  // 开发环境启动
  "dev": "vite",
  // 开发环境构建
  "build:dev": "vue-tsc && vite build --mode development",
  // 生产环境构建
  "build:prod": "vue-tsc && vite build --mode production",
  // 测试环境构建
  "build:test": "vue-tsc && vite build --mode test",
  // 预览环境启动
  "preview": "vite preview"
}
```

### 文件资源目录 📚

```text
SnowAdmin
├─ .github                  # GitHub 配置目录
│  └─ ISSUE_TEMPLATE        # Issue 模板
├─ .husky                   # husky 配置文件
├─ build                    # vite 项目配置目录
├─ node_modules             # 依赖包（自动生成）
├─ public                   # 静态资源文件（不会被打包）
├─ src
│  ├─ api                   # API 接口管理
│  │  └─ modules            # 接口模块
│  ├─ assets                # 静态资源文件
│  │  ├─ fonts              # 字体文件
│  │  ├─ img                # 图片资源
│  │  ├─ logo               # Logo 图标
│  │  └─ svgs               # SVG 图标
│  ├─ components            # 全局组件（code-view/qrcode/svg-icon等）
│  ├─ config                # 全局配置项
│  ├─ directives            # 全局自定义指令（防抖/节流/权限指令）
│  ├─ globals               # 全局函数
│  ├─ hooks                 # 常用 Hooks 封装
│  ├─ lang                  # 语言国际化 i18n
│  ├─ layout                # 框架布局模块
│  │  ├─ components         # 布局组件（Aside/Header/Footer/Menu/Tabs等）
│  │  └─ layout-*           # 布局类型（defaults/head/mixing）
│  ├─ mock                  # 本地数据 mock
│  ├─ router                # 路由管理
│  ├─ store                 # pinia store
│  ├─ style                 # 全局样式文件
│  │  ├─ media              # 响应式样式
│  │  ├─ model              # 动画/滚动条/全局样式
│  │  └─ var                # 主题变量
│  ├─ typings               # 全局 ts 声明
│  ├─ utils                 # 常用工具库
│  ├─ views                 # 项目所有页面
│  │  ├─ about              # 关于页面
│  │  ├─ component          # 组件展示页
│  │  ├─ directive          # 指令展示页
│  │  ├─ error              # 错误页面（401/404/500）
│  │  ├─ file               # 文件管理
│  │  ├─ form               # 表单页
│  │  ├─ functions          # 工具函数演示
│  │  ├─ home               # 首页
│  │  ├─ i18n               # 国际化演示
│  │  ├─ link               # 链接演示
│  │  ├─ login              # 登录页
│  │  ├─ monitor            # 系统监控
│  │  ├─ multilevel         # 多级菜单演示
│  │  ├─ permission         # 权限演示
│  │  ├─ personal           # 个人中心
│  │  ├─ system             # 系统管理
│  │  ├─ table              # 表格页
│  │  └─ thin-preview       # 简洁预览页
│  ├─ App.vue               # 项目主组件
│  ├─ auto-import.d.ts      # 自动导入声明
│  ├─ components.d.ts       # 组件声明
│  ├─ main.ts               # 项目入口文件
│  ├─ style.css             # 全局样式入口
│  └─ vite-env.d.ts         # 指定 ts 识别 vue
├─ .editorconfig            # 代码编辑器配置文件
├─ .env                     # 通用环境配置
├─ .env.development         # 开发环境配置
├─ .env.production          # 生产环境配置
├─ .env.test                # 测试环境配置
├─ .eslintrc-auto-import.json # ESLint 自动导入配置
├─ .gitignore               # 忽略 git 提交
├─ .prettierignore          # 忽略 Prettier 格式化
├─ .prettierrc.cjs          # Prettier 格式化配置
├─ .stylelintignore         # 忽略 Stylelint 校验
├─ .stylelintrc.cjs         # Stylelint 格式化配置
├─ commitlint.config.cjs    # git 提交规范配置
├─ eslint.config.js         # Eslint 校验配置文件
├─ index.html               # 入口 html
├─ LICENSE                  # 开源协议
├─ lint-staged.config.cjs   # lint-staged 配置文件
├─ package.json             # 依赖包管理
├─ pnpm-lock.yaml           # 依赖包版本锁
├─ README.md                # README 介绍
└─ vite.config.ts           # vite 全局配置文件
```

### 浏览器支持 🌎

- 本地开发推荐使用 Chrome 最新版浏览器 [Download](https://www.google.com/intl/zh-CN/chrome/)。
- 生产环境支持现代浏览器，不再支持 IE 浏览器，更多浏览器可以查看 [Can I Use Es Module](https://caniuse.com/?search=ESModule)。

| ![IE](https://i.imgtg.com/2023/04/11/8z7ot.png) | ![Edge](https://i.imgtg.com/2023/04/11/8zr3p.png) | ![Firefox](https://i.imgtg.com/2023/04/11/8zKiU.png) | ![Chrome](https://i.imgtg.com/2023/04/11/8zNrx.png) | ![Safari](https://i.imgtg.com/2023/04/11/8zeGj.png) |
| :---------------------------------------------: | :-----------------------------------------------: | :--------------------------------------------------: | :-------------------------------------------------: | :-------------------------------------------------: |
|                   not support                   |                  last 2 versions                  |                   last 2 versions                    |                   last 2 versions                   |                   last 2 versions                   |

### 鸣谢列表 🧩

- [vue：https://v3.cn.vuejs.org/](https://v3.cn.vuejs.org/)
- [axios：https://github.com/axios/axios](https://github.com/axios/axios)
- [nprogress：https://github.com/rstacruz/nprogress](https://github.com/rstacruz/nprogress)
- [sortablejs：https://github.com/SortableJS/Sortabl](https://github.com/SortableJS/Sortable)
- [sass：https://github.com/sass/sass](https://github.com/sass/sass)
- [typescript：https://github.com/microsoft/TypeScript](https://github.com/microsoft/TypeScript)
- [vite：https://github.com/vitejs/vite](https://github.com/vitejs/vite)
- [wangeditor：https://github.com/wangeditor-team/wangEditor](https://github.com/wangeditor-team/wangEditor)
- [qrcodejs：https://github.com/davidshimjs/qrcodejs](https://github.com/davidshimjs/qrcodejs)
- [print-js：https://github.com/crabbly/Print.js](https://github.com/crabbly/Print.js)
- [eslint：https://github.com/eslint/eslint](https://github.com/eslint/eslint)
- [prettier：https://github.com/prettier/prettier](https://github.com/prettier/prettier)
- [arco.design：https://arco.design/](https://arco.design/)
- [pinyin-pro：https://github.com/zh-lx/pinyin-pro](https://github.com/zh-lx/pinyin-pro)
- [qrcode：https://www.npmjs.com/package/qrcode](https://www.npmjs.com/package/qrcode)
- [jsbarcode：https://github.com/lindell/JsBarcode/wiki/Options](https://github.com/lindell/JsBarcode/wiki/Options)
- [fingerprintjs2：https://www.npmjs.com/package/fingerprintjs2](https://www.npmjs.com/package/fingerprintjs2)

### 特别鸣谢 🎉

- [vue-next-admin：http://vuenextadmin.ccfast.cc/](http://vuenextadmin.ccfast.cc/)
- [Gi-admin：https://lin-97.github.io/gi-demo/](https://lin-97.github.io/gi-demo/#/home)
- [Geeker-Admin：https://admin.spicyboy.cn/](https://admin.spicyboy.cn/#/home/index)

感谢以上框架为本项目提供了帮助，项目中的很多设计思想值得参考和学习，再次感谢。

### 学习 & 交流 🏫

|                                              SnowAdmin 交流群                                               |
| :---------------------------------------------------------------------------------------------------------: |
| ![微信图片_20250726141221](https://github.com/user-attachments/assets/8c1ebf3c-4630-487e-9a54-810fd4d6972c) |

使用过程中有各种问题也可以进群交流，欢迎大家积极参与

### 支持作者 🌟

如果觉得框架不错，或者已经在使用了，希望你可以去 [GitHub](https://github.com/WangFan-io/SnowAdmin) ⭐ Star，这将是对我极大的鼓励与支持。

### 免责条款 📫

您充分了解并同意，您必须为自己使用本服务及注册帐号下的一切行为负责，包括您所发表的任何内容以及由此产生的任何后果。您应对本服务中的内容自行加以判断，并自行承担因使用内容而引起的所有风险。

`SnowAdmin` 对网站上所显示的信息或资料的准确性、内容、完整性、合法性、可靠性、可操作性或可用性不承担任何责任。

`SnowAdmin` 二次开发使用者因为违法而触犯中华人民共和国法律的，一切后果自己负责，`SnowAdmin` 作者不承担任何责任。

本声明未涉及的问题参见国家有关法律法规，当本声明与国家法律法规冲突时，以国家法律法规为准。

### 学习 & 商用 🌈

`SnowAdmin` 是免费和开源的，可免费用于学习、商业使用 。

### 捐赠 💖

|                        微信                        |                        支付宝                        |
| :------------------------------------------------: | :--------------------------------------------------: |
| ![微信](https://ooo.0x0.ooo/2025/07/20/OYtFn1.jpg) | ![支付宝](https://ooo.0x0.ooo/2025/07/20/OYth4I.jpg) |
