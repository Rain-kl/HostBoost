<template>
  <div class="ios-container">
    <!-- 顶部导航栏 -->
    <header class="ios-header">
      <div class="header-content">
        <div class="header-title">
          <h1>HostBoost</h1>
          <p class="header-subtitle">{{ domain || "正在加载..." }}</p>
        </div>
        <button
          @click="showWebDetails = true"
          class="info-button"
          aria-label="网站信息"
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="16" x2="12" y2="12" />
            <line x1="12" y1="8" x2="12.01" y2="8" />
          </svg>
        </button>
      </div>
    </header>

    <!-- 主内容区 -->
    <main class="ios-content">
      <!-- 状态指示器 -->
      <div class="status-indicator" :class="statusClass">
        <span class="status-icon">{{ detectStatus.icon }}</span>
        <span class="status-text">{{ detectStatus.text }}</span>
      </div>

      <!-- 主控制卡片 -->
      <div class="control-card">
        <button
          @click="toggleBoost"
          :disabled="isDetecting"
          class="boost-toggle"
          :class="toggleButtonClass"
        >
          <div v-if="isBoostEnabled" class="pulse-ring"></div>
          <div class="toggle-icon">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
              <path
                d="M12 2L4 6V11C4 16.55 7.84 21.74 12 23C16.16 21.74 20 16.55 20 11V6L12 2Z"
                :fill="isBoostEnabled ? 'currentColor' : 'none'"
                :stroke="isBoostEnabled ? 'none' : 'currentColor'"
                stroke-width="1.5"
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
          </div>
        </button>

        <div class="toggle-label">
          <p class="toggle-title">{{ getToggleTitle() }}</p>
          <p class="toggle-description">{{ getShieldStatusText() }}</p>
        </div>
      </div>

      <!-- CDN 信息卡片 -->
      <transition name="slide-fade">
        <div v-if="isBoostEnabled" class="info-card">
          <div class="info-header">
            <span class="info-title">CDN 节点</span>
            <span class="status-badge">已连接</span>
          </div>
          <div class="info-row">
            <span class="info-label">优选 IP</span>
            <span class="info-value">{{
              optimizedNode.ip || "获取中..."
            }}</span>
          </div>
        </div>
      </transition>
    </main>

    <!-- 网站详情底部抽屉 -->
    <var-popup
      v-model:show="showWebDetails"
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
          <button @click="showWebDetails = false" class="close-button">
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
          <div v-if="loadingWebDetails" class="loading-state">
            <var-loading type="wave" :size="32" />
            <p>正在获取信息...</p>
          </div>

          <!-- 错误状态 -->
          <div v-else-if="webDetailsError" class="error-state">
            <span class="error-icon">⚠️</span>
            <p class="error-message">{{ webDetailsError }}</p>
            <var-button type="primary" size="small" @click="fetchWebDetails"
              >重试</var-button
            >
          </div>

          <!-- 网站信息列表 -->
          <div v-else-if="webDetails" class="details-list">
            <DetailItem
              v-for="item in webDetailsDisplay"
              :key="item.key"
              :icon="item.icon"
              :label="item.label"
              :value="item.value"
            />

            <div v-if="webDetailsDisplay.length === 0" class="empty-state">
              <span class="empty-icon">📭</span>
              <p>暂无可显示的信息</p>
            </div>
          </div>
        </div>
      </div>
    </var-popup>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from "vue";
import { hostApi, toolApi } from "@/api/api-ref.js";
import DetailItem from "./components/DetailItem.vue";

// 状态管理
const domain = ref("");
const isDetecting = ref(true);
const isBoostEnabled = ref(false);
const isBoostSupported = ref(true);
const isBackendError = ref(false);
const isOptimizing = ref(false);
const countdown = ref(3);
const currentTabId = ref(undefined);

// 网站详情状态
const showWebDetails = ref(false);
const webDetails = ref(null);
const loadingWebDetails = ref(false);
const webDetailsError = ref("");

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
const checkCdnSupport = async (domain) => {
  try {
    // 更新状态为检查中
    detectStatus.value = {
      icon: "🔍",
      text: "正在检查 CDN 支持...",
    };

    const response = await toolApi.toolWebDetailsGet(domain);
    if (
      (response.data.code === 200 || response.data.code === "200") &&
      response.data.data
    ) {
      if (
        response.data.data.organization.trim().toLowerCase() === "cloudflare"
      ) {
        detectStatus.value = {
          icon: "🌐",
          text: "已识别为 Cloudflare 节点",
        };
        return true;
      } else {
        // 不是 Cloudflare 节点
        detectStatus.value = {
          icon: "ℹ️",
          text: "该网站不支持加速",
        };
        return false;
      }
    }

    // API 返回数据不正确
    detectStatus.value = {
      icon: "ℹ️",
      text: "该网站不支持加速",
    };
    return false;
  } catch (error) {
    console.error("检查 CDN 支持失败:", error);

    // 显示错误信息
    detectStatus.value = {
      icon: "⚠️",
      text: `CDN 检查失败: ${error.message || "网络错误"}`,
    };

    // 如果是网络错误,可能后端服务有问题
    if (!error.response) {
      isBackendError.value = true;
    }

    return false;
  }
};

