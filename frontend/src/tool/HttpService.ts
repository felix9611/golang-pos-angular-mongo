import { Injectable } from '@angular/core'
import { environment } from '../environments/environment.development'
import { LocalStorageService } from '../state/localStorage.service'
import { resolve } from 'path'
import { rejects } from 'assert'
@Injectable({
  providedIn: 'root'
})
export class HttpService {
    constructor(private localStorageService: LocalStorageService) {}
    private mainUrl = `${environment.apiUrl}`
    private requestHeaders = new Headers()
    private contentType = 'application/json;charset=UTF-8'
    

    async getApi(url: string) {
        const finalUrl = `${this.mainUrl}${url}`

        this.requestHeaders.set('Content-Type', this.contentType)
        if (this.localStorageService.getItem('accessToken')) {
            this.requestHeaders.set('Content-Type', this.localStorageService.getItem('accessToken') || '')
        }
        

        const rawResponse = await fetch(finalUrl, {
            method: 'GET',
            headers: this.requestHeaders
        })
        const content: any = await rawResponse.json()
        return content

    }

    async postApi(url: string, data: any): Promise<any> {
        const finalUrl = `${this.mainUrl}${url}`

        this.requestHeaders.set('Content-Type', this.contentType)
        if (this.localStorageService.getItem('accessToken')) {
            this.requestHeaders.set('Content-Type', this.localStorageService.getItem('accessToken') || '')
        }
        const rawResponse = await fetch(finalUrl, {
                method: 'POST',
                headers: this.requestHeaders,
                body: JSON.stringify(data)
        })
        const content: any = await rawResponse.json()
        return content


        /*
        const rawResponse = await fetch(finalUrl, {
            method: 'POST',
            headers: this.requestHeaders,
            body: JSON.stringify(data)
        })
        const content: any = await rawResponse.json()
        return content
        */

    }

}