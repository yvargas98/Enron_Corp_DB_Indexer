<script>
import Loading from './Loading.vue'
export default {
    props: {
        data: {
            type: Array,
            required: true
        },
        value: {
            type: String,
            required: true
        },
        loading: {
            type: Boolean,
            required: true
        },
        error: {
            type: Boolean,
            required: true
        }
    },
    data() {
        return {
            emails: []
        };
    },
    created() {
        this.emails = this.data;
    },
    updated() {
        this.emails = this.data;
    },
    methods: {
        highlightTextField(text) {
            if (text.length > 200) {
                const queryPosition = text.toLowerCase().indexOf(this.value.toLowerCase());
                const paragraphs = text.split('\n\n');
                const formattedContent = paragraphs.map(paragraph => `${paragraph}`).join('');
                const shortContent = (queryPosition > 40 ? "..." : "") + formattedContent.slice(queryPosition < 40 ? 0 : queryPosition - 40, queryPosition + 85) + "...";
                const highlightedContent = shortContent.replace(new RegExp(this.value, "gi"), `<span style="font-weight: bold;" class="bg-blue-200">$&</span>`);
                return highlightedContent;
            }
            else {
                const highlightedContent = text.replace(new RegExp(this.value, "gi"), `<span style="font-weight: bold;" class="bg-blue-200">$&</span>`);
                return highlightedContent;
            }
        },
        showEmailDetail(email) {
            this.$emit('showDetail', email);
        }
    },
    components: { Loading }
}
</script>

<template>
    <tbody>
        <tr v-if="emails.length === 0 && loading === false && error === false">
            <td class="px-4 py-4 text-sm" colspan="4">
                <div class="text-gray-700 dark:text-gray-500 text-center">
                    No result found.
                </div>
            </td>
        </tr>
        <tr v-if="loading === false && error === true">
            <td class="px-4 py-4 text-sm" colspan="4">
                <div class="text-gray-700 dark:text-gray-500 text-center">
                    Error searching content.
                </div>
            </td>
        </tr>
        <tr v-if="loading === true">
            <td class="px-4 py-4 text-sm" colspan="4">
                <Loading />
            </td>
        </tr>
        <tr v-if="!loading" v-for="(email, i) in emails" :key="i" @click="showEmailDetail(email)"
            class="hover:bg-gray-200 cursor-pointer">
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto"
                v-html="highlightTextField(email.from)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto"
                v-html="highlightTextField(email.to)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto"
                v-html="highlightTextField(email.subject)">
            </td>
            <td class="px-4 py-4 text-sm text-gray-700 dark:text-gray-500 w-64 max-h-40 overflow-y-auto"
                v-html="highlightTextField(email.content)">
            </td>
        </tr>
    </tbody>
</template>

<style></style>