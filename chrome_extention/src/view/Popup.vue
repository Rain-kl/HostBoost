<template>
  <div class="popup-container">
    <!-- 顶部栏 -->
    <header class="header">
      <div class="logo-container">
        <div
          class="logo-ring"
          :class="{ 'logo-ring-active': isDetecting }"
        ></div>
        <span class="logo-text">HostBoost</span>
      </div>
    </header>

    <!-- 主内容区域 -->
    <main class="main-content">
      <!-- 加速状态卡片 (合并 CDN 优选) -->
      <div class="card boost-card">
        <h3 class="card-title">加速状态</h3>

        <!-- 盾牌控制 -->
        <div class="shield-container">
          <button
            class="shield-button"
            :class="{
              'shield-active': isBoostEnabled,
              'shield-disabled': !isBoostSupported && !isBackendError,
              'shield-warning': isBackendError,
            }"
            :disabled="isDetecting"
            @click="toggleBoost"
          >
            <svg
              class="shield-icon"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M12 2L4 6V11C4 16.55 7.84 21.74 12 23C16.16 21.74 20 16.55 20 11V6L12 2Z"
                :fill="isBoostEnabled ? 'currentColor' : 'none'"
                :stroke="isBoostEnabled ? 'none' : 'currentColor'"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                v-if="isBoostEnabled"
                d="M9 12L11 14L15 10"
                stroke="white"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
          <p class="shield-status">
            {{ getShieldStatusText() }}
          </p>
        </div>

        <!-- CDN 优选信息 - 仅在加速开启时显示 -->
        <transition name="slide-fade">
          <div v-if="isBoostEnabled" class="cdn-section">
            <div class="cdn-divider"></div>
            <h4 class="cdn-subtitle">CDN 优选</h4>
            <div class="cdn-info">
              <div class="cdn-detail">
                <span class="cdn-label">最优节点</span>
                <span class="cdn-value">{{ optimizedNode.ip }}</span>
              </div>
              <!--              <div class="cdn-detail">-->
              <!--                <span class="cdn-label">响应时间</span>-->
              <!--                <span class="cdn-value cdn-rtt">{{ optimizedNode.rtt }}ms</span>-->
              <!--              </div>-->
            </div>
            <!--            <div class="latency-bar-container">-->
            <!--              <div-->
            <!--                class="latency-bar"-->
            <!--                :style="{ width: `${latencyPercentage}%` }"-->
            <!--                :class="latencyClass"-->
            <!--              ></div>-->
            <!--            </div>-->
            <!--            <button-->
            <!--              class="reoptimize-button"-->
            <!--              @click="reoptimize"-->
            <!--              :disabled="isOptimizing"-->
            <!--            >-->
            <!--              <span v-if="!isOptimizing">重新优选</span>-->
            <!--              <span v-else>优选中... {{ countdown }}s</span>-->
            <!--            </button>-->
          </div>
        </transition>
      </div>
    </main>

    <!-- 粒子背景 -->
    <div class="particles">
      <div
        v-for="i in 20"
        :key="i"
        class="particle"
        :style="getParticleStyle(i)"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from "vue";
import { hostApi } from "@/api/api-ref.js";

// 状态管理
const domain = ref("");
const isDetecting = ref(true);
const isBoostEnabled = ref(false);
const isBoostSupported = ref(true); // 是否支持加速，默认为true
const isBackendError = ref(false); // 后端服务错误状态
const isOptimizing = ref(false);
const countdown = ref(3);
const currentTabId = ref(undefined); // 当前标签页ID

const detectStatus = ref({
  icon: "🔍",
  text: "正在识别...",
});

const optimizedNode = ref({
  ip: "",
  rtt: 0,
});

// 计算延迟百分比和颜色
const latencyPercentage = computed(() => {
  const rtt = optimizedNode.value.rtt;
  return Math.min((rtt / 200) * 100, 100);
});

const latencyClass = computed(() => {
  const rtt = optimizedNode.value.rtt;
  if (rtt < 50) return "latency-excellent";
  if (rtt < 100) return "latency-good";
  if (rtt < 150) return "latency-fair";
  return "latency-poor";
});

// 检测域名是否支持CDN加速（预留接口，当前版本返回true）
const checkCdnSupport = (domain) => {
  // TODO: 后续版本实现真实的CDN检测逻辑
  // 可以检测域名是否使用Cloudflare、Akamai等CDN服务
  return true;
};

