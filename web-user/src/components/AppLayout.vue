<template>
  <div class="app-layout">
    <router-view />
    <BottomTabBar v-if="showTabBar" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import BottomTabBar from './BottomTabBar.vue'

const route = useRoute()
const showTabBar = computed(() => route.meta.showTabBar === true)
</script>
