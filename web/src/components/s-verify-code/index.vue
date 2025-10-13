<template>
  <span class="s-canvas" @click="changeCode">
    <canvas id="s-canvas" :width="contentWidth" :height="contentHeight"></canvas>
  </span>
</template>

<script setup lang="ts">
defineOptions({ name: "s-verify-code" });
// 通过defineProps接收父组件的值
const props = defineProps({
  identifyCodes: {
    //验证码从该字段中抽取生成
    type: String,
    default: "1234567890"
  },
  fontSizeMin: {
    // 字体最小值
    type: Number,
    default: 25
  },
  fontSizeMax: {
    // 字体最大值
    type: Number,
    default: 35
  },
  backgroundColorMin: {
    // 验证码图片背景色最小值
    type: Number,
    default: 200
  },
  backgroundColorMax: {
    // 验证码图片背景色最大值
    type: Number,
    default: 220
  },
  dotColorMin: {
    // 背景干扰点最小值
    type: Number,
    default: 60
  },
  dotColorMax: {
    // 背景干扰点最大值
    type: Number,
    default: 120
  },
  contentWidth: {
    //容器宽度
    type: Number,
    default: 100
  },
  contentHeight: {
    //容器高度
    type: Number,
    default: 44
  }
});
const emit = defineEmits(["verifyCodeChange"]);
const identifyCode = ref("");
watch(identifyCode, () => {
  drawPic();
});
onMounted(() => {
  drawPic();
  makeCode(props.identifyCodes, 4);
});
// 生成一个随机数
const randomNum = (min: any, max: any) => {
  return Math.floor(Math.random() * (max - min) + min);
};
// 生成一个随机的颜色
const randomColor = (min: any, max: any) => {
  let r = randomNum(min, max);
  let g = randomNum(min, max);
  let b = randomNum(min, max);
  return "rgb(" + r + "," + g + "," + b + ")";
};
const drawPic = () => {
  let canvas = <HTMLCanvasElement>document.getElementById("s-canvas");
  if (!canvas) return;
  let ctx = canvas.getContext("2d");
  if (!ctx) return;
  ctx.textBaseline = "bottom";
  // 绘制背景
  ctx.fillStyle = randomColor(props.backgroundColorMin, props.backgroundColorMax);
  ctx.fillRect(0, 0, props.contentWidth, props.contentHeight);
  // 绘制文字
  for (let i = 0; i < identifyCode.value.length; i++) {
    drawText(ctx, identifyCode.value[i], i);
  }
  drawLine(ctx);
  drawDot(ctx);
};
const drawText = (ctx: any, txt: any, i: any) => {
  ctx.fillStyle = randomColor(50, 160); //随机生成字体颜色
  ctx.font = randomNum(props.fontSizeMin, props.fontSizeMax) + "px SimHei"; //随机生成字体大小
  let x = (i + 1) * (props.contentWidth / (identifyCode.value.length + 1)); // x轴
  let y = randomNum(props.fontSizeMax, props.contentHeight); // y轴
  let deg = randomNum(-30, 30);
  // 修改坐标原点和旋转角度
  ctx.translate(x, y);
  ctx.rotate((deg * Math.PI) / 180);
  ctx.fillText(txt, 0, 0);
  // 恢复坐标原点和旋转角度
  ctx.rotate((-deg * Math.PI) / 180);
  ctx.translate(-x, -y);
};
const drawLine = (ctx: any) => {
  // 绘制干扰线
  for (let i = 0; i < 4; i++) {
    ctx.strokeStyle = randomColor(100, 200);
    ctx.beginPath();
    ctx.moveTo(randomNum(0, props.contentWidth), randomNum(0, props.contentHeight));
    ctx.lineTo(randomNum(0, props.contentWidth), randomNum(0, props.contentHeight));
    ctx.stroke();
  }
};
const drawDot = (ctx: any) => {
  // 绘制干扰点
  for (let i = 0; i < 30; i++) {
    ctx.fillStyle = randomColor(0, 255);
    ctx.beginPath();
    ctx.arc(randomNum(0, props.contentWidth), randomNum(0, props.contentHeight), 1, 0, 2 * Math.PI);
    ctx.fill();
  }
};
/*切换验证码start*/
const changeCode = () => {
  identifyCode.value = "";
  makeCode(props.identifyCodes, 4);
};
const makeCode = (e: any, n: any) => {
  for (let i = 0; i < n; i++) {
    identifyCode.value += e[randomNum(0, e.length)];
  }
  emit("verifyCodeChange", identifyCode.value);
};
/*切换验证码end*/
</script>

<style lang="scss" scoped></style>
