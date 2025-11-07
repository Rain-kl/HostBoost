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
/* macOS 15 Detail Item */
.detail-item {
  display: flex;
  align-items: center;
  gap: var(--macos-space-md);
  padding: var(--macos-space-md) var(--macos-space-lg);
  background: var(--macos-glass-light);
  backdrop-filter: saturate(180%) blur(20px);
  -webkit-backdrop-filter: saturate(180%) blur(20px);
  border-radius: var(--macos-radius-lg);
  border: 1px solid var(--macos-separator-light);
  transition: all var(--macos-transition-fast);
  cursor: default;
}

.detail-item:hover {
  background: var(--macos-glass-medium);
  border-color: var(--macos-separator-medium);
  transform: translateY(-1px);
  box-shadow: var(--macos-shadow-sm);
}

.detail-item:active {
  transform: scale(0.98);
}

.detail-icon {
  font-size: 22px;
  line-height: 1;
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--macos-bg-secondary);
  border-radius: var(--macos-radius-md);
}

.detail-content {
  flex: 1;
  min-width: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--macos-space-md);
}

.detail-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--macos-text-primary);
  flex-shrink: 0;
  letter-spacing: -0.1px;
}

.detail-value {
  font-size: 13px;
  font-weight: 400;
  color: var(--macos-text-secondary);
  text-align: right;
  word-break: break-word;
  font-family: "SF Mono", Monaco, "Courier New", monospace;
  line-height: 1.4;
}
</style>
