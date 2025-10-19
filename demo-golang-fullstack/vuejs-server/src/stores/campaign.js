import { defineStore } from "pinia";

import axiosInstance, { axiosWithToken } from "@/plugins/axios";
import router from "@/router";
import { toastSuccess, toastError } from "@/helpers/toastHelper";

export const useCampaignStore = defineStore("campaign", {
  state: () => ({
    campaigns: [],
    loading: false,
    error: null,
    success: null,
  }),

  actions: {
    async fetchAllCampaigns() {
      this.loading = true;
      try {
        const response = await axiosWithToken.get("/campaigns");
        this.campaigns = response.data.data || [];
        return this.campaigns;
      } catch (error) {
        this.error = error.response?.data?.message || "Failed to fetch top campaigns";
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },

    async fetchTopCampaigns() {
      this.loading = true;
      try {
        const response = await axiosWithToken.get("/campaigns/top");
        this.campaigns = response.data.data || [];
        return this.campaigns;
      } catch (error) {
        this.error = error.response?.data?.message || "Failed to fetch top campaigns";
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },
    
    async fetchCampaignById(id) {
      this.loading = true;
      try {
        const response = await axiosWithToken.get(`/campaigns/${id}`);
        return response.data.data;
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async fetchMyCampaigns() {
      this.loading = true;
      try {
        const response = await axiosWithToken.get(`/campaigns/me`);
        return response.data.data;
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async fetchMyCampaignById(id) {
      this.loading = true;
      try {
        const response = await axiosWithToken.get(`/campaigns/${id}/me`);
        return response.data.data;
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async createMyCampaign(payload) {
      this.loading = true;
      try {
        const { data } = await axiosWithToken.post("/campaigns", payload);
        toastSuccess("Create campaign successfully");
        router.push({ name: "dashboard" });
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async updateMyCampaign(id, payload) {
      this.loading = true;
      try {
        const { data } = await axiosWithToken.put(`/campaigns/${id}`, payload);
        toastSuccess("Update campaign successfully");
        router.push({ name: "dashboard" }); 
        return data.data;
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },

    async uploadMyCampaignImage(id, payload) {
      this.loading = true;
      try {
        const formData = new FormData();
        formData.append("campaign_image", payload.file);
        formData.append("is_primary", payload.is_primary); // true/false

        const response = await axiosWithToken.post(
          `/campaigns/${id}/images`,
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );

        toastSuccess("Upload image successfully");
        console.log(response.data.data);
        return response.data.data;
      } catch (error) {
        toastError(error.response?.data.message);
      } finally {
        this.loading = false;
      }
    },
  },
});