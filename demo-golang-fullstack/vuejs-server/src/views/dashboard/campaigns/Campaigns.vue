<script setup>
import { onMounted, ref } from "vue";
import { useCampaignStore } from "@/stores/campaign";

const campaignStore = useCampaignStore();
const { fetchMyCampaigns } = campaignStore;

const campaigns = ref({});

const getMyCampaigns = async () => {
  const data = await fetchMyCampaigns();
  campaigns.value = data;
};

const campaignImageResolver = (image) => {
  const BASE_URL = import.meta.env.VITE_BACKEND_URL;
  return image
    ? `${BASE_URL}/uploads/campaigns/${image}`
    : "https://placehold.co/150x150.png?text=No+Image";
};

onMounted(async () => {
  await getMyCampaigns();
});
</script>

<template>
  <section class="container mx-auto pt-8">
    <!-- Headbar -->
    <div class="flex justify-between items-center mb-6">
      <div class="w-3/4 mr-6">
        <h2 class="text-4xl text-gray-900 mb-2 font-medium">Dashboard</h2>
        <ul class="flex mt-2">
          <li class="mr-6">
            <RouterLink to="/dashboard" class="text-gray-800 font-bold">
              Your Campaigns
            </RouterLink>
          </li>
          <li class="mr-6">
            <RouterLink
              to="/dashboard/transactions"
              class="text-gray-800 hover:text-gray-800"
            >
              Your Transactions
            </RouterLink>
          </li>
        </ul>
      </div>

      <div class="w-1/4 text-right">
        <RouterLink
          to="/dashboard/campaigns/create"
          class="bg-orange-300 hover:bg-orange-500 text-white font-bold py-4 px-4 rounded inline-flex items-center"
        >
          + Create Campaign
        </RouterLink>
      </div>
    </div>
    <hr />

    <!-- Content -->
    <div class="block mt-2 mb-2">
      <div
        v-for="campaign in campaigns"
        :key="campaign.id"
        class="w-full lg:max-w-full lg:flex mb-4"
      >
        <div
          class="border h-48 lg:h-auto lg:w-48 flex-none bg-cover rounded-t lg:rounded-t-none lg:rounded-l text-center overflow-hidden"
          :style="{
            backgroundImage: `url(${campaignImageResolver(
              campaign.image_name
            )})`,
            backgroundPosition: 'center',
            backgroundColor: '#bbb',
          }"
        ></div>

        <router-link
          :to="`/dashboard/campaigns/${campaign.id}`"
          class="w-full border-r border-b border-l border-gray-400 lg:border-l-0 lg:border-t lg:border-gray-400 bg-white rounded-b lg:rounded-b-none lg:rounded-r p-8 flex flex-col justify-between leading-normal"
        >
          <div class="mb-8">
            <div class="text-gray-900 font-bold text-xl mb-1">
              {{ campaign.title }}
            </div>
            <p class="text-sm text-gray-600 flex items-center mb-2">
              Rp.
              {{ new Intl.NumberFormat().format(campaign.goal_amount) }}
              &middot;
              {{
                (
                  (campaign.current_amount / campaign.goal_amount) *
                  100
                ).toFixed(0)
              }}%
            </p>
            <p class="text-gray-700 text-base">
              {{ campaign.short_description }}
            </p>
          </div>
          <div class="flex items-center">
            <button class="bg-green-500 text-white py-2 px-4 rounded">
              Detail
            </button>
          </div>
        </router-link>
      </div>
    </div>
  </section>
</template>
