import { Directive } from "vue";
import { useUserInfoStore } from "@/store/modules/user-info";

/**
 * 检测指令绑定值是否为空
 * @param value 指令绑定值
 * @returns {Array<string>} 指令绑定值数组
 */
const bindingValueEmpty = (value: unknown): Array<string> => {
  // 处理未定义或空值情况
  if (!value) throw new Error("v-hasPerm 指令需要配置权限标识");

  // 标准化为数组格式
  // 如果 value 是一个数组，则直接返回该数组
  // 如果 value 不是数组，则将其转换为包含单个字符串元素的数组
  return Array.isArray(value) ? (value as string[]) : [String(value)];
};

/**
 * 检查自定义指令权限
 * @param {HTMLElement} el dom元素
 * @param {unknown} bindingValue 指令绑定值
 */
const checkPermissions = (el: HTMLElement, bindingValue: unknown) => {
  try {
    // 检测自定义指令值并转化为数组格式
    const requiredPermissions = bindingValueEmpty(bindingValue);

    // 超级管理员标识
    const all_permission = "*:*:*";

    // 获取用户权限标识-Array[string]
    let { permissions } = useUserInfoStore().account;

    // 如果是超级管理员则放行
    if (permissions.includes(all_permission)) return;

    // 是否有权限
    const hasPermissions = requiredPermissions.some((perm: string) => permissions.includes(perm));

    // 无权限、父节点存在时，删除当前节点
    if (!hasPermissions && el.parentNode) el.parentNode.removeChild(el);
  } catch (error) {
    console.error(`权限指令错误: ${error}`);
    // 删除当前节点
    if (el.parentNode) el.parentNode.removeChild(el);
  }
};

const hasPerm: Directive = {
  mounted: (el, binding) => checkPermissions(el, binding.value),
  updated: (el, binding) => checkPermissions(el, binding.value)
};

export default hasPerm;
