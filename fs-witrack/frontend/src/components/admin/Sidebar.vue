<script setup>
import { HomeIcon, TicketIcon } from '@heroicons/vue/24/solid';
</script>

<template>
    <aside class="w-64 bg-white border-r border-gray-100">
        <nav class="mt-6">
            <RouterLink :to="{ name: 'admin.dashboard' }"
                class="flex items-center px-6 py-3 text-gray-600 hover:bg-gray-50 hover:border-l-4 hover:border-gray-200"
                :class="{
                    'bg-blue-50 border-l-4 border-blue-600': $route.name === 'admin.dashboard',
                }">
                <HomeIcon class="w-6 h-6 mr-2 text-gray-500" />
                Dashboard
            </RouterLink>
            <RouterLink :to="{ name: 'admin.tickets' }"
                class="flex items-center px-6 py-3 text-gray-600 hover:bg-gray-50 hover:border-l-4 hover:border-gray-200"
                :class="{
                    'bg-blue-50 border-l-4 border-blue-600': $route.name === 'admin.tickets',
                }">
                <TicketIcon class="w-6 h-6 mr-2 text-gray-500" />
                Tiket
            </RouterLink>
        </nav>
    </aside>
</template>