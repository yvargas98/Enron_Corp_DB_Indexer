<script>
export default {
    props: ["onClick"]
}
</script>

<script setup>
defineProps({
    onClick: {
        type: Function,
        required: true
    }
})
</script>
<template>
    <button @click="onClick"
        class="h-9 flex items-center rounded-lg border-2 border-gray-400 bg-gray-300 text-gray-700 shadow-lg hover:bg-gray-400 focus:outline-none focus:ring focus:ring-gray-400 ml-2 px-4 py-2">Submit</button>
</template>

<style></style>