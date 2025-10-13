<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-alert type="success">工具类位置: src/utils/tree-tools.ts</a-alert>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>判断code在树形数据中是否存在: treeCodeExist</div>
        <div>返回true / false，示例: treeCodeExist(treeData, "ABC123", { key: 'id', childrenKey: "children" })</div>
        <div>
          使用方式：
          <a-typography-text code
            >treeCodeExist(树数据, key值, { key: key名, childrenKey: 下级绑定的key比如children })</a-typography-text
          >
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>返回指定名称的所有父节点: findPathOfParentNode</div>
        <div>示例: findPathOfParentNode(treeData,key, keyValue)</div>
        <div>
          使用方式：
          <a-typography-text code>findPathOfParentNode(树数据, 指定的key, 指定key绑定的值)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>数组扁平化-多维转一维: arrayFlattened</div>
        <div>返回转一维的数组，示例: arrayFlattened(treeData, 'children')</div>
        <div>
          使用方式：
          <a-typography-text code>arrayFlattened(树数据, 下级的键名如children)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>乱序一维数组转树形结构: buildTreeOptimized</div>
        <div>返回根据parentId将乱序的一维数组转树形结构，示例: buildTreeOptimized(nodes)</div>
        <div>
          使用方式：
          <a-typography-text code>buildTreeOptimized(乱序的一维数组)</a-typography-text>
        </div>
      </a-space>
      <a-divider />
      <a-space direction="vertical" fill>
        <div>根据指定key找到对应节点: findCategoryById</div>
        <div>返回指定id找到对应的节点，示例: findCategoryById(tree, key, keyValue)</div>
        <div>
          使用方式：
          <a-typography-text code>findCategoryById(树数据, 指定的key, 指定的key绑定的值)</a-typography-text>
        </div>
      </a-space>
    </div>
  </div>
</template>
