<script >
import MamuroHeader from "../components/MamuroHeader.vue";
import Input from "../components/Input.vue";
import Button from "../components/Button.vue";
import TableHeaders from "../components/TableHeaders.vue";
import TableBody from "../components/TableBody.vue";
import Pagination from "../components/Pagination.vue";
import EmailDetails from "../components/EmailDetails.vue";

import { SearchServices } from "../services/service.js"

export default {
    components: {
    MamuroHeader,
    Input,
    Button,
    TableHeaders,
    TableBody,
    Pagination,
    EmailDetails
},
    created() {
        this.searchService = new SearchServices()
    },
    data() {
        return {
            value: "",
            offset: 1,
            limit: 10,
            totalResults: 0,
            emails: [],
            showDetail: false,
            selectedEmail: null,
            loading: false,
            error: false,
        }
    },
    methods: {
        async searchContents() {
            this.loading = true
            const requestValues = {
                stream: import.meta.env.VITE_INDEX_NAME,
                value: this.value,
                from: this.offset,
                size: this.limit
            }
            if (this.value != "") {
                const content = await this.searchService.searchContents(requestValues)
                const formatEmails = content.hits.map((item) => ({
                    content: item._source.content,
                    date: item._source.date,
                    from: item._source.from,
                    subject: item._source.subject,
                    to: item._source.to,
                }))
                this.setEmails(formatEmails)
                this.totalResults = formatEmails.length;
            }
            else {
                this.setEmails([]);
                this.totalResults = 0;
            }
        },
        setEmails(newEmails) {
            this.emails = [...newEmails]
            this.loading = false;
        },
        created() {
            this.emails = [""]
        },
        async updateOffset(newOffset) {
            this.offset = newOffset;
            await this.searchContents().catch(() => {
                this.loading = false
                this.error = true
            });
        },
        handleSearch() {
            this.offset = 1;
            this.searchContents().catch(() => {
                this.loading = false
                this.error = true
            });
        },
        showEmailDetail(email) {
            this.selectedEmail = email;
            this.showDetail = true;
        },
        closeEmailDetail() {
            this.showDetail = false;
            this.selectedEmail = null;
        }
    }
}
</script>

<template>
    <MamuroHeader />
    <section class="container px-4 mx-auto">
        <div v-show="!showDetail" class="mt-6 md:flex md:items-center md:justify-between">
            <div class="relative flex items-center mt-4 md:mt-0">
                <span class="absolute">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                        stroke="currentColor" class="w-5 h-5 mx-3 text-gray-400 dark:text-gray-600">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
                    </svg>
                </span>
                <Input @enter="handleSearch" v-model="this.value" />
                <Button :onClick="handleSearch"></Button>
            </div>
        </div>
        <div v-show="!showDetail">
            <div class="flex flex-col mt-6">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
                        <div class="overflow-hidden border border-gray-200 dark:border-gray-700 md:rounded-lg">
                            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                                <TableHeaders />
                                <TableBody :data="this.emails" :value="this.value" @showDetail="showEmailDetail" :loading="this.loading" :error="this.error"/>
                            </table>
                        </div>
                    </div>
                </div>
                <Pagination :offset="offset" :limit="limit" :totalResults="totalResults" @update:offset="updateOffset">
                </Pagination>
            </div>
        </div>
        <EmailDetails v-if="showDetail" :selectedEmail="selectedEmail" @closeDetail="closeEmailDetail"/>
    </section>
</template>

<style scoped>
.bg-email-search {
    background-image: url('../assets/email_search.jpeg');
}
</style>