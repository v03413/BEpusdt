/**
 * 菜单处理hooks，内置多种菜单处理场景
 * @returns 菜单处理方法
 */
export const useMenuMethod = () => {
  /**
   * 多级菜单的显示隐藏
   * @param {Menu.MenuOptions} item 菜单的item
   * @returns type:1为目录，并且不隐藏，返回true，否则返回false
   */
  const menuShow = (item: Menu.MenuOptions) => {
    if (item.meta.type == 1 && !item.meta.hide) return true;
    return false;
  };
  /**
   * 单级菜单的显示隐藏
   * @param {Menu.MenuOptions} item 菜单的item
   * @returns type:2为菜单，并且单级菜单不隐藏，返回true，否则返回false
   */
  const aMenuShow = (item: Menu.MenuOptions) => {
    if (item.meta.type == 2 && !item.meta.hide) return true;
    return false;
  };

  return {
    menuShow,
    aMenuShow
  };
};
