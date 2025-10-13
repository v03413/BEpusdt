<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space direction="vertical" :size="10">
        <a-typography-text id="newbie-title-1"> SnowAdmin </a-typography-text>
        <a-typography-text id="newbie-title-2"> 清晰的项目架构：结构清晰，优雅易懂。 </a-typography-text>
        <a-typography-text id="newbie-title-3"> 严格的代码规范：严格遵循开发设计规范，保证代码的规范性。 </a-typography-text>
        <a-typography-text id="newbie-title-4"> 内置国际化方案：轻松实现多语言支持。</a-typography-text>
        <a-typography-text id="newbie-title-5"> 清新丰富的主题配置：内置多样的主题配置。</a-typography-text>
      </a-space>
      <div class="margin-top">
        <a-space>
          <a-button type="primary" @click="openNewbie">打开指引</a-button>
          <a-button type="primary" @click="openSystemNewbie">系统指引</a-button>
        </a-space>
      </div>
      <div class="margin-top">
        新手指引，基于<a-link href="https://driver.employleague.cn/" target="_blank">driver.js</a-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { driver } from "driver.js";
import "driver.js/dist/driver.css";

const openNewbie = () => {
  const driverObj = driver({
    allowClose: true,
    doneBtnText: "结束",
    nextBtnText: "下一步",
    prevBtnText: "上一步",
    steps: [
      {
        element: "#newbie-title-1",
        popover: {
          title: "SnowAdmin",
          description: "清新优雅的中后台管理模板",
          side: "right"
        }
      },
      {
        element: "#newbie-title-2",
        popover: {
          title: "清晰的项目架构",
          description: "结构清晰，优雅易懂",
          side: "right"
        }
      },
      {
        element: "#newbie-title-3",
        popover: {
          title: "严格的代码规范",
          description: "严格遵循开发设计规范，保证代码的规范性",
          side: "right"
        }
      },
      {
        element: "#newbie-title-4",
        popover: {
          title: "内置国际化方案",
          description: "轻松实现多语言支持",
          side: "right"
        }
      },
      {
        element: "#newbie-title-5",
        popover: {
          title: "清新丰富的主题配置",
          description: "内置多样的主题配置",
          side: "right"
        }
      }
    ]
  });
  driverObj.drive();
};

let systemList = ref<any>([
  {
    element: "#system-collapsed",
    popover: {
      title: "折叠菜单",
      description: "菜单的展开收起",
      side: "right"
    }
  },
  {
    element: "#system-breadcrumb",
    popover: {
      title: "面包屑",
      description: "页面路径",
      side: "right"
    }
  },
  {
    element: "#system-language",
    popover: {
      title: "国际化",
      description: "切换系统语言",
      side: "left"
    }
  },
  {
    element: "#system-dark",
    popover: {
      title: "黑暗模式",
      description: "切换黑暗和白天模式",
      side: "left"
    }
  },
  {
    element: "#system-notice",
    popover: {
      title: "通知",
      description: "系统通知信息",
      side: "left"
    }
  },
  {
    element: "#system-fullscreen",
    popover: {
      title: "全屏",
      description: "系统全屏显示",
      side: "left"
    }
  },
  {
    element: "#system-settings",
    popover: {
      title: "系统设置",
      description: "修改系统的渲染配置",
      side: "left"
    }
  },
  {
    element: "#system-theme",
    popover: {
      title: "主题设置",
      description: "修改系统的主题配置",
      side: "left"
    }
  },
  {
    element: "#system-my-setting",
    popover: {
      title: "我的",
      description: "个人中心设置",
      side: "left"
    }
  },
  {
    element: "#system-tabs-setting",
    popover: {
      title: "Tabs栏功能",
      description: "Tabs栏的关闭、刷新等功能",
      side: "left"
    }
  },
  {
    element: "#system-tabs-refresh",
    popover: {
      title: "页面刷新",
      description: "刷新当前窗口页面",
      side: "left"
    }
  }
]);
const openSystemNewbie = () => {
  const driverObj = driver({
    allowClose: true,
    doneBtnText: "结束",
    nextBtnText: "下一步",
    prevBtnText: "上一步",
    steps: systemList.value
  });
  driverObj.drive();
};
</script>

<style lang="scss" scoped>
.margin-top {
  margin-top: $padding;
}
</style>
