<script setup lang="ts">
import { ref, onMounted } from "vue";
import { hostApi, optApi } from "@/api/api-ref.ts";
import type { HostVo, OptVo } from "@/api";

// 响应式数据
const hostList = ref<HostVo[]>([]);
const currentOpt = ref<OptVo | null>(null);
const loading = ref(false);
const newDomain = ref("");

// 加载 Host 列表
const loadHostList = async () => {
  loading.value = true;
  try {
    const response = await hostApi.hostListGet();
    if (response.data.code === 200 && response.data.data) {
      hostList.value = response.data.data.list || [];
    }
  } catch (error) {
    console.error("加载列表失败:", error);
  } finally {
    loading.value = false;
  }
};

// 添加 Host
const addHost = async () => {
  if (!newDomain.value.trim()) {
    return;
  }

  loading.value = true;
  try {
    const response = await hostApi.hostPost({
      domain: newDomain.value.trim(),
    });

    if (response.data.code === "200") {
      newDomain.value = "";
      await loadHostList(); // 重新加载列表
    } else {
    }
  } catch (error) {
    console.error("添加失败:", error);
  } finally {
    loading.value = false;
  }
};

// 删除 Host
const deleteHost = async (domain: string) => {
  if (!confirm(`确定要删除 ${domain} 吗?`)) {
    return;
  }

  loading.value = true;
  try {
    const response = await hostApi.hostDelete({ domain });

    if (response.data.code === "200") {
      await loadHostList(); // 重新加载列表
    } else {
    }
  } catch (error) {
    console.error("删除失败:", error);
  } finally {
    loading.value = false;
  }
};

// 获取优选 IP
const loadCurrentOpt = async (type: string = "github") => {
  loading.value = true;
  try {
    const response = await optApi.optGet(type);
    currentOpt.value = response.data;
  } catch (error) {
    console.error("获取优选IP失败:", error);
  } finally {
    loading.value = false;
  }
};

// 切换优选 IP
const changeOpt = async (type: string = "github") => {
  loading.value = true;
  try {
    const response = await optApi.optChangeGet(type);

    if (response.data.code === "200") {
      await loadCurrentOpt(type); // 重新加载优选IP
    } else {
    }
  } catch (error) {
    console.error("切换失败:", error);
  } finally {
    loading.value = false;
  }
};

// 页面加载时获取数据
onMounted(() => {
  loadHostList();
  loadCurrentOpt();
});
</script>

<template>
  <div class="host-management">
    <h2>HostBoost 管理面板</h2>

    <!-- 优选 IP 信息 -->
    <div class="opt-section">
      <h3>当前优选 IP</h3>
      <div v-if="currentOpt" class="opt-info">
        <p><strong>IP:</strong> {{ currentOpt.ip }}</p>
        <p><strong>延迟:</strong> {{ currentOpt.delay }}ms</p>
        <p><strong>速率:</strong> {{ currentOpt.rate }}</p>
        <button @click="changeOpt('github')" :disabled="loading">
          切换优选 IP
        </button>
      </div>
      <div v-else>
        <p>暂无优选 IP 信息</p>
      </div>
    </div>

    <!-- 添加 Host 表单 -->
    <div class="add-host-section">
      <h3>添加 Host</h3>
      <div class="input-group">
        <input
          v-model="newDomain"
          type="text"
          placeholder="输入域名，如: github.com"
          @keyup.enter="addHost"
        />
        <button @click="addHost" :disabled="loading">添加</button>
      </div>
    </div>

    <!-- Host 列表 -->
    <div class="host-list-section">
      <h3>Host 列表</h3>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="hostList.length === 0" class="empty">暂无数据</div>
      <ul v-else class="host-list">
        <li v-for="host in hostList" :key="host.domain" class="host-item">
          <div class="host-info">
            <strong>{{ host.domain }}</strong>
            <span class="host-ip">{{ host.ip }}</span>
          </div>
          <button
            @click="deleteHost(host.domain)"
            :disabled="loading"
            class="btn-delete"
          >
            删除
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.host-management {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
    Ubuntu, sans-serif;
}

h2 {
  margin-bottom: 20px;
  color: #333;
}

h3 {
  margin: 15px 0 10px;
  color: #555;
  font-size: 16px;
}

.opt-section,
.add-host-section,
.host-list-section {
  margin-bottom: 30px;
  padding: 15px;
  background: #f5f5f5;
  border-radius: 8px;
}

.opt-info p {
  margin: 5px 0;
}

.input-group {
  display: flex;
  gap: 10px;
}

input[type="text"] {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

button {
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
}

button:hover:not(:disabled) {
  background: #0056b3;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-delete {
  background: #dc3545;
}

.btn-delete:hover:not(:disabled) {
  background: #c82333;
}

.host-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.host-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 8px;
}

.host-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.host-ip {
  color: #666;
  font-size: 12px;
}

.loading,
.empty {
  text-align: center;
  color: #999;
  padding: 20px;
}
</style>
