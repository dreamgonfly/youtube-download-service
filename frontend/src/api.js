import Axios from 'axios'

const BASE_URL = process.env.REACT_APP_BACKEND_URL

const axios = Axios.create({
    baseURL: BASE_URL,
})

const preview = (id) => {
    return axios.get(`/preview/` + id)
}

const updateThumbnail = (id, url, name) => {
    return axios.post(`/update-thumbnail/` + id, {
        URL: url,
        Name: name
    })
}

const play = (id, format, filename) => {
    const url = new URL(`${BASE_URL}/play/${id}`)
    url.searchParams.append("format", format)
    url.searchParams.append("filename", filename)
    return axios.get(url.href)
}

const composeDownloadLink = (id, format, filename) => {
    const url = new URL(`${BASE_URL}/download/${id}`)
    url.searchParams.append("format", format)
    url.searchParams.append("filename", filename)
    return url.href
}

const api = { preview, updateThumbnail, play, composeDownloadLink }

export default api
