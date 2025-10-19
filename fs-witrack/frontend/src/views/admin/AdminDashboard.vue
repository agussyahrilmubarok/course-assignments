<script setup>
import { onMounted, watch } from "vue";
import { storeToRefs } from "pinia";
import { Chart } from "chart.js/auto";
import {
  TagIcon,
  ClockIcon,
  CheckCircleIcon,
  ChartBarIcon,
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon,
} from "@heroicons/vue/24/outline";
import { useDashboardStore } from "@/stores/dashboard";
import { useTicketStore } from "@/stores/ticket";

const dashboardStore = useDashboardStore();
const { statistic } = storeToRefs(dashboardStore);
const { fetchStatistics } = dashboardStore;

const ticketStore = useTicketStore();
const { tickets } = storeToRefs(ticketStore);
const { fetchTickets } = ticketStore;

let chart = null;

// update chart when statistic changes
watch(
  statistic,
  () => {
    if (statistic.value && chart) {
      chart.data.datasets[0].data = [
        statistic.value.statusDistribution?.open ?? 0,
        statistic.value.statusDistribution?.onprogress ?? 0,
        statistic.value.statusDistribution?.resolved ?? 0,
        statistic.value.statusDistribution?.rejected ?? 0,
      ];
      chart.update();
    }
  },
  { deep: true }
);

onMounted(async () => {
  await fetchTickets();
  await fetchStatistics();

  const statusCtx = document.getElementById("statusChart")?.getContext("2d");

  if (statusCtx && statistic.value) {
    chart = new Chart(statusCtx, {
      type: "doughnut",
      data: {
        labels: ["Open", "In Progress", "Resolved", "Rejected"],
        datasets: [
          {
            data: [
              statistic.value.statusDistribution?.open ?? 0,
              statistic.value.statusDistribution?.onprogress ?? 0,
              statistic.value.statusDistribution?.resolved ?? 0,
              statistic.value.statusDistribution?.rejected ?? 0,
            ],
            backgroundColor: ["#3B82F6", "#F59E0B", "#10B981", "#EF4444"],
          },
        ],
      },
      options: {
        responsive: true,
        plugins: {
          legend: {
            position: "bottom",
          },
        },
        cutout: "70%",
      },
    });
  }
});
</script>

<template>
  <div>
    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Total Tickets -->
      <div class="stat-card bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Tickets</p>
            <h3 class="text-2xl font-bold text-gray-800 mt-1">
              {{ statistic?.totalTickets }}
            </h3>
          </div>
          <div class="p-3 bg-blue-50 rounded-lg">
            <TagIcon class="w-6 h-6 text-blue-600" />
          </div>
        </div>
        <div class="mt-4 flex items-center text-sm">
          <span class="text-green-500 flex items-center">
            <ArrowTrendingUpIcon class="w-4 h-4 mr-1" /> 12%
          </span>
          <span class="text-gray-500 ml-2">vs last month</span>
        </div>
      </div>

      <!-- Active Tickets -->
      <div class="stat-card bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Active Tickets</p>
            <h3 class="text-2xl font-bold text-gray-800 mt-1">
              {{ statistic?.activeTickets }}
            </h3>
          </div>
          <div class="p-3 bg-yellow-50 rounded-lg">
            <ClockIcon class="w-6 h-6 text-yellow-600" />
          </div>
        </div>
        <div class="mt-4 flex items-center text-sm">
          <span class="text-red-500 flex items-center">
            <ArrowTrendingDownIcon class="w-4 h-4 mr-1" /> 12%
          </span>
          <span class="text-gray-500 ml-2">vs last month</span>
        </div>
      </div>

      <!-- Resolved -->
      <div class="stat-card bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Resolved</p>
            <h3 class="text-2xl font-bold text-gray-800 mt-1">
              {{ statistic?.resolvedTickets }}
            </h3>
          </div>
          <div class="p-3 bg-green-50 rounded-lg">
            <CheckCircleIcon class="w-6 h-6 text-green-600" />
          </div>
        </div>
        <div class="mt-4 flex items-center text-sm">
          <span class="text-green-500 flex items-center">
            <ArrowTrendingUpIcon class="w-4 h-4 mr-1" /> 12%
          </span>
          <span class="text-gray-500 ml-2">vs last month</span>
        </div>
      </div>

      <!-- Avg Resolution Time -->
      <div class="stat-card bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Average Resolution Time</p>
            <h3 class="text-2xl font-bold text-gray-800 mt-1">
              {{ statistic?.avgResolutionTime }}
            </h3>
          </div>
          <div class="p-3 bg-purple-50 rounded-lg">
            <ChartBarIcon class="w-6 h-6 text-purple-600" />
          </div>
        </div>
        <div class="mt-4 flex items-center text-sm">
          <span class="text-green-500 flex items-center">
            <ArrowTrendingUpIcon class="w-4 h-4 mr-1" /> 12%
          </span>
          <span class="text-gray-500 ml-2">vs last month</span>
        </div>
      </div>
    </div>

    <!-- Chart Section -->
    <div class="mt-8 bg-white rounded-xl shadow-sm p-6 border border-gray-100 flex justify-center">
    <div class="w-full max-w-sm">
        <h3 class="text-lg font-bold text-gray-800 mb-4 text-center">
        Ticket Status Distribution
        </h3>
        <div class="flex justify-center">
        <canvas id="statusChart" class="w-64 h-64"></canvas>
        </div>
    </div>
    </div>
  </div>
</template>
