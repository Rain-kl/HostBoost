<template>
  <div>
    <h2 class="text-2xl">当前域名</h2>
    <p v-if="domain">{{ domain }}</p>
    <p v-else>正在获取...</p>

  </div>
</template>

<script setup>
import {ref, onMounted, watch} from "vue";
import {hostApi} from "@/api/api-ref.js";

const domain = ref("");

const getHost = async (domain) => {
  const rsp = await hostApi.hostGet(domain)
  console.log(rsp);
}

onMounted(() => {
  chrome.tabs.query({active: true, currentWindow: true}, (tabs) => {
    const tab = tabs[0];
    if (tab?.url) {
      try {
        domain.value = new URL(tab.url).hostname;
      } catch {
        domain.value = "无法解析域名";
      }
    } else {
      domain.value = "未获取到当前标签页";
    }
  });
});

watch(domain, (newVal, oldVal) => {
  console.log('count 发生变化:', oldVal, '→', newVal)
  getHost(newVal)
})
</script>