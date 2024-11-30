import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useServerStore = defineStore('server', () => {
  const url = ref("http://localhost")
  const port = ref("8080")
  const token = ref("123")
  return { url, port, token }
})
