<script setup>
import { ref, onMounted, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { debounce } from 'lodash';
import { DateTime } from 'luxon';
import { PlusIcon, MagnifyingGlassIcon, ClockIcon, ChatBubbleLeftIcon } from '@heroicons/vue/24/solid';
import { useTicketStore } from '@/stores/ticket';

const ticketStore = useTicketStore();
const { fetchMyTickets } = ticketStore;
const { tickets } = storeToRefs(ticketStore);

const filters = ref({
    search: '',
    status: '',
    priority: '',
});

// Watch filters dan fetch ticket with debounce
watch(
    filters,
    debounce(async () => {
        await fetchMyTickets(filters.value);
    }, 300),
    { deep: true },
);

onMounted(async () => {
    await fetchMyTickets();
});

const capitalize = (str) => {
    if (!str) return '';
    return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
};
</script>

<template>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row items-center justify-between mb-8 gap-y-4 sm:gap-y-0">
        <div class="text-center sm:text-start">
            <h1 class="text-2xl font-bold text-gray-800">My Tickets</h1>
            <p class="text-sm text-gray-500 sm:mt-1">View and manage all of your submitted tickets efficiently.</p>
        </div>

        <RouterLink :to="{ name: 'app.tickets.create' }"
            class="w-full sm:w-auto inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700">
            <PlusIcon class="w-4 h-4 mr-2" />
            Create New Ticket
        </RouterLink>
    </div>

    <!-- Search & Filters -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 mb-6 p-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="relative">
                <input v-model="filters.search" type="text" placeholder="Search..."
                    class="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" />
                <MagnifyingGlassIcon class="w-4 h-4 text-gray-400 absolute left-3 top-2.5" />
            </div>

            <select v-model="filters.status"
                class="border border-gray-200 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
                <option value="">All Status</option>
                <option value="OPEN">Open</option>
                <option value="ONPROGRESS">On Progress</option>
                <option value="RESOLVED">Resolved</option>
                <option value="REJECTED">Rejected</option>
            </select>

            <select v-model="filters.priority"
                class="border border-gray-200 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
                <option value="">All Priority</option>
                <option value="LOW">Low</option>
                <option value="MEDIUM">Medium</option>
                <option value="HIGH">High</option>
                <option value="CRITICAL">Critical</option>
            </select>

            <select v-model="filters.date"
                class="border border-gray-200 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500">
                <option value="">All Dates</option>
                <option value="TODAY">Today</option>
                <option value="WEEK">This Week</option>
                <option value="MONTH">This Month</option>
            </select>
        </div>
    </div>

    <!-- Tickets List -->
    <div v-if="tickets.length" class="space-y-4">
        <div v-for="ticket in tickets" :key="ticket.code"
            class="bg-white rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
            <RouterLink :to="{ name: 'app.tickets.create', params: { code: ticket.code } }" class="block p-6">
                <!-- Header Ticket -->
                <div
                    class="flex flex-col sm:flex-row items-start sm:items-center sm:justify-between gap-y-1 sm:gap-y-0 space-x-3">
                    <h3 class="text-lg font-semibold text-gray-800">{{ ticket.title }}</h3>

                    <div class="flex items-center space-x-2">
                        <!-- Status -->
                        <span class="px-3 py-1 text-xs font-medium rounded-lg" :class="{
                            'text-blue-700 bg-blue-100': ticket.status.toLowerCase() === 'open',
                            'text-yellow-700 bg-yellow-100': ticket.status.toLowerCase() === 'onprogress',
                            'text-green-700 bg-green-100': ticket.status.toLowerCase() === 'resolved',
                            'text-red-700 bg-red-100': ticket.status.toLowerCase() === 'rejected',
                        }">
                            {{ capitalize(ticket.status) }}
                        </span>

                        <!-- Priority -->
                        <span class="px-3 py-1 text-xs font-medium rounded-lg" :class="{
                            'text-red-700 bg-red-100': ticket.priority.toLowerCase() === 'high',
                            'text-yellow-700 bg-yellow-100': ticket.priority.toLowerCase() === 'medium',
                            'text-green-700 bg-green-100': ticket.priority.toLowerCase() === 'low',
                        }">
                            {{ capitalize(ticket.priority) }}
                        </span>
                    </div>
                </div>

                <!-- Info Tanggal -->
                <p class="text-sm text-gray-500 mt-2 sm:mt-1">
                    #{{ ticket.code }} | Created at
                    {{ DateTime.fromISO(ticket.createdAt).toFormat('HH:mm, dd MMMM yyyy') }}
                </p>

                <!-- Deskripsi -->
                <p class="text-base text-gray-700 mt-2 sm:mt-1">{{ ticket.description }}</p>

                <!-- Footer -->
                <div class="mt-4 flex flex-col sm:flex-row items-start sm:items-center space-x-4 text-gray-500 text-sm">
                    <div class="shrink-0 flex items-center">
                        <ChatBubbleLeftIcon class="w-4 h-4 mr-1" />
                        <span>{{ ticket.totalReplies }} replies</span>
                    </div>
                    <div class="shrink-0 flex items-center">
                        <ClockIcon class="w-4 h-4 mr-1" />
                        <span>
                            Last updated
                            {{ DateTime.fromISO(ticket.updatedAt).toFormat('dd MMMM yyyy, HH:mm') }}
                        </span>
                    </div>
                </div>
            </RouterLink>
        </div>
    </div>

    <div v-else class="text-center text-gray-400 py-8">
        No data found
    </div>
</template>
