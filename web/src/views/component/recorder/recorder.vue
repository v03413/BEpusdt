<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space direction="vertical" fill>
        <a-alert>本地录音需要https或localhost安全环境，是否为安全环境：{{ isSecureEnvironment() }}</a-alert>
        <a-alert type="success"
          >基于<a-link href="https://github.com/xiangyuecn/Recorder" target="_blank">Recorder.js</a-link
          >，采用PCM实时帧回调的方式，用于流式语音传输</a-alert
        >
        <a-space direction="vertical" fill>
          <!-- 波形绘制区域 -->
          <div class="waveview-box">
            <s-recorder-pcm ref="RecorderPCM" @change="onChange" />
          </div>
          <a-space>
            <a-button type="primary" @click="recOpen">打开录音权限</a-button>
            <a-button type="primary" @click="recStart">开始录音</a-button>
            <a-button @click="recStop">结束录音</a-button>
            <a-button @click="audioList.length = 0">清空列表</a-button>
          </a-space>
        </a-space>
        <audio ref="audioElement" controls v-if="audioList.length"></audio>
        <a-list v-if="audioList.length">
          <template #header>音频帧列表</template>
          <a-list-item v-for="(item, index) in audioList" :key="index">
            <a-space>
              <span v-if="!item.close">PCM音频实时帧 - {{ index }}</span>
              <span v-else>PCM完整音频</span>
              <a-button-group>
                <a-button @click="playAudio(item)">播放</a-button>
              </a-button-group>
            </a-space>
          </a-list-item>
        </a-list>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { isSecureEnvironment } from "@/utils/index";

const RecorderPCM = ref();

const recOpen = () => {
  audioList.value.length = 0;
  RecorderPCM.value.recOpen();
};

const recStart = () => {
  RecorderPCM.value.recStart();
};

const recStop = () => {
  RecorderPCM.value.recStop();
};

interface AudioInfo {
  arrayBuffer: ArrayBuffer;
  close: boolean;
}

// 音频帧列表
const audioList = ref<AudioInfo[]>([]);

// 音频帧实时回调
const onChange = async (e: AudioInfo) => {
  if (!e.close) {
    audioList.value.push(e);
  } else {
    // 合并音频帧列表
    let complete = mergePcmChunks(audioList.value.map(item => item.arrayBuffer));
    audioList.value.push({ arrayBuffer: complete, close: true });
  }
};

// 合并多个PCM分片
const mergePcmChunks = (chunks: ArrayBuffer[]) => {
  // 计算总长度
  const totalLength = chunks.reduce((sum, chunk) => sum + chunk.byteLength, 0);

  // 创建新的ArrayBuffer
  const mergedBuffer = new ArrayBuffer(totalLength);
  const mergedView = new Uint8Array(mergedBuffer);

  // 按顺序复制每个分片
  let offset = 0;
  chunks.forEach(chunk => {
    mergedView.set(new Uint8Array(chunk), offset);
    offset += chunk.byteLength;
  });

  return mergedBuffer;
};

const audioElement = ref<any>();
const playAudio = (e: AudioInfo) => {
  const wavBuffer = pcmToWav(e.arrayBuffer);
  const audioBlob = new Blob([wavBuffer], { type: "audio/wav" });

  // 为Blob创建URL对象
  const audioUrl = URL.createObjectURL(audioBlob);

  // 设置audio元素的src
  audioElement.value.src = audioUrl;

  // 开始播放
  audioElement.value.play();

  // 播放结束后释放URL对象
  audioElement.value.onended = function () {
    URL.revokeObjectURL(audioUrl);
  };
};

// PCM转WAV函数
const pcmToWav = (pcmArrayBuffer: ArrayBuffer, sampleRate = 16000, numChannels = 1, bitDepth = 16) => {
  const buffer = new ArrayBuffer(44 + pcmArrayBuffer.byteLength);
  const view = new DataView(buffer);

  // 写入WAV文件头
  writeString(view, 0, "RIFF");
  view.setUint32(4, 36 + pcmArrayBuffer.byteLength, true);
  writeString(view, 8, "WAVE");

  // 写入fmt子块
  writeString(view, 12, "fmt ");
  view.setUint32(16, 16, true); // 子块大小
  view.setUint16(20, 1, true); // 音频格式（1表示PCM）
  view.setUint16(22, numChannels, true); // 声道数
  view.setUint32(24, sampleRate, true); // 采样率
  view.setUint32(28, sampleRate * numChannels * (bitDepth / 8), true); // 字节率
  view.setUint16(32, numChannels * (bitDepth / 8), true); // 块对齐
  view.setUint16(34, bitDepth, true); // 位深度

  // 写入data子块
  writeString(view, 36, "data");
  view.setUint32(40, pcmArrayBuffer.byteLength, true);

  // 写入PCM数据
  const pcmData = new Uint8Array(pcmArrayBuffer);
  const resultData = new Uint8Array(buffer);
  resultData.set(pcmData, 44);

  return buffer;
};

const writeString = (view: DataView, offset: number, string: string) => {
  for (let i = 0; i < string.length; i++) {
    view.setUint8(offset + i, string.charCodeAt(i));
  }
};
</script>

<style lang="scss" scoped>
.waveview-box {
  display: inline-block;
  width: 400px;
  height: 100px;
  border: 1px solid $color-border-2;
}
</style>
