# 收银台模板开发指南

本文档面向开发者，介绍收银台模板的**目录结构与开发规范**，帮助你完成两类任务：

- **修改默认模板**：在官方内置模板的基础上调整样式或逻辑
- **新增自定义模板**：从零创建一套独立的收银台主题

---

## 模板目录结构

每套模板放在 `static/checkout/` 下的一个子目录中，目录名即模板的唯一标识。结构如下：

```
<模板目录名>/
├── checkout.json       # 【必须】模板元数据
├── views/
│   └── checkout.html   # 【必须】主模板 HTML
└── assets/             # 【可选】静态资源
    ├── css/
    ├── js/
    ├── img/
    └── locales/
```

BEpusdt 启动时会扫描该目录，自动注册所有合法模板，并在后台「收银台模板」下拉框中展示。

---

## checkout.json

模板的元数据配置文件，用于在后台展示模板信息：

```json
{
  "name": "模板显示名称",
  "author": "作者或团队",
  "desc": "一句话描述",
  "link": "模板主页或仓库地址"
}
```

`name` 为必填项，其余可选。

---

## views/checkout.html

标准的 Go `html/template` 文件。服务端只注入一个变量：

```
{{ .trade_id }}   // 当前订单的交易 ID
```

**所有订单数据均通过前端 AJAX 请求获取**，不在服务端渲染。可以用任意前端框架或原生 JS 实现交互逻辑，参考官方模板
`static/checkout/official/assets/js/checkout.js` 的写法。

---

## assets 静态资源

`assets/` 下的文件会被挂载到 `/checkout/<模板目录名>/assets/` 路由，在 HTML 中**必须使用绝对路径**引用：

```html

<link rel="stylesheet" href="/checkout/my-theme/assets/css/style.css">
<script src="/checkout/my-theme/assets/js/main.js"></script>
```

第三方 JS 库建议本地化放入 `assets/js/`，避免依赖外部 CDN。

---

## 任务一：修改默认模板

默认模板 `official` 的源文件位于 `static/checkout/official/`，直接在该目录下修改即可。

**1. 克隆源码**

```bash
git clone https://github.com/v03413/BEpusdt.git
cd BEpusdt
```

**2. 按需修改**

- 品牌色 → `static/checkout/official/assets/css/checkout.css` 中的 CSS 变量
- Logo / 图标 → 替换 `static/checkout/official/assets/img/` 下的图片并更新引用
- 页面结构 → 编辑 `static/checkout/official/views/checkout.html`
- 交互逻辑 → 编辑 `static/checkout/official/assets/js/checkout.js`
- 多语言文案 → 编辑 `static/checkout/official/assets/locales/zh.json`、`en.json`

**3. 重新编译**（见[部署与验证](#部署与验证)）

---

## 任务二：新增自定义模板

**1. 创建目录**

```bash
mkdir -p my-theme/views my-theme/assets/{css,js,img}
```

**2. 编写 checkout.json**

```json
{
  "name": "My Theme",
  "author": "yourname",
  "desc": "简洁风格收银台",
  "link": ""
}
```

**3. 编写 views/checkout.html**

从 `{{ .trade_id }}` 获取交易 ID，通过 AJAX 调用后端接口拉取订单数据，轮询检测支付状态，完成后跳转。可直接参考官方模板的
`checkout.js` 逻辑复用。

静态资源路径前缀为 `/checkout/my-theme/assets/`。

**4. 放置静态资源**

将 CSS、JS、图片等放入 `assets/` 对应子目录。

---

## 部署与验证

`static/checkout/` 目录通过 Go `//go:embed` 在**编译时**打包进程序二进制。这意味着：

- 新增或修改模板后，必须**重新编译**才能生效
- 运行时无法热加载，不存在"上传到服务器某个路径"这种方式

### 步骤

**1. 将模板目录放入源码**

把你的模板目录放到 `static/checkout/` 下：

```
static/checkout/
├── official/        # 内置官方模板
└── my-theme/        # 你的模板
    ├── checkout.json
    ├── views/
    │   └── checkout.html
    └── assets/
```

**2. 重新编译**

```bash
go build -o bepusdt ./main
```

或使用项目已有的构建脚本/Dockerfile。

**3. 验证注册结果**

启动程序后查看日志，出现以下内容说明模板注册成功：

```
前台收银模板注册成功： my-theme
```

**4. 在后台切换模板**

进入 **系统管理 → 基本设置 → API 设置**，在「收银台模板」中选择对应模板名称并保存。

---

## 注意事项

- `checkout.json` 缺失或 JSON 格式有误会导致程序启动失败
- `views/checkout.html` 缺失时该模板会被跳过，不影响其他模板加载
- 目录名会出现在 URL 中，建议用小写字母与连字符，如 `my-theme`
