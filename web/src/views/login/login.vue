<template>
  <div class="container">
    <div class="login">
      <LoginBanner v-if="isPc" />
      <div class="login_box">
        <div class="login_title">Welcome Back</div>
        <div class="login_title_desc">一款更好用的个人 USDT/USDC 收款网关，让您的收款更加便捷！</div>
        <LoginForm />
        <!-- <div class="author">by 兔子先森</div> -->
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import LoginBanner from "@/views/login/components/login-banner.vue";
import LoginForm from "@/views/login/components/login-form.vue";
import { useDevicesSize } from "@/hooks/useDevicesSize";
const { isPc } = useDevicesSize();
</script>

<style lang="scss" scoped>
.container {
  position: relative;
  height: 100vh;
  overflow: hidden;
  .login {
    position: absolute;
    top: 50%;
    left: 50%;
    display: flex;
    align-items: center;
    max-width: 1000px;
    height: 500px;
    box-shadow: 0 0 8px 1px $color-fill-2;
    transform: translate(-50%, -50%);
    .login_box {
      position: relative;
      box-sizing: border-box;
      width: 350px;
      height: 100%;
      padding: 40px 30px 30px;
      .login_title {
        margin-bottom: $margin-text;
        font-size: $font-size-title-2;
        color: $color-text-1;
      }
      .login_title_desc {
        font-size: $font-size-body-1;
        color: $color-text-3;
      }
      .author {
        position: absolute;
        bottom: 30px;
        font-size: $font-size-body-1;
        color: $color-text-4;
      }
    }
  }
}
</style>
