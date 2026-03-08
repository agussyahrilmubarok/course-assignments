<script setup>
import { ref } from 'vue';
import { storeToRefs } from 'pinia';
import { WifiIcon, ChevronDownIcon, UserIcon } from '@heroicons/vue/24/solid';

import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const { user } = storeToRefs(authStore);
const { logout } = authStore;

const showUserMenu = ref(false);

const handleLogout = async () => {
    await logout()
}

const toggleUserMenu = () => {
    showUserMenu.value = !showUserMenu.value
};

</script>

<template>
    <nav class="bg-white shadow-sm z-10">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">

                <!-- Logo -->
                <div class="flex items-center">
                    <RouterLink :to="{ name: 'app.dashboard' }" class="flex items-center">
                        <WifiIcon class="w-8 h-8 text-blue-600" />
                        <span class="ml-2 text-xl font-bold text-blue-600">WiTrack</span>
                    </RouterLink>
                </div>

                <!-- Right Menu -->
                <div class="flex items-center space-x-4">

                    <!-- TODO: Notification -->
                    <!-- <button @click="clearNotifications"
                        class="relative p-2 text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-full">
                        <i data-feather="bell" class="w-6 h-6"></i>
                        <span v-if="unreadCount > 0"
                            class="absolute -top-0.5 -right-0.5 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                            {{ unreadCount }}
                        </span>
                    </button> -->

                    <!-- User Menu -->
                    <div class="relative">
                        <button @click="toggleUserMenu()"
                            class="flex items-center bg-gray-50 px-4 py-2 rounded-full hover:bg-gray-100">
                            <img :src="`https://ui-avatars.com/api/?name=${user?.name}&background=0D8ABC&color=fff`"
                                alt="Profile" class="w-8 h-8 rounded-full" />
                            <span class="ml-2 text-sm font-medium text-gray-700">{{ user?.name }}</span>
                            <ChevronDownIcon class="w-4 h-4 ml-2 text-gray-500" />
                        </button>

                        <!-- Dropdown menu -->
                        <div v-if="showUserMenu"
                            class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-100 py-1 z-50">
                            <a href="#" class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">
                                <UserIcon class="w-5 h-5 mr-2 text-gray-500" />
                                Profile
                            </a>
                            <div class="border-t border-gray-100 my-1"></div>
                            <button type="button" @click="handleLogout"
                                class="w-full px-4 py-2 text-sm text-red-600 hover:bg-gray-50 cursor-pointer">
                                Logout
                            </button>
                        </div>
                    </div>

                </div>
            </div>
        </div>
    </nav>
</template>
