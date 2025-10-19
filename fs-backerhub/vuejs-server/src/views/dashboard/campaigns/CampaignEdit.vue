<script setup>
import { ref, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useRoute } from "vue-router";
import { useCampaignStore } from "@/stores/campaign";

const route = useRoute();
const campaignStore = useCampaignStore();
const { fetchMyCampaignById, updateMyCampaign } = campaignStore;
const { loading } = storeToRefs(campaignStore);

const form = ref({
  title: "",
  short_description: "",
  description: "",
  goal_amount: 0,
  perks: "",
});

const handleUpdate = async () => {
  await updateMyCampaign(route.params.id, form.value);
};

onMounted(async () => {
  const data = await fetchMyCampaignById(route.params.id);
  if (data) {
    form.value = {
      title: data.title,
      short_description: data.short_description,
      description: data.description,
      goal_amount: data.goal_amount,
      perks: data.perks,
    };
  }
});
</script>

<template>
  <section class="container mx-auto pt-8">
    <div class="flex justify-between items-center">
      <div class="w-full mr-6">
        <h2 class="text-4xl text-gray-900 mb-2 font-medium">Dashboard</h2>
      </div>
    </div>

    <form @submit.prevent="handleUpdate">
      <div class="flex justify-between items-center">
        <div class="w-3/4 mr-6">
          <h3 class="text-2xl text-gray-900 mb-4">
            Edit Campaign "{{ form.title }}"
          </h3>
        </div>
        <div class="w-1/4 text-right">
          <button
            type="submit"
            :disabled="loading"
            :class="[
              'bg-green-500 hover:bg-green-600 text-white font-bold px-4 py-1 rounded inline-flex items-center',
              loading ? 'cursor-not-allowed' : '',
            ]"
          >
            <span v-if="loading">Updating...</span>
            <span v-else>Update</span>
          </button>
        </div>
      </div>

      <div class="block mb-2">
        <div class="w-full lg:max-w-full lg:flex mb-4">
          <div
            class="w-full border border-gray-400 bg-white rounded p-8 flex flex-col justify-between leading-normal"
          >
            <form class="w-full">
              <div class="flex flex-wrap -mx-3 mb-6">
                <!-- Campaign Name -->
                <div class="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
                  >
                    Campaign Name
                  </label>
                  <input
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                    type="text"
                    placeholder="e.g., Build a Gunpla for My Wife"
                    v-model="form.title"
                  />
                </div>

                <!-- Price -->
                <div class="w-full md:w-1/2 px-3">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
                  >
                    Goal Amount
                  </label>
                  <input
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                    type="number"
                    placeholder="e.g., 200000"
                    v-model.number="form.goal_amount"
                  />
                </div>

                <!-- Short Description -->
                <div class="w-full px-3">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 mt-3"
                  >
                    Short Description
                  </label>
                  <input
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                    type="text"
                    placeholder="A short description of your campaign"
                    v-model="form.short_description"
                  />
                </div>

                <!-- Perks -->
                <div class="w-full px-3">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
                  >
                    What Backers Will Get
                  </label>
                  <input
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                    type="text"
                    placeholder="e.g., T-shirt, Mug, Poster"
                    v-model="form.perks"
                  />
                </div>

                <!-- Description -->
                <div class="w-full px-3">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
                  >
                    Description
                  </label>
                  <textarea
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                    placeholder="Write a detailed description of your campaign"
                    v-model="form.description"
                  ></textarea>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </form>
  </section>
</template>
