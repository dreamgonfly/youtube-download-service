const reconstructYoutubeURL = (id) => {
    return `https://www.youtube.com/watch?v=${id}`
}

const isValidURL = (input) => {
    let isValid = true
    try {
        new URL(input)
    } catch (error) {
        isValid = false
    }
    return isValid
}

const extractVideoIdFromURL = (url) => {
    const u = new URL(url)
    return u.searchParams.get('v')
}

const url = { extractVideoIdFromURL, isValidURL, reconstructYoutubeURL }

export default url