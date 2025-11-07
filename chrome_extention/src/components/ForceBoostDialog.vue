<template>
  <var-popup
    v-model:show="visible"
    position="center"
    :close-on-click-overlay="true"
    :overlay-style="{ background: 'rgba(0, 0, 0, 0.4)' }"
    class="macos-dialog-popup"
  >
    <div class="macos-dialog">
      <!-- Icon -->
      <div class="dialog-icon-wrapper">
        <svg
          class="dialog-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path
            d="M12 2L4 6V11C4 16.55 7.84 21.74 12 23C16.16 21.74 20 16.55 20 11V6L12 2Z"
          />
          <path d="M12 8V13" stroke-width="2" />
          <circle cx="12" cy="16" r="1" fill="currentColor" />
        </svg>
      </div>

      <!-- Content -->
      <div class="dialog-content">
        <h3 class="dialog-title">强制开启加速</h3>
        <p class="dialog-message">
          该网站未被识别为支持的加速节点。 强制开启可能无法正常工作。
          是否确认开启?
        </p>
      </div>

      <!-- Actions -->
      <div class="dialog-actions">
        <button
          @click="handleCancel"
          class="dialog-button dialog-button-secondary"
        >
          取消
        </button>
        <button
          @click="handleConfirm"
          class="dialog-button dialog-button-primary"
        >
          确认开启
        </button>
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
.macos-dialog-popup :deep(.var-popup__content) {
  background: transparent;
}

.macos-dialog {
  width: 300px;
  background: var(--macos-glass-heavy);
  backdrop-filter: saturate(180%) blur(40px);
  -webkit-backdrop-filter: saturate(180%) blur(40px);
  border-radius: var(--macos-radius-xl);
  border: 1px solid var(--macos-separator-light);
  padding: var(--macos-space-2xl);
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: var(--macos-shadow-xl);
  animation: dialogAppear 0.3s var(--macos-spring);
}

@keyframes dialogAppear {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.dialog-icon-wrapper {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: color-mix(in srgb, var(--macos-warning) 15%, transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--macos-space-lg);
  animation: iconBounce 0.5s var(--macos-spring);
}

@keyframes iconBounce {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.15);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.dialog-icon {
  width: 32px;
  height: 32px;
  color: var(--macos-warning);
}

.dialog-content {
  text-align: center;
  margin-bottom: var(--macos-space-xl);
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.3px;
  color: var(--macos-text-primary);
  margin: 0 0 var(--macos-space-sm);
}

.dialog-message {
  font-size: 13px;
  line-height: 1.5;
  color: var(--macos-text-secondary);
  margin: 0;
  font-weight: 400;
}

.dialog-actions {
  width: 100%;
  display: flex;
  gap: var(--macos-space-sm);
}

.dialog-button {
  flex: 1;
  height: 40px;
  border-radius: var(--macos-radius-md);
  border: none;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--macos-transition-fast);
}

.dialog-button-secondary {
  background: var(--macos-bg-secondary);
  color: var(--macos-text-primary);
}

.dialog-button-secondary:hover {
  background: var(--macos-bg-tertiary);
  transform: translateY(-1px);
}

.dialog-button-secondary:active {
  transform: translateY(0);
}

.dialog-button-primary {
  background: var(--macos-warning);
  color: white;
  box-shadow: 0 2px 8px
    color-mix(in srgb, var(--macos-warning) 30%, transparent);
}

.dialog-button-primary:hover {
  background: color-mix(in srgb, var(--macos-warning) 90%, white);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px
    color-mix(in srgb, var(--macos-warning) 40%, transparent);
}

.dialog-button-primary:active {
  transform: translateY(0);
}
</style>
