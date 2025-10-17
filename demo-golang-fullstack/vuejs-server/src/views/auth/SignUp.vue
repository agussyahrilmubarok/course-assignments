<script setup>
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { RouterLink } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const authStore = useAuthStore();
const { signUp } = authStore;
const { loading } = storeToRefs(authStore);

const form = ref({
  name: null,
  occupation: null,
  email: null,
  password: null,
});

const handleSubmit = async () => {
  await signUp(form.value);
};
</script>

<template>
  <div class="h-screen flex justify-center items-center">
    <!-- Background Image (Desktop Only) -->
    <div
      class="hidden md:block lg:w-1/3 bg-white h-full auth-background rounded-tr-3xl rounded-br-3xl shadow-lg"
    ></div>
    <!-- Begin:Form Section -->
    <div class="w-auto md:w-2/4 lg:w-2/3 flex justify-center items-center">
      <div
        class="w-full lg:w-1/2 px-10 py-8 bg-white/10 backdrop-blur-md rounded-2xl shadow-xl border border-white/20"
      >
        <form @submit.prevent="handleSubmit">
          <!-- Title -->
          <h2 class="font-bold mb-6 text-3xl text-white text-center">
            Create an Account
          </h2>

          <!-- Full Name -->
          <div class="mb-4">
            <label class="block text-white text-sm mb-2">Name</label>
            <input
              type="text"
              v-model="form.name"
              id="name"
              name="name"
              placeholder="John Doe"
              class="w-full px-4 py-2 rounded-lg bg-white/20 text-white placeholder-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-400"
            />
          </div>

          <!-- Occupation -->
          <div class="mb-4">
            <label class="block text-white text-sm mb-2">Occupation</label>
            <input
              type="text"
              v-model="form.occupation"
              id="occupation"
              name="occupation"
              placeholder="Your profession"
              class="w-full px-4 py-2 rounded-lg bg-white/20 text-white placeholder-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-400"
            />
          </div>

          <!-- Email -->
          <div class="mb-4">
            <label class="block text-white text-sm mb-2">Email</label>
            <input
              type="email"
              v-model="form.email"
              id="email"
              name="email"
              placeholder="example@mail.com"
              class="w-full px-4 py-2 rounded-lg bg-white/20 text-white placeholder-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-400"
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
              class="w-full px-4 py-2 rounded-lg bg-white/20 text-white placeholder-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-400"
            />
          </div>

          <!-- Sign Up Button -->
          <button
            type="submit"
            :disabled="loading"
            :class="[
              'w-full py-3 bg-orange-500 hover:bg-orange-600 text-white font-medium rounded-lg transition-all duration-300 shadow-md',
              loading ? 'cursor-not-allowed' : '',
            ]"
          >
            <span v-if="loading">Loading...</span>
            <span v-else>Sign Up</span>
          </button>
        </form>

        <!-- Already Have Account -->
        <p class="mt-6 text-sm text-white text-center">
          Already have an account?
          <RouterLink to="/auth" class="text-orange-400 hover:underline">
            Sign In
          </RouterLink>
        </p>
      </div>
    </div>
    <!-- End:Form Section -->
  </div>
</template>

<style scoped>
.auth-background {
  background-image: url("@/assets/images/sign-up-background.jpg");
  background-position: center;
  background-size: cover;
}
</style>
