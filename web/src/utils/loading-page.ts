/**
 * 全局加载 loading-page
 * @method start 创建 loading
 * @method done 移除 loading
 */
export const loadingPage = {
  // 开始渲染loading
  start: () => {
    // 获取顶层body
    // 将新创建的dc-loader元素（div）插入到body元素的子元素列表中的指定位置（在指定元素之前）
    // 插入的位置是作为body元素的第一个子元素，即页面的最顶部位置
    const bodyDom: Element = document.body;
    const div = document.createElement("div");
    div.className = "loading-page";
    const loader = document.createElement("div");
    loader.className = "dc-loader";
    div.appendChild(loader);
    bodyDom.insertBefore(div, bodyDom.firstChild);
  },
  // 结束渲染loading
  done: (time: number = 0) => {
    setTimeout(() => {
      // 找到第一个匹配对象
      // 找到loading-page的父节点，移除loading-page
      const dom = document.querySelector(".loading-page");
      dom?.parentNode?.removeChild(dom);
    }, time);
  }
};
