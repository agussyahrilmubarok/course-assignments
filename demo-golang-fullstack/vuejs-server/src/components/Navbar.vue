<script setup>
import { computed } from "vue";
import { storeToRefs } from "pinia";
import { useAuthStore } from "@/stores/auth";
import Logo from "@/assets/images/logo.svg";

const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

const menus = [
  { name: "Home", link: "/" },
  { name: "Campaigns", link: "/campaigns" },
];

const currentUser = computed(() => {
  const BASE_URL = import.meta.env.VITE_BACKEND_URL;
  if (!user.value) {
    return {
      name: "Guest",
      image_name: "https://i.pravatar.cc/100",
    };
  }

  return {
    name: user.value.name,
    image_name: user.value.image_name
      ? `${BASE_URL}/uploads/users/${user.value.image_name}`
      : "https://i.pravatar.cc/100",
  };
});

const loggedIn = computed(() => authStore.isAuthenticated);

const logout = () => {
  authStore.signOut();
};
</script>

<template>
  <header class="flex items-center w-full px-6 py-4">
    <!-- Logo -->
    <div class="h-[54px] pr-5">
      <img :src="Logo" alt="Logo" class="h-full" />
    </div>

    <!-- Menu -->
    <ul class="flex items-center">
      <li v-for="menu in menus" :key="menu.name">
        <a :href="menu.link" class="text-white hover:text-teal-500 text-lg px-4 py-3">
          {{ menu.name }}
        </a>
      </li>
    </ul>

    <!-- Auth Buttons -->
    <ul v-if="!loggedIn" class="flex ml-auto items-center">
      <li>
        <RouterLink to="/auth/register"
          class="inline-block bg-transparent border border-white hover:bg-white hover:bg-opacity-25 text-white font-light w-40 text-center px-6 py-1 text-lg rounded-full mr-4">
          Sign Up
        </RouterLink>
      </li>
      <li>
        <RouterLink to="/auth"
          class="inline-block bg-transparent border border-white hover:bg-white hover:bg-opacity-25 text-white font-light w-40 text-center px-6 py-1 text-lg rounded-full">
          Sign In
        </RouterLink>
      </li>
    </ul>

    <!-- Dropdown User -->
    <div v-else class="flex ml-auto relative">
      <div class="dropdown inline-block">
        <button class="bg-white text-gray-700 font-semibold py-2 px-4 rounded inline-flex items-center">
          <img :src="currentUser.image_name" alt="User Avatar" class="h-8 w-8 rounded-full mr-2" />
          <span class="mr-1">{{ currentUser.name }}</span>
          <svg class="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
            <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
          </svg>
        </button>
        <ul class="dropdown-menu absolute hidden bg-white text-gray-700 pt-1 shadow w-48 right-0">
          <li>
            <RouterLink to="/dashboard" class="block px-4 py-2 hover:bg-gray-100 hover:text-orange-500">My Dashboard
            </RouterLink>
          </li>
          <li>
            <button @click="logout" class="w-full text-left block px-4 py-2 hover:bg-gray-100 hover:text-orange-500">
              Logout
            </button>
          </li>
        </ul>
      </div>
    </div>
  </header>
</template>

<style scoped>
.dropdown:hover .dropdown-menu {
  display: block;
}
</style>
