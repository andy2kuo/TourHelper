import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // State
  const userId = ref(null)
  const userPreferences = ref({
    maxDistance: 50, // 最大距離 (公里)
    preferredTags: [], // 偏好的標籤
    budget: 'medium' // 預算範圍: low, medium, high
  })
  const currentLocation = ref({
    latitude: null,
    longitude: null
  })

  // Actions
  function setUserId(id) {
    userId.value = id
  }

  function setLocation(lat, lon) {
    currentLocation.value = {
      latitude: lat,
      longitude: lon
    }
  }

  function updatePreferences(preferences) {
    userPreferences.value = { ...userPreferences.value, ...preferences }
  }

  function resetPreferences() {
    userPreferences.value = {
      maxDistance: 50,
      preferredTags: [],
      budget: 'medium'
    }
  }

  return {
    userId,
    userPreferences,
    currentLocation,
    setUserId,
    setLocation,
    updatePreferences,
    resetPreferences
  }
})
