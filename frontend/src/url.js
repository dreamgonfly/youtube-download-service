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

export default { extractVideoIdFromURL, isValidURL }