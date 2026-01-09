<template>
  <div class="doc-container">
    <div class="card">
      <h2>项目文档</h2>
      <div class="divider"></div>
      <p>
        仓库地址：
        <a href="https://github.com/v03413/BEpusdt" target="_blank" class="link"> https://github.com/v03413/BEpusdt </a>
      </p>
      <p style="margin-top: 12px">
        API对接：
        <a href="https://github.com/v03413/BEpusdt/blob/main/docs/api/api.md" target="_blank" class="link">
          https://github.com/v03413/BEpusdt/blob/main/docs/api/api.md
        </a>
      </p>
    </div>
  </div>
</template>

<style scoped>
.doc-container {
  padding: 20px;
}

.card {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

h2 {
  margin: 0 0 10px 0;
  font-size: 18px;
  color: #303133;
}

.divider {
  height: 1px;
  background-color: #dcdfe6;
  margin: 16px 0;
}

p {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.link {
  color: #409eff;
  text-decoration: none;
}

.link:hover {
  color: #66b1ff;
}
</style>
