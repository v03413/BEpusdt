<template>
  <div>
    <div class="login_form_box">
      <a-form :rules="rules" :model="form" layout="vertical" @submit="onSubmit">
        <a-form-item field="username" :hide-asterisk="true">
          <a-input v-model="form.username" allow-clear placeholder="请输入账号">
            <template #prefix>
              <icon-user />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item field="password" :hide-asterisk="true">
          <a-input-password v-model="form.password" allow-clear placeholder="请输入密码">
            <template #prefix>
              <icon-lock />
            </template>
          </a-input-password>
        </a-form-item>
        <a-form-item field="remember">
          <div class="remember">
            <a-checkbox v-model="form.remember">记住密码</a-checkbox>
            <div class="forgot-password" @click="handleForgotPassword">忘记密码</div>
          </div>
        </a-form-item>
        <a-form-item>
          <a-button long type="primary" html-type="submit">登录</a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUserInfoStore } from "@/store/modules/user-info";
import { loginAPI } from "@/api/modules/user/index";

let userStores = useUserInfoStore();
const router = useRouter();

// 记住密码的存储key
const REMEMBER_KEY = "login_remember_info";

// 简单的加密解密函数
const encrypt = (str: string) => {
  return btoa(encodeURIComponent(str));
};

const decrypt = (str: string) => {
  try {
    return decodeURIComponent(atob(str));
  } catch {
    return "";
  }
};

const form = ref({
  username: "",
  password: "",
  verifyCode: null,
  remember: false
});
const rules = ref({
  username: [
    {
      required: true,
      message: "请输入账号"
    }
  ],
  password: [
    {
      required: true,
      message: "请输入密码"
    }
  ]
});

// 组件挂载时读取记住的密码
onMounted(() => {
  const savedInfo = localStorage.getItem(REMEMBER_KEY);
  if (savedInfo) {
    try {
      const { username, password, remember } = JSON.parse(savedInfo);
      form.value.username = decrypt(username);
      form.value.password = decrypt(password);
      form.value.remember = remember;
    } catch (error) {
      console.error("读取记住的密码失败:", error);
      localStorage.removeItem(REMEMBER_KEY);
    }
  }
});

// 提交表单
const onSubmit = async ({ errors }: any) => {
  if (errors) return;
  onLogin();
};

// 登录
const onLogin = async () => {
  // 处理记住密码
  if (form.value.remember) {
    // 保存加密后的账号密码
    const rememberInfo = {
      username: encrypt(form.value.username),
      password: encrypt(form.value.password),
      remember: true
    };
    localStorage.setItem(REMEMBER_KEY, JSON.stringify(rememberInfo));
  } else {
    // 不记住密码时清除存储
    localStorage.removeItem(REMEMBER_KEY);
  }

  // 登录
  let res = await loginAPI(form.value);

  userStores.token = res.data.token;

  // 加载用户信息
  await userStores.setAccount();
  // // 加载路由信息
  // await routeStore.initSetRouter();

  arcoMessage("success", "登录成功");
  // 跳转首页
  router.replace("/home");
  // 设置字典
  // useSystemStore().setDictData();
};

// 忘记密码
const handleForgotPassword = () => {
  window.open("https://github.com/v03413/BEpusdt", "_blank");
};
</script>

<style lang="scss" scoped>
.login_form_box {
  margin-top: 28px;
  .verifyCode {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }
  .remember {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    .forgot-password {
      color: $color-primary;
      cursor: pointer;
    }
  }
}
</style>
