<script>
export default {
    props: {
        data: {
            type: Array,
            required: true
        },
        value: {
            type: String,
            required: true
        }
    },
    data() {
        return {
            emails: []
        }
    },
    created() {
        this.emails = this.data
    },
    updated() {
        this.emails = this.data
    },
    methods: {
        highlightTextField(text) {
            if(text.length > 200) {
                const queryPosition = text.toLowerCase().indexOf(this.value.toLowerCase());
                const paragraphs = text.split('\n\n');
                const formattedContent = paragraphs.map(paragraph => `${paragraph}`).join('');
                const shortContent = (queryPosition > 40 ? "..." : "") + formattedContent.slice(queryPosition < 40 ? 0 : queryPosition - 40, queryPosition + 85) + "..."
                const highlightedContent = shortContent.replace(this.value, `<span style="font-weight: bold;" class="bg-blue-200">${this.value}</span>`);
                return highlightedContent
            }
            else {
                const highlightedContent = text.replace(this.value, `<span style="font-weight: bold;" class="bg-blue-200">${this.value}</span>`);
                return highlightedContent
            }
        }
    }
}
</script>

<template>
    <tbody>
        <tr v-if="emails.length === 0">
            <td class="px-4 py-4 text-sm" colspan="4">
                <div class="text-gray-700 dark:text-gray-500 text-center">
                    No result found.
                </div>
            </td>
        </tr>
        <tr v-for="(email, i) in emails" :key="i">
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto" v-html="highlightTextField(email.from)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto" v-html="highlightTextField(email.to)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto" v-html="highlightTextField(email.subject)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto" v-html="highlightTextField(email.content)">
            </td>
        </tr>
    </tbody>
</template>

<style>
</style>