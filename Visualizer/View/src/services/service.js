const main_url = window.parent.location.href;
const url = `${main_url}api/default`

export class SearchServices {
    constructor() {
    }

    async searchContents(values) {
        try {
            const searchUrl = new URL(`${url}/_search`);
            requestOption = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: {
                    stream: values.stream,
                    value: values.value,
                    from: values.from,
                    size: values.size
                }
            }

            const response = await fetch(searchUrl.toString(), requestOption);
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