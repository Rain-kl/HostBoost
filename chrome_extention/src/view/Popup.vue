<template>
  <div class="macos-container">
    <!-- Header -->
    <header class="macos-header">
      <div class="header-content">
        <div class="header-info">
          <h1 class="app-title">HostBoost</h1>
          <p class="domain-text">{{ domain || "正在加载..." }}</p>
        </div>
        <div class="header-actions">
          <button
            @click="openDnsClearPage"
            class="icon-button"
            aria-label="清理DNS缓存"
            title="清理DNS缓存"
          >
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M3 6h18" />
              <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
              <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
              <line x1="10" y1="11" x2="10" y2="17" />
              <line x1="14" y1="11" x2="14" y2="17" />
            </svg>
          </button>
          <button
            @click="showWebDetails = true"
            class="icon-button"
            aria-label="网站信息"
            title="网站信息"
          >
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="16" x2="12" y2="12" />
              <line x1="12" y1="8" x2="12.01" y2="8" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="macos-main">
      <!-- Status Badge -->
      <div class="status-badge" :class="statusClass">
        <span class="status-icon">{{ detectStatus.icon }}</span>
        <span class="status-label">{{ detectStatus.text }}</span>
      </div>

      <!-- Control Center -->
      <div class="control-center">
        <button
          @click="toggleBoost"
          :disabled="isDetecting"
          class="boost-button"
          :class="toggleButtonClass"
        >
          <div v-if="isBoostEnabled" class="active-ring"></div>
          <div class="boost-icon-wrapper">
            <svg
              class="boost-icon"
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
              <path
                v-if="isBoostEnabled"
                d="M9 12L11 14L15 10"
                class="checkmark"
                stroke-width="2"
              />
            </svg>
          </div>
        </button>

        <div class="control-info">
          <h2 class="control-title">{{ getToggleTitle() }}</h2>
          <p class="control-description">{{ getShieldStatusText() }}</p>
        </div>
      </div>

      <!-- CDN Info Card -->
      <transition name="macos-fade">
        <div v-if="isBoostEnabled" class="info-card">
          <div class="card-header">
            <span class="card-title">CDN 节点</span>
            <span class="badge badge-success">已解析</span>
          </div>
          <div class="card-content">
            <div class="info-item">
              <span class="info-key">优选 IP</span>
              <span class="info-value">{{
                optimizedNode.ip || "获取中..."
              }}</span>
            </div>
          </div>
          <div class="card-footer">
            <button
              @click="changeOptimizedIP"
              :disabled="isChangingIP"
              class="action-button"
              title="当前IP效果不好时，更换为新的优选IP"
            >
              <svg
                v-if="!isChangingIP"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path
                  d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"
                />
              </svg>
              <span v-if="isChangingIP" class="loading-spinner"></span>
              <span>{{ isChangingIP ? "更换中..." : "更换优选IP" }}</span>
            </button>
          </div>
        </div>
      </transition>
    </main>

    <!-- Dialogs -->
    <ForceBoostDialog
      v-model:show="showForceBoostDialog"
      @confirm="handleForceBoost"
      @cancel="handleCancelForceBoost"
    />

    <WebDetailsDrawer v-model:show="showWebDetails" :domain="domain" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from "vue";
import { hostApi, toolApi, optApi } from "@/api/api-ref.js";
import ForceBoostDialog from "@/components/ForceBoostDialog.vue";
import WebDetailsDrawer from "@/components/WebDetailsDrawer.vue";

// 状态管理
const domain = ref("");
const isDetecting = ref(true);
const isBoostEnabled = ref(false);
const isBoostSupported = ref(true);
const isBackendError = ref(false);
const currentTabId = ref(undefined);
const isForceBoost = ref(false); // 标记是否是强制开启的加速
const isChangingIP = ref(false); // 标记是否正在更换优选IP
const currentType = ref(""); // 保存当前 host 的 type，用于调用 /opt/change 接口

// 三连击检测相关状态
const clickCount = ref(0);
const clickTimer = ref(null);
const showForceBoostDialog = ref(false);

// 网站详情状态
const showWebDetails = ref(false);

const detectStatus = ref({
  icon: "🔍",
  text: "正在识别...",
});

const optimizedNode = ref({
  ip: "",
  rtt: 0,
});

