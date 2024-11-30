<script setup>
import { useServerStore } from '@/stores/server';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toast-notification';
const $toast = useToast();
const router = useRouter();

const command = ref("python3");
const commandError = ref("" | null);
const args = ref("long_task.py");
const argsError = ref("" | null);

function handleSubmit() {
    const serverStore = useServerStore()

    const loading = ref(false)

    args.value = args.value.trim()
    let parsedArgs = args.value.split(",")
    if (parsedArgs.length == 1 && parsedArgs[0] == "") {
        parsedArgs = []
    }

    const url = serverStore.url + (serverStore.port ? (":" + serverStore.port) : "") + "/jobs"
    if (command.value == "") {
        commandError.value = "Command is required"
        return
    } else {
        commandError.value = null
    }
    $toast.open({
        message: "Creating Command",
        type: "info",
        duration: 2000
    })
    loading.value = true
    fetch(url, {
        method: "POST",
        headers: {
            "Authorization": serverStore.token
        },
        body: JSON.stringify({
            command: command.value,
            args: parsedArgs
        })
    }).then(response => {
        if (response.ok) {
     
            router.push({ name: "home" })
            return
        } else {
            throw new Error(response.status)
        }
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
            loading.value = false
        })



}
</script>
<template>
    <div>
        <div>
            <label for="command" class="block text-sm/6 font-medium text-gray-900">Command</label>
            <div class="mt-2">
                <input type="text" name="command" id="command" v-model="command"
                    class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                    placeholder="you@example.com" />
            </div>
            <p v-if="commandError" class="mt-2 text-sm text-red-600" id="email-error">{{ commandError }}</p>
        </div>

        <div>
            <label for="args" class="block text-sm/6 font-medium text-gray-900 mt-3">Arguments</label>
            <div class="mt-2">
                <input type="text" name="args" id="args" v-model="args"
                    class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6"
                    placeholder="you@example.com" aria-describedby="email-description" />
            </div>
            <p class="mt-2 text-sm text-gray-500" id="email-description">Comma separated</p>
            <p v-if="argsError" class="mt-2 text-sm text-red-600" id="email-error">{{ argsError }}</p>
        </div>
        <button type="button" @click="handleSubmit"
            class="block mt-3 rounded-md bg-indigo-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Create</button>

    </div>
</template>