// API 调用 - 检查域名状态
const getHost = async (domain) => {
  try {
    // 先调用 hostGet 接口查询状态
    const response = await hostApi.hostGet(domain);

    isDetecting.value = false;
    isBackendError.value = false; // 能收到响应，清除后端错误状态

    // 检测域名是否支持加速
    isBoostSupported.value = checkCdnSupport(domain);

    // 如果查询成功(code === 200)，说明已有加速记录，直接开启加速
    if (response.data.code === 200 && response.data.data) {
      isBoostEnabled.value = true;

      detectStatus.value = {
        icon: "🌐",
        text: "已识别为 Cloudflare 节点",
      };

      // 从 API 响应中获取优化节点信息
      if (response.data.data.ip) {
        optimizedNode.value = {
          ip: response.data.data.ip,
          rtt: 0,
        };
      }
    } else {
      // 查询失败或无记录（但服务端有响应）
      isBoostEnabled.value = false;

      detectStatus.value = {
        icon: "🌐",
        text: "可加速网站",
      };
    }
  } catch (error) {
    console.error("查询域名状态失败:", error);
    isDetecting.value = false;

    // 检测域名是否支持加速
    isBoostSupported.value = checkCdnSupport(domain);

    // 只有在网络错误（无法连接、超时等）时才设置后端错误状态
    // 如果error.response存在，说明服务端有响应，不是网络问题
    if (!error.response) {
      // 网络错误：ERR_CONNECTION_REFUSED, ECONNREFUSED, timeout等
      isBackendError.value = true;
      isBoostEnabled.value = false;

      detectStatus.value = {
        icon: "⚠️",
        text: "后端服务未启动",
      };
    } else {
      // 服务端有响应但返回错误（如404, 500等）
      isBackendError.value = false;
      isBoostEnabled.value = false;

      detectStatus.value = {
        icon: "❌",
        text: `服务错误: ${error.response.status}`,
      };
    }
  }
};

// 切换加速状态
const toggleBoost = async () => {
  if (isDetecting.value) {
    return;
  }

  // 如果是后端错误状态，点击后重新检查后端状态
  if (isBackendError.value) {
    isDetecting.value = true;
    isBackendError.value = false;
    await getHost(domain.value);
    return;
  }

  if (!isBoostSupported.value) {
    return;
  }

  try {
    const hostData = {
      domain: domain.value,
    };

    if (!isBoostEnabled.value) {
      // 开启加速 - 调用 hostPost
      const response = await hostApi.hostPost(hostData);

      if (response.data.code === 200) {
        isBoostEnabled.value = true;
        isBackendError.value = false; // 清除后端错误状态
        console.log("加速已开启:", response.data);

        // 再次调用 hostGet 获取完整的 CDN IP 等信息
        try {
          const getResponse = await hostApi.hostGet(domain.value);

          if (getResponse.data.code === 200 && getResponse.data.data) {
            console.log("获取 CDN 信息成功:", getResponse.data);

            // 更新优化节点信息
            if (getResponse.data.data.ip) {
              optimizedNode.value = {
                ip: getResponse.data.data.ip,
                rtt: 0,
              };
            }
          }
        } catch (getError) {
          console.error("获取 CDN 信息失败:", getError);

          // 如果 hostPost 返回了 IP，使用它作为备选
          if (response.data.data && response.data.data.ip) {
            optimizedNode.value = {
              ip: response.data.data.ip,
              rtt: 0,
            };
          }
        }

        // 等待1秒后重载当前网页，刷新DNS缓存
        setTimeout(() => {
          if (currentTabId.value) {
            chrome.tabs.reload(currentTabId.value, { bypassCache: true });
            console.log("已重载当前网页，刷新DNS缓存");
          }
        }, 1000);
      } else {
        // 服务端有响应但返回错误
        isBackendError.value = false;
        console.error("开启加速失败:", response.data.message);
      }
    } else {
      // 关闭加速 - 调用 hostDelete
      const response = await hostApi.hostDelete(hostData);

      if (response.data.code === 200) {
        isBoostEnabled.value = false;
        isBackendError.value = false; // 清除后端错误状态
        console.log("加速已关闭:", response.data);
      } else {
        // 服务端有响应但返回错误
        isBackendError.value = false;
        console.error("关闭加速失败:", response.data.message);
      }
    }
  } catch (error) {
    console.error("切换加速状态失败:", error);

    // 只有在网络错误时才设置后端错误状态
    if (!error.response) {
      isBackendError.value = true;
      isBoostEnabled.value = false;
      detectStatus.value = {
        icon: "⚠️",
        text: "后端服务未启动",
      };
    } else {
      // 服务端有响应但返回错误
      isBackendError.value = false;
      detectStatus.value = {
        icon: "❌",
        text: `操作失败: ${error.response.status}`,
      };
    }
  }
};

