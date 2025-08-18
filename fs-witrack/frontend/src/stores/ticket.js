import { defineStore } from "pinia";

import router from "@/router";
import { axiosWithToken } from "@/plugins/axios";
import { handleError } from "@/helpers/errorHelper";
import { toastSuccess, toastError } from "@/helpers/toastHelper";

export const useTicketStore = defineStore("ticket", {
  state: () => ({
    tickets: [],
    loading: false,
    error: null,
    success: null,
  }),

  actions: {
    async fetchTickets(params) {
      this.loading = true;
      try {
        const response = await axiosWithToken.get(`/tickets`, { params });
        this.tickets = response.data;
      } catch (error) {
        this.error = handleError(error);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },

    async fetchMyTickets(params) {
      this.loading = true;
      try {
        const response = await axiosWithToken.get(`/tickets/me`, { params });
        this.tickets = response.data;
      } catch (error) {
        this.error = handleError(error);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },

    async fetchTicketByCode(code) {
      this.loading = true;
      try {
        const response = await axiosWithToken.get(`/tickets/${code}`);
        return response.data;
      } catch (error) {
        this.error = handleError(error);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },

    async createTicket(payload) {
      this.loading = true;
      try {
        const response = await axiosWithToken.post("/tickets", payload);
        toastSuccess(`Ticket code ${response.data.code} has been created`);
        router.push({ name: "app.dashboard" });
      } catch (error) {
        this.error = handleError(error);
        toastError(this.error);
      } finally {
        this.loading = false;
      }
    },
  },
});
