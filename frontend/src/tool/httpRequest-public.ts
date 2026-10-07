import { environment } from "../environments/environment.development"


const contentType = 'application/json;charset=UTF-8'

export const postApi = async (url: string, data: any) => {
    const requestHeaders = new Headers()
    requestHeaders.append('Content-Type', contentType)

    const finalUrl = `${environment.apiUrl}${url}`
    const rawResponse = await fetch(finalUrl, {
        method: 'POST',
        headers: {
            'Content-Type': contentType
        },
        body: JSON.stringify(data)
    })
    const content: any = await rawResponse.json()
    return content
}

export const getApi = async (url: string) => {
    const finalUrl = `${environment.apiUrl}${url}`
    const rawResponse = await fetch(finalUrl, {
        method: 'GET',
        headers: {
            'Content-Type': contentType
        }
    })
    const content: any = await rawResponse.json()
    return content
}
export const getApiFromOutside = async (url: string, headers?: any) => {
    const rawResponse = await fetch(url, {
        method: 'GET',
    /*    headers: {
            'Content-Type': contentType
        } */
    })
    const content: any = await rawResponse.json()
    return content
}


export const deleteApi = async (url: string, data: any) => {
    const finalUrl = `${environment.apiUrl}${url}`
    return new Promise((resolve, rejects) => {
        fetch(url, {
            method: 'DELETE',
            headers: {
                'Content-Type': contentType
            }
        })
    })
}

export const putApi = async (url: string, data: any) => {
    const finalUrl = `${environment.apiUrl}${url}`
    return new Promise((resolve, rejects) => {
        fetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': contentType
            },
            body: JSON.stringify(data)
        })
    })
}