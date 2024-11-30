<script setup>
import { useServerStore } from '@/stores/server';
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
const router = useRoute();
const { id } = router.params;

const allLogs = ref(true);
const lastLines = ref(10);
const loading = ref(false);
const logs = ref("" | null);

const intervalInMS = ref(1500);
const toggleInterval = ref(false);
var intervalRef = null

onMounted(() => {
    fetchLogs();
});

onUnmounted(() => {
    if (intervalRef) clearInterval(intervalRef);
});

function handleInterval() {
    if (toggleInterval.value) {
        intervalRef = setInterval(() => {
            fetchLogs();
        }, intervalInMS.value);
    } else {
        if (intervalRef) clearInterval(intervalRef);
        else console.log("Interval not set")
    }
}


function fetchLogs() {
    const url = getDownloadUrl();
    loading.value = true

    fetch(url, {
        method: "GET",
        headers: {
            "Authorization": "123"
        }
    }).then((response) => response.text())
        .then(data => {
            loading.value = false
            logs.value = data
        })
        .catch(e => {
            console.log(e)
            loading.value = false
            logs.value = null
        })
}
function getDownloadUrl() {
    const serverStore = useServerStore();
    return serverStore.url + (serverStore.port ? (":" + serverStore.port) : "") + "/jobs/" + id + "/logs" + (allLogs.value ? "" : "?lines=" + lastLines.value);
}

</script>
<template>
    <div class="flex justify-between items-center">
        <div class="flex space-x-3 items-center h-4">
            <div class="flex gap-3">
                <div class="flex h-6 shrink-0 items-center">
                    <div class="group grid size-4 grid-cols-1">
                        <input v-model="allLogs" id="AllLogs" aria-describedby="AllLogs-description" name="AllLogs"
                            type="checkbox"
                            class="col-start-1 row-start-1 appearance-none rounded border border-gray-300 bg-white checked:border-indigo-600 checked:bg-indigo-600 indeterminate:border-indigo-600 indeterminate:bg-indigo-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto" />
                        <svg class="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-[:disabled]:stroke-gray-950/25"
                            viewBox="0 0 14 14" fill="none">
                            <path class="opacity-0 group-has-[:checked]:opacity-100" d="M3 8L6 11L11 3.5"
                                stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                            <path class="opacity-0 group-has-[:indeterminate]:opacity-100" d="M3 7H11" stroke-width="2"
                                stroke-linecap="round" stroke-linejoin="round" />
                        </svg>
                    </div>
                </div>
                <div class="text-sm/6">
                    <label for="AllLogs" class="font-medium text-gray-900">All Logs</label>
                </div>
            </div>
            <div v-if="!allLogs">
                <div>
                    <input type="number" name="lastLines" id="lastLines" aria-label="lastLines" v-model="lastLines"
                        class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                        placeholder="# of most recent lines" />
                </div>
            </div>
        </div>

        <div class="flex space-x-3 items-center h-4">
            <div class="flex gap-3">
                <div class="flex h-6 shrink-0 items-center">
                    <div class="group grid size-4 grid-cols-1">
                        <input v-model="toggleInterval" id="toggleInterval"
                            aria-describedby="toggleInterval-description" name="toggleInterval" @change="handleInterval"
                            type="checkbox"
                            class="col-start-1 row-start-1 appearance-none rounded border border-gray-300 bg-white checked:border-indigo-600 checked:bg-indigo-600 indeterminate:border-indigo-600 indeterminate:bg-indigo-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto" />
                        <svg class="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-[:disabled]:stroke-gray-950/25"
                            viewBox="0 0 14 14" fill="none">
                            <path class="opacity-0 group-has-[:checked]:opacity-100" d="M3 8L6 11L11 3.5"
                                stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                            <path class="opacity-0 group-has-[:indeterminate]:opacity-100" d="M3 7H11" stroke-width="2"
                                stroke-linecap="round" stroke-linejoin="round" />
                        </svg>
                    </div>
                </div>
                <div class="text-sm/6">
                    <label for="toggleInterval" class="font-medium text-gray-900">Toggle Interval Fetching</label>
                </div>
            </div>
        </div>


        <div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none">
            <button type="button" @click="fetchLogs"
                class="block rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Refresh</button>
        </div>
    </div>
    <div>
        <label for="comment" class="block text-sm/6 font-medium text-gray-900">Logs</label>
        <div class="mt-2">
            <textarea rows="15" name="comment" id="comment" 
                class="block w-full resize-y rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6">{{ loading ? "Loading..." : logs }}</textarea>
        </div>
    </div>
</template>