// 状态管理
let currentDomain = "";
let isActive = false;
let isLoading = false;

// DOM 元素
const elements = {
  currentDomain: null,
  statusCard: null,
  statusIcon: null,
  statusTitle: null,
  statusSubtitle: null,
  toggleButton: null,
  toggleLabel: null,
  toggleSwitch: null,
  errorMessage: null,
  errorText: null,
};

// 初始化
document.addEventListener("DOMContentLoaded", async () => {
  // 获取 DOM 元素
  elements.currentDomain = document.getElementById("currentDomain");
  elements.statusCard = document.getElementById("statusCard");
  elements.statusIcon = document.getElementById("statusIcon");
  elements.statusTitle = document.getElementById("statusTitle");
  elements.statusSubtitle = document.getElementById("statusSubtitle");
  elements.toggleButton = document.getElementById("toggleButton");
  elements.toggleLabel = document.getElementById("toggleLabel");
  elements.toggleSwitch = document.getElementById("toggleSwitch");
  elements.errorMessage = document.getElementById("errorMessage");
  elements.errorText = document.getElementById("errorText");

  // 绑定事件
  elements.toggleButton.addEventListener("click", handleToggle);

  // 获取当前标签页信息并检查状态
  await initializeApp();
});

/**
 * 初始化应用
 */
async function initializeApp() {
  try {
    // 获取当前标签页
    const [tab] = await chrome.tabs.query({
      active: true,
      currentWindow: true,
    });

    if (!tab || !tab.url) {
      showError("无法获取当前页面信息");
      return;
    }

    // 解析域名
    const url = new URL(tab.url);
    currentDomain = url.hostname;

    // 显示域名
    elements.currentDomain.textContent = currentDomain;

    // 检查状态
    await checkHostStatus();
  } catch (error) {
    console.error("初始化失败:", error);
    showError("初始化失败: " + error.message);
  }
}

/**
 * 检查 host 状态
 */
async function checkHostStatus() {
  setLoading(true);
  hideError();

  try {
    const result = await api.getHost(currentDomain);

    if (result.success && result.data) {
      // 检查返回的数据是否包含该域名
      const hasHost =
        result.data.domain === currentDomain ||
        (result.data.code === 200 && result.data.data);
      setActiveStatus(hasHost);
    } else {
      // 获取失败或数据为空,视为未开启
      setActiveStatus(false);
    }
  } catch (error) {
    console.error("检查状态失败:", error);
    setActiveStatus(false);
    showError("检查状态失败: " + error.message);
  } finally {
    setLoading(false);
  }
}

/**
 * 处理开关切换
 */
async function handleToggle() {
  if (isLoading) return;

  setLoading(true);
  hideError();

  try {
    if (isActive) {
      // 关闭 - 删除 host
      const result = await api.deleteHost(currentDomain);

      if (result.success) {
        setActiveStatus(false);
        showSuccess("已关闭");
      } else {
        showError("关闭失败: " + (result.error || "未知错误"));
      }
    } else {
      // 开启 - 添加 host
      const result = await api.addHost(currentDomain);

      if (result.success) {
        setActiveStatus(true);
        showSuccess("已开启");
      } else {
        showError("开启失败: " + (result.error || "未知错误"));
      }
    }
  } catch (error) {
    console.error("切换状态失败:", error);
    showError("操作失败: " + error.message);
  } finally {
    setLoading(false);
  }
}

/**
 * 设置激活状态
 */
function setActiveStatus(active) {
  isActive = active;

  // 更新状态图标
  elements.statusIcon.className =
    "status-icon " + (active ? "active" : "inactive");

  // 更新状态文本
  elements.statusTitle.textContent = active ? "已开启" : "未开启";
  elements.statusSubtitle.textContent = active
    ? `${currentDomain} 已加入加速列表`
    : `${currentDomain} 未加入加速列表`;

  // 更新开关
  elements.toggleLabel.textContent = active ? "关闭加速" : "开启加速";

  if (active) {
    elements.toggleSwitch.classList.add("active");
  } else {
    elements.toggleSwitch.classList.remove("active");
  }

  // 更新图标 SVG
  updateStatusIcon(active);
}

/**
 * 更新状态图标
 */
function updateStatusIcon(active) {
  if (active) {
    elements.statusIcon.innerHTML = `
            <svg width="40" height="40" viewBox="0 0 40 40">
                <circle cx="20" cy="20" r="18" fill="none" stroke="currentColor" stroke-width="2"/>
                <path d="M12 20 L17 25 L28 14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
        `;
  } else {
    elements.statusIcon.innerHTML = `
            <svg width="40" height="40" viewBox="0 0 40 40">
                <circle cx="20" cy="20" r="18" fill="none" stroke="currentColor" stroke-width="2"/>
                <path d="M15 15 L25 25 M25 15 L15 25" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
        `;
  }
}

/**
 * 设置加载状态
 */
function setLoading(loading) {
  isLoading = loading;

  if (loading) {
    elements.toggleButton.disabled = true;
    elements.statusIcon.className = "status-icon loading";
    elements.statusIcon.innerHTML = `
            <svg width="40" height="40" viewBox="0 0 40 40">
                <circle cx="20" cy="20" r="18" fill="none" stroke="currentColor" stroke-width="2"/>
                <path d="M20 2 A18 18 0 0 1 38 20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
        `;
  } else {
    elements.toggleButton.disabled = false;
    updateStatusIcon(isActive);
  }
}

/**
 * 显示错误信息
 */
function showError(message) {
  elements.errorText.textContent = message;
  elements.errorMessage.style.display = "flex";
}

/**
 * 隐藏错误信息
 */
function hideError() {
  elements.errorMessage.style.display = "none";
}

/**
 * 显示成功提示(短暂显示后自动隐藏)
 */
function showSuccess(message) {
  elements.errorText.textContent = message;
  elements.errorMessage.style.display = "flex";
  elements.errorMessage.style.background = "rgba(16, 185, 129, 0.15)";
  elements.errorMessage.style.borderColor = "rgba(16, 185, 129, 0.3)";

  setTimeout(() => {
    hideError();
    // 重置样式
    elements.errorMessage.style.background = "rgba(239, 68, 68, 0.15)";
    elements.errorMessage.style.borderColor = "rgba(239, 68, 68, 0.3)";
  }, 2000);
}
