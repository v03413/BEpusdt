<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-alert type="success">工具类位置: src/utils/verify-tools.ts</a-alert>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>处理输入内容，只允许输入数字和2位小数👉️ {{ text01 }}</div>
        <div>
          <a-input
            :style="{ width: '320px' }"
            placeholder="请输入"
            allow-clear
            v-model="text01"
            @input="onVerifyNumberIntegerAndFloat"
          />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>输入框限制只允许正整数输入👉️ {{ text02 }}</div>
        <div>
          <a-input
            :style="{ width: '320px' }"
            placeholder="请输入"
            allow-clear
            v-model="text02"
            @input="onVerifiyNumberInteger"
          />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          手机号码格式验证，是否正确👉️
          {{ verifyPhone(text03) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text03" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          国内电话号码格式验证👉️
          {{ verifyTelPhone(text04) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text04" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          登录账号验证 (字母开头，允许5-16字节，允许字母数字下划线)👉️
          {{ verifyAccount(text05) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text05" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          密码 (以字母开头，长度在6~16之间，只能包含字母、数字和下划线)👉️
          {{ verifyPassword(text06) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text06" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          强密码 (字母+数字+特殊字符，长度在6-16之间)👉️
          {{ verifyPasswordPowerful(text07) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text07" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          密码强度(长度6~16)👉️
          {{ verifyPasswordStrength(text08) || "无" }}
        </div>
        <div>弱：纯数字，纯字母，纯特殊字符</div>
        <div>中：字母+数字，字母+特殊字符，数字+特殊字符</div>
        <div>强：字母+数字+特殊字符</div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text08" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          校验IP地址是否正确👉️
          {{ verifyIPAddress(text09) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text09" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          校验邮箱是否正确👉️
          {{ verifyEmail(text10) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text10" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          身份证验证👉️
          {{ verifyIdCard(text11) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text11" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          中文姓名验证👉️
          {{ verifyFullName(text12) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text12" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          邮政编码验证👉️
          {{ verifyPostalCode(text13) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text13" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          url正确性验证👉️
          {{ verifyUrl(text14) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text14" />
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>
          车牌号验证👉️
          {{ verifyCarNum(text15) }}
        </div>
        <div>
          <a-input :style="{ width: '320px' }" placeholder="请输入" allow-clear v-model="text15" />
        </div>
      </a-space>
      <a-divider />
      <div>基于<a-link href="https://gitee.com/lyt-top/vue-next-admin" target="_blank">vue-next-admin</a-link>的验证工具</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  verifyNumberIntegerAndFloat,
  verifiyNumberInteger,
  verifyPhone,
  verifyTelPhone,
  verifyAccount,
  verifyPassword,
  verifyPasswordPowerful,
  verifyPasswordStrength,
  verifyIPAddress,
  verifyEmail,
  verifyIdCard,
  verifyFullName,
  verifyPostalCode,
  verifyUrl,
  verifyCarNum
} from "@/utils/verify-tools.ts";
const text01 = ref<string>("");
const onVerifyNumberIntegerAndFloat = (e: string) => {
  text01.value = verifyNumberIntegerAndFloat(e);
};

const text02 = ref<string>("");
const onVerifiyNumberInteger = (e: string) => {
  text02.value = verifiyNumberInteger(e);
};

const text03 = ref<string>("");
const text04 = ref<string>("");
const text05 = ref<string>("");
const text06 = ref<string>("");
const text07 = ref<string>("");
const text08 = ref<string>("");
const text09 = ref<string>("");
const text10 = ref<string>("");
const text11 = ref<string>("");
const text12 = ref<string>("");
const text13 = ref<string>("");
const text14 = ref<string>("");
const text15 = ref<string>("");
</script>
