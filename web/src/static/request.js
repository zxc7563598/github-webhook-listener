import axios from "axios";
import config from "./config";

const service = axios.create({
  baseURL: config.baseUrl,
  timeout: 10000, // 请求超时时间毫秒
  headers: {
    "Content-Type": "application/json",
  },
});

// 添加请求拦截器
service.interceptors.request.use(
  async function (config) {
    return config;
  },
  function (error) {
    return Promise.reject(error);
  },
);

// 添加响应拦截器
service.interceptors.response.use(
  function (response) {
    // switch (String(response.data.code)) {
    // 	case '900006':
    // }
    return response.data;
  },
  function (error) {
    return Promise.reject(error);
  },
);
export default service;
