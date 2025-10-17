import { defineStore } from "pinia";
import { useStorage } from "@vueuse/core";
import Cookies from "js-cookie";

import axiosInstance, { axiosWithToken } from "@/plugins/axios";
import router from "@/router";
import { toastSuccess, toastError } from "@/helpers/toastHelper";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: useStorage("backerhub_user", null, localStorage, {
      serializer: {
        read: (v) => (v ? JSON.parse(v) : null),
        write: (v) => JSON.stringify(v),
      },
    }),
    loading: false,
    error: null,
    success: null,
  }),

  getters: {
    token: () => Cookies.get("backerhub_token"),
    isAuthenticated: (state) => !!state.user && !!Cookies.get("backerhub_token"),
  },

  actions: {
    async signUp(payload) {
      this.loading = true;
      try {
        const { data } = await axiosInstance.post("/auth/sign-up", payload);
        toastSuccess("Sign up successfully");
        router.push({ name: "signin" });
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async signIn(payload) {
      this.loading = true;
      try {
        const { data } = await axiosInstance.post("/auth/sign-in", payload);
        this.setSession(data.data.token);
        toastSuccess("Sign in successfully");
        router.push({ name: "dashboard" });
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async signOut() {
      try {
        this.loading = true;
        this.error = null;
        this.success = null;

        this.clearSession();
        toastSuccess("Signed out successfully");

        router.push({ name: "signin" });
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async checkAuth() {
      try {
        const { data } = await axiosWithToken.get("/profiles/me");
        this.user = data.data;
      } catch (error) {
        this.clearSession();
      }
    },

    setSession(token) {
      Cookies.set("backerhub_token", token);
      this.checkAuth();
    },

    clearSession() {
      Cookies.remove("backerhub_token");
      this.user = null;
    },
  },
});