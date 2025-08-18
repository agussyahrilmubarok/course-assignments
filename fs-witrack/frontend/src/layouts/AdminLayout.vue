<script setup>
import Navbar from '@/components/app/Navbar.vue'
import Sidebar from '@/components/admin/Sidebar.vue';
</script>

<template>
    <div class="flex flex-col h-screen">
        <Navbar />

        <main class="flex h-screen">
            <Sidebar />

            <div class="flex-1 overflow-x-hidden overflow-y-auto">
                <div class="p-6 space-y-6">
                    <router-view></router-view>
                </div>
            </div>

        </main>
    </div>
</template>