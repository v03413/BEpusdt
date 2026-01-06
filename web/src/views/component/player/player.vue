<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space direction="vertical">
        <div class="player-box">
          <div id="mse"></div>
        </div>
        <div>采用开源视频播放器：<a-link href="https://h5player.bytedance.com/" target="_blank">西瓜播放器</a-link></div>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import Player from "xgplayer";
import "xgplayer/dist/index.min.css";

onMounted(() => {
  new Player({
    id: "mse",
    autoplay: true,
    volume: 0.3,
    lang: "zh",
    url: "//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-360p.mp4",
    poster: "//lf3-static.bytednsdoc.com/obj/eden-cn/nupenuvpxnuvo/xgplayer_doc/poster.jpg",
    playsinline: true,
    thumbnail: {
      pic_num: 44,
      width: 160,
      height: 90,
      col: 10,
      row: 10,
      urls: ["//lf3-static.bytednsdoc.com/obj/eden-cn/nupenuvpxnuvo/xgplayer_doc/xgplayer-demo-thumbnail.jpg"]
    },
    danmu: {
      comments: [
        {
          duration: 15000,
          id: "1",
          start: 3000,
          txt: "长弹幕长弹幕长弹幕长弹幕长弹幕",
          style: {
            //弹幕自定义样式
            color: "#ff9500",
            fontSize: "20px",
            border: "solid 1px #ff9500",
            borderRadius: "50px",
            padding: "5px 11px",
            backgroundColor: "rgba(255, 255, 255, 0.1)"
          }
        }
      ],
      area: {
        start: 0,
        end: 1
      }
    },
    height: "100%",
    width: "100%"
  });
});
</script>

<style lang="scss" scoped>
.player-box {
  width: 600px;
  height: 400px;
}
</style>
