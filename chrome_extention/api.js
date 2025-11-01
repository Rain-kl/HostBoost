// API 配置
const API_BASE_URL = "http://localhost:15920"; // 请根据实际情况修改

// API 服务类
class HostBoostAPI {
  constructor(baseUrl = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  /**
   * 获取指定域名的 host 配置
   * @param {string} domain - 域名
   * @returns {Promise<Object>} 响应数据
   */
  async getHost(domain) {
    try {
      const url = new URL("/host", this.baseUrl);
      if (domain) {
        url.searchParams.append("domain", domain);
      }

      const response = await fetch(url.toString(), {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data,
      };
    } catch (error) {
      console.error("获取 host 失败:", error);
      return {
        success: false,
        error: error.message,
      };
    }
  }

  /**
   * 新增 host 配置
   * @param {string} domain - 域名
   * @returns {Promise<Object>} 响应数据
   */
  async addHost(domain) {
    try {
      const response = await fetch(`${this.baseUrl}/host`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ domain }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data,
      };
    } catch (error) {
      console.error("添加 host 失败:", error);
      return {
        success: false,
        error: error.message,
      };
    }
  }

  /**
   * 删除 host 配置
   * @param {string} domain - 域名
   * @returns {Promise<Object>} 响应数据
   */
  async deleteHost(domain) {
    try {
      const response = await fetch(`${this.baseUrl}/host`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ domain }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data,
      };
    } catch (error) {
      console.error("删除 host 失败:", error);
      return {
        success: false,
        error: error.message,
      };
    }
  }

  /**
   * 获取所有 host 列表
   * @returns {Promise<Object>} 响应数据
   */
  async getHostList() {
    try {
      const response = await fetch(`${this.baseUrl}/host/list`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return {
        success: true,
        data: data,
      };
    } catch (error) {
      console.error("获取 host 列表失败:", error);
      return {
        success: false,
        error: error.message,
      };
    }
  }
}

// 导出 API 实例
const api = new HostBoostAPI();
