/**
 * 路由path路径与文件夹名称相同，找文件可以浏览器地址快速查找，方便定位文件
 *
 * component
 * 模块加载路径，除了layout外，全部默认路径为src/views文件夹内的文件
 * 不要以"/"开头，结尾不要带文件扩展名，如:".vue"
 * 举例：views文件夹下的home文件夹下的home.vue
 * component: "home/home"
 *
 * id的具体配置规则：
 * 父节点: id: '01'
 *  |-- 子节点1: '0101'
 *  |    |-- 子节点1-1: '010101'
 *  |    |-- 子节点1-2: '010102'
 *  |-- 子节点2: '0102'
 *  |    |-- 子节点2-1: '010201'
 *  |    |-- 子节点2-2: '010202'
 * 当然，你们也可以自行设计id结构，只要保证id的唯一性即可，最后接口返回的是一个树形结构
 * 注意：节点id不能与parentId相同，否则在遍历路由树的时候会报错
 * 在这种层级编码设计模式下，​不可能自然产生节点ID与parentId相同的情况，但还是需要注意手动录入的情况
 *
 *
 * 路由meta对象参数，我们通常将属性放到meta对象中
 * meta: {
 *   title:     国际化key，无对应Key则直接展示
 *   isFull:    是否全屏显示，开启条件： true \ false
 *   hide:      是否隐藏此路由，不会显示在菜单树，可以访问
 *   disable:   是否停用此路由-不显示且不可访问
 *   keepAlive: 是否缓存组件
 *   affix:     固定在Tabs栏上-开启后Tabs栏无关闭按钮
 *   link:      是否是超链接菜单，开启外链条件：1、 link：链接地址不为空  2、iframe: false
 *   iframe:    是否内嵌外链窗口，开启条件：1、iframe：true  2、link：链接地址不为空
 *   roles:     路由权限标识，取角色权限，超级管理员：admin；普通角色：common
 *   svgIcon:   svg菜单图标，优先级高于icon，取src/assets/svgs内的svg文件
 *   icon:      普通icon菜单图标，默认取arco.design图标
 *   sort:      排序字段
 *   type:      1目录 2菜单 3按钮
 * }
 */

/**
 * 路由的层级设置
 * layout为框架布局，顶层路由
 * layout.children下的路由为systemMenu菜单，例如：首页、系统设置、权限管理等
 */
