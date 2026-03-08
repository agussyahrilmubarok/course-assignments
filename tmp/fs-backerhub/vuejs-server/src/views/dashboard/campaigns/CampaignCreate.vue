<script setup>
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { useCampaignStore } from "@/stores/campaign";

const campaignStore = useCampaignStore();
const { createMyCampaign } = campaignStore;
const { loading } = storeToRefs(campaignStore);

const form = ref({
  title: "",
  short_description: "",
  description: "",
  goal_amount: 0,
  perks: "",
});

const handleSubmit = async () => {
  await createMyCampaign(form.value);
};
</script>

<template>
  <section class="container mx-auto pt-8">
    <div class="flex justify-between items-center">
      <div class="w-full mr-6">
        <h2 class="text-4xl text-gray-900 mb-2 font-medium">Dashboard</h2>
      </div>
    </div>

    <form @submit.prevent="handleSubmit">
      <div class="flex justify-between items-center">
        <div class="w-3/4 mr-6">
          <h3 class="text-2xl text-gray-900 mb-4">Create New Projects</h3>
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
            <span v-if="loading">Saving...</span>
            <span v-else>Save</span>
          </button>
        </div>
      </div>

      <div class="block mb-2">
        <div class="w-full lg:max-w-full lg:flex mb-4">
          <div
            class="w-full border border-gray-400 bg-white rounded p-8 flex flex-col justify-between leading-normal"
          >
            <div class="w-full">
              <div class="flex flex-wrap -mx-3 mb-6">
                <!-- Campaign Name -->
                <div class="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
                  >
                    Campaign Name
                  </label>
                  <input
                    type="text"
                    v-model="form.title"
                    id="title"
                    name="title"
                    placeholder="Example: Save the Cats"
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                  />
                </div>

                <!-- Goal Amount -->
                <div class="w-full md:w-1/2 px-3">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
                  >
                    Goal Amount
                  </label>
                  <input
                    type="number"
                    v-model.number="form.goal_amount"
                    id="goal_amount"
                    name="goal_amount"
                    placeholder="Example: 200000"
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
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
                    type="text"
                    v-model="form.short_description"
                    id="short_description"
                    name="short_description"
                    placeholder="Short description of your project"
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                  />
                </div>

                <!-- Perks -->
                <div class="w-full px-3">
                  <label
                    class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2"
                  >
                    What will backers get
                  </label>
                  <input
                    type="text"
                    v-model="form.perks"
                    id="perks"
                    name="perks"
                    placeholder="Example: Stickers, T-shirt, Mug"
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
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
                    v-model="form.description"
                    id="description"
                    name="description"
                    placeholder="Write a long description for your project"
                    class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500"
                  ></textarea>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </form>
  </section>
</template>
