<script setup>
import { ref, onMounted, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { debounce } from 'lodash';
import { DateTime } from 'luxon';
import { MagnifyingGlassIcon, ChatBubbleLeftRightIcon } from '@heroicons/vue/24/solid';
import { useTicketStore } from '@/stores/ticket';

const ticketStore = useTicketStore();
const { fetchTickets } = ticketStore;
const { tickets } = storeToRefs(ticketStore);

const filters = ref({
    search: '',
    status: '',
    priority: '',
    date: '',
});

watch(
    filters,
    debounce(async () => {
        await fetchTickets(filters.value);
    }, 300),
    { deep: true },
);

onMounted(async () => {
    await fetchTickets(filters.value);
});

const capitalize = (str) => {
    if (!str) return '';
    return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
};
</script>

<template>
    <div class="p-6">
        <!-- Search & Filters -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 mb-6">
            <div class="p-6">
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
            </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm border border-gray-100">
            <div class="overflow-x-auto w-full">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Ticket ID
                            </th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Title
                            </th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Reporter
                            </th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Status
                            </th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Priority
                            </th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Date
                            </th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Action
                            </th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-100">
                        <tr v-if="!tickets || tickets.length === 0">
                            <td colspan="7" class="px-6 py-10 text-center text-sm text-gray-500">
                                No data found
                            </td>
                        </tr>

                        <tr v-else v-for="ticket in tickets" :key="ticket.code" class="hover:bg-gray-50">
                            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                                #{{ ticket.code }}
                            </td>
                            <td class="px-6 py-4">
                                <div class="text-sm text-gray-800">
                                    {{ ticket.title }}
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <img :src="`https://ui-avatars.com/api/?name=${ticket.user.name}&background=0D8ABC&color=fff`"
                                        :alt="ticket.user.name" class="w-6 h-6 rounded-full" />
                                    <span class="ml-2 text-sm text-gray-800">{{
                                        ticket.user.fullName
                                    }}</span>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span class="px-3 py-1 text-xs font-medium rounded-full" :class="{
                                    'text-blue-700 bg-blue-100': ticket.status.toLowerCase() === 'open',
                                    'text-yellow-700 bg-yellow-100':
                                        ticket.status.toLowerCase() === 'onprogress',
                                    'text-green-700 bg-green-100': ticket.status.toLowerCase() === 'resolved',
                                    'text-red-700 bg-red-100': ticket.status.toLowerCase() === 'rejected',
                                }">
                                    {{ capitalize(ticket.status) }}
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span class="px-3 py-1 text-xs font-medium rounded-full" :class="{
                                    'text-red-700 bg-red-100': ticket.priority.toLowerCase() === 'high',
                                    'text-yellow-700 bg-yellow-100':
                                        ticket.priority.toLowerCase() === 'medium',
                                    'text-green-700 bg-green-100': ticket.priority.toLowerCase() === 'low',
                                }">
                                    {{ capitalize(ticket.priority) }}
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                {{ DateTime.fromISO(ticket.createdAt).toFormat('HH:mm, dd MMMM yyyy') }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm">
                                <RouterLink :to="{
                                    name: 'admin.tickets.detail',
                                    params: { code: ticket.code },
                                }"
                                    class="flex items-center px-3 py-2 bg-blue-500 text-white text-sm leading-4 font-medium rounded-lg">
                                    <ChatBubbleLeftRightIcon class="w-4 h-4 mr-2" />
                                    Reply
                                </RouterLink>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>