export const systemMenu = [
  // layout-二级路由(主要渲染页面)
  {
    id: "01", // 路由id
    parentId: "0", // 顶层路由 parentId: 0
    path: "/home", // 路由path
    name: "home", // 路由name
    component: "home/home", // 路由跳转的文件路径，默认在src/views内，这里就是src/views/home/home.vue
    meta: {
      title: "home", // 国际化key，无对应Key则直接展示
      hide: false, // 是否隐藏此路由-不显示但可访问
      disable: false, // 是否停用此路由-不显示且不可访问
      keepAlive: false, // 是否缓存组件
      affix: true, // 固定在Tabs栏上-开启后Tabs栏无关闭按钮
      link: "", // 是否外链
      iframe: false, // 是否内嵌外链窗口
      isFull: false, // 是否全屏显示
      roles: ["admin", "common"], // 路由角色权限
      svgIcon: "home", // svg菜单图标，优先级高于icon，取src/assets/svgs内的svg文件
      icon: "", // 普通icon菜单图标，默认取arco.design图标
      sort: 1, // 排序字段
      type: 2 // type 1目录 2菜单 3按钮
    },
    children: null // 存储子节点字段，默认为null
  },
  {
    id: "02",
    parentId: "0",
    path: "/file",
    name: "file",
    redirect: "/file/document-library",
    meta: {
      title: "file",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "folder-menu",
      sort: 2,
      type: 1
    },
    children: null
  },
  {
    id: "0201",
    parentId: "02",
    path: "/file/document-library",
    name: "document-library",
    component: "file/document-library/document-library",
    meta: {
      title: "document-library",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      isFull: false,
      link: "",
      iframe: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "03",
    parentId: "0",
    path: "/table",
    name: "table",
    redirect: "/table/common-table",
    meta: {
      title: "table",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "table",
      sort: 3,
      type: 1
    },
    children: null
  },
  {
    id: "0301",
    parentId: "03",
    path: "/table/common-table",
    name: "common-table",
    component: "table/common-table/common-table",
    meta: {
      title: "common-table",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "0302",
    parentId: "03",
    path: "/table/custom-table",
    name: "custom-table",
    component: "table/custom-table/custom-table",
    meta: {
      title: "custom-table",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "04",
    parentId: "0",
    path: "/form",
    name: "form",
    redirect: "/form/common-form",
    meta: {
      title: "form",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "form",
      sort: 4,
      type: 1
    },
    children: null
  },
  {
    id: "0401",
    parentId: "04",
    path: "/form/common-form",
    name: "common-form",
    component: "form/common-form/common-form",
    meta: {
      title: "common-form",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "0402",
    parentId: "04",
    path: "/form/step-form",
    name: "step-form",
    component: "form/step-form/step-form",
    meta: {
      title: "step-form",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "05",
    parentId: "0",
    path: "/multilevel",
    name: "multilevel",
    redirect: "/multilevel/second-1",
    meta: {
      title: "multilevel",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "switch",
      sort: 5,
      type: 1
    },
    children: null
  },
  {
    id: "0501",
    parentId: "05",
    path: "/multilevel/second-1",
    name: "second-1",
    component: "multilevel/second/second-1",
    meta: {
      title: "second-1",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "0502",
    parentId: "05",
    path: "/multilevel/second-2",
    name: "second-2",
    redirect: "/multilevel/third-1",
    meta: {
      title: "second-2",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 1
    },
    children: null
  },
  {
    id: "050201",
    parentId: "0502",
    path: "/multilevel/third-2",
    name: "third-2",
    component: "multilevel/third/third-2",
    meta: {
      title: "third-2",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "050202",
    parentId: "0502",
    path: "/multilevel/third-1",
    name: "third-1",
    component: "multilevel/third/third-1",
    meta: {
      title: "third-1",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "050203",
    parentId: "0502",
    path: "/multilevel/third-3",
    name: "third-3",
    component: "multilevel/third/third-3",
    meta: {
      title: "third-3",
      link: "",
      hide: false,
      keepAlive: true,
      affix: false,
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "06",
    parentId: "0",
    path: "/component",
    name: "component",
    redirect: "/component/player",
    meta: {
      title: "component",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "classify",
      sort: 6,
      type: 1
    },
    children: null
  },
  {
    id: "0601",
    parentId: "06",
    path: "/component/player",
    name: "player",
    component: "component/player/player",
    meta: {
      title: "player",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "0602",
    parentId: "06",
    path: "/component/print",
    name: "print",
    component: "component/print/print",
    meta: {
      title: "print",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "0603",
    parentId: "06",
    path: "/component/draggable",
    name: "draggable",
    component: "component/draggable/draggable",
    meta: {
      title: "draggable",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "0604",
    parentId: "06",
    path: "/component/editor",
    name: "editor",
    component: "component/editor/editor",
    meta: {
      title: "editor",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 4,
      type: 2
    },
    children: null
  },
  {
    id: "0605",
    parentId: "06",
    path: "/component/newbie",
    name: "newbie",
    component: "component/newbie/newbie",
    meta: {
      title: "newbie",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 5,
      type: 2
    },
    children: null
  },
  {
    id: "0606",
    parentId: "06",
    path: "/component/icon-selector",
    name: "icon-selector",
    component: "component/icon-selector/icon-selector",
    meta: {
      title: "icon-selector",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 6,
      type: 2
    },
    children: null
  },
  {
    id: "0607",
    parentId: "06",
    path: "/component/user-center",
    name: "user-center",
    component: "component/user-center/user-center",
    meta: {
      title: "user-center",
      hide: true,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 7,
      type: 2
    },
    children: null
  },
  {
    id: "0608",
    parentId: "06",
    path: "/component/fingerprintjs2",
    name: "fingerprintjs2",
    component: "component/fingerprintjs2/fingerprintjs2",
    meta: {
      title: "fingerprintjs2",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 8,
      type: 2
    },
    children: null
  },
  {
    id: "0609",
    parentId: "06",
    path: "/component/barcode",
    name: "barcode",
    component: "component/barcode/barcode",
    meta: {
      title: "barcode",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 9,
      type: 2
    },
    children: null
  },
  {
    id: "0610",
    parentId: "06",
    path: "/component/qrcode",
    name: "qrcode",
    component: "component/qrcode/qrcode",
    meta: {
      title: "qrcode",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 10,
      type: 2
    },
    children: null
  },
  {
    id: "0611",
    parentId: "06",
    path: "/component/pinyin",
    name: "pinyin",
    component: "component/pinyin/pinyin",
    meta: {
      title: "pinyin",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 11,
      type: 2
    },
    children: null
  },
  {
    id: "0612",
    parentId: "06",
    path: "/component/recorder",
    name: "recorder",
    component: "component/recorder/recorder",
    meta: {
      title: "recorder",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 12,
      type: 2
    },
    children: null
  },
  {
    id: "0613",
    parentId: "06",
    path: "/component/virtual-list",
    name: "virtual-list",
    component: "component/virtual-list/index",
    meta: {
      title: "virtual-list",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 13,
      type: 2
    },
    children: null
  },
  {
    id: "0614",
    parentId: "06",
    path: "/component/common-layouts",
    name: "common-layouts",
    component: "component/common-layouts/index",
    meta: {
      title: "common-layouts",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 14,
      type: 2
    },
    children: null
  },
  {
    id: "07",
    parentId: "0",
    path: "/directive",
    name: "directive",
    redirect: "/directive/anti-shake",
    meta: {
      title: "directive",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "directives",
      sort: 7,
      type: 1
    },
    children: null
  },
  {
    id: "0701",
    parentId: "07",
    path: "/directive/anti-shake",
    name: "anti-shake",
    component: "directive/anti-shake/anti-shake",
    meta: {
      title: "anti-shake",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "0702",
    parentId: "07",
    path: "/directive/throttle",
    name: "throttle",
    component: "directive/throttle/throttle",
    meta: {
      title: "throttle",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "0703",
    parentId: "07",
    path: "/directive/test-directive",
    name: "test-directive",
    component: "directive/test-directive/test-directive",
    meta: {
      title: "test-directive",
      hide: true,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  // {
  //   id: "08",
  //   parentId: "0",
  //   path: "/personal",
  //   name: "personal",
  //   redirect: "/personal/userinfo",
  //   meta: {
  //     title: "personal",
  //     hide: false,
  //     disable: false,
  //     keepAlive: true,
  //     affix: false,
  //     link: "",
  //     iframe: false,
  //     isFull: false,
  //     roles: ["admin"],
  //     svgIcon: "user",
  //     sort: 8,
  //     type: 1
  //   },
  //   children: null
  // },
  // {
  //   id: "0801",
  //   parentId: "08",
  //   path: "/personal/userinfo",
  //   name: "userinfo",
  //   component: "personal/userinfo/userinfo",
  //   meta: {
  //     title: "userinfo",
  //     hide: false,
  //     disable: false,
  //     keepAlive: true,
  //     affix: false,
  //     link: "",
  //     iframe: false,
  //     isFull: false,
  //     roles: ["admin"],
  //     icon: "icon-menu",
  //     sort: 1,
  //     type: 2
  //   },
  //   children: null
  // },
  // {
  //   id: "0802",
  //   parentId: "08",
  //   path: "/personal/user-settings",
  //   name: "user-settings",
  //   component: "personal/user-settings/user-settings",
  //   meta: {
  //     title: "user-settings",
  //     hide: false,
  //     disable: false,
  //     keepAlive: true,
  //     affix: false,
  //     link: "",
  //     iframe: false,
  //     isFull: false,
  //     roles: ["admin"],
  //     icon: "icon-menu",
  //     sort: 2,
  //     type: 2
  //   },
  //   children: null
  // },
  {
    id: "09",
    parentId: "0",
    path: "/functions",
    name: "functions",
    redirect: "/functions/routing-operation",
    meta: {
      title: "functions",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "functions",
      sort: 9,
      type: 1
    },
    children: null
  },
  {
    id: "0901",
    parentId: "09",
    path: "/functions/routing-operation",
    name: "routing-operation",
    component: "functions/routing-operation/index",
    meta: {
      title: "routing-operation",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "090101",
    parentId: "0901",
    path: "/functions/routing-operation/common-route",
    name: "common-route",
    component: "functions/routing-operation/common-route",
    meta: {
      title: "common-route",
      hide: true,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "switch",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "090102",
    parentId: "0901",
    path: "/functions/routing-operation/dynamic-route/:name/:text",
    name: "dynamic-route",
    component: "functions/routing-operation/dynamic-route",
    meta: {
      title: "dynamic-route",
      hide: true,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "switch",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "0902",
    parentId: "09",
    path: "/functions/common-tools",
    name: "common-tools",
    component: "functions/common-tools/common-tools",
    meta: {
      title: "common-tools",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "0903",
    parentId: "09",
    path: "/functions/tree-tools",
    name: "tree-tools",
    component: "functions/tree-tools/tree-tools",
    meta: {
      title: "tree-tools",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "0904",
    parentId: "09",
    path: "/functions/file-tools",
    name: "file-tools",
    component: "functions/file-tools/file-tools",
    meta: {
      title: "file-tools",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 4,
      type: 2
    },
    children: null
  },
  {
    id: "0905",
    parentId: "09",
    path: "/functions/verify-tools",
    name: "verify-tools",
    component: "functions/verify-tools/verify-tools",
    meta: {
      title: "verify-tools",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 5,
      type: 2
    },
    children: null
  },
  {
    id: "10",
    parentId: "0",
    path: "/system",
    name: "system",
    redirect: "/system/account",
    meta: {
      title: "system",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin", "common"],
      svgIcon: "set",
      sort: 10,
      type: 1
    },
    children: null
  },
  {
    id: "1001",
    parentId: "10",
    path: "/system/account",
    name: "account",
    component: "system/account/account",
    meta: {
      title: "account",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "1002",
    parentId: "10",
    path: "/system/role",
    name: "role",
    component: "system/role/role",
    meta: {
      title: "role",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "1003",
    parentId: "10",
    path: "/system/menu",
    name: "menu",
    component: "system/menu/menu",
    meta: {
      title: "menu",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "1004",
    parentId: "10",
    path: "/system/division",
    name: "division",
    component: "system/division/division",
    meta: {
      title: "division",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 4,
      type: 2
    },
    children: null
  },
  {
    id: "1005",
    parentId: "10",
    path: "/system/dictionary",
    name: "dictionary",
    component: "system/dictionary/dictionary",
    meta: {
      title: "dictionary",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 5,
      type: 2
    },
    children: null
  },
  {
    id: "1006",
    parentId: "10",
    path: "/system/log",
    name: "log",
    component: "system/log/log",
    meta: {
      title: "log",
      hide: true,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 6,
      type: 2
    },
    children: null
  },
  {
    id: "1007",
    parentId: "10",
    path: "/system/userinfo",
    name: "userinfo",
    component: "system/userinfo/userinfo",
    meta: {
      title: "userinfo",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin", "common"],
      icon: "icon-menu",
      sort: 7,
      type: 2
    },
    children: null
  },
  {
    id: "11",
    parentId: "0",
    path: "/disable-menu",
    name: "disable-menu",
    component: "disable-menu/disable-menu",
    meta: {
      title: "disable-menu",
      hide: false,
      disable: true,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "switch",
      sort: 11,
      type: 2
    },
    children: null
  },
  {
    id: "12",
    parentId: "0",
    path: "/hide-menu",
    name: "hide-menu",
    component: "hide-menu/hide-menu",
    meta: {
      title: "hide-menu",
      hide: true,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "switch",
      sort: 12,
      type: 2
    },
    children: null
  },
  {
    id: "13",
    parentId: "0",
    path: "/permission",
    name: "permission",
    redirect: "/permission/btn-perm",
    meta: {
      title: "permission",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin", "common"],
      svgIcon: "defend",
      sort: 13,
      type: 1
    },
    children: null
  },
  {
    id: "1301",
    parentId: "13",
    path: "/permission/btn-perm",
    name: "btn-perm",
    component: "permission/btn-perm/btn-perm",
    meta: {
      title: "btn-perm",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin", "common"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "1302",
    parentId: "13",
    path: "/permission/admin-page",
    name: "admin-page",
    component: "permission/admin-page/admin-page",
    meta: {
      title: "admin-page",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "1303",
    parentId: "13",
    path: "/permission/common-page",
    name: "common-page",
    component: "permission/common-page/common-page",
    meta: {
      title: "common-page",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["common"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "14",
    parentId: "0",
    path: "/link",
    name: "link",
    redirect: "/link/internal",
    meta: {
      title: "link",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "link",
      sort: 14,
      type: 1
    },
    children: null
  },
  {
    id: "1401",
    parentId: "14",
    path: "/link/internal",
    name: "internal",
    redirect: "/link/internal/uigradients",
    meta: {
      title: "internal",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 1
    },
    children: null
  },
  {
    id: "140101",
    parentId: "1401",
    path: "/link/internal/uigradients",
    name: "uigradients",
    component: "link/internal/internal",
    meta: {
      title: "uigradients",
      hide: false,
      keepAlive: true,
      disable: false,
      affix: false,
      link: "https://uigradients.com/#HoneyDew", // 链接
      iframe: true, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "140102",
    parentId: "1401",
    path: "/link/internal/color-taking-tool",
    name: "color-taking-tool",
    component: "link/internal/internal",
    meta: {
      title: "color-taking-tool",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "https://photokit.com/colors/eyedropper/?lang=zh", // 链接
      iframe: true, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "140103",
    parentId: "1401",
    path: "/link/internal/grid-generator",
    name: "grid-generator",
    component: "link/internal/internal",
    meta: {
      title: "grid-generator",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "https://cssgrid-generator.netlify.app/", // 链接
      iframe: true, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "140104",
    parentId: "1401",
    path: "/link/internal/gaodemap",
    name: "gaodemap",
    component: "link/internal/internal",
    meta: {
      title: "amap",
      hide: false,
      keepAlive: true,
      disable: false,
      affix: false,
      link: "http://115.190.79.132:82/", // 链接
      iframe: true, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 4,
      type: 2
    },
    children: null
  },
  {
    id: "1402",
    parentId: "14",
    path: "/link/external",
    name: "external",
    redirect: "/link/external/link-vue",
    meta: {
      title: "external",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 1
    },
    children: null
  },
  {
    id: "140201",
    parentId: "1402",
    path: "/link/external/SnowAdmin-Docs",
    name: "SnowAdmin-Docs",
    component: "link/external/external",
    meta: {
      title: "SnowAdmin-Docs",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "http://115.190.79.132:81/", // 链接
      iframe: false, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 5,
      type: 2
    },
    children: null
  },
  {
    id: "140202",
    parentId: "1402",
    path: "/link/external/vue",
    name: "vue",
    component: "link/external/external",
    meta: {
      title: "vue",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "https://cn.vuejs.org/", // 链接
      iframe: false, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "140203",
    parentId: "1402",
    path: "/link/external/vite",
    name: "vite",
    component: "link/external/external",
    meta: {
      title: "vite",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "https://www.vitejs.net/", // 链接
      iframe: false, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "140204",
    parentId: "1402",
    path: "/link/external/github",
    name: "github",
    component: "link/external/external",
    meta: {
      title: "github",
      hide: false,
      disable: true,
      keepAlive: true,
      affix: false,
      link: "https://github.com/WangFan-io/SnowAdmin", // 链接
      iframe: false, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "140205",
    parentId: "1402",
    path: "/link/external/juejin",
    name: "juejin",
    component: "link/external/external",
    meta: {
      title: "juejin",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "https://juejin.cn/user/1728883023940600", // 链接
      iframe: false, // 区分是否内链 true内链 false外链
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 4,
      type: 2
    },
    children: null
  },
  {
    id: "15",
    parentId: "0",
    path: "/monitor",
    name: "monitor",
    redirect: "/monitor/onlineuser",
    meta: {
      title: "system-monitor",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "financial-statement",
      sort: 15,
      type: 1
    },
    children: null
  },
  {
    id: "1501",
    parentId: "15",
    path: "/monitor/onlineuser",
    name: "onlineuser",
    component: "monitor/onlineuser/index",
    meta: {
      title: "onlineuser",
      hide: false,
      keepAlive: true,
      disable: false,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 1,
      type: 2
    },
    children: null
  },
  {
    id: "1502",
    parentId: "15",
    path: "/monitor/crontab",
    name: "crontab",
    component: "monitor/crontab/index",
    meta: {
      title: "crontab",
      hide: false,
      keepAlive: true,
      disable: false,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 2,
      type: 2
    },
    children: null
  },
  {
    id: "1503",
    parentId: "15",
    path: "/monitor/crontab-logs",
    name: "crontab-logs",
    component: "monitor/crontab-logs/index",
    meta: {
      title: "crontab-logs",
      hide: false,
      keepAlive: true,
      disable: false,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      icon: "icon-menu",
      sort: 3,
      type: 2
    },
    children: null
  },
  {
    id: "16",
    parentId: "0",
    path: "/SnowAdmin-Thin",
    name: "thin-preview",
    component: "link/external/external",
    meta: {
      title: "thin-preview",
      hide: false,
      disable: false,
      keepAlive: false,
      affix: true,
      link: "http://115.190.79.132:83/#/login",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "propaganda",
      icon: "",
      sort: 16,
      type: 2
    },
    children: null
  },
  {
    id: "17",
    parentId: "0",
    path: "/i18n",
    name: "i18n",
    component: "i18n/show/index",
    meta: {
      title: "i18n",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin"],
      svgIcon: "earth",
      sort: 17,
      type: 2
    },
    children: null
  },
  {
    id: "18",
    parentId: "0",
    path: "/about",
    name: "about",
    component: "about/about",
    meta: {
      title: "about",
      hide: false,
      disable: false,
      keepAlive: true,
      affix: false,
      link: "",
      iframe: false,
      isFull: false,
      roles: ["admin", "common"],
      svgIcon: "about",
      sort: 18,
      type: 2
    },
    children: null
  }
];

// 权限数据
export const permissionData = [
  {
    id: "130101",
    parentId: "1301",
    meta: {
      title: "新增按钮",
      hide: false,
      disable: false,
      roles: ["admin"],
      permission: "sys:btn:add",
      sort: 1,
      type: 3
    }
  },
  {
    id: "130102",
    parentId: "1301",
    meta: {
      title: "编辑按钮",
      hide: false,
      disable: false,
      roles: ["admin"],
      permission: "sys:btn:edit",
      sort: 2,
      type: 3
    }
  },
  {
    id: "130103",
    parentId: "1301",
    meta: {
      title: "删除按钮",
      hide: false,
      disable: false,
      roles: ["admin"],
      permission: "sys:btn:delete",
      sort: 3,
      type: 3
    }
  },
  {
    id: "130201",
    parentId: "1301",
    meta: {
      title: "新增按钮-2",
      hide: false,
      disable: false,
      roles: ["admin", "common"],
      permission: "common:btn:add",
      sort: 1,
      type: 3
    }
  },
  {
    id: "130202",
    parentId: "1301",
    meta: {
      title: "编辑按钮-2",
      hide: false,
      disable: false,
      roles: ["admin", "common"],
      permission: "common:btn:edit",
      sort: 2,
      type: 3
    }
  },
  {
    id: "130203",
    parentId: "1301",
    meta: {
      title: "删除按钮-2",
      hide: false,
      disable: false,
      roles: ["admin", "common"],
      permission: "common:btn:delete",
      sort: 3,
      type: 3
    }
  }
];
