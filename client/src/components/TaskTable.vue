<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-base font-semibold text-gray-900">Tasks</h1>
      </div>
      <div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
        <button type="button" @click="fetchTableData"
          class="block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Refresh</button>
      </div>
    </div>
    <div class="mt-8 flow-root">
      <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
          <div v-if="loading">Loading...</div>
          <table v-else class="min-w-full divide-y divide-gray-300">
            <thead>
              <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-0">ID</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Command</th>
                <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Arguments</th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                  <span class="sr-only">Logs</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="job in jobs" :key="job.email">
                <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">{{ job.id }}</td>
                <!-- <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500 uppercase font-bold"
                  :class="getStatClass(job.status)">{{ job.status }}</td> -->
                  <td>
                    <span :class="getStatusBadge(job.status)">{{ job.status }}</span>
                  </td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ job.command }}</td>
                <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ job.args }}</td>
                <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                  <span v-if="job.status=='running'">
                    <span @click="handleKill(job.id)" class="text-indigo-600 hover:text-indigo-900"> Kill</span> |
                  </span>
                  <RouterLink :to="`/logs/${job.id}`" class="text-indigo-600 hover:text-indigo-900">Logs<span
                      class="sr-only">, {{ job.name }}</span></RouterLink>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useServerStore } from '@/stores/server';
import { ref, onMounted } from 'vue';
import { useToast } from 'vue-toast-notification';
const $toast = useToast();

const jobs = ref([])
const loading = ref(false)
onMounted(() => {
  fetchTableData()
})

function handleKill(id) {
  const serverStore = useServerStore()
  const url = serverStore.url + (serverStore.port ? (":" + serverStore.port) : "") + "/jobs/"+id+"/kill"
  console.log("Killing job with id: " + id)
  fetch(url, {
    method: "GET",
    headers: {
      "Authorization": serverStore.token
    }
  }).then(response => {
    if (response.ok) {
      $toast.open({
          message: "Task Killed",
          type: "success",
          duration: 3000
        })
        fetchTableData()
    } else {
      throw new Error(response.status)
    }
  })
    .then(data => {
      loading.value = false
      jobs.value = data.jobs
    })
    .catch(e => {
      console.log(e)
      if (e.message === "401") {
        $toast.open({
          message: "Unauthorized",
          type: "error",
          duration: 3000
        })
      } else if (e.message === "400") {
        $toast.open({
          message: "Bad Request",
          type: "error",
          duration: 3000
        })
      } else if (e.message === "500") {
        $toast.open({
          message: "Internal Server Error",
          type: "error",
          duration: 3000
        })
      } else if (e.message === "Failed to fetch") {
        $toast.open({
          message: "Cannot Reach Server",
          type: "error",
          duration: 3000
        })
      }
      jobs.value = []
      loading.value = false
    })


}

function getStatusBadge(status) {
  if (status == "running") {
    return "inline-flex items-center rounded-md bg-green-50 px-2 py-1 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20"
  } else if (status == "error") {
    return "inline-flex items-center rounded-md bg-red-50 px-2 py-1 text-xs font-medium text-red-700 ring-1 ring-inset ring-red-600/10"
  }else if (status == "killed") {
    return "inline-flex items-center rounded-md bg-yellow-50 px-2 py-1 text-xs font-medium text-yellow-800 ring-1 ring-inset ring-yellow-600/20"
  } else if (status == "pending") {
    return "inline-flex items-center rounded-md bg-gray-50 px-2 py-1 text-xs font-medium text-gray-600 ring-1 ring-inset ring-gray-500/10"
  } else if (status == "finished") {
    return "inline-flex items-center rounded-md bg-indigo-50 px-2 py-1 text-xs font-medium text-indigo-700 ring-1 ring-inset ring-indigo-700/10"
  } else {
    return "inline-flex items-center rounded-md bg-gray-50 px-2 py-1 text-xs font-medium text-gray-600 ring-1 ring-inset ring-gray-500/10"
  }
}

function getStatClass(stat) {
  if (stat == "running" || stat == "new") {
    return "b text-green-800"
  } else if (stat == "error" || stat == "killed") {
    return "text-red-800"
  } else if (stat == "pending") {
    return "text-yellow-800"
  } else if (stat == "finished") {
    return "text-blue-800"
  } else {
    return "text-gray-800"
  }
}

function fetchTableData() {
  const serverStore = useServerStore()
  const url = serverStore.url + (serverStore.port ? (":" + serverStore.port) : "") + "/jobs"
  loading.value = true

  fetch(url, {
    method: "GET",
    headers: {
      "Authorization": serverStore.token
    }
  }).then(response => {
    if (response.ok) {
      return response.json()
    } else {
      throw new Error(response.status)
    }
  })
    .then(data => {
      loading.value = false
      jobs.value = data.jobs
    })
    .catch(e => {
      console.log(e)
      if (e.message === "401") {
        $toast.open({
          message: "Unauthorized",
          type: "error",
          duration: 3000
        })
      } else if (e.message === "400") {
        $toast.open({
          message: "Bad Request",
          type: "error",
          duration: 3000
        })
      } else if (e.message === "500") {
        $toast.open({
          message: "Internal Server Error",
          type: "error",
          duration: 3000
        })
      } else if (e.message === "Failed to fetch") {
        $toast.open({
          message: "Cannot Reach Server",
          type: "error",
          duration: 3000
        })
      }
      jobs.value = []
      loading.value = false
    })
}
</script>