// 获取盾牌状态文本
const getShieldStatusText = () => {
  if (isDetecting.value) {
    return "正在识别...";
  }
  if (isBackendError.value) {
    return "点击重新检查后端服务";
  }
  if (!isBoostSupported.value) {
    return "该网站不支持加速";
  }
  return isBoostEnabled.value ? "加速已开启" : "点击盾牌开启加速";
};

// 重新优选
const reoptimize = async () => {
  if (isOptimizing.value) return;

  isOptimizing.value = true;
  countdown.value = 3;

  try {
    // 重新调用 hostPost 和 hostGet 获取最新的优化节点
    const hostData = {
      domain: domain.value,
    };

    await hostApi.hostPost(hostData);
    const response = await hostApi.hostGet(domain.value);

    if (
      response.data.code === 200 &&
      response.data.data &&
      response.data.data.ip
    ) {
      optimizedNode.value = {
        ip: response.data.data.ip,
        rtt: 0,
      };
      isBackendError.value = false; // 清除后端错误状态
    }
  } catch (error) {
    console.error("重新优选失败:", error);

    // 只有在网络错误时才设置后端错误状态
    if (!error.response) {
      isBackendError.value = true;
      detectStatus.value = {
        icon: "⚠️",
        text: "后端服务未启动",
      };
    } else {
      // 服务端有响应但返回错误
      isBackendError.value = false;
      detectStatus.value = {
        icon: "❌",
        text: `重新优选失败: ${error.response.status}`,
      };
    }
  } finally {
    isOptimizing.value = false;
  }
};

// 粒子动画样式
const getParticleStyle = (index) => {
  const x = Math.random() * 100;
  const y = Math.random() * 100;
  const delay = Math.random() * 5;
  const duration = 3 + Math.random() * 4;

  return {
    left: `${x}%`,
    top: `${y}%`,
    animationDelay: `${delay}s`,
    animationDuration: `${duration}s`,
  };
};

// 获取当前域名
onMounted(() => {
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    const tab = tabs[0];
    if (tab?.url) {
      try {
        domain.value = new URL(tab.url).hostname;
        currentTabId.value = tab.id; // 保存当前标签页ID
      } catch {
        domain.value = "无法解析域名";
      }
    } else {
      domain.value = "未获取到当前标签页";
    }
  });
});

watch(domain, (newVal) => {
  if (newVal) {
    isDetecting.value = true;
    getHost(newVal);
  }
});
</script>

