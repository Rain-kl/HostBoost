<template>
  <div class="detail-item">
    <div class="detail-icon">{{ icon }}</div>
    <div class="detail-content">
      <span class="detail-label">{{ label }}</span>
      <span class="detail-value">{{ displayValue }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  icon: {
    type: String,
    required: true,
  },
  label: {
    type: String,
    required: true,
  },
  value: {
    type: [String, Number],
    default: "",
  },
});

const displayValue = computed(() => {
  if (props.value === undefined || props.value === null || props.value === "") {
    return "暂无数据";
  }
  return String(props.value);
});
</script>

<style scoped>
.detail-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 14px;
  border: 0.5px solid rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
}

.dark .detail-item {
  background: rgba(58, 58, 60, 0.6);
  border-color: rgba(255, 255, 255, 0.06);
}

.detail-item:active {
  background: rgba(0, 0, 0, 0.05);
  transform: scale(0.98);
}

.dark .detail-item:active {
  background: rgba(255, 255, 255, 0.1);
}

.detail-icon {
  font-size: 24px;
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.detail-content {
  flex: 1;
  min-width: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.detail-label {
  font-size: 15px;
  color: #1d1d1f;
  font-weight: 500;
  flex-shrink: 0;
}

.dark .detail-label {
  color: #f5f5f7;
}

.detail-value {
  font-size: 14px;
  color: #86868b;
  text-align: right;
  word-break: break-all;
  font-family: "SF Mono", Monaco, "Courier New", monospace;
}

.dark .detail-value {
  color: #98989d;
}
</style>
