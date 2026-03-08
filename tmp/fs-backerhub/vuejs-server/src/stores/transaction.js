import { defineStore } from "pinia";

import axiosInstance, { axiosWithToken } from "@/plugins/axios";
import router from "@/router";
import { toastSuccess, toastError } from "@/helpers/toastHelper";

export const useTransactionStore = defineStore("transaction", {
  state: () => ({
    transactions: [],
    loading: false,
    error: null,
    success: null,
  }),

  actions: {
    async donationCampaign(payload) {
      this.loading = true;
      try {
        const { data } = await axiosWithToken.post(
          "/transactions/donation",
          payload
        );
        if (data.data.payment_url) {
          window.open(data.data.payment_url, "_blank");
        }
        return data.data.transaction_id;
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async fetchTransactionsByCampaign(id) {
      this.loading = true;
      try {
        const { data } = await axiosWithToken.get(`/transactions/campaign/${id}`);
        return data.data
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async fetchTransactionsByUser(id) {
      this.loading = true;
      try {
        const { data } = await axiosWithToken.get("/transactions/me");
        return data.data
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },
  },
});
