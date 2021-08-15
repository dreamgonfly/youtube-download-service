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
    if (url.includes("youtube.com")) {
        return u.searchParams.get('v')
    }
    if (url.includes("youtu.be")) {
        return u.pathname.substring(1)  // Remove leftmost slash
    }
    return null
}

const url = { extractVideoIdFromURL, isValidURL, reconstructYoutubeURL }

export default url