// 计算延迟百分比和颜色
computed(() => {
  const rtt = optimizedNode.value.rtt;
  return Math.min((rtt / 200) * 100, 100);
});
computed(() => {
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
      isForceBoost.value = false; // 清除强制标记

      // 保存 type 用于后续更换 IP
      if (response.data.data.type) {
        currentType.value = response.data.data.type;
      }

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
      isForceBoost.value = false; // 清除强制标记

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

  // 如果网站不支持加速，检测三连击
  if (!isBoostSupported.value) {
    handleUnsupportedClick();
    return;
  }

  // 执行加速开关逻辑
  await performBoostToggle();
};

// 处理不支持加速时的点击
const handleUnsupportedClick = () => {
  clickCount.value++;

  // 清除之前的定时器
  if (clickTimer.value) {
    clearTimeout(clickTimer.value);
  }

  // 检测是否达到三次点击
  if (clickCount.value >= 3) {
    clickCount.value = 0;
    showForceBoostDialog.value = true;
    return;
  }

  // 设置1秒后重置计数器
  clickTimer.value = setTimeout(() => {
    clickCount.value = 0;
  }, 1000);
};

// 执行加速开关逻辑
const performBoostToggle = async () => {
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

            // 保存 type 用于后续更换 IP
            if (getResponse.data.data.type) {
              currentType.value = getResponse.data.data.type;
            }

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

        // 如果是强制开启的加速，关闭后恢复原始状态
        if (isForceBoost.value) {
          isForceBoost.value = false;
          isBoostSupported.value = false;
          detectStatus.value = {
            icon: "ℹ️",
            text: "该网站不支持加速",
          };
        } else {
          detectStatus.value = {
            icon: "🌐",
            text: "可加速网站",
          };
        }
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

// 处理强制加速确认
const handleForceBoost = async () => {
  // 临时标记为支持加速，执行开启逻辑
  const originalSupported = isBoostSupported.value;
  isBoostSupported.value = true;
  isForceBoost.value = true; // 标记为强制开启

  try {
    await performBoostToggle();
    // 如果成功开启，更新状态
    detectStatus.value = {
      icon: "✅",
      text: "已强制开启加速",
    };
  } catch (error) {
    // 如果失败，恢复原状态
    isBoostSupported.value = originalSupported;
    isForceBoost.value = false;
    console.error("强制加速失败:", error);
  }
};

// 处理取消强制加速
const handleCancelForceBoost = () => {
  clickCount.value = 0;
  console.log("用户取消了强制加速");
};

// 打开 DNS 清理页面
const openDnsClearPage = () => {
  chrome.tabs.create({ url: "chrome://net-internals/#dns" });
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
// 更换优选 IP
const changeOptimizedIP = async () => {
  if (isChangingIP.value) return;

  // 检查是否有 type 参数
  if (!currentType.value) {
    console.error("缺少 type 参数，无法更换优选 IP");
    detectStatus.value = {
      icon: "❌",
      text: "更换失败: 缺少必要参数",
    };
    return;
  }

  isChangingIP.value = true;

  try {
    // 调用 /opt/change 接口更换优选 IP，传递 type 参数
    const response = await optApi.optChangeGet(currentType.value);

    if (response.data.code === 200 || response.data.code === "200") {
      console.log("更换优选 IP 成功:", response.data);

      // 更换成功后，重新获取当前的优选 IP 信息
      try {
        const getResponse = await hostApi.hostGet(domain.value);

        if (getResponse.data.code === 200 && getResponse.data.data) {
          // 更新 type（可能会变化）
          if (getResponse.data.data.type) {
            currentType.value = getResponse.data.data.type;
          }

          // 更新 IP 信息
          if (getResponse.data.data.ip) {
            optimizedNode.value = {
              ip: getResponse.data.data.ip,
              rtt: 0,
            };
          }

          // 显示成功提示
          detectStatus.value = {
            icon: "✅",
            text: "已更换为新的优选 IP",
          };

          // 3秒后恢复状态提示
          setTimeout(() => {
            if (isBoostEnabled.value) {
              detectStatus.value = {
                icon: "✅",
                text: "加速已启用",
              };
            }
          }, 3000);
        }
      } catch (getError) {
        console.error("获取新的优选 IP 信息失败:", getError);
      }

      isBackendError.value = false;
    } else {
      console.error("更换优选 IP 失败:", response.data);
      detectStatus.value = {
        icon: "❌",
        text: `更换失败: ${response.data.message || "未知错误"}`,
      };
    }
  } catch (error) {
    console.error("更换优选 IP 失败:", error);

    // 判断是否是网络错误
    if (!error.response) {
      isBackendError.value = true;
      detectStatus.value = {
        icon: "⚠️",
        text: "后端服务未启动",
      };
    } else {
      isBackendError.value = false;
      const errorData = error.response?.data;
      const errorMsg = errorData?.message || error.message || "未知错误";
      detectStatus.value = {
        icon: "❌",
        text: `更换失败: ${errorMsg}`,
      };
    }
  } finally {
    isChangingIP.value = false;
  }
};

// 粒子动画样式
// 验证是否为有效域名
const isValidDomain = (hostname) => {
  if (!hostname) return false;

  // 过滤特殊页面
  const invalidPatterns = [
    "newtab",
    "extensions",
    "settings",
    "chrome",
    "about:",
    "edge:",
    "localhost",
    "127.0.0.1",
    "0.0.0.0",
    "::1",
  ];

  // 检查是否匹配无效模式
  const lowerHostname = hostname.toLowerCase();
  if (invalidPatterns.some((pattern) => lowerHostname.includes(pattern))) {
    return false;
  }

  // 检查是否为IP地址(本地网络)
  const ipv4Pattern = /^(\d{1,3}\.){3}\d{1,3}$/;
  const ipv6Pattern = /^([0-9a-fA-F]{0,4}:){2,7}[0-9a-fA-F]{0,4}$/;
  if (ipv4Pattern.test(hostname) || ipv6Pattern.test(hostname)) {
    // 检查是否为本地IP
    if (
      hostname.startsWith("192.168.") ||
      hostname.startsWith("10.") ||
      hostname.startsWith("172.") ||
      hostname === "127.0.0.1" ||
      hostname === "::1"
    ) {
      return false;
    }
  }

  // 检查是否包含点(.)，基本的域名格式
  return hostname.includes(".");


};

// 获取当前域名
onMounted(() => {
  chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
    const tab = tabs[0];
    if (tab?.url) {
      try {
        const url = new URL(tab.url);
        domain.value = url.hostname;
        currentTabId.value = tab.id; // 保存当前标签页ID

        // 验证域名有效性
        if (!isValidDomain(domain.value)) {
          // 如果不是有效域名，直接标记为不支持
          isDetecting.value = false;
          isBoostSupported.value = false;
          isBoostEnabled.value = false;
          detectStatus.value = {
            icon: "ℹ️",
            text: "该网站不支持加速",
          };
        }
      } catch {
        domain.value = "无法解析域名";
        isDetecting.value = false;
        isBoostSupported.value = false;
        detectStatus.value = {
          icon: "⚠️",
          text: "无法解析域名",
        };
      }
    } else {
      domain.value = "未获取到当前标签页";
      isDetecting.value = false;
      isBoostSupported.value = false;
      detectStatus.value = {
        icon: "⚠️",
        text: "未获取到当前标签页",
      };
    }
  });
});

watch(domain, (newVal) => {
  if (newVal && isValidDomain(newVal)) {
    isDetecting.value = true;
    getHost(newVal);
  }
});
</script>

<style scoped>
/* Container */
.macos-container {
  width: 360px;
  min-height: 480px;
  background: var(--macos-bg-primary);
  display: flex;
  flex-direction: column;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display",
    "SF Pro Text", sans-serif;
  color: var(--macos-text-primary);
  overflow: hidden;
}

/* Header */
.macos-header {
  padding: var(--macos-space-lg) var(--macos-space-lg) var(--macos-space-md);
  border-bottom: 1px solid var(--macos-separator-light);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--macos-space-md);
}

