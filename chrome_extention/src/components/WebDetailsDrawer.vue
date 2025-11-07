<template>
  <var-popup
    v-model:show="isShow"
    position="bottom"
    :close-on-click-overlay="true"
    :safe-area-inset-bottom="true"
    class="web-details-popup"
  >
    <div class="popup-content">
      <!-- 拖动条 -->
      <div class="popup-handle">
        <div class="handle-bar"></div>
      </div>

      <!-- 弹窗标题 -->
      <div class="popup-header">
        <h2>网站信息</h2>
        <button @click="close" class="close-button">
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>

      <!-- 内容区域 -->
      <div class="popup-body">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <var-loading type="wave" :size="32" />
          <p>正在获取信息...</p>
        </div>

        <!-- 错误状态 -->
        <div v-else-if="error" class="error-state">
          <span class="error-icon">⚠️</span>
          <p class="error-message">{{ error }}</p>
          <var-button type="primary" size="small" @click="handleRetry"
            >重试</var-button
          >
        </div>

        <!-- 网站信息列表 -->
        <div v-else-if="details" class="details-list">
          <DetailItem
            v-for="item in displayItems"
            :key="item.key"
            :icon="item.icon"
            :label="item.label"
            :value="item.value"
          />

          <div v-if="displayItems.length === 0" class="empty-state">
            <span class="empty-icon">📭</span>
            <p>暂无可显示的信息</p>
          </div>
        </div>
      </div>
    </div>
  </var-popup>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { toolApi } from "@/api/api-ref.ts";
import DetailItem from "./DetailItem.vue";

const props = defineProps({
  show: {
    type: Boolean,
    required: true,
  },
  domain: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["update:show"]);

// 状态管理
const details = ref(null);
const loading = ref(false);
const error = ref("");

// 双向绑定 show
const isShow = computed({
  get: () => props.show,
  set: (value) => emit("update:show", value),
});

// 获取网站详情
const fetchWebDetails = async () => {
  if (
    !props.domain ||
    props.domain === "无法解析域名" ||
    props.domain === "未获取到当前标签页"
  ) {
    error.value = "无效的域名";
    return;
  }

  loading.value = true;
  error.value = "";
  details.value = null;

  try {
    const response = await toolApi.toolWebDetailsGet(props.domain);
    if (
      (response.data.code === 200 || response.data.code === "200") &&
      response.data.data
    ) {
      details.value = response.data.data;
    } else {
      error.value = response.data.message || "获取网站信息失败";
    }
  } catch (err) {
    console.error("获取网站信息失败:", err);
    error.value = err.message || "网络请求失败";
  } finally {
    loading.value = false;
  }
};

// 重试处理
const handleRetry = () => {
  fetchWebDetails();
};

// 关闭抽屉
const close = () => {
  isShow.value = false;
};

// 当打开弹窗时自动获取网站信息
watch(
  () => props.show,
  (newVal) => {
    if (newVal && !details.value && !loading.value) {
      fetchWebDetails();
    }
  }
);

// 网站详情展示数据（处理字段不存在的情况）
const displayItems = computed(() => {
  if (!details.value) return [];

  const data = details.value;
  const items = [
    { key: "ip", icon: "🌐", label: "IP 地址", value: data.ip },
    { key: "country", icon: "🌍", label: "国家", value: data.country },
    {
      key: "country_code",
      icon: "🏳️",
      label: "国家代码",
      value: data.country_code,
    },
    { key: "region", icon: "📍", label: "地区", value: data.region },
    {
      key: "region_code",
      icon: "🗺️",
      label: "地区代码",
      value: data.region_code,
    },
    { key: "city", icon: "🏙️", label: "城市", value: data.city },
    {
      key: "organization",
      icon: "🏢",
      label: "组织",
      value: data.organization,
    },
    { key: "isp", icon: "📡", label: "ISP", value: data.isp },
    { key: "asn", icon: "🔢", label: "ASN", value: data.asn },
    {
      key: "asn_organization",
      icon: "🏛️",
      label: "ASN 组织",
      value: data.asn_organization,
    },
    { key: "timezone", icon: "🕐", label: "时区", value: data.timezone },
    {
      key: "offset",
      icon: "⏱️",
      label: "时区偏移",
      value: data.offset ? `UTC+${data.offset / 3600}` : undefined,
    },
    { key: "latitude", icon: "🧭", label: "纬度", value: data.latitude },
    { key: "longitude", icon: "🧭", label: "经度", value: data.longitude },
    {
      key: "continent_code",
      icon: "🌏",
      label: "洲代码",
      value: data.continent_code,
    },
  ];

  // 过滤掉值为 undefined, null, 或空字符串的项
  return items.filter((item) => {
    const value = item.value;
    return value !== undefined && value !== null && value !== "";
  });
});
</script>

<style scoped>
/* 底部弹窗样式 */
.web-details-popup :deep(.var-popup) {
  border-radius: 20px 20px 0 0;
  background: #f5f5f7;
  max-height: 70vh;
}

.dark .web-details-popup :deep(.var-popup) {
  background: #1c1c1e;
}

.popup-content {
  padding: 0;
}

.popup-handle {
  padding: 8px 0 12px;
  display: flex;
  justify-content: center;
}

.handle-bar {
  width: 36px;
  height: 5px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.dark .handle-bar {
  background: rgba(255, 255, 255, 0.3);
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px 16px;
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.1);
}

.dark .popup-header {
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.popup-header h2 {
  font-size: 20px;
  font-weight: 700;
  color: #1d1d1f;
  margin: 0;
  letter-spacing: -0.3px;
}

.dark .popup-header h2 {
  color: #f5f5f7;
}

.close-button {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.05);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #86868b;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dark .close-button {
  background: rgba(255, 255, 255, 0.1);
  color: #98989d;
}

.close-button:hover {
  background: rgba(0, 0, 0, 0.1);
  transform: scale(1.05);
}

.close-button:active {
  transform: scale(0.95);
}

.popup-body {
  padding: 20px;
  max-height: calc(70vh - 80px);
  overflow-y: auto;
}

/* 加载/错误/空状态 */
.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 20px;
  gap: 16px;
}

.loading-state p,
.empty-state p {
  font-size: 14px;
  color: #86868b;
  margin: 0;
}

.error-icon,
.empty-icon {
  font-size: 48px;
}

.error-message {
  font-size: 14px;
  color: #ff3b30;
  text-align: center;
  margin: 0;
}

.dark .error-message {
  color: #ff453a;
}

/* 详情列表 */
.details-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 滚动条样式 */
.popup-body::-webkit-scrollbar {
  width: 6px;
}

.popup-body::-webkit-scrollbar-track {
  background: transparent;
}

.popup-body::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.dark .popup-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}
</style>
