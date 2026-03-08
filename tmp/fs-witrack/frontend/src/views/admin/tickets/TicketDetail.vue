<script setup>
import { ref, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useRoute } from 'vue-router';
import { DateTime } from 'luxon';
import { ArrowLeftIcon, ClockIcon, ArrowDownTrayIcon, PaperAirplaneIcon } from '@heroicons/vue/24/outline';
import { useAuthStore } from '@/stores/auth';
import { useTicketStore } from '@/stores/ticket';

const ticketStore = useTicketStore();
const { fetchTicketByCode, updateTicketByCode, createTicketReply } = ticketStore;
const { success, error, loading } = storeToRefs(ticketStore);
const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

const route = useRoute();

const ticket = ref({});

const form = ref({
    content: '',
});

const fetchTicketDetail = async () => {
    const response = await fetchTicketByCode(route.params.code);
    ticket.value = response;
    form.value.status = response.status;
}

const handleUpdateStatus = async () => {
    if (!ticket.value.code) return;
    await updateTicketByCode(ticket.value.code, {
        title: ticket.value.title,
        description: ticket.value.description,
        status: ticket.value.status,
        priority: ticket.value.priority,
    });
    await fetchTicketDetail();
};

const handleUpdatePriority = async () => {
    if (!ticket.value.code) return;
    await updateTicketByCode(ticket.value.code, {
        title: ticket.value.title,
        description: ticket.value.description,
        status: ticket.value.status,
        priority: ticket.value.priority,
    });
    await fetchTicketDetail();
};

const handleSubmit = async () => {
    error.value = null
    await createTicketReply(route.params.code, form.value)

    await fetchTicketDetail()
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
    <div class="min-h-screen bg-gray-50 flex flex-col">
        <!-- Back link -->
        <div class="mb-6">
            <RouterLink :to="{ name: 'admin.tickets' }"
                class="inline-flex items-center text-sm text-gray-600 hover:text-gray-800">
                <ArrowLeftIcon class="w-4 h-4 mr-2" />
                Back to Ticket List
            </RouterLink>
        </div>

        <!-- Ticket Overview -->
        <div class="bg-white mb-2 rounded-xl shadow-sm border border-gray-100">
            <div class="p-6">
                <div class="flex items-start justify-between">
                    <div>
                        <h1 class="text-2xl font-bold text-gray-800">
                            {{ ticket.title }}
                        </h1>
                        <div class="mt-2 flex items-center space-x-4">
                            <!-- Status -->
                            <select v-model="ticket.status" @change="handleUpdateStatus"
                                class="px-3 py-1 text-xs font-medium rounded-lg border border-gray-200 focus:outline-none focus:ring-1 focus:ring-blue-500"
                                :class="{
                                    'text-blue-700 bg-blue-100': ticket.status?.toLowerCase() === 'open',
                                    'text-yellow-700 bg-yellow-100': ticket.status?.toLowerCase() === 'onprogress',
                                    'text-green-700 bg-green-100': ticket.status?.toLowerCase() === 'resolved',
                                    'text-red-700 bg-red-100': ticket.status?.toLowerCase() === 'rejected',
                                }">
                                <option value="OPEN">Open</option>
                                <option value="ONPROGRESS">On Progress</option>
                                <option value="RESOLVED">Resolved</option>
                                <option value="REJECTED">Rejected</option>
                            </select>
                            <!-- Priority -->
                            <select v-model="ticket.priority" @change="handleUpdatePriority"
                                class="px-3 py-1 text-xs font-medium rounded-lg border border-gray-200 focus:outline-none focus:ring-1 focus:ring-blue-500"
                                :class="{
                                    'text-red-700 bg-red-100': ticket.priority?.toLowerCase() === 'high',
                                    'text-yellow-700 bg-yellow-100': ticket.priority?.toLowerCase() === 'medium',
                                    'text-green-700 bg-green-100': ticket.priority?.toLowerCase() === 'low',
                                }">
                                <option value="LOW">Low</option>
                                <option value="MEDIUM">Medium</option>
                                <option value="HIGH">High</option>
                                <option value="CRITICAL">Critical</option>
                            </select>
                            <span class="text-sm text-gray-500"> #{{ ticket.code }} </span>
                            <div class="text-sm text-gray-500 flex items-center space-x-1">
                                <ClockIcon class="w-4 h-4 inline-block" />
                                <span>
                                    Created At
                                    {{ DateTime.fromISO(ticket.createdAt).toFormat('HH:mm, dd MMMM yyyy') }}
                                </span>
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

        <!-- Main Content -->
        <div class="flex-1 space-y-6">

            <!-- Replies -->
            <div class="bg-white rounded-xl shadow-sm border border-gray-100">
                <div v-if="!ticket.replies || ticket.replies.length === 0"
                    class="p-6 text-center text-sm text-gray-500">
                    No replies yet
                </div>
                <div v-else v-for="reply in ticket.replies" :key="reply.id" :class="[
                    'p-6 border-b border-gray-100',
                    reply.user.id === user?.id ? 'bg-gray-100' : 'bg-white'
                ]">
                    <div class="flex items-start space-x-4">
                        <img :src="`https://ui-avatars.com/api/?name=${reply.user.fullName}&background=0D8ABC&color=fff`"
                            alt="profile" class="w-8 h-8 rounded-full" />
                        <div class="flex-1">
                            <div class="flex items-start justify-between">
                                <div class="text-left">
                                    <p class="text-sm font-medium text-gray-700">
                                        {{ reply.user.fullName }}
                                    </p>
                                    <p class="text-xs text-gray-500">
                                        {{ reply.user.email }}
                                    </p>
                                </div>
                                <p class="text-xs text-gray-500 whitespace-nowrap ml-4">
                                    {{ DateTime.fromISO(reply.createdAt).toFormat('HH:mm, dd MMM yyyy') }}
                                </p>
                            </div>
                        </div>
                    </div>
                    <div class="mt-3 text-sm text-gray-800">
                        <p>{{ reply.content }}</p>
                    </div>

                </div>
            </div>

        </div>

        <!-- Reply Section -->
        <div class="p-6 mt-2 rounded-xl border-t shadow-sm border-gray-100 bg-white">
            <form @submit.prevent="handleSubmit" class="space-y-4">
                <!-- Content -->
                <div class="group">
                    <textarea placeholder="Write your reply here.." rows="4" v-model="form.content"
                        class="w-full px-4 py-3 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                        :class="{ 'border-red-500 ring-red-500': error?.content }"></textarea>
                    <p class="mt-1 text-xs text-red-500" v-if="error?.content">
                        {{ error.content }}
                    </p>
                </div>

                <div class="flex items-center justify-between">
                    <!-- Attachment -->
                    <button type="button"
                        class="flex items-center gap-2 h-9 px-4 border border-gray-200 rounded-lg text-sm text-gray-600 hover:bg-gray-50">
                        <ArrowDownTrayIcon class="w-4 h-4" />
                        <span>Attachment</span>
                    </button>

                    <!-- Send -->
                    <button type="submit" :class="[
                        'flex items-center gap-2 h-9 px-6 text-white rounded-lg text-sm',
                        loading ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'
                    ]">
                        <PaperAirplaneIcon class="w-4 h-4" />
                        <span v-if="loading">Loading...</span>
                        <span v-else>Send</span>
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>