.header-info {
  flex: 1;
  min-width: 0;
}

.app-title {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.3;
  letter-spacing: -0.4px;
  margin: 0;
  color: var(--macos-text-primary);
}

.domain-text {
  font-size: 13px;
  color: var(--macos-text-secondary);
  margin: 2px 0 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: var(--macos-space-sm);
  flex-shrink: 0;
}

.icon-button {
  width: 32px;
  height: 32px;
  border-radius: var(--macos-radius-md);
  background: var(--macos-bg-secondary);
  border: none;
  color: var(--macos-text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--macos-transition-fast);
}

.icon-button:hover {
  background: var(--macos-bg-tertiary);
  color: var(--macos-text-primary);
  transform: scale(1.05);
}

.icon-button:active {
  transform: scale(0.95);
}

/* Main Content */
.macos-main {
  flex: 1;
  padding: var(--macos-space-lg);
  display: flex;
  flex-direction: column;
  gap: var(--macos-space-lg);
  overflow-y: auto;
}

/* Status Badge */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--macos-space-sm);
  padding: var(--macos-space-sm) var(--macos-space-md);
  border-radius: var(--macos-radius-lg);
  font-size: 13px;
  font-weight: 500;
  transition: all var(--macos-transition-normal);
}

.status-badge.status-active {
  background: color-mix(in srgb, var(--macos-success) 15%, transparent);
  color: var(--macos-success);
}

.status-badge.status-idle {
  background: color-mix(in srgb, var(--macos-accent) 15%, transparent);
  color: var(--macos-accent);
}

.status-badge.status-error {
  background: color-mix(in srgb, var(--macos-error) 15%, transparent);
  color: var(--macos-error);
}

.status-badge.status-disabled {
  background: var(--macos-bg-secondary);
  color: var(--macos-text-tertiary);
}

.status-icon {
  font-size: 16px;
  line-height: 1;
}

.status-label {
  line-height: 1;
}

