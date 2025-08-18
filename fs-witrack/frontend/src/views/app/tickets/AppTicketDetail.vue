<script setup>
import { ref, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useRoute } from 'vue-router';
import { DateTime } from 'luxon';
import { ArrowLeftIcon, ClockIcon, ArrowDownTrayIcon } from '@heroicons/vue/24/outline';

import { useTicketStore } from '@/stores/ticket';

const ticketStore = useTicketStore();
const { fetchTicketByCode } = ticketStore;
const { success, error, loading } = storeToRefs(ticketStore);

const route = useRoute();

const ticket = ref({});

const form = ref({
    content: '',
});

const fetchTicketDetail = async () => {
    const response = await fetchTicketByCode(route.params.code);
    console.log(response);
    ticket.value = response;
    form.value.status = response.status;
}

onMounted(async () => {
    await fetchTicketDetail();
});

const capitalize = (str) => {
    if (!str) return '';
    return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
};
</script>

<template>
    <div class="mb-6">
        <RouterLink :to="{ name: 'app.dashboard' }"
            class="inline-flex items-center text-sm text-gray-600 hover:text-gray-800">
            <ArrowLeftIcon class="w-4 h-4 mr-2" />
            Back to Ticket List
        </RouterLink>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-100 mb-6">
        <div class="p-6">
            <div class="flex items-start justify-between">

                <div>
                    <h1 class="text-2xl font-bold text-gray-800">
                        {{ ticket.title }}
                    </h1>
                    <div class="mt-2 flex items-center space-x-4">
                        <span class="px-3 py-1 text-xs font-medium rounded-lg" :class="{
                            'text-blue-700 bg-blue-100': ticket.status?.toLowerCase() === 'open',
                            'text-yellow-700 bg-yellow-100': ticket.status?.toLowerCase() === 'onprogress',
                            'text-green-700 bg-green-100': ticket.status?.toLowerCase() === 'resolved',
                            'text-red-700 bg-red-100': ticket.status?.toLowerCase() === 'rejected',
                        }">
                            {{ capitalize(ticket.status) }}
                        </span>
                        <span class="px-3 py-1 text-xs font-medium rounded-lg" :class="{
                            'text-red-700 bg-red-100': ticket.priority?.toLowerCase() === 'high',
                            'text-yellow-700 bg-yellow-100': ticket.priority?.toLowerCase() === 'medium',
                            'text-green-700 bg-green-100': ticket.priority?.toLowerCase() === 'low',
                        }">
                            {{ capitalize(ticket.priority) }}
                        </span>
                        <span class="text-sm text-gray-500"> #{{ ticket.code }} </span>
                        <div class="text-sm text-gray-500 flex items-center space-x-1">
                            <ClockIcon class="w-4 h-4 inline-block" />
                            <span>Created At {{ DateTime.fromISO(ticket.createdAt).toFormat('HH:mm, dd MMMM yyyy')
                            }}</span>
                        </div>
                    </div>
                    <p class="mt-4 text-gray-700 text-sm leading-relaxed">
                        {{ ticket.description }}
                    </p>
                </div>

                <button type="button"
                    class="px-4 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50 flex items-center">
                    <ArrowDownTrayIcon class="w-4 h-4 inline-block mr-2" />
                    Attachment
                </button>
            </div>
        </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-100">
        <!-- Additional content here -->
    </div>
</template>
