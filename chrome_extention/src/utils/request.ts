/**
 * 请求工具函数
 * 可以根据需要自定义请求逻辑
 */
export default function request<T = any>(url: string, options?: RequestInit): Promise<T> {
  return fetch(url, options).then((response) => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  });
}
