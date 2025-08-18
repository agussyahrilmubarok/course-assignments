<script setup>
import { ref } from 'vue';
import { storeToRefs } from 'pinia';
import { ArrowLeftIcon, PaperAirplaneIcon } from '@heroicons/vue/24/solid';

import { useTicketStore } from '@/stores/ticket';

const ticketStore = useTicketStore();
const { createTicket } = ticketStore;
const { loading, error } = storeToRefs(ticketStore);

const form = ref({
    title: '',
    description: '',
    status: 'OPEN',
    priority: '',
})

const handleSubmit = async () => {
    await createTicket(form.value);
}
</script>

<template>
    <div class="min-h-screen bg-gray-50 flex flex-col">

        <div class="mb-6">
            <RouterLink :to="{ name: 'app.dashboard' }"
                class="inline-flex items-center text-sm text-gray-600 hover:text-gray-800">
                <ArrowLeftIcon class="w-4 h-4 mr-2" />
                Back to Ticket List
            </RouterLink>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-gray-100">

            <div class="p-6 border-b border-gray-100">
                <h1 class="text-2xl font-bold text-gray-800">Create New Ticket</h1>
                <p class="text-sm text-gray-500 mt-1">Fill out the form below to submit a new ticket</p>
            </div>

            <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
                <div>
                    <label for="title" class="block text-sm font-medium text-gray-700 mb-2">Title</label>
                    <input v-model="form.title" type="text" id="title"
                        placeholder="e.g., Wi-Fi connectivity issue in office"
                        class="w-full px-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                        :class="{ 'border-red-500 ring-red-500': error?.title }" />
                    <div v-if="error?.title" class="flex items-center mt-2">
                        <p class="text-xs text-red-500">
                            {{ error.title }}
                        </p>
                    </div>
                </div>

                <div>
                    <label for="description" class="block text-sm font-medium text-gray-700 mb-2">Description</label>
                    <textarea v-model="form.description" id="description"
                        placeholder="Provide a detailed description of the issue, including steps to reproduce or affected areas"
                        class="w-full px-4 py-3 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                        :class="{ 'border-red-500 ring-red-500': error?.description }"></textarea>
                    <div v-if="error?.description" class="flex items-center mt-2">
                        <p class="text-xs text-red-500">
                            {{ error.description }}
                        </p>
                    </div>
                </div>

                <div>
                    <label for="priority" class="block text-sm font-medium text-gray-700 mb-2">Priority</label>
                    <select v-model="form.priority" id="priority"
                        class="w-full px-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
                        :class="{ 'border-red-500 ring-red-500': error?.priority }">
                        <option value="" disabled>Select priority</option>
                        <option value="LOW">Low</option>
                        <option value="MEDIUM">Medium</option>
                        <option value="HIGH">High</option>
                        <option value="CRITICAL">Critical</option>
                    </select>
                    <div v-if="error?.priority" class="flex items-center mt-2">
                        <p class="text-xs text-red-500">
                            {{ error.priority }}
                        </p>
                    </div>
                </div>

                <div class="flex items-center justify-end">
                    <button type="submit" :class="[
                        'px-6 py-2 text-white rounded-lg transition flex items-center justify-center gap-2 text-sm font-medium',
                        loading ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'
                    ]">
                        <PaperAirplaneIcon class="w-4 h-4" />
                        <span v-if="loading">Loading...</span>
                        <span v-else>Submit</span>
                    </button>
                </div>
            </form>
        </div>

    </div>
</template>
