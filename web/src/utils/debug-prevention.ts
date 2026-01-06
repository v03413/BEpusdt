/**
 * 阻止键盘事件打开控制台
 */
class KeydownControl {
  /**
   * 阻止键盘事件打开控制台
   * 主要的组合键有四种：
   * 1、ctrl+shift+i
   * 2、F12
   * 3、ctrl+shift+c
   * 4、shift+f10
   * @param e 键盘事件
   */
  private keydown = (e: KeyboardEvent) => {
    const code = e.code; // 具体按键
    const ctrl = e.ctrlKey; // Control键是否按下
    const shift = e.shiftKey; // Shift键是否按下
    // ctrl+shift+i
    const isCSI = ctrl && shift && code === "KeyI";
    // F12
    const isF12 = code === "F12";
    // ctrl+shift+c
    const isCSC = ctrl && shift && code === "KeyC";
    // shift+f10
    const isSF10 = shift && code === "F10";
    // 禁止打开控制台
    if (isF12 || isCSI || isCSC || isSF10) {
      e.preventDefault();
    }
  };
  // 监听键盘事件
  start() {
    document.addEventListener("keydown", this.keydown);
  }
  // 移除键盘事件监听
  stop() {
    document.removeEventListener("keydown", this.keydown);
  }
}

/**
 * 阻止鼠标事件打开控制台
 */
class RightMouseControl {
  /**
   * 阻止鼠标事件打开控制台
   * @param e 鼠标事件
   */
  private rightClick = (e: MouseEvent) => {
    e.preventDefault();
  };
  // 监听鼠标右键
  start() {
    // 禁用鼠标右键菜单
    document.addEventListener("contextmenu", this.rightClick);
  }
  // 移除鼠标右键监听
  stop() {
    document.removeEventListener("contextmenu", this.rightClick);
  }
}

/**
 * 使用debugger关键字阻止打开控制台
 */
class DebugProtector {
  private isActive = false; // 是否开启防调试

  start() {
    if (this.isActive) return;
    this.isActive = true; // 开启防调试
    this.asyncCheck(); // 异步检查
  }

  private asyncCheck() {
    if (!this.isActive) return;
    // 直接使用debugger关键字，避免eval的安全风险
    debugger;

    // 异步调度避免栈溢出
    setTimeout(() => {
      this.asyncCheck();
    }, 200); // 保持0.2秒间隔
  }

  stop() {
    this.isActive = false;
  }
}

/**
 * 防调试开关
 * 1、阻止键盘事件打开控制台
 * 2、阻止鼠标事件打开控制台
 * 3、使用debugger关键字阻止打开控制台
 */
class DebugControl {
  private modules: any[] = [];

  constructor() {
    this.modules.push(new KeydownControl(), new RightMouseControl(), new DebugProtector());
  }
  /**
   * 开启防调试
   */
  start() {
    this.modules.forEach(m => m.start?.());
  }
  /**
   * 关闭防调试
   */
  stop() {
    this.modules.forEach(m => m.stop?.());
  }
}

export { DebugControl };
