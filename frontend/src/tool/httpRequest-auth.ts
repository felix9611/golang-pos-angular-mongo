import { environment } from "../environments/environment.development"

const contentType = 'application/json;charset=UTF-8'

// localStorage on read in browser
const getAccessToken = () => {
    if (typeof window !== 'undefined' && localStorage) {
        return localStorage.getItem('accessToken') || ''
    }
    return ''
}

export const postApiWithAuth = async (url: string, data: any) => {
    const requestHeaders = new Headers()
    requestHeaders.append('Content-Type', contentType)

    const accessToken = getAccessToken()
    if (accessToken) {
        requestHeaders.append('Authorization', accessToken)
    }

    const finalUrl = `${environment.apiUrl}${url}`
    const rawResponse = await fetch(finalUrl, {
        method: 'POST',
        headers: requestHeaders,
        body: JSON.stringify(data)
    })
    return rawResponse.json()
}

export const getApiWithAuth = async (url: string) => {
    const requestHeaders = new Headers()
    requestHeaders.append('Content-Type', contentType)

    const accessToken = getAccessToken()
    if (accessToken) {
        requestHeaders.append('Authorization', accessToken)
    }

    const finalUrl = `${environment.apiUrl}${url}`
    const rawResponse = await fetch(finalUrl, {
        method: 'GET',
        headers: requestHeaders
    })
    return rawResponse.json()
}

export const deleteApiWithAuth = async (url: string) => {
    const requestHeaders = new Headers()
    requestHeaders.append('Content-Type', contentType)

    const accessToken = getAccessToken()
    if (accessToken) {
        requestHeaders.append('Authorization', accessToken)
    }

    const finalUrl = `${environment.apiUrl}${url}`
    const rawResponse = await fetch(finalUrl, {
        method: 'DELETE',
        headers: requestHeaders
    })
    return rawResponse.json()
}

export const putApiWithAuth = async (url: string, data: any) => {
    const requestHeaders = new Headers()
    requestHeaders.append('Content-Type', contentType)

    const accessToken = getAccessToken()
    if (accessToken) {
        requestHeaders.append('Authorization', accessToken)
    }

    const finalUrl = `${environment.apiUrl}${url}`
    const rawResponse = await fetch(finalUrl, {
        method: 'PUT',
        headers: requestHeaders,
        body: JSON.stringify(data)
    })
    return rawResponse.json()
}