<style scoped>
/* 容器基础样式 - 亮色模式 */
.popup-container {
  width: 380px;
  min-height: 400px;
  background: linear-gradient(135deg, #f0f4ff 0%, #e5edff 50%, #fef3f2 100%);
  color: #1f2937;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display", "Segoe UI",
    sans-serif;
  position: relative;
  overflow: hidden;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  transition: min-height 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* 动态光感背景 */
.popup-container::before {
  content: "";
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(
      circle at 30% 30%,
      rgba(102, 126, 234, 0.15) 0%,
      transparent 50%
    ),
    radial-gradient(
      circle at 70% 70%,
      rgba(245, 158, 11, 0.12) 0%,
      transparent 50%
    ),
    radial-gradient(
      circle at 50% 50%,
      rgba(139, 92, 246, 0.08) 0%,
      transparent 60%
    );
  animation: lightFlow 15s ease-in-out infinite;
  pointer-events: none;
}

@keyframes lightFlow {
  0%,
  100% {
    transform: translate(0, 0) rotate(0deg);
    opacity: 0.8;
  }
  33% {
    transform: translate(10%, 10%) rotate(120deg);
    opacity: 1;
  }
  66% {
    transform: translate(-10%, 5%) rotate(240deg);
    opacity: 0.9;
  }
}

/* 顶部栏 - 增强毛玻璃效果 */
.header {
  padding: 20px 24px;
  backdrop-filter: blur(30px) saturate(180%);
  -webkit-backdrop-filter: blur(30px) saturate(180%);
  background: rgba(255, 255, 255, 0.7);
  border-bottom: 1px solid rgba(102, 126, 234, 0.2);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  position: relative;
  z-index: 2;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-ring {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #f59e0b 100%);
  position: relative;
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.logo-ring::after {
  content: "";
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #f59e0b 100%);
  opacity: 0.3;
  filter: blur(8px);
  z-index: -1;
}

.logo-ring-active {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.05);
  }
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.5px;
  background: linear-gradient(135deg, #667eea 0%, #f59e0b 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 主内容区域 */
.main-content {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
  z-index: 1;
}

/* 卡片通用样式 - 增强毛玻璃效果 */
.card {
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(30px) saturate(180%);
  -webkit-backdrop-filter: blur(30px) saturate(180%);
  border-radius: 20px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08), 0 0 0 1px rgba(102, 126, 234, 0.1);
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  position: relative;
}

/* 卡片动态光感效果 */
.card::before {
  content: "";
  position: absolute;
  inset: 0;
  border-radius: 20px;
  padding: 1px;
  background: linear-gradient(
    135deg,
    rgba(102, 126, 234, 0.3) 0%,
    rgba(245, 158, 11, 0.2) 50%,
    rgba(139, 92, 246, 0.3) 100%
  );
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.card:hover {
  background: rgba(255, 255, 255, 0.75);
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12),
    0 0 0 1px rgba(102, 126, 234, 0.2);
}

.card:hover::before {
  opacity: 1;
}

.card-title {
  font-size: 14px;
  font-weight: 600;
  color: rgba(102, 126, 234, 0.7);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 12px;
}

/* 域名显示 */
.domain-display {
  padding: 12px 16px;
  background: linear-gradient(
    135deg,
    rgba(102, 126, 234, 0.1) 0%,
    rgba(139, 92, 246, 0.08) 100%
  );
  border-radius: 12px;
  margin-bottom: 12px;
  border: 1px solid rgba(102, 126, 234, 0.2);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
}

.domain-text {
  font-size: 16px;
  font-weight: 500;
  color: #667eea;
}

.detect-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 10px;
  border: 1px solid rgba(102, 126, 234, 0.15);
}

.status-icon {
  font-size: 20px;
}

.status-text {
  font-size: 14px;
  color: #4b5563;
  font-weight: 500;
}

/* 加速状态卡片 */
.detect-status {
  margin-bottom: 12px;
}

.shield-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 16px 0;
}

.shield-button {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: none;
  background: linear-gradient(
    135deg,
    rgba(239, 68, 68, 0.15) 0%,
    rgba(220, 38, 38, 0.1) 100%
  );
  color: #ef4444;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  border: 2px solid rgba(239, 68, 68, 0.3);
  box-shadow: 0 4px 16px rgba(239, 68, 68, 0.15);
}

.shield-button::after {
  content: "";
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  background: radial-gradient(
    circle,
    rgba(239, 68, 68, 0.2) 0%,
    transparent 70%
  );
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: -1;
}

.shield-button:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 8px 24px rgba(239, 68, 68, 0.25);
}

.shield-button:hover:not(:disabled)::after {
  opacity: 1;
}

.shield-button.shield-active {
  background: linear-gradient(
    135deg,
    rgba(34, 197, 94, 0.15) 0%,
    rgba(22, 163, 74, 0.1) 100%
  );
  color: #22c55e;
  border-color: rgba(34, 197, 94, 0.4);
  box-shadow: 0 4px 16px rgba(34, 197, 94, 0.2);
}

.shield-button.shield-active::after {
  background: radial-gradient(
    circle,
    rgba(34, 197, 94, 0.25) 0%,
    transparent 70%
  );
}

.shield-button.shield-active:hover:not(:disabled) {
  box-shadow: 0 8px 24px rgba(34, 197, 94, 0.3);
}

.shield-button.shield-disabled,
.shield-button:disabled {
  background: linear-gradient(
    135deg,
    rgba(156, 163, 175, 0.1) 0%,
    rgba(107, 114, 128, 0.08) 100%
  );
  color: #9ca3af;
  border-color: rgba(156, 163, 175, 0.2);
  cursor: not-allowed;
  opacity: 0.6;
  box-shadow: none;
}

