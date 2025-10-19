<script setup>
import { ref, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useRoute } from "vue-router";
import { useCampaignStore } from "@/stores/campaign";
import { useTransactionStore } from "@/stores/transaction";

const campaignStore = useCampaignStore();
const { fetchMyCampaignById, uploadMyCampaignImage } = campaignStore;

const route = useRoute();
const campaign = ref({});
const fileInput = ref(null);
const formImage = ref({
  selectedFile: null,
  isPrimary: false, // checkbox state
});

const getCampaign = async () => {
  const data = await fetchMyCampaignById(route.params.id);
  campaign.value = data;
};

const selectFile = () => {
  formImage.value.selectedFile = fileInput.value?.files[0] || null;
};

const upload = async () => {
  if (!formImage.value.selectedFile) return;

  try {
    await uploadMyCampaignImage(route.params.id, {
      file: formImage.value.selectedFile,
      is_primary: formImage.value.isPrimary,
    });

    await getCampaign();
  } catch (err) {
    toastError("Failed to upload image");
  } finally {
    fileInput.value.value = "";
    formImage.value.selectedFile = null;
    formImage.value.isPrimary = false;
  }
};

const transactionStore = useTransactionStore();
const { fetchTransactionsByCampaign } = transactionStore;

const transactions = ref([]);

const getTransactions = async () => {
  const data = await fetchTransactionsByCampaign(route.params.id);
  transactions.value = data || [];
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
    return b.is_primary - a.is_primary; // true (1) grather than false (0)
  });
});

onMounted(async () => {
  await getCampaign();
  await getTransactions();
});
</script>

<template>
  <!-- Campaign details -->
  <section class="container mx-auto pt-8">
    <div class="flex justify-between items-center">
      <div class="w-full mr-6">
        <h2 class="text-4xl text-gray-900 mb-2 font-medium">Dashboard</h2>
      </div>
    </div>

    <div class="flex justify-between items-center">
      <div class="w-3/4 mr-6">
        <h3 class="text-2xl text-gray-900 mb-4">Campaign Details</h3>
      </div>
      <div class="w-1/4 text-right">
        <router-link
          :to="`/dashboard/campaigns/${campaign.id}/edit`"
          class="inline-block"
        >
          <button
            class="bg-green-500 text-white font-bold px-4 py-1 rounded inline-flex items-center hover:bg-green-600 transition"
          >
            Edit
          </button>
        </router-link>
      </div>
    </div>

    <!-- Campaign Card -->
    <div class="block mb-2">
      <div class="w-full lg:max-w-full lg:flex mb-4">
        <div
          class="w-full border border-gray-400 bg-white rounded p-8 flex flex-col justify-between leading-normal"
        >
          <div>
            <div class="text-gray-900 font-bold text-xl mb-2">
              {{ campaign.title }}
            </div>

            <p class="text-sm font-bold mb-1">Short Description</p>
            <p class="text-gray-700">{{ campaign.short_description }}</p>

            <p class="text-sm font-bold mb-1 mt-4">Description</p>
            <p class="text-gray-700">{{ campaign.description }}</p>

            <p class="text-sm font-bold mb-1 mt-4">What Will Funders Get</p>
            <ul class="list-disc ml-5">
              <li v-for="(perk, index) in perkList" :key="index">
                {{ perk }}
              </li>
            </ul>

            <p class="text-sm font-bold mb-1 mt-4">Goal Amount</p>
            <p class="text-4xl text-gray-700">
              Rp.
              {{ new Intl.NumberFormat().format(campaign.goal_amount) }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Gallery -->
    <div class="flex justify-between items-center">
      <div class="w-2/4 mr-6">
        <h3 class="text-2xl text-gray-900 mb-4 mt-5">Gallery</h3>
      </div>
      <div class="w-2/4 text-right space-x-2">
        <input
          type="file"
          ref="fileInput"
          @change="selectFile"
          class="border p-1 rounded overflow-hidden"
        />
        <label class="inline-flex items-center ml-2">
          <input type="checkbox" v-model="formImage.isPrimary" class="mr-1" />
          Set as Primary
        </label>
        <button
          @click="upload"
          :disabled="!formImage.selectedFile"
          class="px-4 py-2 rounded inline-flex items-center font-bold text-white transition bg-green-500 hover:bg-green-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          Upload
        </button>
      </div>
    </div>

    <div
      class="flex space-x-4 overflow-x-auto pb-2 scrollbar-thin scrollbar-thumb-gray-400 scrollbar-track-gray-200"
    >
      <div
        v-for="image in sortedImages"
        :key="image.image_name"
        class="relative flex-none w-64 h-40 border rounded overflow-hidden"
        :class="{
          'border-green-500 border-4': image.is_primary,
          'border-gray-400': !image.is_primary,
        }"
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

    <!-- Transactions -->
    <div class="flex justify-between items-center">
      <div class="w-3/4 mr-6">
        <h3 class="text-2xl text-gray-900 mb-4 mt-5">Transaction History</h3>
      </div>
    </div>

    <div class="block mb-2">
      <div
        v-for="transaction in transactions"
        :key="transaction.id"
        class="w-full lg:max-w-full lg:flex mb-4"
      >
        <div
          class="w-full border border-gray-400 bg-white rounded p-8 flex flex-col justify-between leading-normal"
        >
          <div>
            <div class="text-gray-900 font-bold text-xl mb-1">
              {{ transaction.user_name }}
            </div>
            <p class="text-sm text-gray-600 mb-2">
              {{ transaction.user_email }}
            </p>
            <p class="text-sm text-gray-600 mb-2">
              Rp. {{ new Intl.NumberFormat().format(transaction.amount) }} ·
              {{
                new Date(transaction.created_at).toLocaleString("en-US", {
                  dateStyle: "medium",
                  timeStyle: "short",
                })
              }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
