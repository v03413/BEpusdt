import { Directive } from "vue";

const custom: Directive = {
  mounted(el, binding) {
    el.__onClick__ = () => {
      let { goodsId, event } = binding.value;
      event(goodsId);
    };
    el.addEventListener("click", el.__onClick__);
  }
};

export default custom;
