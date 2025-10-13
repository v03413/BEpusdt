/**
 * px转响应式px
 * 示例：$px2px(20, 1920)，在1920px的设计稿上元素大小为20px，若在其它窗口大小下，会根据窗口大小自动缩放
 * @param {number} px 设计稿元素大小
 * @param {number} viewport 设计稿窗口大小
 * @returns {number} 响应式px
 */
export const px2px = (px: number, viewport: number): number => {
  const viewportWidth: number = (window.innerWidth || document.documentElement.clientWidth) as number;
  let originalWidth = px;
  let originalViewport = viewport;
  // 计算新视口下的宽度
  let newWidth = (originalWidth / originalViewport) * viewportWidth;
  return parseInt(String(newWidth));
};
