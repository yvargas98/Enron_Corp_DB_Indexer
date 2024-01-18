<template>
    <div v-if="selectedEmail" class="mt-20 mb-20 bg-gray-200 shadow-md rounded-lg container px-20 py-20">
        <h2 class="text-3xl font-semibold tracking-wide mb-6">
            {{ selectedEmail.subject }}
        </h2>
        <div class="grid grid-cols-1">
            <div class="mb-6">
                <p class="text-gray-700"><strong>Date:</strong> {{ selectedEmail.date }}</p>
            </div>
            <div class="mb-6">
                <p class="text-gray-700"><strong>From:</strong> {{ selectedEmail.from }}</p>
            </div>
            <div class="mb-6">
                <p class="text-gray-700"><strong>To:</strong> {{ selectedEmail.to }}</p>
            </div>
        </div>
        <p class="mb-6">
            <strong class=" text-gray-700 font-bold">Content:</strong>
            <br />
            {{ selectedEmail.content }}
        </p>
        <button @click="goBack"
            class="flex items-center rounded-lg border-2 border-gray-400 bg-gray-300 text-gray-700 px-6 py-2 rounded-md hover:bg-blue-700 shadow-lg hover:bg-gray-400 focus:outline-none focus:ring focus:ring-gray-400 ml-2 px-4 py-2">
            Back
        </button>
    </div>
</template>
  
<script>
export default {
    props: {
        selectedEmail: {
            type: Object,
            required: true,
        },
    },
    methods: {
        goBack() {
            this.$emit('closeDetail');
        },
    },
};
</script>
  
<style scoped>
/* Añade estilos personalizados aquí si lo deseas */
</style>