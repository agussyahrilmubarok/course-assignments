<script setup>
import { ref, onMounted } from "vue";
import { useCampaignStore } from "@/stores/campaign";
import HeroImage from "@/assets/images/hero-image@2x.png";
import LineStep from "@/assets/images/line-step.svg";
import LineStep1 from "@/assets/images/step-1-illustration.svg";
import LineStep2 from "@/assets/images/step-2-illustration.svg";
import LineStep3 from "@/assets/images/step-3-illustration.svg";

const campaignStore = useCampaignStore();
const { fetchTopCampaigns } = campaignStore;

const campaigns = ref({});

const getTopCampaigns = async () => {
  const data = await fetchTopCampaigns();
  campaigns.value = data;
};

const campaignImageResolver = (image) => {
  const BASE_URL = import.meta.env.VITE_BACKEND_URL;
  if (!image) {
    return "https://placehold.co/400x250.png?text=No+Image";
  }
  return `${BASE_URL}/uploads/campaigns/${image}`;
};

onMounted(async () => {
  await getTopCampaigns();
});
</script>

<template>
  <!-- HeroSection -->
  <section class="landing-hero pt-5 pb-5 bg-[#3B41E3]">
    <div class="container mx-auto relative">
      <div class="flex items-center pt-10 px-5 md:px-0">
        <div class="w-1/2">
          <h1 class="text-4xl text-white mb-5">
            We helps <u class="hero-underline">startup</u> to
            <br />
            getting started &
            <u class="hero-underline">funding</u> <br />
            their truly needs
          </h1>
          <p class="text-white text-xl font-light mb-8">
            Fund the best idea to become <br />
            a real product and be the contributor
          </p>
          <button
            @click="$router.push({ path: '/auth/register' })"
            class="block bg-orange-600 hover:bg-green-600 text-white font-semibold px-12 py-3 text-xl rounded-full"
          >
            Create a Campaign
          </button>
        </div>
        <div class="w-1/2 flex justify-center">
          <img :src="HeroImage" alt="project" />
        </div>
      </div>
    </div>
  </section>

  <!-- StepSection -->
  <section class="container mx-auto pt-24">
    <div class="flex justify-between items-center mb-10">
      <div class="w-auto">
        <h2 class="text-3xl text-gray-900 mb-8">
          Only 3 steps to execute <br />
          your bright ideas
        </h2>
      </div>
    </div>
    <div class="flex">
      <div class="w-full px-56 mb-5">
        <img :src="LineStep" alt="" class="w-full" />
      </div>
    </div>
    <div class="flex justify-between items-center text-center">
      <div class="w-1/3">
        <figure class="flex justify-center items-center">
          <img :src="LineStep1" alt="" class="h-30 mb-8" />
        </figure>
        <div class="step-content">
          <h3 class="font-medium">Create an Account</h3>
          <p class="font-light">
            Sign Up account and start <br />funding project
          </p>
        </div>
      </div>
      <div class="w-1/3">
        <figure class="flex justify-center items-center -mt-24">
          <img :src="LineStep2" alt="" class="h-30 mb-8" />
        </figure>
        <div class="step-content">
          <h3 class="font-medium">Open Project</h3>
          <p class="font-light">
            Choose some project idea, <br />
            and start funding
          </p>
        </div>
      </div>
      <div class="w-1/3">
        <figure class="flex justify-center items-center -mt-48">
          <img :src="LineStep3" alt="" class="h-30 mb-8" />
        </figure>
        <div class="step-content">
          <h3 class="font-medium">Execute</h3>
          <p class="font-light">
            Time to makes dream <br />
            comes true
          </p>
        </div>
      </div>
    </div>
  </section>

  <!-- Top Campaigns -->
  <section
    v-if="campaigns && campaigns.length > 0"
    class="container mx-auto pt-24"
  >
    <div class="flex justify-between items-center">
      <div class="w-auto">
        <h2 class="text-3xl text-gray-900 mb-8">
          New projects you can <br />
          taken care of
        </h2>
      </div>
      <div class="w-auto mt-5">
        <a class="text-gray-900 hover:underline text-md font-medium" href="">
          View All
        </a>
      </div>
    </div>

    <div class="grid grid-cols-3 gap-4 mt-3">
      <div
        v-for="campaign in campaigns"
        :key="campaign.id"
        class="card-project w-full p-5 border border-gray-500 rounded-20"
      >
        <div class="item">
          <figure class="item-image">
            <img
              :src="campaignImageResolver(campaign.image_name)"
              alt="Campaign Image"
              class="rounded-20 w-full h-56 object-cover"
            />
          </figure>
          <div class="item-meta">
            <h4 class="text-2xl font-medium text-gray-900 mt-5">
              {{ campaign.title }}
            </h4>
            <p class="text-md font-light text-gray-900 h-12">
              {{ campaign.short_description }}
            </p>
            <div class="relative pt-4 progress-bar">
              <div class="overflow-hidden h-2 mb-4 flex rounded bg-gray-200">
                <div
                  :style="`width: ${
                    (campaign.current_amount / campaign.goal_amount) * 100
                  }%`"
                  class="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-purple-600 transition-all"
                ></div>
              </div>
            </div>
            <div class="flex progress-info">
              <div>
                {{
                  Math.round(
                    (campaign.current_amount / campaign.goal_amount) * 100
                  )
                }}%
              </div>
              <div class="ml-auto font-semibold">
                Rp {{ new Intl.NumberFormat().format(campaign.goal_amount) }}
              </div>
            </div>
          </div>
          <button
            @click="
              $router.push({ name: 'campaign', params: { id: campaign.id } })
            "
            class="text-center mt-5 button-cta block w-full bg-orange-500 hover:bg-green-500 text-white font-semibold px-6 py-2 text-lg rounded-full"
          >
            Fund Now
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
