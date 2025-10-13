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

const getCryptoColor = (crypto: string) => {
  const cryptoColorMap: Record<string, string> = {
    USDT: "blue",
    USDC: "green",
    TRX: "red"
  };
  return cryptoColorMap[crypto] || "#0fc6c2";
};

export { getFiatFlag, getCryptoColor };
