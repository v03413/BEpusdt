import { ref } from "vue";

export interface WalletDetail {
  id: number;
  name: string;
  address: string;
  trade_type: string;
  remark: string;
  other_notify: number;
  status: number;
  created_at?: string;
  updated_at?: string;
}

export const useWalletDetail = () => {
  const detailVisible = ref(false);
  const detailData = ref<WalletDetail>({
    id: 0,
    name: "",
    address: "",
    trade_type: "",
    remark: "",
    other_notify: 0,
    status: 0
  });

  // 显示详情
  const showDetail = (record: WalletDetail) => {
    detailData.value = { ...record };
    detailVisible.value = true;
  };

  // 关闭详情
  const closeDetail = () => {
    detailVisible.value = false;
    detailData.value = {
      id: 0,
      name: "",
      address: "",
      trade_type: "",
      remark: "",
      other_notify: 0,
      status: 0
    };
  };

  return {
    detailVisible,
    detailData,
    showDetail,
    closeDetail
  };
};
