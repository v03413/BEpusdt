<template>
  <div class="snow-page">
    <div class="snow-inner">
      <div class="my-avatar">
        <a-avatar :size="100">
          <img alt="avatar" :src="myImage" />
        </a-avatar>
        <div class="my-name">
          <div class="my-title title-size-1">兔子先森</div>
          <div class="my-local"><icon-location />浙江杭州</div>
        </div>
      </div>
      <a-grid :cols="{ xs: 1, sm: 1, md: 1, lg: 2, xl: 3 }" :col-gap="16" :row-gap="16">
        <a-grid-item>
          <div class="title-size-1">👋关于我</div>
          <div class="type-row" v-for="item in list.aboutMe" :key="item.label">
            <div class="column-title-size-1">{{ item.label }}：</div>
            <div v-for="(type, index) in item.value" :key="index">
              <a-tag :color="item.theme">
                {{ type }}
              </a-tag>
            </div>
          </div>
        </a-grid-item>
        <a-grid-item>
          <div class="title-size-1">💬联系我</div>
          <div class="type-row" v-for="item in list.callMe" :key="item.label">
            <div class="column-title-size-2">{{ item.label }}：</div>
            <div :class="item.link ? 'link-hover' : ''" @click="onLink(item)">{{ item.value }}</div>
          </div>
        </a-grid-item>
        <a-grid-item>
          <div class="title-size-1">🐾其它</div>
          <div class="type-row" v-for="item in list.other" :key="item.label">
            <div class="column-title-size-2">{{ item.label }}：</div>
            <div>{{ item.value }}</div>
          </div>
          <a-image width="100" :src="officialAccount" />
        </a-grid-item>
        <a-grid-item>
          <a-image default-scale height="220" :src="otherImage" />
        </a-grid-item>
        <a-grid-item>
          <div class="title-size-2">✨Hello 朋友们✨</div>
          <div class="type-row" v-for="item in list.introduce" :key="item.label">
            <div class="text-ellipsis">{{ item.label }}</div>
          </div>
        </a-grid-item>
      </a-grid>
    </div>
  </div>
</template>

<script setup lang="ts">
import myImage from "@/assets/img/my-image.jpg";
import otherImage from "@/assets/img/other-image.jpg";
import officialAccount from "@/assets/img/official-account.png";

const list = reactive({
  aboutMe: [
    { label: "🐦我用语雀", value: ["知识管理", "记录生活"], theme: "green" },
    { label: "💫个性性格", value: ["天秤座", "INFJ-T"], theme: "blue" },
    { label: "💼职业工作", value: ["互联网", "前端开发"], theme: "orangered" },
    { label: "🏠现居城市", value: ["杭州"], theme: "magenta" },
    { label: "💙兴趣爱好", value: ["读书", "音乐", "狗狗🐕"], theme: "purple" }
  ],
  callMe: [
    { label: "📧邮箱", value: "2547096351@qq.com", link: false },
    { label: "📠微信", value: "fanction-w", link: false },
    { label: "🛰️github", value: "兔子先森", link: true, url: "https://github.com/WangFan-io" },
    { label: "📚思否", value: "兔子先森", link: true, url: "https://segmentfault.com/blog/dcodes" },
    { label: "✒️掘金", value: "兔子先森Ace", link: true, url: "https://juejin.cn/user/1728883023940600/posts" }
  ],
  other: [
    { label: "🧭公众号", value: "DCodes" },
    { label: "🌟公众号二维码", value: "" }
  ],
  introduce: [
    { label: "👋欢迎来到兔子先森的project，记录思考成长，一起进步" },
    { label: "🔖长期主义者一枚，持续学习拓展思维及方法论" },
    { label: "🤔行业瞬息万变，时刻关注行业最新资讯" },
    { label: "🤝职场相处沟通都需要技巧，持续记录学习职场二三事" },
    { label: "🌈人生多彩，除了工作还有很大的世界供我们探索，Let's go!" }
  ]
});

const onLink = (item: any) => {
  if (!item.link) return;
  window.open(item.url);
};
</script>

<style lang="scss" scoped>
.border {
  border: 1px solid red;
}
.my-avatar {
  display: flex;
  align-items: center;
  margin-bottom: calc($margin * 2);
  .my-name {
    height: 60px;
    margin-left: $margin;
    .my-title {
      height: 30px;
    }
    .my-local {
      display: flex;
      align-items: end;
      height: 30px;
      font-size: $font-size-body-3;
    }
  }
}
.type-row {
  display: flex;
  column-gap: $margin;
  align-items: center;
  height: 25px;
  margin-top: $margin;
  .column-title-size-1 {
    min-width: 90px;
  }
  .column-title-size-2 {
    min-width: 68px;
  }
  .text-ellipsis {
    white-space: nowrap;
  }
}
.title-size-1 {
  font-size: $font-size-title-1;
  font-weight: bold;
}
.title-size-2 {
  font-size: $font-size-title-2;
  font-weight: bold;
}
.link-hover {
  color: $color-primary;
  cursor: pointer;
}
</style>
