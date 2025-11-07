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
/* macOS 15 Bottom Drawer */
.web-details-popup :deep(.var-popup) {
  border-radius: var(--macos-radius-2xl) var(--macos-radius-2xl) 0 0;
  background: var(--macos-bg-primary);
  max-height: 75vh;
}

.popup-content {
  padding: 0;
  display: flex;
  flex-direction: column;
}

/* Drag Handle */
.popup-handle {
  padding: var(--macos-space-md) 0;
  display: flex;
  justify-content: center;
  cursor: grab;
}

.popup-handle:active {
  cursor: grabbing;
}

.handle-bar {
  width: 40px;
  height: 4px;
  background: var(--macos-text-tertiary);
  border-radius: var(--macos-radius-sm);
  opacity: 0.4;
  transition: all var(--macos-transition-fast);
}

.popup-handle:hover .handle-bar {
  opacity: 0.6;
  width: 48px;
}

/* Header */
.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--macos-space-xl) var(--macos-space-lg);
  border-bottom: 1px solid var(--macos-separator-light);
}

.popup-header h2 {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.3px;
  color: var(--macos-text-primary);
  margin: 0;
}

.close-button {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--macos-bg-secondary);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--macos-text-secondary);
  cursor: pointer;
  transition: all var(--macos-transition-fast);
}

.close-button:hover {
  background: var(--macos-bg-tertiary);
  color: var(--macos-text-primary);
  transform: scale(1.05);
}

.close-button:active {
  transform: scale(0.95);
}

/* Body */
.popup-body {
  flex: 1;
  padding: var(--macos-space-xl);
  overflow-y: auto;
  overscroll-behavior: contain;
}

/* States */
.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--macos-space-3xl) var(--macos-space-xl);
  gap: var(--macos-space-lg);
  min-height: 200px;
}

.loading-state p,
.empty-state p {
  font-size: 14px;
  color: var(--macos-text-secondary);
  margin: 0;
  font-weight: 400;
}

.error-icon,
.empty-icon {
  font-size: 48px;
  opacity: 0.8;
}

.error-message {
  font-size: 14px;
  color: var(--macos-error);
  text-align: center;
  margin: 0;
  line-height: 1.5;
}

/* Details List */
.details-list {
  display: flex;
  flex-direction: column;
  gap: var(--macos-space-sm);
}

/* Scrollbar */
.popup-body::-webkit-scrollbar {
  width: 6px;
}

.popup-body::-webkit-scrollbar-track {
  background: transparent;
}

.popup-body::-webkit-scrollbar-thumb {
  background: var(--macos-text-tertiary);
  border-radius: var(--macos-radius-sm);
  opacity: 0.3;
}

.popup-body::-webkit-scrollbar-thumb:hover {
  opacity: 0.5;
}
</style>
