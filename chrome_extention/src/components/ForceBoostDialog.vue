<template>
  <var-popup
    v-model:show="visible"
    position="center"
    :close-on-click-overlay="true"
    :overlay-style="{ background: 'rgba(0, 0, 0, 0.5)' }"
    class="force-boost-dialog"
  >
    <div class="dialog-container">
      <!-- 图标 -->
      <div class="dialog-icon">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
          <path
            d="M12 2L4 6V11C4 16.55 7.84 21.74 12 23C16.16 21.74 20 16.55 20 11V6L12 2Z"
            stroke="#FF9500"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
          />
          <path
            d="M12 8V13"
            stroke="#FF9500"
            stroke-width="2"
            stroke-linecap="round"
          />
          <circle cx="12" cy="16" r="1" fill="#FF9500" />
        </svg>
      </div>

      <!-- 标题 -->
      <h3 class="dialog-title">强制开启加速</h3>

      <!-- 描述 -->
      <p class="dialog-description">
        该网站未被识别为支持的加速节点。<br />
        强制开启可能无法正常工作。<br />
        是否确认开启?
      </p>

      <!-- 按钮组 -->
      <div class="dialog-actions">
        <var-button
          type="default"
          size="large"
          block
          @click="handleCancel"
          class="cancel-button"
        >
          取消
        </var-button>
        <var-button
          type="primary"
          size="large"
          block
          @click="handleConfirm"
          class="confirm-button"
        >
          确认开启
        </var-button>
      </div>
    </div>
  </var-popup>
</template>

<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  show: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["update:show", "confirm", "cancel"]);

const visible = ref(props.show);

watch(
  () => props.show,
  (newVal) => {
    visible.value = newVal;
  }
);

watch(visible, (newVal) => {
  emit("update:show", newVal);
});

const handleConfirm = () => {
  emit("confirm");
  visible.value = false;
};

const handleCancel = () => {
  emit("cancel");
  visible.value = false;
};
</script>

<style scoped>
.force-boost-dialog :deep(.var-popup__content) {
  background: transparent;
}

.dialog-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: saturate(180%) blur(20px);
  border-radius: 20px;
  padding: 32px 24px 24px;
  width: 320px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dark .dialog-container {
  background: rgba(44, 44, 46, 0.95);
}

.dialog-icon {
  margin-bottom: 16px;
  animation: bounce 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

@keyframes bounce {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.dialog-title {
  font-size: 20px;
  font-weight: 700;
  color: #1d1d1f;
  margin: 0 0 12px;
  letter-spacing: -0.3px;
  text-align: center;
}

.dark .dialog-title {
  color: #f5f5f7;
}

.dialog-description {
  font-size: 14px;
  color: #86868b;
  line-height: 1.6;
  text-align: center;
  margin: 0 0 24px;
}

.dialog-actions {
  width: 100%;
  display: flex;
  gap: 12px;
}

.cancel-button {
  background: rgba(0, 0, 0, 0.05) !important;
  color: #1d1d1f !important;
  border: none !important;
  font-weight: 500 !important;
  transition: all 0.2s ease !important;
}

.dark .cancel-button {
  background: rgba(255, 255, 255, 0.1) !important;
  color: #f5f5f7 !important;
}

.cancel-button:hover {
  background: rgba(0, 0, 0, 0.1) !important;
  transform: translateY(-1px);
}

.cancel-button:active {
  transform: translateY(0);
}

.confirm-button {
  background: linear-gradient(135deg, #ff9500 0%, #ff9f0a 100%) !important;
  border: none !important;
  font-weight: 600 !important;
  box-shadow: 0 4px 12px rgba(255, 149, 0, 0.3) !important;
  transition: all 0.2s ease !important;
}

.confirm-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(255, 149, 0, 0.4) !important;
}

.confirm-button:active {
  transform: translateY(0);
}

/* 按钮组响应式 */
@media (max-width: 360px) {
  .dialog-container {
    width: 280px;
    padding: 24px 20px 20px;
  }

  .dialog-actions {
    flex-direction: column;
  }
}
</style>
