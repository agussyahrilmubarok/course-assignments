import { defineStore } from "pinia";
import { useStorage } from "@vueuse/core";
import Cookies from "js-cookie";

import axiosInstance, { axiosWithToken } from "@/plugins/axios";
import router from "@/router";
import { handleError } from "@/helpers/errorHelper";
import { toastSuccess, toastError } from "@/helpers/toastHelper";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: useStorage("witrack_user", null),
    loading: false,
    error: null,
    success: null,
  }),

  getters: {
    token: () => Cookies.get("witrack_token"),
    isAuthenticated: (state) => !!state.user && !!Cookies.get("witrack_token"),
    isAdmin: (state) => state.user?.roles?.includes("ROLE_ADMIN") ?? false,
  },

  actions: {
    async register(payload) {
      this.loading = true;
      try {
        const { data } = await axiosInstance.post("/auth/sign-up", payload);
        this.setSession(data.token, data.user);

        toastSuccess("Sign up successfully");

        router.push({
          name: this.isAdmin ? "admin.dashboard" : "app.dashboard",
        });
      } catch (error) {
        this.error = handleError(error);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },

    async login(payload) {
      this.loading = true;
      try {
        const { data } = await axiosInstance.post("/auth/sign-in", payload);
        this.setSession(data.token, data.user);

        toastSuccess("Sign in successfully");

        router.push({
          name: this.isAdmin ? "admin.dashboard" : "app.dashboard",
        });
      } catch (error) {
        this.error = handleError(error);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },

    async logout() {
      try {
        this.loading = true;
        this.error = null;
        this.success = null;

        this.clearSession();
        toastSuccess("Signed out successfully");

        router.push({ name: "login" });
      } catch (error) {
        this.error = handleError(error);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },
    async checkAuth() {
      try {
        const { data } = await axiosWithToken.get("/users/profiles/me");
        this.user = data.user;
      } catch (error) {
        this.clearSession();
        this.error = handleError(error);
      }
    },

    setSession(token, user) {
      Cookies.set("witrack_token", token);
      this.user = user;
    },

    clearSession() {
      Cookies.remove("witrack_token");
      this.user = null;
    },
  },
});
