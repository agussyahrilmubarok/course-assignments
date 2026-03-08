<script setup>
import { ref, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { useCampaignStore } from "@/stores/campaign";

const campaignStore = useCampaignStore();
const { fetchAllCampaigns } = campaignStore;

const campaigns = ref([]);
const searchQuery = ref("");

const getAllCampaigns = async () => {
  const data = await fetchAllCampaigns();
  campaigns.value = data || [];
};

const campaignImageResolver = (image) => {
  const BASE_URL = import.meta.env.VITE_BACKEND_URL;
  if (!image) {
    return "https://placehold.co/400x250?text=No+Image";
  }
  return `${BASE_URL}/uploads/campaigns/${image}`;
};

const filteredCampaigns = computed(() => {
  if (!searchQuery.value) return campaigns.value;
  return campaigns.value.filter((c) =>
    c.title.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

onMounted(async () => {
  await getAllCampaigns();
});
</script>

<template>
  <!-- Search Section -->
  <section class="bg-gray-100 py-10">
    <div class="container mx-auto px-5">
      <h2 class="text-3xl font-semibold text-gray-900 mb-4">Search Campaign</h2>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Type campaign..."
        class="w-full border border-gray-300 rounded-md px-5 py-3 focus:outline-none focus:ring-2 focus:ring-orange-400"
      />
    </div>
  </section>

  <!-- Campaign List -->
  <section
    v-if="filteredCampaigns && filteredCampaigns.length > 0"
    class="container mx-auto pt-12 px-5"
  >
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div
        v-for="campaign in filteredCampaigns"
        :key="campaign.id"
        class="card-project w-full p-5 border border-gray-300 rounded-20 shadow hover:shadow-lg transition-all"
      >
        <figure>
          <img
            :src="campaignImageResolver(campaign.image_name)"
            alt="Campaign Image"
            class="rounded-20 w-full h-56 object-cover"
          />
        </figure>
        <div class="mt-4">
          <h4 class="text-xl font-medium text-gray-900 mb-1">
            {{ campaign.title }}
          </h4>
          <p class="text-sm font-light text-gray-700 h-12">
            {{ campaign.short_description }}
          </p>

          <!-- Progress Bar -->
          <div class="relative pt-4">
            <div class="overflow-hidden h-2 mb-2 flex rounded bg-gray-200">
              <div
                :style="`width: ${
                  (campaign.current_amount / campaign.goal_amount) * 100
                }%`"
                class="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-purple-600 transition-all"
              ></div>
            </div>
            <div class="flex justify-between text-sm text-gray-600">
              <span>
                {{
                  Math.round(
                    (campaign.current_amount / campaign.goal_amount) * 100
                  )
                }}%
              </span>
              <span class="font-semibold">
                Rp {{ new Intl.NumberFormat().format(campaign.goal_amount) }}
              </span>
            </div>
          </div>

          <button
            @click="
              $router.push({ name: 'campaign', params: { id: campaign.id } })
            "
            class="text-center mt-5 w-full bg-orange-500 hover:bg-green-500 text-white font-semibold px-6 py-2 text-lg rounded-full"
          >
            Fund Now
          </button>
        </div>
      </div>
    </div>
  </section>

  <!-- Empty state -->
  <section v-else class="container mx-auto py-20 text-center text-gray-600">
    Campaign not found.
  </section>
</template>
