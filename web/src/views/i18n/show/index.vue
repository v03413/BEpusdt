<template>
  <div class="snow-page">
    <div class="snow-inner">
      <p>{{ $t(`system.switch-language-to-preview`) }}</p>
      <br />
      <div>
        <a-date-picker style="width: 200px" v-model="form.time" />
      </div>
      <br />
      <div>
        <a-time-picker type="time-range" style="width: 252px" v-model="form.timeRange" />
      </div>
      <br />
      <div>
        <a-range-picker style="width: 360px" show-time format="YYYY-MM-DD HH:mm" v-model="form.date" />
      </div>
      <br />
      <div>
        <a-pagination :total="50" show-total show-jumper show-page-size />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const form = reactive({
  time: "",
  timeRange: [],
  date: []
});
</script>

<style lang="scss" scoped></style>
