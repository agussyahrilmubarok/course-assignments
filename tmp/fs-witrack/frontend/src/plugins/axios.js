import Cookies from "js-cookie";
import axios from "axios";

const BASE_URL = import.meta.env.VITE_BACKEND_URL;

const axiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest",
    Accept: "application/json",
  },
});

export const axiosWithToken = axios.create({
  baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest",
    Accept: "application/json",
  },
});

axiosWithToken.interceptors.request.use((config) => {
  const token = Cookies.get("witrack_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default axiosInstance;
