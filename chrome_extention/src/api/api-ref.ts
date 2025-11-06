/**
 * API 使用示例文件
 * 展示如何使用 OpenAPI Generator 生成的 API 代码
 */

import { Configuration, HostApi, OptApi, ToolApi } from './index';
import type { HostPostRequest, OptRequest } from './models';

// ============================================
// 1. 配置 API 客户端
// ============================================

// 创建配置对象 - 基础配置
const config = new Configuration({
  basePath: 'http://localhost:15920', // 你的 API 基础路径
  // 如果需要认证,可以添加以下配置:
  // apiKey: 'your-api-key',
  // accessToken: 'your-access-token',
  // 或者使用函数动态获取 token:
  // accessToken: async () => {
  //   return await getTokenFromStorage();
  // }
});

// 创建 API 实例
const hostApi = new HostApi(config);
const optApi = new OptApi(config);
const toolApi = new ToolApi(config);

// ============================================
// 2. HostApi 使用示例
// ============================================

/**
 * 示例 1: 查询单个 Host
 */
export async function getHostExample() {
  try {
    const response = await hostApi.hostGet('github.com');
    console.log('查询成功:', response.data);
    return response.data;
  } catch (error) {
    console.error('查询失败:', error);
    throw error;
  }
}

/**
 * 示例 2: 查询 Host 列表
 */
export async function getHostListExample() {
  try {
    const response = await hostApi.hostListGet();
    console.log('列表查询成功:', response.data);
    
    // 访问返回的数据 (注意: QueryHostListResponse.code 是 number 类型)
    if (response.data.code === 200) {
      const hosts = response.data.data?.list || [];
      console.log(`共查询到 ${hosts.length} 条记录`);
      return hosts;
    }
  } catch (error) {
    console.error('列表查询失败:', error);
    throw error;
  }
}

/**
 * 示例 3: 添加 Host
 */
export async function addHostExample() {
  try {
    const hostData: HostPostRequest = {
      domain: 'github.com'
    };
    
    const response = await hostApi.hostPost(hostData);
    console.log('添加成功:', response.data);
    
    // 注意: BaseResponse.code 是 string 类型
    if (response.data.code === '200') {
      console.log('Host 添加成功!');
    }
    return response.data;
  } catch (error) {
    console.error('添加失败:', error);
    throw error;
  }
}

/**
 * 示例 4: 删除 Host
 */
export async function deleteHostExample() {
  try {
    const hostData: HostPostRequest = {
      domain: 'github.com'
    };
    
    const response = await hostApi.hostDelete(hostData);
    console.log('删除成功:', response.data);
    
    if (response.data.code === '200') {
      console.log('Host 删除成功!');
    }
    return response.data;
  } catch (error) {
    console.error('删除失败:', error);
    throw error;
  }
}

// ============================================
// 3. OptApi 使用示例
// ============================================

/**
 * 示例 5: 获取优选 IP
 */
export async function getOptExample() {
  try {
    const response = await optApi.optGet('github');
    console.log('获取优选IP成功:', response.data);
    
    // response.data 就是 OptVo 对象
    const optData = response.data;
    console.log(`优选IP: ${optData.ip}, 延迟: ${optData.delay}ms, 速率: ${optData.rate}`);
    return optData;
  } catch (error) {
    console.error('获取优选IP失败:', error);
    throw error;
  }
}

/**
 * 示例 6: 切换优选 IP
 */
export async function changeOptExample() {
  try {
    const response = await optApi.optChangeGet('github');
    console.log('切换优选IP成功:', response.data);
    
    if (response.data.code === '200') {
      console.log('切换成功!');
    }
    return response.data;
  } catch (error) {
    console.error('切换优选IP失败:', error);
    throw error;
  }
}

/**
 * 示例 7: 上报优选 IP (OptReportPost)
 */
export async function reportOptExample() {
  try {
    const optRequest: OptRequest = {
      type: 'github',
      data: [
        {
          ip: '20.205.243.166',
          delay: '50',
          rate: '100'
        },
        {
          ip: '20.205.243.167',
          delay: '60',
          rate: '95'
        }
      ]
    };
    
    const response = await optApi.optReportPost(optRequest);
    console.log('上报优选IP成功:', response.data);
    
    if (response.data.code === '200') {
      console.log('上报成功!');
    }
    return response.data;
  } catch (error) {
    console.error('上报优选IP失败:', error);
    throw error;
  }
}

// ============================================
// 4. 在 Chrome 扩展中使用的高级示例
// ============================================

/**
 * 示例 9: 带错误处理和重试的请求
 */
export async function getHostWithRetry(host: string, maxRetries = 3) {
  let lastError;
  
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await hostApi.hostGet(host);
      return response.data;
    } catch (error) {
      console.warn(`请求失败 (${i + 1}/${maxRetries})`, error);
      lastError = error;
      
      // 等待一段时间后重试
      if (i < maxRetries - 1) {
        await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
      }
    }
  }
  
  throw lastError;
}

/**
 * 示例 10: 批量操作 - 批量添加 Hosts
 */
export async function batchAddHosts(domains: string[]) {
  const results = await Promise.allSettled(
    domains.map(domain => hostApi.hostPost({ domain }))
  );
  
  const succeeded = results.filter(r => r.status === 'fulfilled').length;
  const failed = results.filter(r => r.status === 'rejected').length;
  
  console.log(`批量添加完成: 成功 ${succeeded}, 失败 ${failed}`);
  
  return {
    total: domains.length,
    succeeded,
    failed,
    results
  };
}

// ============================================
// 5. 导出配置好的 API 实例供全局使用
// ============================================

export { hostApi, optApi, config,toolApi };

// 如果需要重新配置 API (例如切换环境)
export function reconfigureApi(newBasePath: string) {
  const newConfig = new Configuration({
    basePath: newBasePath,
    // 保留其他配置...
  });
  
  return {
    hostApi: new HostApi(newConfig),
    optApi: new OptApi(newConfig)
  };
}
