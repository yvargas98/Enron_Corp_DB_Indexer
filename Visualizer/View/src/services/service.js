const main_url = window.parent.location.href;
const url = `${main_url}api/default`

export class SearchServices {
    constructor() {
        this.requestOption = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        }
    }

    async searchContents(values) {
        try {
            const searchUrl = new URL(`${url}/_search`);
            Object.keys(values).forEach((key) => {
                searchUrl.searchParams.append(key, values[key]);
            });

            const response = await fetch(searchUrl.toString(), this.requestOption);
            if (!response.ok) {
                throw new Error(response.status);
            }

            const data = await response.json();
            return data;
        } catch (error) {
            console.log(error);
            throw error;
        }
    }
}