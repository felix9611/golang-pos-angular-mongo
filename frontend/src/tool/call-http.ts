import { environment } from '../environments/environment.development'
const requestHeaders = new Headers()
const contentType = 'application/json;charset=UTF-8'

export async function callPostApi(url: string, data: any) {
    const finalUrl = `${environment.apiUrl}${url}`

    requestHeaders.set('Content-Type', contentType)

    if (localStorage.getItem('accessToken')) {
        const accessToken = localStorage.getItem('accessToken')
        requestHeaders.set('Content-Type', accessToken || '')
    }
    const rawResponse = await fetch(finalUrl, {
        method: 'POST',
        headers: requestHeaders,
        body: JSON.stringify(data)
    })
    const content: any = await rawResponse.json()
    return content

}