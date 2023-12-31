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
        console.log("this.data", this.data)
    },
    updated() {
        this.emails = this.data
    },
    methods: {
        formatContent(content) {
            const queryPosition = content.toLowerCase().indexOf(this.value.toLowerCase());
            const paragraphs = content.split('\n\n');
            const formattedContent = paragraphs.map(paragraph => `${paragraph}`).join('');
            const shortContent = (queryPosition > 40 ? "..." : "") + formattedContent.slice(queryPosition < 40 ? 0 : queryPosition - 40, queryPosition + 85) + "..."
            const highlightedContent = shortContent.replace(this.value, `<span style="font-weight: bold;">${this.value}</span>`);
            return highlightedContent
        }
    }
}
</script>

<template>
    <tbody>
        <tr v-for="(email, i) in emails" :key="i">
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500" v-html="formatContent(email.from)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500" v-html="formatContent(email.to)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500" v-html="formatContent(email.subject)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500" v-html="formatContent(email.content)">
            </td>
        </tr>
    </tbody>
</template>

<style>
</style>