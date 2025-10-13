/**
 * https://prazdevs.github.io/pinia-plugin-persistedstate/zh/guide/config.html
 * @description: 持久化配置
 * @param {string} key 存储的key
 * @param {Array<string>} paths 需要持久化的字段
 * @return {*}
 */
const persistedstateConfig = (key: string, paths?: Array<string>) => {
  return {
    key,
    storage: localStorage, // 存储方式，使用localStorage，也可以使用sessionStorage
    paths
  };
};

export default persistedstateConfig;
