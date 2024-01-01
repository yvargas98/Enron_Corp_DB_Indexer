<script>
export default {
    data() {
        return {
            value: ""
        }
    },
    methods: {
        handleEnter() {
            this.$emit("enter", this.value);
        }
    }
}
</script>

<template>
    <input v-model="value" @keyup.enter="handleEnter" type="text" placeholder="Search"
        class="block w-full h-9 py-1.5 pr-5 text-gray-700 bg-white border border-gray-200 rounded-lg md:w-80 placeholder-gray-400/70 pl-11 rtl:pr-11 rtl:pl-5 dark:bg-gray-900 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 dark:focus:border-blue-300 focus:ring-blue-300 focus:outline-none focus:ring focus:ring-opacity-40">
</template>

<style></style>