<script setup>
import { onMounted, ref } from "vue";
import { useTransactionStore } from "@/stores/transaction";

const transactionStore = useTransactionStore();
const transactions = ref([]);

const getTransactions = async () => {
  const data = await transactionStore.fetchTransactionsByUser();
  if (data) {
    transactions.value = data;
  }
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString("en-US", {
    dateStyle: "medium",
    timeStyle: "short",
  });
};

const campaignImageResolver = (image) => {
  const BASE_URL = import.meta.env.VITE_BACKEND_URL;
  return image
    ? `${BASE_URL}/uploads/campaigns/${image}`
    : "https://placehold.co/400x250.png?text=No+Image";
};

onMounted(() => {
  getTransactions();
});
</script>

<template>
  <section class="container mx-auto pt-8">
    <div class="flex justify-between items-center mb-6">
      <div class="w-3/4 mr-6">
        <h2 class="text-4xl text-gray-900 mb-2 font-medium">Dashboard</h2>
        <ul class="flex mt-2">
          <li class="mr-6">
            <RouterLink
              to="/dashboard"
              class="text-gray-800 hover:text-gray-800"
            >
              Your Campaigns
            </RouterLink>
          </li>
          <li class="mr-6">
            <RouterLink
              to="/dashboard/transactions"
              class="text-gray-800 font-bold"
            >
              Your Transactions
            </RouterLink>
          </li>
        </ul>
      </div>
    </div>
    <hr />

    <!-- Content -->
    <div class="block mt-2 mb-2">
      <div
        v-for="transaction in transactions"
        :key="transaction.id"
        class="w-full lg:max-w-full lg:flex mb-4"
      >
        <!-- Campaign Thumbnail -->
        <div
          class="h-48 lg:h-auto lg:w-48 flex-none bg-cover rounded-t lg:rounded-t-none lg:rounded-l text-center overflow-hidden"
          :style="{
            backgroundColor: '#bbb',
            backgroundPosition: 'center',
            backgroundImage: `url(${campaignImageResolver(
              transaction.campaign_image_name
            )})`,
          }"
        ></div>

        <!-- Transaction Details -->
        <div
          class="w-full border-r border-b border-l border-gray-400 lg:border-l-0 lg:border-t lg:border-gray-400 bg-white rounded-b lg:rounded-b-none lg:rounded-r p-8 flex flex-col justify-between leading-normal"
        >
          <div>
            <div class="text-gray-900 font-bold text-xl mb-1">
              {{ transaction.campaign_name }}
            </div>
            <p class="text-sm text-gray-600 flex items-center mb-2">
              Rp.
              {{ new Intl.NumberFormat().format(transaction.amount) }} &middot;
              {{ formatDate(transaction.created_at) }} &middot;
              {{ transaction.status }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
