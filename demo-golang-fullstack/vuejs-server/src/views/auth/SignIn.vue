<script setup>
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { RouterLink } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const authStore = useAuthStore();
const { signIn } = authStore;
const { loading } = storeToRefs(authStore);

const form = ref({
  email: null,
  password: null,
});

const handleSubmit = async () => {
  await signIn(form.value);
};
</script>

<template>
  <div class="h-screen flex justify-center items-center">
    <!-- Background Image (Desktop Only) -->
    <div
      class="hidden md:block lg:w-1/3 bg-white h-full auth-background rounded-tr-3xl rounded-br-3xl shadow-lg"
    ></div>

    <!-- End:Form Section -->
    <div class="w-auto md:w-2/4 lg:w-2/3 flex justify-center items-center">
      <div
        class="w-full lg:w-1/2 px-10 py-8 bg-white/10 backdrop-blur-md rounded-2xl shadow-xl border border-white/20"
      >
        <form @submit.prevent="handleSubmit">
          <!-- Title -->
          <h2 class="font-bold mb-6 text-3xl text-white text-center">
            Sign In
          </h2>

          <!-- Email -->
          <div class="mb-4">
            <label class="block text-white text-sm mb-2">Email</label>
            <input
              type="email"
              v-model="form.email"
              id="email"
              name="email"
              placeholder="example@mail.com"
              class="w-full px-4 py-2 rounded-lg bg-white/20 text-white placeholder-gray-300 focus:outline-none focus:ring-2 focus:ring-sky-400"
            />
          </div>

          <!-- Password -->
          <div class="mb-6">
            <label class="block text-white text-sm mb-2">Password</label>
            <input
              type="password"
              v-model="form.password"
              id="password"
              name="password"
              placeholder="••••••••"
              class="w-full px-4 py-2 rounded-lg bg-white/20 text-white placeholder-gray-300 focus:outline-none focus:ring-2 focus:ring-sky-400"
            />
          </div>

          <!-- Sign In Button -->
          <button
            type="submit"
            :disabled="loading"
            :class="[
              'w-full py-3 bg-orange-500 hover:bg-orange-600 text-white font-medium rounded-lg transition-all duration-300 shadow-md',
              loading ? 'cursor-not-allowed' : '',
            ]"
          >
            <span v-if="loading">Loading...</span>
            <span v-else>Sign In</span>
          </button>
        </form>

        <!-- Forgot Password -->
        <p class="mt-4 text-white text-sm text-center">
          <a href="/forgot-password" class="text-white hover:underline">
            Forgot your password?
          </a>
        </p>

        <!-- Don't Have Account -->
        <p class="mt-6 text-sm text-white text-center">
          Don’t have an account?
          <RouterLink
            to="/auth/register"
            class="text-orange-400 hover:underline"
          >
            Create one
          </RouterLink>
        </p>
      </div>
    </div>
    <!-- End:Form Section -->
  </div>
</template>

<style scoped>
.auth-background {
  background-image: url("@/assets/images/sign-in-background.jpg");
  background-position: center;
  background-size: cover;
}
</style>
