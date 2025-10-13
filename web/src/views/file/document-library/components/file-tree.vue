<template>
  <div class="container-page">
    <a-input class="search-tree" v-model="searchKey" placeholder="请输入搜索关键词">
      <template #prefix>
        <icon-search />
      </template>
    </a-input>
    <div class="tree-box">
      <a-spin class="tree-loading" :loading="loading" tip="加载中...">
        <!-- TODO: 滚动条边距和宽度需要调整 -->
        <a-scrollbar style="height: 100%; overflow: auto" outer-class="scrollbar">
          <a-tree ref="aTreeRef" :data="treeData" :default-expand-all="true" :show-line="true" @select="onNode">
            <template #title="node">
              <span class="tree-title">{{ node.title }}</span>
            </template>
            <template #icon="node">
              <s-svg-icon name="folder-close" :size="16" v-if="!node.isLeaf && !node.expanded"></s-svg-icon>
              <s-svg-icon name="folder-open" :size="16" v-if="!node.isLeaf && node.expanded"></s-svg-icon>
              <s-svg-icon name="txt" :size="16" v-if="node.isLeaf"></s-svg-icon>
            </template>
          </a-tree>
        </a-scrollbar>
      </a-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { findPathOfParentNode } from "@/utils/tree-tools";
import { getDocumentLibraryTreeAPI } from "@/api/modules/file/index";

const emit = defineEmits(["onNode"]);

const onNode = (selectedKeys: Array<string>) => {
  let list = findPathOfParentNode(treeData.value, "key", selectedKeys[0]);
  emit("onNode", list || []);
};

const searchKey = ref<string>("");
const treeData = computed(() => {
  if (!searchKey.value) return sourceTree.value;
  return searchData(searchKey.value);
});

// 搜索树
const searchData = (keyword: string) => {
  const loop = (tree: SourceTree[]) => {
    const result: SourceTree[] = []; // 存储搜索结果
    tree.forEach((item: SourceTree) => {
      // 匹配节点
      if (item.title.toLowerCase().indexOf(keyword.toLowerCase()) > -1) {
        result.push({ ...item });
      } else if (item.children) {
        // 不匹配，但是有children，深层遍历
        // filterData 存储遍历结果
        const filterData = loop(item.children);
        // 如果在深层遍历中，匹配到
        if (filterData.length) {
          // 将结果存储到children中
          result.push({
            ...item,
            children: filterData
          });
        }
      }
    });
    return result;
  };
  // 开始递归
  return loop(sourceTree.value);
};

interface SourceTree {
  title: string;
  key: string;
  children?: SourceTree[];
}
const aTreeRef = ref();
const loading = ref<boolean>(false);
const sourceTree = ref<SourceTree[]>([]);
const getDocumentLibraryTree = async () => {
  try {
    loading.value = true;
    let { data } = await getDocumentLibraryTreeAPI();
    loading.value = false;
    sourceTree.value = data;
    nextTick(() => {
      aTreeRef.value.expandAll(true);
    });
  } finally {
    loading.value = false;
  }
};
getDocumentLibraryTree();
</script>

<style lang="scss" scoped>
.container-page {
  box-sizing: border-box;
  height: 100%;
  padding: $padding;
  .search-tree {
    height: 32px;
    margin-bottom: $margin;
  }
  .tree-box {
    height: calc(100% - 32px - $margin);
    .tree-loading {
      width: 100%;
      height: 100%;
    }
    .scrollbar {
      height: 100%;
      .tree-title {
        white-space: nowrap;
      }
    }
  }
}
</style>
