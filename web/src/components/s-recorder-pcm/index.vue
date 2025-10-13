<template>
  <div style="width: 100%; height: 100%" ref="recwave"></div>
</template>

<script setup lang="ts">
//必须引入的核心
import Recorder from "recorder-core";
//引入pcm格式支持文件
import "recorder-core/src/engine/pcm";
//可选的插件支持项，这个是波形可视化插件
import "recorder-core/src/extensions/waveview";
//ts import 提示：npm包内已自带了.d.ts声明文件（不过是any类型）

const emit = defineEmits(["change"]);

let wave: any; // 用于绘制波形
const recwave = ref(null); // 绑定到dom元素上
let rec: any; // Recorder实例
let send_chunk: any; // 上次分割点数据
let testSampleRate = 16000; // 采样率

//重置环境，每次开始录音时必须先调用此方法，清理环境
const RealTimeSendReset = () => {
  send_chunk = null;
};

// 打开录音
const recOpen = () => {
  RealTimeSendReset();
  // 创建录音对象
  rec = Recorder({
    type: "unknown", //这里特意使用unknown格式，方便清理内存
    onProcess: (buffers: any, powerLevel: any, _: any, bufferSampleRate: any) => {
      // 所有的pcm数据queue，缓存采样率，是否结束
      RealTimeSendTry(buffers, bufferSampleRate, false);
      if (wave) {
        wave.input(buffers[buffers.length - 1], powerLevel, bufferSampleRate);
      }
    }
  });
  if (!rec) {
    alert("当前浏览器不支持录音功能！");
    return;
  }
  // 打开录音，获得权限
  rec.open(
    () => {
      if (recwave.value) {
        // 创建音频可视化图形绘制对象
        wave = Recorder.WaveView({ elem: recwave.value });
      }
    },
    (msg: any, isUserNotAllow: any) => {
      // 用户拒绝了录音权限，或者浏览器不支持录音
      console.log((isUserNotAllow ? "UserNotAllow，" : "") + "无法录音:" + msg);
    }
  );
};

const RealTimeSendTry = (buffers: any, bufferSampleRate: any, isClose: any) => {
  //提取出新的pcm数据
  let pcm = new Int16Array(0);
  if (buffers.length > 0) {
    //【关键代码】借用SampleData函数进行数据的连续处理，采样率转换是顺带的，得到新的pcm数据
    // send_chunk为上次分割点
    let chunk = Recorder.SampleData(buffers, bufferSampleRate, testSampleRate, send_chunk);
    send_chunk = chunk; // 保存本次分割点，用于下次使用
    pcm = chunk.data; //此时的pcm就是原始的音频16位pcm数据（小端LE），直接保存即为16位pcm文件、加个wav头即为wav文件、丢给mp3编码器转一下码即为mp3文件
  }

  //没有指定固定的帧大小，直接把pcm发送出去即可
  TransferUpload(pcm, isClose);
  return;
};

//=====数据传输函数==========
const TransferUpload = (pcmFrame: any, isClose: any) => {
  if (isClose && pcmFrame.length == 0) {
    // 这里就是数据发送完成，录音结束的位置，可以停止ws了
    emit("change", { arrayBuffer: pcmFrame.buffer, close: true });
    // ws.send(arrayBuffer,true)
    return; //如果不需要处理最后一帧数据，直接return不做任何处理
  }
  // 二进制数据
  let arrayBuffer = pcmFrame.buffer;

  emit("change", { arrayBuffer, close: false });
  //可以实现
  //WebSocket send(arrayBuffer) ...
  //WebRTC send(arrayBuffer) ...
  //XMLHttpRequest send(arrayBuffer) ...
};

// 开始录音
const recStart = () => {
  if (!rec) return console.error("未打开录音");
  rec.start();
  // 已开始录音
};
// 结束录音
const recStop = () => {
  if (!rec) return console.error("未打开录音");
  rec.close(); // 关闭录音，释放录音资源
  rec = null;
  RealTimeSendTry([], 0, true); //最后一次发送
};

defineExpose({
  recOpen,
  recStart,
  recStop
});
</script>

<style lang="scss" scoped></style>
