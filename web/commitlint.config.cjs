module.exports = {
  ignores: [commit => commit.includes("init")],
  extends: ["@commitlint/config-conventional"],
  rules: {
    "body-leading-blank": [2, "always"],
    "footer-leading-blank": [1, "always"],
    "header-max-length": [2, "always", 108],
    "subject-empty": [2, "never"],
    "type-empty": [2, "never"],
    "type-enum": [
      2,
      "always",
      [
        "feat",
        "fix",
        "style",
        "perf",
        "refactor",
        "test",
        "chore",
        "docs",
        "workflow",
        "build",
        "ci",
        "revert",
        "wip",
        "release",
        "deps",
        "merge",
        "sync"
        // feat：交新功能
        // fix：bug修复
        // style： 调整代码格式，未修改代码逻辑
        // perf: 性能优化
        // refactor：代码重构，既没修复bug也没有添加新功能
        // test：增加测试
        // chore: 其他修改（不在上述类型中的修改）
        // docs：修改文档或注释
        // workflow: 工作流
        // build: 项目打包构建相关的配置修改
        // ci: 持续集成相关
        // revert: 恢复上一次提交（回滚）
        // wip: work in progress 工作中 还没完成
        // release: 发版
        // deps: 依赖相关的修改
        // merge	代码合并
        // sync	同步主线或分支的Bug
        // scope	说明 commit 影响的范围，比如数据层、控制层、视图层等等
      ]
    ]
  }
};