// API 调用 - 检查域名状态
const getHost = async (domain) => {
  try {
    // 优先调用 hostGet 接口查询状态
    const response = await hostApi.hostGet(domain);

    isBackendError.value = false; // 能收到响应，清除后端错误状态
    isDetecting.value = false;

    // 如果查询成功(code === 200 或 code === "200")，说明已有加速记录，直接开启加速
    if (
      (response.data.code === 200 || response.data.code === "200") &&
      response.data.data
    ) {
      isBoostEnabled.value = true;
      isBoostSupported.value = true; // 已经加速说明肯定支持

      // 从 API 响应中获取优化节点信息
      if (response.data.data.ip) {
        optimizedNode.value = {
          ip: response.data.data.ip,
          rtt: 0,
        };
      }

      // 已经加速，不需要再检查 CDN 支持
      detectStatus.value = {
        icon: "✅",
        text: "加速已启用",
      };
    } else {
      // 查询失败或无记录（但服务端有响应），需要检测域名是否支持加速
      isBoostEnabled.value = false;

      // 检测域名是否支持加速
      isBoostSupported.value = await checkCdnSupport(domain);

      // 根据是否支持加速显示不同的状态
      if (isBoostSupported.value) {
        detectStatus.value = {
          icon: "🌐",
          text: "可加速网站",
        };
      } else {
        detectStatus.value = {
          icon: "ℹ️",
          text: "该网站不支持加速",
        };
      }
    }
  } catch (error) {
    console.error("查询域名状态失败:", error);

    isDetecting.value = false;

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
      // 服务端有响应但返回错误（如404, 500等），需要检测域名是否支持加速
      isBackendError.value = false;
      isBoostEnabled.value = false;

      // 检测域名是否支持加速
      isBoostSupported.value = await checkCdnSupport(domain);

      const errorData = error.response?.data;
      const errorCode = errorData?.code || error.response.status;
      const errorMsg = errorData?.message || error.message || "未知错误";
      detectStatus.value = {
        icon: "❌",
        text: `查询失败 [${errorCode}]: ${errorMsg}`,
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
        const errorMsg = response.data.message || "未知错误";
        console.error("开启加速失败:", response.data);
        detectStatus.value = {
          icon: "❌",
          text: `开启失败 [${response.data.code}]: ${errorMsg}`,
        };
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
        const errorMsg = response.data.message || "未知错误";
        console.error("关闭加速失败:", response.data);
        detectStatus.value = {
          icon: "❌",
          text: `关闭失败 [${response.data.code}]: ${errorMsg}`,
        };
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
      const errorData = error.response?.data;
      const errorCode = errorData?.code || error.response.status;
      const errorMsg = errorData?.message || error.message || "未知错误";
      detectStatus.value = {
        icon: "❌",
        text: `操作失败 [${errorCode}]: ${errorMsg}`,
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
    return "点击重新检查服务";
  }
  if (!isBoostSupported.value) {
    return "该网站不支持加速";
  }
  return isBoostEnabled.value ? "加速服务已启用" : "轻触以启用加速";
};

// 获取切换按钮标题
const getToggleTitle = () => {
  if (isBoostEnabled.value) return "已加速";
  if (isBackendError.value) return "服务异常";
  if (!isBoostSupported.value) return "不支持";
  return "未加速";
};

// 计算状态类名
const statusClass = computed(() => {
  if (isBackendError.value) return "status-error";
  if (isBoostEnabled.value) return "status-active";
  if (!isBoostSupported.value) return "status-disabled";
  return "status-idle";
});

// 计算按钮类名
const toggleButtonClass = computed(() => {
  if (isBoostEnabled.value) return "toggle-active";
  if (isBackendError.value) return "toggle-error";
  if (!isBoostSupported.value) return "toggle-disabled";
  return "toggle-idle";
});

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

// 网站详情相关方法
const fetchWebDetails = async () => {
  if (
    !domain.value ||
    domain.value === "无法解析域名" ||
    domain.value === "未获取到当前标签页"
  ) {
    webDetailsError.value = "无效的域名";
    return;
  }

  loadingWebDetails.value = true;
  webDetailsError.value = "";
  webDetails.value = null;

  try {
    const response = await toolApi.toolWebDetailsGet(domain.value);
    if (
      (response.data.code === 200 || response.data.code === "200") &&
      response.data.data
    ) {
      webDetails.value = response.data.data;
    } else {
      webDetailsError.value = response.data.message || "获取网站信息失败";
    }
  } catch (error) {
    console.error("获取网站信息失败:", error);
    webDetailsError.value = error.message || "网络请求失败";
  } finally {
    loadingWebDetails.value = false;
  }
};

// 当打开弹窗时自动获取网站信息
watch(showWebDetails, (newVal) => {
  if (newVal && !webDetails.value && !loadingWebDetails.value) {
    fetchWebDetails();
  }
});

const closeWebDetails = () => {
  // 弹窗关闭时可选择清理数据
  // webDetails.value = null;
  // webDetailsError.value = "";
};

// 网站详情展示数据（处理字段不存在的情况）
const webDetailsDisplay = computed(() => {
  if (!webDetails.value) return [];

  const details = webDetails.value;
  const items = [
    { key: "ip", icon: "🌐", label: "IP 地址", value: details.ip },
    { key: "country", icon: "🌍", label: "国家", value: details.country },
    {
      key: "country_code",
      icon: "🏳️",
      label: "国家代码",
      value: details.country_code,
    },
    { key: "region", icon: "📍", label: "地区", value: details.region },
    {
      key: "region_code",
      icon: "🗺️",
      label: "地区代码",
      value: details.region_code,
    },
    { key: "city", icon: "🏙️", label: "城市", value: details.city },
    {
      key: "organization",
      icon: "🏢",
      label: "组织",
      value: details.organization,
    },
    { key: "isp", icon: "📡", label: "ISP", value: details.isp },
    { key: "asn", icon: "🔢", label: "ASN", value: details.asn },
    {
      key: "asn_organization",
      icon: "🏛️",
      label: "ASN 组织",
      value: details.asn_organization,
    },
    { key: "timezone", icon: "🕐", label: "时区", value: details.timezone },
    {
      key: "offset",
      icon: "⏱️",
      label: "时区偏移",
      value: details.offset ? `UTC+${details.offset / 3600}` : undefined,
    },
    { key: "latitude", icon: "🧭", label: "纬度", value: details.latitude },
    { key: "longitude", icon: "🧭", label: "经度", value: details.longitude },
    {
      key: "continent_code",
      icon: "🌏",
      label: "洲代码",
      value: details.continent_code,
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
/* iOS 16 风格设计 */
.ios-container {
  width: 380px;
  height: 600px;
  background: linear-gradient(180deg, #f5f5f7 0%, #ffffff 100%);
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display",
    "SF Pro Text", "Helvetica Neue", Arial, sans-serif;
}

.dark .ios-container {
  background: linear-gradient(180deg, #1c1c1e 0%, #000000 100%);
}

/* 顶部导航栏 */
.ios-header {
  padding: 16px 20px 12px;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: saturate(180%) blur(20px);
  border-bottom: 0.5px solid rgba(0, 0, 0, 0.08);
}

.dark .ios-header {
  background: rgba(28, 28, 30, 0.72);
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-title h1 {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: #1d1d1f;
  margin: 0;
}

.dark .header-title h1 {
  color: #f5f5f7;
}

.header-subtitle {
  font-size: 13px;
  color: #86868b;
  margin: 2px 0 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 260px;
}

.info-button {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.05);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #007aff;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .info-button {
  background: rgba(255, 255, 255, 0.1);
  color: #0a84ff;
}

.info-button:hover {
  background: rgba(0, 122, 255, 0.1);
  transform: scale(1.05);
}

.info-button:active {
  transform: scale(0.95);
}

/* 主内容区 */
.ios-content {
  padding: 20px;
  overflow-y: auto;
  height: calc(600px - 72px);
}

/* 状态指示器 */
.status-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  margin-bottom: 16px;
  transition: all 0.3s ease;
  border: 0.5px solid rgba(0, 0, 0, 0.04);
}

.dark .status-indicator {
  background: rgba(58, 58, 60, 0.6);
  border-color: rgba(255, 255, 255, 0.06);
}

.status-icon {
  font-size: 20px;
}

.status-text {
  font-size: 14px;
  font-weight: 500;
  color: #1d1d1f;
  flex: 1;
}

.dark .status-text {
  color: #f5f5f7;
}

.status-active {
  background: linear-gradient(
    135deg,
    rgba(52, 199, 89, 0.15) 0%,
    rgba(48, 209, 88, 0.1) 100%
  );
  border-color: rgba(52, 199, 89, 0.2);
}

.status-error {
  background: linear-gradient(
    135deg,
    rgba(255, 149, 0, 0.15) 0%,
    rgba(255, 159, 10, 0.1) 100%
  );
  border-color: rgba(255, 149, 0, 0.2);
}

.status-disabled {
  opacity: 0.6;
}

/* 主控制卡片 */
.control-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: saturate(180%) blur(20px);
  border-radius: 24px;
  padding: 32px 24px;
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08), 0 0 0 0.5px rgba(0, 0, 0, 0.04);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dark .control-card {
  background: rgba(58, 58, 60, 0.7);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3),
    0 0 0 0.5px rgba(255, 255, 255, 0.1);
}

/* 加速切换按钮 */
.boost-toggle {
  position: relative;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  background: linear-gradient(135deg, #f5f5f7 0%, #e8e8ed 100%);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.dark .boost-toggle {
  background: linear-gradient(135deg, #3a3a3c 0%, #2c2c2e 100%);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.boost-toggle:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.boost-toggle:not(:disabled):hover {
  transform: scale(1.05);
}

.boost-toggle:not(:disabled):active {
  transform: scale(0.98);
}

.toggle-active {
  background: linear-gradient(135deg, #34c759 0%, #30d158 100%);
  box-shadow: 0 12px 32px rgba(52, 199, 89, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.toggle-error {
  background: linear-gradient(135deg, #ff9500 0%, #ff9f0a 100%);
  box-shadow: 0 12px 32px rgba(255, 149, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.toggle-disabled {
  background: linear-gradient(135deg, #c7c7cc 0%, #d1d1d6 100%);
  opacity: 0.6;
}

.dark .toggle-disabled {
  background: linear-gradient(135deg, #48484a 0%, #3a3a3c 100%);
}

.toggle-icon {
  color: #86868b;
  transition: all 0.3s ease;
}

.toggle-active .toggle-icon {
  color: white;
}

.toggle-error .toggle-icon {
  color: white;
}

.dark .toggle-icon {
  color: #98989d;
}

/* 脉冲环 */
.pulse-ring {
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  background: radial-gradient(
    circle,
    rgba(52, 199, 89, 0.3) 0%,
    transparent 70%
  );
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(1.1);
  }
}

/* 切换标签 */
.toggle-label {
  text-align: center;
}

.toggle-title {
  font-size: 20px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0 0 4px;
  letter-spacing: -0.3px;
}

.dark .toggle-title {
  color: #f5f5f7;
}

.toggle-description {
  font-size: 13px;
  color: #86868b;
  margin: 0;
}

/* 信息卡片 */
.info-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: saturate(180%) blur(20px);
  border-radius: 18px;
  padding: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06), 0 0 0 0.5px rgba(0, 0, 0, 0.04);
}

.dark .info-card {
  background: rgba(58, 58, 60, 0.7);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3),
    0 0 0 0.5px rgba(255, 255, 255, 0.1);
}

.info-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.info-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d1d1f;
}

.dark .info-title {
  color: #f5f5f7;
}

.status-badge {
  font-size: 11px;
  font-weight: 600;
  color: #34c759;
  background: rgba(52, 199, 89, 0.15);
  padding: 4px 10px;
  border-radius: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.info-label {
  font-size: 13px;
  color: #86868b;
}

.info-value {
  font-size: 13px;
  font-weight: 500;
  font-family: "SF Mono", Monaco, "Courier New", monospace;
  color: #1d1d1f;
}

.dark .info-value {
  color: #f5f5f7;
}

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

/* 滑入滑出动画 */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(-12px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(12px);
}

/* 滚动条样式 */
.ios-content::-webkit-scrollbar,
.popup-body::-webkit-scrollbar {
  width: 6px;
}

.ios-content::-webkit-scrollbar-track,
.popup-body::-webkit-scrollbar-track {
  background: transparent;
}

.ios-content::-webkit-scrollbar-thumb,
.popup-body::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.dark .ios-content::-webkit-scrollbar-thumb,
.dark .popup-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}
</style>
