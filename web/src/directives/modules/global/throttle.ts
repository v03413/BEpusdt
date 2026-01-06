import { Directive } from "vue";

const throttle: Directive = {
  mounted(el, binding) {
    if (typeof binding.value != "function") {
      return console.warn("The throttling instruction must be bound to a function");
    }
    let countDown = 0;
    let timer: any = null;
    el.__onClick__ = () => {
      if (timer) return;
      if (countDown == 0) {
        binding.value();
        countDown = 1;
      }
      timer = setInterval(() => {
        countDown--;
        if (countDown <= 0) {
          clearInterval(timer);
          timer = null;
        }
      }, 1000);
    };
    el.addEventListener("click", el.__onClick__);
  },
  beforeUnmount(el) {
    el.removeEventListener("click", el.__onClick__);
  }
};

export default throttle;
