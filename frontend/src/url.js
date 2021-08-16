export const getVideoId = (input) => {
    if (isValidURL(input)) {
        return extractVideoIdFromURL(input)
    } else {
        return input
    }
}


export const reconstructYoutubeURL = (id) => {
    return `https://www.youtube.com/watch?v=${id}`
}

export const isValidURL = (input) => {
    try {
        new URL(input)
    } catch (error) {
        return false
    }
    return true
}

export const extractVideoIdFromURL = (url) => {
    const u = new URL(url)
    if (url.includes("youtube.com")) {
        return u.searchParams.get('v')
    }
    if (url.includes("youtu.be")) {
        return u.pathname.substring(1)  // Remove leftmost slash
    }
    return null
}
