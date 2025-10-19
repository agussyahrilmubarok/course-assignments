<script setup>
import { ref, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useRoute } from "vue-router";
import { toastSuccess, toastError } from "@/helpers/toastHelper";
import { useCampaignStore } from "@/stores/campaign";
import { useTransactionStore } from "@/stores/transaction";

const campaignStore = useCampaignStore();
const { fetchCampaignById } = campaignStore;

const route = useRoute();
const campaign = ref({});
const defaultImage = ref("");

const getCampaign = async () => {
  const data = await fetchCampaignById(route.params.id);
  if (data) {
    campaign.value = data;

    const primaryImg = data.campaign_images?.find((img) => img.is_primary);

    if (primaryImg) {
      defaultImage.value = campaignImageResolver(primaryImg.image_name);
    } else if (data.campaign_images?.length > 0) {
      defaultImage.value = campaignImageResolver(
        data.campaign_images[0].image_name
      );
    } else {
      defaultImage.value = "https://placehold.co/400x250.png?text=No+Image";
    }
  }
};

function changeImage(url) {
  defaultImage.value = campaignImageResolver(url);
}

const transactionStore = useTransactionStore();
const { donationCampaign } = transactionStore;

const formFund = ref({
  amount: 0,
  campaign_id: route.params.id,
});

const fund = async () => {
  if (!formFund.value.amount || formFund.value.amount <= 0) {
    toastError("Please enter a valid amount");
    return;
  }

  try {
    const transactionId = await donationCampaign(formFund.value);
    if (transactionId) {
      toastSuccess("Fund created, please complete payment in new tab");
    }
  } catch (err) {
    toastError("Funding failed");
  }
};

const perkList = computed(() => {
  return (campaign.value.perks || "")
    .split(",")
    .map((p) => p.trim())
    .filter((p) => p.length > 0);
});

const campaignImageResolver = (image) => {
  const BASE_URL = import.meta.env.VITE_BACKEND_URL;
  return image
    ? `${BASE_URL}/uploads/campaigns/${image}`
    : "https://placehold.co/400x250.png?text=No+Image";
};

const sortedImages = computed(() => {
  if (!campaign.value?.campaign_images) return [];
  return [...campaign.value.campaign_images].sort((a, b) => {
    return b.is_primary - a.is_primary;
  });
});

const progress = computed(() => {
  if (!campaign.value.goal_amount) return 0;
  return Math.round(
    (campaign.value.current_amount / campaign.value.goal_amount) * 100
  );
});

onMounted(async () => {
  await getCampaign();
});
</script>

<template>
  <section class="container project-container mx-auto">
    <div class="flex mt-3">
      <!-- Left: Main Image + Thumbnails -->
      <div class="w-3/4 mr-6">
        <div class="bg-white p-3 mb-3 border border-gray-400 rounded-2xl">
          <figure class="item-image">
            <img
              :src="defaultImage"
              alt=""
              class="rounded-2xl w-full h-[400px] object-cover"
            />
          </figure>
        </div>

        <!-- Thumbnails -->
        <div
          class="flex flex-nowrap space-x-4 overflow-x-auto pb-2 scrollbar-thin scrollbar-thumb-gray-400 scrollbar-track-gray-200"
        >
          <div
            v-for="image in sortedImages"
            :key="image.image_name"
            class="relative flex-none w-64 h-40 border rounded overflow-hidden"
            :class="{
              'border-green-500 border-4': image.is_primary,
              'border-gray-400': !image.is_primary,
            }"
            @click="changeImage(image.image_name)"
          >
            <img
              :src="campaignImageResolver(image.image_name)"
              alt="campaign image"
              class="w-full h-full object-cover rounded"
            />

            <div
              v-if="image.is_primary"
              class="absolute top-2 left-2 bg-green-500 text-white text-xs px-2 py-1 rounded"
            >
              Primary
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Project Leader + Funding -->
      <div class="w-1/4">
        <div
          class="bg-white w-full p-5 border border-gray-400 rounded-2xl sticky"
          style="top: 15px"
        >
          <h3 class="font-bold">Support This Campaign</h3>

          <h4 class="mt-5 font-semibold">What will you get:</h4>
          <ul class="list-disc ml-5 mt-3">
            <li v-for="(perk, index) in perkList" :key="index">
              {{ perk }}
            </li>
          </ul>

          <!-- Input fund -->
          <input
            type="number"
            class="border border-gray-500 block w-full px-6 py-3 mt-4 rounded-full text-gray-800"
            placeholder="Amount in Rp"
            v-model.number="formFund.amount"
            @keyup.enter="fund"
          />
          <button
            @click="fund"
            class="text-center mt-3 block w-full bg-[#FF872E] hover:bg-[#1ABC9C] text-white font-medium px-6 py-3 text-md rounded-full"
          >
            Fund Now
          </button>
        </div>
      </div>
    </div>
  </section>

  <!-- Campaign Details -->
  <section class="container mx-auto pt-8">
    <h2 class="text-4xl text-gray-900 mb-2 font-medium">
      {{ campaign.title }}
    </h2>
    <p class="font-light text-xl mb-5">{{ campaign.short_description }}</p>

    <!-- Progress bar -->
    <div class="relative progress-bar">
      <div
        class="overflow-hidden mb-4 text-xs flex rounded-full bg-gray-200 h-6"
      >
        <div
          :style="'width: ' + progress + '%'"
          class="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-[#3B41E3]"
        ></div>
      </div>
    </div>

    <div class="flex progress-info mb-6">
      <div class="text-2xl">{{ progress }}%</div>
      <div class="ml-auto font-semibold text-2xl">
        Rp {{ new Intl.NumberFormat().format(campaign.goal_amount) }}
      </div>
    </div>

    <p class="font-light text-xl mb-5">{{ campaign.description }}</p>
  </section>
</template>
