import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import HomeLayout from "@/layouts/HomeLayout.vue";
import Home from "@/views/Home.vue";
import SearchCampaign from "@/views/campaign/SearchCampaign.vue";
import Campaign from "@/views/campaign/Campaign.vue";
import Transaction from "@/views/transaction/Transaction.vue";
import AuthLayout from "@/layouts/AuthLayout.vue";
import SignUp from "@/views/auth/SignUp.vue";
import SignIn from "@/views/auth/SignIn.vue";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import Dashboard from "@/views/dashboard/Dashboard.vue";
import Campaigns from "@/views/dashboard/campaigns/Campaigns.vue";
import CampaignCreate from "@/views/dashboard/campaigns/CampaignCreate.vue";
import CampaignDetail from "@/views/dashboard/campaigns/CampaignDetail.vue";
import CampaignEdit from "@/views/dashboard/campaigns/CampaignEdit.vue";
import Transactions from "@/views/dashboard/transactions/Transactions.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "",
      component: HomeLayout,
      children: [
        {
          path: "",
          name: "home",
          component: Home,
          meta: { title: "", requiresUnauth: true },
        },
        {
          path: "campaigns",
          name: "campaigns.search",
          component: SearchCampaign,
          meta: { title: "Campaigns", requiresUnauth: true },
        },
        {
          path: "campaign/:id",
          name: "campaign",
          component: Campaign,
          meta: { title: "Campaign", requiresUnauth: true },
        },
        {
          path: "transaction",
          name: "transaction",
          component: Transaction,
          meta: { title: "Transaction", requiresAuth: true },
        },
      ],
    },
    {
      path: "/auth",
      component: AuthLayout,
      children: [
        {
          path: "",
          name: "signin",
          component: SignIn,
          meta: { title: "Sign In", requiresUnauth: true },
        },
        {
          path: "register",
          name: "signup",
          component: SignUp,
          meta: { title: "Sign Up", requiresUnauth: true },
        },
      ],
    },
    {
      path: "/dashboard",
      component: DashboardLayout,
      children: [
        {
          path: "",
          name: "dashboard",
          component: Dashboard,
          meta: { title: "Dashboard", requiresAuth: true },
        },
        {
          path: "campaigns",
          name: "campaigns",
          component: Campaigns,
          meta: { title: "Campaigns", requiresAuth: true },
        },
        {
          path: "campaigns/create",
          name: "campaigns.create",
          component: CampaignCreate,
          meta: { title: "Create Campaign", requiresAuth: true },
        },
        {
          path: "campaigns/:id",
          name: "campaigns.detail",
          component: CampaignDetail,
          meta: { title: "Detail Campaign", requiresAuth: true },
        },
        {
          path: "campaigns/:id/edit",
          name: "campaigns.edit",
          component: CampaignEdit,
          meta: { title: "Edit Campaign", requiresAuth: true },
        },
        {
          path: "transactions",
          name: "transactions",
          component: Transactions,
          meta: { title: "Transactions", requiresAuth: true },
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
      if (to.name !== "signin") {
        return next({ name: "signin" });
      }
    }
  }

  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      return next({ name: "signin" });
    }
    return next();
  }

  return next();
});

router.afterEach((to) => {
  const defaultTitle = "BackerHub";
  document.title = to.meta.title
    ? `${to.meta.title} | ${defaultTitle}`
    : defaultTitle;
});

export default router;
