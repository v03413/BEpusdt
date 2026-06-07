// 获取法币对应的旗帜emoji
const getFiatFlag = (fiat: string) => {
  const fiatFlagMap: Record<string, string> = {
    CNY: "🇨🇳",
    USD: "🇺🇸",
    JPY: "🇯🇵",
    EUR: "🇪🇺",
    GBP: "🇬🇧"
  };
  return fiatFlagMap[fiat] || "🌍";
};

export const getCryptoColor = (crypto: string): string => {
  const colorMap: Record<string, string> = {
    USDT: "green",
    USDC: "blue",
    TRX: "red",
    ETH: "purple",
    BNB: "orange",
    GRAM: "#0088CC"
  };
  return colorMap[crypto] ?? "gray";
};

export { getFiatFlag };
