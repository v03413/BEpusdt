<template>
  <div class="doc-container">
    <!-- 作者声明 -->
    <div class="notice-card">
      <div class="notice-title">
        <span>❗️</span>
        <span>作者声明</span>
      </div>
      <p class="notice-text">
        本项目为研究学习区块链的开源项目，不提供任何形式的收费服务(谨防诈骗)，不鼓励任何衍生金融属性的交易行为，不负责任何使用本项目进行的三方行为；使用过程中遇见问题请提
        <a href="https://github.com/v03413/BEpusdt/issues" target="_blank" class="notice-link">issue</a>
        或群里交流，开源项目请自重！
      </p>
    </div>

    <div class="cards-row">
      <!-- 社区交流 -->
      <div class="card">
        <div class="card-title">🏝️ 社区交流</div>
        <div class="divider"></div>
        <ul class="link-list">
          <li>
            <span class="label">Telegram 群组</span>
            <a href="https://t.me/BEpusdtChat" target="_blank" class="link">https://t.me/BEpusdtChat</a>
          </li>
          <li>
            <span class="label">Telegram 频道</span>
            <a href="https://t.me/BEpusdtChannel" target="_blank" class="link">https://t.me/BEpusdtChannel</a>
          </li>
        </ul>
      </div>

      <!-- 项目文档 -->
      <div class="card">
        <div class="card-title">📄 项目文档</div>
        <div class="divider"></div>
        <ul class="link-list">
          <li>
            <span class="label">仓库地址</span>
            <a href="https://github.com/v03413/BEpusdt" target="_blank" class="link">https://github.com/v03413/BEpusdt</a>
          </li>
          <li>
            <span class="label">API 对接</span>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/api/api.md" target="_blank" class="link"
              >https://github.com/v03413/BEpusdt/blob/main/docs/api/api.md</a
            >
          </li>
          <li>
            <span class="label">第三方程序对接说明</span>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/api/other.md" target="_blank" class="link"
              >https://github.com/v03413/BEpusdt/blob/main/docs/api/other.md</a
            >
          </li>
          <li>
            <span class="label">交易类型完整列表</span>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/trade-type.md" target="_blank" class="link"
              >https://github.com/v03413/BEpusdt/blob/main/docs/trade-type.md</a
            >
          </li>
        </ul>
      </div>

      <!-- 常见问题 -->
      <div class="card">
        <div class="card-title">❓ 常见问题</div>
        <div class="divider"></div>
        <ul class="link-list faq-list">
          <li>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/faq/server.md" target="_blank" class="link"
              >服务器配置性能选型推荐 ⚡️</a
            >
          </li>
          <li>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/faq/bandwidth.md" target="_blank" class="link"
              >服务器流量消耗说明解释</a
            >
          </li>
          <li>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/faq/login-reset.md" target="_blank" class="link"
              >后台入口账密忘记重置教程</a
            >
          </li>
          <li>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/faq/telegram-chat-id.md" target="_blank" class="link"
              >Telegram 通知 Chat ID 获取教程</a
            >
          </li>
          <li>
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/faq/rpc-endpoint.md" target="_blank" class="link"
              >区块 RPC 节点稳定性说明指南 ‼️</a
            >
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.doc-container {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 作者声明 */
.notice-card {
  background: rgba(var(--warning-1), 0.3);
  border: 1px solid rgb(var(--warning-5));
  border-radius: 6px;
  padding: 16px 20px;
}

.notice-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 16px;
  font-weight: 600;
  color: rgb(var(--warning-7));
  margin-bottom: 10px;
}

.notice-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.8;
  color: $color-text-1;
  word-break: break-word;
}

.notice-link {
  color: rgb(var(--warning-7));
  text-decoration: underline;

  &:hover {
    color: rgb(var(--warning-9));
  }
}

/* 两列布局 */
.cards-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

/* 卡片 */
.card {
  background: $color-bg-2;
  border: 1px solid $color-border-2;
  border-radius: 6px;
  padding: 20px;
  box-shadow: $shadow2-center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: $color-text-1;
}

.divider {
  height: 1px;
  background-color: $color-border-2;
  margin: 12px 0;
}

/* 链接列表 */
.link-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;

  li {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
}

.label {
  font-size: 12px;
  font-weight: 500;
  color: $color-text-3;
}

.link {
  font-size: 13px;
  color: $color-link;
  text-decoration: none;
  word-break: break-all;

  &:hover {
    color: rgb(var(--link-5));
  }
}

.faq-list {
  li {
    flex-direction: row;
    align-items: center;
    gap: 8px;

    &::before {
      content: "·";
      color: $color-text-3;
      font-size: 18px;
      line-height: 1;
      flex-shrink: 0;
    }
  }
}
</style>