/* Control Center */
.control-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--macos-space-lg);
  padding: var(--macos-space-2xl) 0;
}

.boost-button {
  position: relative;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  transition: all var(--macos-transition-normal);
  display: flex;
  align-items: center;
  justify-content: center;
}

.boost-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.boost-button.toggle-idle {
  background: var(--macos-bg-secondary);
  box-shadow: var(--macos-shadow-md);
}

.boost-button.toggle-idle:hover:not(:disabled) {
  background: var(--macos-bg-tertiary);
  box-shadow: var(--macos-shadow-lg);
  transform: scale(1.05);
}

.boost-button.toggle-active {
  background: linear-gradient(
    135deg,
    var(--macos-success) 0%,
    color-mix(in srgb, var(--macos-success) 85%, white) 100%
  );
  box-shadow: 0 8px 24px
      color-mix(in srgb, var(--macos-success) 40%, transparent),
    var(--macos-shadow-lg);
}

.boost-button.toggle-active:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 12px 32px
      color-mix(in srgb, var(--macos-success) 50%, transparent),
    var(--macos-shadow-xl);
}

.boost-button.toggle-error {
  background: linear-gradient(
    135deg,
    var(--macos-error) 0%,
    color-mix(in srgb, var(--macos-error) 85%, white) 100%
  );
  box-shadow: 0 8px 24px color-mix(in srgb, var(--macos-error) 40%, transparent),
    var(--macos-shadow-lg);
}

.boost-button.toggle-disabled {
  background: var(--macos-bg-secondary);
  box-shadow: var(--macos-shadow-sm);
}

.boost-button:active:not(:disabled) {
  transform: scale(0.95);
}

.active-ring {
  position: absolute;
  inset: -12px;
  border-radius: 50%;
  border: 2px solid var(--macos-success);
  opacity: 0.3;
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
    opacity: 0.3;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.1;
  }
}

.boost-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
}

.boost-icon {
  width: 56px;
  height: 56px;
}

.boost-button.toggle-active .boost-icon {
  color: white;
}

.boost-button.toggle-idle .boost-icon,
.boost-button.toggle-disabled .boost-icon {
  color: var(--macos-text-secondary);
}

.boost-button.toggle-error .boost-icon {
  color: white;
}

.boost-icon .checkmark {
  stroke-dasharray: 100;
  stroke-dashoffset: 100;
  animation: checkmark 0.4s ease-out 0.2s forwards;
}

@keyframes checkmark {
  to {
    stroke-dashoffset: 0;
  }
}

.control-info {
  text-align: center;
}

.control-title {
  font-size: 22px;
  font-weight: 600;
  letter-spacing: -0.3px;
  margin: 0;
  color: var(--macos-text-primary);
}

.control-description {
  font-size: 14px;
  color: var(--macos-text-secondary);
  margin: 4px 0 0;
  font-weight: 400;
}

/* Info Card */
.info-card {
  background: var(--macos-glass-light);
  backdrop-filter: saturate(180%) blur(20px);
  -webkit-backdrop-filter: saturate(180%) blur(20px);
  border-radius: var(--macos-radius-xl);
  border: 1px solid var(--macos-separator-light);
  padding: var(--macos-space-lg);
  box-shadow: var(--macos-shadow-sm);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--macos-space-md);
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--macos-text-primary);
}

.badge {
  padding: 4px 10px;
  border-radius: var(--macos-radius-sm);
  font-size: 12px;
  font-weight: 500;
}

.badge-success {
  background: color-mix(in srgb, var(--macos-success) 15%, transparent);
  color: var(--macos-success);
}

.card-content {
  margin-bottom: var(--macos-space-md);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--macos-space-md);
  padding: var(--macos-space-sm) 0;
}

.info-key {
  font-size: 14px;
  color: var(--macos-text-secondary);
  font-weight: 400;
}

.info-value {
  font-size: 14px;
  color: var(--macos-text-primary);
  font-weight: 500;
  font-family: "SF Mono", Monaco, "Courier New", monospace;
  text-align: right;
  word-break: break-all;
}

.card-footer {
  display: flex;
  gap: var(--macos-space-sm);
  padding-top: var(--macos-space-sm);
  border-top: 1px solid var(--macos-separator-light);
}

.action-button {
  flex: 1;
  height: 36px;
  border-radius: var(--macos-radius-md);
  background: var(--macos-accent);
  color: white;
  border: none;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--macos-space-sm);
  cursor: pointer;
  transition: all var(--macos-transition-fast);
}

.action-button:hover:not(:disabled) {
  background: var(--macos-accent-secondary);
  transform: translateY(-1px);
  box-shadow: var(--macos-shadow-md);
}

.action-button:active:not(:disabled) {
  transform: translateY(0);
}

.action-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

</style>
