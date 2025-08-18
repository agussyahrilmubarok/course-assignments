import { createRouter, createWebHistory } from "vue-router";

import { useAuthStore } from "@/stores/auth";
import AuthLayout from "@/layouts/AuthLayout.vue";
import Login from "@/views/auth/Login.vue";
import Register from "@/views/auth/Register.vue";
import AppLayout from "@/layouts/AppLayout.vue";
import AppDashboard from "@/views/app/AppDashboard.vue";
import AppTicketCreate from "@/views/app/tickets/AppTicketCreate.vue";
import AppTicketDetail from "@/views/app/tickets/AppTicketDetail.vue";
import AdminLayout from "@/layouts/AdminLayout.vue";
import AdminDashboard from "@/views/admin/AdminDashboard.vue";
import TicketList from "@/views/admin/tickets/TicketList.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/auth",
      component: AuthLayout,
      children: [
        {
          path: "",
          name: "login",
          component: Login,
          meta: { title: "Sign In", requiresUnauth: true },
        },
        {
          path: "register",
          name: "register",
          component: Register,
          meta: { title: "Sign Up", requiresUnauth: true },
        },
      ],
    },
    {
      path: "",
      component: AppLayout,
      children: [
        {
          path: "",
          name: "app.dashboard",
          component: AppDashboard,
          meta: { title: "Dashboard", requiresAuth: true },
        },
        {
          path: "tickets/create",
          name: "app.tickets.create",
          component: AppTicketCreate,
          meta: { title: "Create Ticket", requiresAuth: true },
        },
        {
          path: "tickets/:code",
          name: "app.tickets.detail",
          component: AppTicketDetail,
          meta: { title: "Detail Ticket", requiresAuth: true },
        },
      ],
    },
    {
      path: "/admin",
      component: AdminLayout,
      children: [
        {
          path: "dashboard",
          name: "admin.dashboard",
          component: AdminDashboard,
          meta: { title: "Admin Dashboard", requiresAuth: true },
        },
        {
          path: "tickets",
          name: "admin.tickets",
          component: TicketList,
          meta: {
            requiresAuth: true,
            title: "Tickets",
          },
        },
      ],
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();

  // Load user if we already have a token but no user
  if (authStore.token && !authStore.user) {
    try {
      await authStore.checkAuth();
    } catch {
      return next({ name: "login" });
    }
  }

  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      return next({ name: "login" });
    }
    return next();
  }

  if (to.meta.requiresUnauth && authStore.isAuthenticated) {
    return next({
      name: authStore.isAdmin ? "admin.dashboard" : "app.dashboard",
    });
  }

  return next();
});

router.afterEach((to) => {
  const defaultTitle = "WiTrack";
  document.title = to.meta.title
    ? `${to.meta.title} | ${defaultTitle}`
    : defaultTitle;
});

export default router;
