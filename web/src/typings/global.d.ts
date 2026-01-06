/* ts命名空间 */
/* 路由-Menu */
declare namespace Menu {
  interface MenuOptions {
    id: string;
    parentId: string;
    fullPath?: stirng;
    path: string;
    name: string;
    redirect?: string;
    component?: string | (() => Promise<unknown>);
    meta: MetaType;
    children?: MenuOptions[] | null | any;
  }
  interface MetaType {
    title: string;
    hide: boolean;
    disable: boolean;
    keepAlive: boolean;
    affix: boolean;
    link?: string;
    iframe?: boolean;
    roles: Array<string>;
    icon?: string;
    svgIcon?: string;
    sort?: number;
    type?: number;
  }
}
/* tabs菜单 */
declare namespace TabsMenu {
  interface list {
    list: Array<{ type: number; icon: string; title: string }>;
  }
  interface menu {
    type: number;
    icon: string;
    title: string;
  }
}

/* viteEnv */
interface ViteEnv {
  VITE_GLOB_APP_TITLE: string;
  VITE_IMG_BASE_URL: string;
  VITE_APP_OPEN_MOCK: boolean | string;
  VITE_USER_NODE_ENV: "development" | "production" | "test";
  VITE_PUBLIC_PATH: string;
  VITE_APP_BASE_URL: string;
}
