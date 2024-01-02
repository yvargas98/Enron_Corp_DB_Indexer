<template>
	<div class="flex items-center justify-center mt-6">
		<button @click="goToPage(offset - 1)" :disabled="offset === 1"
			class="px-3 py-1 mr-2 text-sm font-medium text-gray-500 bg-gray-200 rounded-md cursor-pointer">Prev</button>
		<span class="text-sm font-medium text-gray-700">{{ offset }}</span>
		<button @click="goToPage(offset + 1)" :disabled="totalResults < limit"
			class="px-3 py-1 ml-2 text-sm font-medium text-gray-500 bg-gray-200 rounded-md cursor-pointer">Next</button>
	</div>
</template>
  
<script>
export default {
	props: {
		offset: {
			type: Number,
			required: true,
		},
		limit: {
			type: Number,
			required: true,
		},
		totalResults: {
			type: Number,
			required: true,
		},
	},
	methods: {
		goToPage(page) {
			if (page > 0 && this.totalResults == this.limit) {
				this.$emit('update:offset', page);
			}
		},
	},
};
</script>
  
<style></style>
