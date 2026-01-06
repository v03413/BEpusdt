<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space direction="vertical">
        <div>
          <div style="border: 1px solid #cccccc">
            <Toolbar style="border-bottom: 1px solid #cccccc" :editor="editorRef" :default-config="toolbarConfig" :mode="mode" />
            <Editor
              style="height: 500px; overflow-y: hidden"
              v-model="valueHtml"
              :default-config="editorConfig"
              :mode="mode"
              @on-created="handleCreated"
            />
          </div>
        </div>
        <div>
          采用开源富文本编辑器：<a-link href="https://github.com/wangeditor-team/wangEditor" target="_blank">wangEditor</a-link>
        </div>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import "@wangeditor/editor/dist/css/style.css"; // 引入 css
import { Editor, Toolbar } from "@wangeditor/editor-for-vue";

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef();

// 内容 HTML
const valueHtml = ref("<p>hello</p>");

// 模拟 ajax 异步获取内容
onMounted(() => {
  setTimeout(() => {
    valueHtml.value = "<p>模拟ajax异步获取的内容</p>";
  }, 1500);
});

const toolbarConfig = {};
const editorConfig = { placeholder: "请输入内容..." };

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor == null) return;
  editor.destroy();
});
const mode = "default";
const handleCreated = (editor: any) => {
  editorRef.value = editor; // 记录 editor 实例，重要！
};
</script>

<style lang="scss" scoped></style>
