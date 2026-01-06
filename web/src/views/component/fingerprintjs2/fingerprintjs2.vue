<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-card title="浏览器检测-通过设备浏览器信息获取浏览器指纹">
        <a-descriptions :column="1" :data="browserInfo" bordered />
      </a-card>
      <br />
      <a-card title="移动设备">
        <a-descriptions :column="1" :data="mobileFinger" bordered />
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import Fingerprint2 from "fingerprintjs2";

interface IComponents {
  key: string;
  value: any;
}
const browserInfo = ref<any>([]);
const getBrowserInfo = () => {
  // 浏览器指纹
  Fingerprint2.get((components: any) => {
    // 参数只有回调函数时，默认浏览器指纹依据所有配置信息进行生成
    const values = components.map((component: any) => component.value); // 配置的值的数组
    const murmur = Fingerprint2.x64hash128(values.join(""), 31).toUpperCase(); // 生成浏览器指纹
    let arr = [];
    components.forEach((el: IComponents) => {
      switch (el.key) {
        case "userAgent":
          arr.push({ label: "用户代理(User-Agent)", value: el.value });
          break;
        case "language":
          arr.push({ label: "语言(Language)", value: el.value });
          break;
        case "colorDepth":
          arr.push({ label: "屏幕颜色深度(ColorDepth)", value: el.value });
          break;
        case "deviceMemory":
          arr.push({ label: "设备内存(DviceMemory)", value: `此设备至少拥有 ${el.value}GiB 的 RAM` });
          break;
        case "hardwareConcurrency":
          arr.push({ label: "硬件并发(HardwareConcurrency)", value: el.value });
          break;
        case "screenResolution":
          arr.push({ label: "屏幕分辨率(ScreenResolution)", value: el.value.join("*") });
          break;
        case "availableScreenResolution":
          arr.push({ label: "可用屏幕分辨率(AvailableScreenResolution)", value: el.value.join("*") });
          break;
        case "timezoneOffset":
          arr.push({ label: "到UTC的时区差异-分钟(TimezoneOffset)", value: el.value });
          break;
        case "timezone":
          arr.push({ label: "时区(Timezone)", value: el.value });
          break;
        case "sessionStorage":
          arr.push({ label: "会话存储支持(SessionStorage)", value: el.value ? "是(true)" : "否(false)" });
          break;
        case "localStorage":
          arr.push({ label: "本地存储支持(LocalStorage)", value: el.value ? "是(true)" : "否(false)" });
          break;
        case "indexedDb":
          arr.push({ label: "本地数据库支持(IndexedDb)", value: el.value ? "是(true)" : "否(false)" });
          break;
        case "openDatabase":
          arr.push({ label: "是否支持Web SQL(OpenDatabase)", value: el.value ? "是(true)" : "否(false)" });
          break;
        case "platform":
          arr.push({ label: "系统平台(Platform)", value: el.value });
          break;
        default:
          break;
      }
    });
    arr.push({ label: "浏览器指纹(BrowserFingerprint)", value: murmur });
    browserInfo.value = arr;
  });
};

const mobileFinger = ref<any>([]);
const getMobileFinger = () => {
  Fingerprint2.get((components: any) => {
    // 由于Fingerprint2是通过浏览器设备信息生成的浏览器指纹，
    // 大多数用户不会修改手机配置，同型号产品可能会出现指纹相同，这里需要加上唯一标识
    components.push({
      key: "ip",
      value: "192.168.1.1" //通过接口获取的到ip
    });
    const values = components.map((component: any) => component.value); // 配置的值的数组
    const murmur = Fingerprint2.x64hash128(values.join(""), 31).toUpperCase(); // 生成浏览器指纹
    mobileFinger.value.push({ label: "移动设备指纹(需要唯一标识)", value: murmur });
  });
};

getBrowserInfo();
getMobileFinger();
</script>