.shield-button.shield-disabled::after,
.shield-button:disabled::after {
  display: none;
}

.shield-button.shield-disabled:hover,
.shield-button:disabled:hover {
  transform: none;
  box-shadow: none;
}

.shield-button.shield-warning {
  background: linear-gradient(
    135deg,
    rgba(245, 158, 11, 0.15) 0%,
    rgba(217, 119, 6, 0.1) 100%
  );
  color: #f59e0b;
  border-color: rgba(245, 158, 11, 0.4);
  box-shadow: 0 4px 16px rgba(245, 158, 11, 0.2);
  cursor: pointer;
  opacity: 1;
}

.shield-button.shield-warning::after {
  display: block;
  background: radial-gradient(
    circle,
    rgba(245, 158, 11, 0.25) 0%,
    transparent 70%
  );
}

.shield-button.shield-warning:hover {
  transform: scale(1.05);
  box-shadow: 0 8px 24px rgba(245, 158, 11, 0.3);
}

.shield-button.shield-warning:hover::after {
  opacity: 1;
}

.shield-icon {
  width: 48px;
  height: 48px;
  transition: all 0.3s ease;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.shield-status {
  font-size: 14px;
  font-weight: 500;
  color: #4b5563;
}

/* CDN 优选部分 - 集成在加速状态卡片内 */
.cdn-section {
  margin-top: 20px;
  padding-top: 20px;
}

.cdn-divider {
  width: 100%;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(102, 126, 234, 0.2) 50%,
    transparent 100%
  );
  margin-bottom: 16px;
}

.cdn-subtitle {
  font-size: 13px;
  font-weight: 600;
  color: rgba(102, 126, 234, 0.7);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 12px;
}

.cdn-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.cdn-detail {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 10px;
  border: 1px solid rgba(102, 126, 234, 0.15);
}

.cdn-label {
  font-size: 13px;
  color: #6b7280;
  font-weight: 500;
}

.cdn-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.cdn-rtt {
  color: #667eea;
}

/* 延迟条 */
.latency-bar-container {
  height: 8px;
  background: rgba(156, 163, 175, 0.15);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 16px;
  border: 1px solid rgba(156, 163, 175, 0.1);
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.05);
}

.latency-bar {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s cubic-bezier(0.34, 1.56, 0.64, 1),
    background-color 0.3s ease;
  position: relative;
  box-shadow: 0 0 8px currentColor;
}

.latency-bar::after {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.3) 50%,
    transparent 100%
  );
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.latency-excellent {
  background: linear-gradient(90deg, #22c55e 0%, #16a34a 100%);
  color: #22c55e;
}

.latency-good {
  background: linear-gradient(90deg, #3b82f6 0%, #2563eb 100%);
  color: #3b82f6;
}

.latency-fair {
  background: linear-gradient(90deg, #f59e0b 0%, #d97706 100%);
  color: #f59e0b;
}

.latency-poor {
  background: linear-gradient(90deg, #ef4444 0%, #dc2626 100%);
  color: #ef4444;
}

/* 重新优选按钮 */
.reoptimize-button {
  width: 100%;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  color: #ffffff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
  position: relative;
  overflow: hidden;
}

.reoptimize-button::before {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.2) 0%,
    transparent 100%
  );
  opacity: 0;
  transition: opacity 0.3s ease;
}

.reoptimize-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
}

.reoptimize-button:hover:not(:disabled)::before {
  opacity: 1;
}

.reoptimize-button:active:not(:disabled) {
  transform: translateY(0);
}

.reoptimize-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 粒子背景 - 亮色模式 */
.particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
  z-index: 0;
}

.particle {
  position: absolute;
  width: 3px;
  height: 3px;
  background: radial-gradient(
    circle,
    rgba(102, 126, 234, 0.4) 0%,
    rgba(102, 126, 234, 0.1) 100%
  );
  border-radius: 50%;
  animation: float linear infinite;
  box-shadow: 0 0 4px rgba(102, 126, 234, 0.3);
}

@keyframes float {
  0% {
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-100vh) translateX(20px);
    opacity: 0;
  }
}

/* 过渡动画 */
.slide-fade-enter-active {
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 1, 1);
}

.slide-fade-enter-from {
  transform: translateY(-20px);
  opacity: 0;
}

.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}
</style>
