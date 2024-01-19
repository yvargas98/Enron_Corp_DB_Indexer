import axios from "axios";

const main_url = window.parent.location.href;

export class SearchServices {
    constructor() {
        this.apiclient = axios.create({
            baseURL: `${main_url}api`,
            headers: {
                'Content-Type': 'application/json'
            }
        })
    }

    async searchContents(values) {
        try {

            const response = await this.apiclient.post('/_search', values);
            
            if (response.status !== 200) {
                throw new Error(`API request failed with status: ${response.status}`);
            }
            return response.data;
        } catch (error) {
            console.log('Search API error:', error);
            throw error;
        }
    }
}