<script setup>
import { ref } from 'vue';
import { storeToRefs } from 'pinia';
import { RouterLink } from "vue-router";
import { EyeIcon, EyeSlashIcon } from "@heroicons/vue/24/outline";

import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const { login } = authStore;
const { loading, error } = storeToRefs(authStore);

const form = ref({
    email: null,
    password: null,
});

const handleSubmit = async () => {
    await login(form.value);
}

const showPassword = ref(false);
</script>

<template>
    <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Email -->
        <div>
            <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
            <div class="mt-1 relative">
                <input v-model="form.email" type="email" id="email" name="email" placeholder="Your email address"
                    class="w-full px-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                    :class="{ 'border-red-500 ring-red-500': error?.email }" />
            </div>
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                <i data-feather="user" class="w-4 h-4 text-gray-400"></i>
            </div>
            <p class="mt-1 text-xs text-red-500" v-if="error?.email">
                {{ error.email }}
            </p>
        </div>

        <!-- Password -->
        <div>
            <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
            <div class="mt-1 relative">
                <input v-model="form.password" :type="showPassword ? 'text' : 'password'" id="password" name="password"
                    placeholder="Your password"
                    class="w-full px-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                    :class="{ 'border-red-500 ring-red-500': error?.password }" />
                <button type="button" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400"
                    @click="showPassword = !showPassword">
                    <EyeSlashIcon v-if="showPassword" class="w-5 h-5" />
                    <EyeIcon v-else class="w-5 h-5" />
                </button>
            </div>
            <p class="mt-1 text-xs text-red-500" v-if="error?.password">
                {{ error.password }}
            </p>
        </div>

        <!-- Submit -->
        <div>
            <button type="submit" :disabled="loading" :class="[
                'w-full flex justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white',
                loading ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'
            ]">
                <span v-if="loading">Loading...</span>
                <span v-else>Sign In</span>
            </button>
        </div>
    </form>

    <!-- Divider -->
    <div class="mt-6">
        <div class="relative">
            <div class="absolute inset-0 flex items-center">
                <div class="w-full border-t border-gray-200"></div>
            </div>
            <div class="relative flex justify-center text-sm">
                <span class="px-2 bg-white text-gray-500">Or</span>
            </div>
        </div>
    </div>

    <!-- Do not have an account -->
    <div class="mt-6 text-center text-sm">
        <span class="text-gray-600">Do not have an account? </span>
        <RouterLink to="/auth/register" class="font-medium text-blue-600 hover:text-blue-500">
            Sign up
        </RouterLink>
    </div>
</template>