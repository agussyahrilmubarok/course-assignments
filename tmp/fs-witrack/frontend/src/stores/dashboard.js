import { defineStore } from "pinia";
import { axiosWithToken } from "@/plugins/axios";
import { handleError } from "@/helpers/errorHelper";
import { toastSuccess, toastError } from "@/helpers/toastHelper";

export const useDashboardStore = defineStore("dashboard", {
  state: () => ({
    statistic: null,
    loading: false,
    error: null,
    success: null,
  }),

  actions: {
    async fetchStatistics() {
      this.loading = true;
      this.error = null;

      try {
        const { data } = await axiosWithToken.get("dashboards/statistics");
        console.log(data);
        this.statistic = data;
      } catch (err) {
        this.error = handleError(err);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },
  },
});
