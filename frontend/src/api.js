import Axios from 'axios';

const BASE_URL = 'http://localhost:8000/';

const axios = Axios.create({
    baseURL: BASE_URL,
});


const hello = () => {
    return axios.get(`/hello`)
}


const preview = (id) => {
    return axios.get(`/preview/` + id)
}

const updateThumbnail = (id, url, name) => {
    return axios.post(`/update-thumbnail/` + id, {
        URL: url,
        Name: name
    })
}

const composeDownloadLink = (id, format) => {
    return `${BASE_URL}/download/${id}?format=${format}`
}

export default { preview, hello, composeDownloadLink, updateThumbnail }