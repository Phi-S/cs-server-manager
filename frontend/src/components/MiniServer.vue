<script setup lang="ts">
import {onMounted, ref} from "vue";

let displayMapAndPlayers = ref(false)
let started = ref(false)

let posts = ref<Post[]>()
type Post = {
    userId: number,
    id: number,
    title: string,
    body: string
}

function onClick() {
    displayMapAndPlayers.value = displayMapAndPlayers.value!
}

onMounted(() => {
    fetch("https://jsonplaceholder.typicode.com/posts", {
        method: "get",
        headers: {
            'content-type': 'application/json'
        }
    }).then(res => {
        if (!res.ok) {
            throw new Error(res.statusText + " with json " + res.json());
        }

        return res.json()
    }).then(json => {
        const temp = json as Post[]
        console.log(temp)
    })
})

</script>

<template>
    <li v-if="posts !== undefined" v-for="post of posts">
        <p><strong>{{ post }}</strong></p>
        <p></p>
    </li>


    <div class="container bg-success rounded-2" style="height: 2.5rem; width: 30rem">
        <div class="row align-items-center text-nowrap flex-nowrap">
            <div class="col-10 btn text-start">
                asdf

                <div>

                </div>
            </div>
            <div class="col-2 text-end">
                <button v-if="started" class="btn bi-stop"></button>
                <button v-else-if="true" class="spinner-grow">
                </button>
                <button v-else class="btn bi-play"></button>
            </div>
        </div>
    </div>

    <button @click="onClick">display toggle</button>
</template>