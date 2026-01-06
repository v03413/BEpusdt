<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-alert type="success">工具类位置: src/utils/file-tools.ts</a-alert>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>获取麦克风权限: getMicPower</div>
        <div>返回Promis的boolean，示例: getMicPower()</div>
        <div>
          使用方式：
          <a-typography-text code>getMicPower()</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>十六进制字符串转换为字节数组: hexStringToByteArray</div>
        <div>返回Uint8Array，示例: hexStringToByteArray(hexString)</div>
        <div>
          使用方式：
          <a-typography-text code>hexStringToByteArray(十六进制字符串)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>Uint8Array转换为ArrayBuffer: byteArrayToArrayBuffer</div>
        <div>返回ArrayBuffer，示例: byteArrayToArrayBuffer(Uint8Array)</div>
        <div>
          使用方式：
          <a-typography-text code>byteArrayToArrayBuffer(你的Uint8Array数据)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>base64转blob: base64ToBlob</div>
        <div>返回处理后的blob，示例: base64ToBlob(base64)</div>
        <div>
          使用方式：
          <a-typography-text code>base64ToBlob(base64字符串, MIME类型-可不传, 分片大小-可不传)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>base64转ArrayBuffer: base64ToArrayBuffer</div>
        <div>返回处理后的ArrayBuffer，示例: base64ToArrayBuffer(base64)</div>
        <div>
          使用方式：
          <a-typography-text code>base64ToArrayBuffer(base64字符串)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>blob转Uint8Array: blobToUint8Array</div>
        <div>返回处理后的Uint8Array，示例: blobToUint8Array(base64)</div>
        <div>
          使用方式：
          <a-typography-text code>blobToUint8Array(你的blob数据)</a-typography-text>
        </div>
      </a-space>
    </div>
  </div>
</template>
