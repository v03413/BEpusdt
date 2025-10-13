// 用于全局使用globalProperties
import { getCurrentInstance, ComponentInternalInstance } from "vue";
export default () => (getCurrentInstance() as ComponentInternalInstance).appContext.config.globalProperties;
