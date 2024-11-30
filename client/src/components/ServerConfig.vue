<script setup>
import { useServerStore } from '@/stores/server';

const serverStore = useServerStore();
</script>
<template>
    <div class="divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow">
        <div class="px-4 py-5 sm:px-6">
            <h3 class="text-base font-semibold text-gray-900">Server Info</h3>
        </div>
        <div class="px-4 py-5 sm:p-6">
            <div class="grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
                <div class="sm:col-span-2 sm:col-start-1">
                    <label for="url" class="block text-sm/6 font-medium text-gray-900">URL / IP</label>
                    <div class="mt-2">
                        <input type="text" name="url" id="url" v-model="serverStore.url"
                            class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm/6" />
                    </div>
                </div>
                <div class="sm:col-span-2">
                    <label for="port" class="block text-sm/6 font-medium text-gray-900">Port</label>
                    <div class="mt-2">
                        <input type="text" name="port" id="port" v-model="serverStore.port"
                            class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm/6" />
                    </div>
                </div>

                <div class="sm:col-span-2">
                    <label for="token" class="block text-sm/6 font-medium text-gray-900">Token</label>
                    <div class="mt-2">
                        <input type="text" name="token" id="token" v-model="serverStore.token"
                            class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm/6" />
                    </div>
                </div>
            </div>

        </div>
    </div>

</template>