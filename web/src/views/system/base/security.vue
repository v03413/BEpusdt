<template>
  <div>
    <a-row align="center" :gutter="[0, 16]">
      <a-col :span="24">
        <a-card title="安全设置">
          <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
            <a-form-item field="admin_secure" extra="长度限制：8-18位，必须以 / 开头" label="安全入口">
              <a-input v-model="form.admin_secure" placeholder="请输入安全入口" allow-clear />
            </a-form-item>

            <a-form-item field="admin_username" label="管理账号">
              <div class="username-input-wrapper">
                <a-input v-model="form.admin_username" placeholder="请输入管理账号" allow-clear />
                <a-button type="text" @click="showPasswordModal" class="password-btn">修改密码</a-button>
              </div>
            </a-form-item>

            <a-form-item>
              <a-button type="primary" html-type="submit">保存设置</a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>

    <!-- 修改密码弹窗 -->
    <a-modal v-model:visible="passwordModalVisible" title="修改密码" @ok="handlePasswordSubmit" @cancel="handlePasswordCancel">
      <a-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef">
        <a-form-item field="password" label="当前密码">
          <a-input-password v-model="passwordForm.password" placeholder="请输入当前密码" allow-clear />
        </a-form-item>

        <a-form-item field="new_password" label="新密码">
          <a-input-password v-model="passwordForm.new_password" placeholder="请输入新密码" allow-clear />
        </a-form-item>

        <a-form-item field="confirm_password" label="重复新密码">
          <a-input-password v-model="passwordForm.confirm_password" placeholder="请再次输入新密码" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { Message } from "@arco-design/web-vue";
import { setPasswordAPI } from "@/api/modules/user";
import { setsConfAPI } from "@/api/modules/conf/index";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;

// 基础设置表单
const form = ref({
  admin_secure: "",
  admin_username: ""
});

// 密码修改表单
const passwordForm = ref({
  password: "",
  new_password: "",
  confirm_password: ""
});

// 密码弹窗相关
const passwordModalVisible = ref(false);
const passwordFormRef = ref();

// 基础设置验证规则
const rules = {
  admin_secure: [
    {
      required: true,
      message: "安全入口不能为空"
    },
    {
      validator: (value: string, cb: any) => {
        if (!value.startsWith("/")) {
          cb("安全入口必须以 / 开头");
        } else if (value.length < 8) {
          cb("安全入口长度不能少于8位");
        } else if (value.length > 18) {
          cb("安全入口长度不能超过18位");
        } else if (!/^\/[a-zA-Z0-9]+$/.test(value)) {
          cb("安全入口只能包含字母和数字");
        } else {
          cb();
        }
      }
    }
  ],
  admin_username: [
    {
      required: true,
      message: "管理账号不能为空"
    }
  ]
};

// 密码修改验证规则
const passwordRules = {
  current_password: [
    {
      required: true,
      message: "当前密码不能为空"
    }
  ],
  new_password: [
    {
      required: true,
      message: "新密码不能为空"
    },
    {
      minLength: 6,
      message: "新密码长度不能少于6位"
    }
  ],
  confirm_password: [
    {
      required: true,
      message: "请重复输入新密码"
    },
    {
      validator: (value: string, cb: any) => {
        if (value !== passwordForm.value.new_password) {
          cb("两次输入的密码不一致");
        } else {
          cb();
        }
      }
    }
  ]
};

// 保存基础设置
const onSubmit = async ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;

  try {
    const response = await setsConfAPI([
      { key: "admin_username", value: form.value.admin_username },
      { key: "admin_secure", value: form.value.admin_secure }
    ]);

    if (response && response.code === 200) {
      Message.success("设置保存成功！");
      emit("refresh");
    } else {
      Message.error(response?.msg || "设置保存失败");
    }
  } catch (error: any) {
    Message.error(error);
  }
};

// 显示密码修改弹窗
const showPasswordModal = () => {
  passwordModalVisible.value = true;
  // 重置密码表单
  passwordForm.value = {
    password: "",
    new_password: "",
    confirm_password: ""
  };
};

// 提交密码修改
const handlePasswordSubmit = async () => {
  try {
    const valid = await passwordFormRef.value?.validate();
    if (!valid) {
      const response = await setPasswordAPI({
        password: passwordForm.value.password,
        new_password: passwordForm.value.new_password,
        confirm_password: passwordForm.value.confirm_password
      });

      if (response && response.code === 200) {
        Message.success("密码修改成功，请重新登录！");
        passwordModalVisible.value = false;
        emit("refresh");
      } else {
        Message.error(response?.msg || "密码修改失败");
      }
    }
  } catch (error: any) {
    console.error("密码修改失败:", error);
    Message.error("密码修改失败，请稍后重试");
  }
};

// 取消密码修改
const handlePasswordCancel = () => {
  passwordModalVisible.value = false;
  passwordForm.value = {
    password: "",
    new_password: "",
    confirm_password: ""
  };
};

// 监听数据变化，更新表单
watch(
  () => data.value,
  () => {
    if (data.value) {
      form.value.admin_username = data.value.admin_username || "";
      form.value.admin_secure = data.value.admin_secure || "";
    }
  },
  { immediate: true }
);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}

.username-input-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;

  .password-btn {
    flex-shrink: 0;
    margin-left: 8px;
  }
}
</style>
