import { Injectable } from '@angular/core'
import { BehaviorSubject } from 'rxjs'
import { getApiWithAuth } from '../tool/httpRequest-auth'

@Injectable({ providedIn: 'root' })
export class GlobalDataService {
    private placeListData: any[] = []
    private placeListDataSubject = new BehaviorSubject<any[]>(this.placeListData)
    placeListData$ = this.placeListDataSubject.asObservable()

    setPlaceListData(data: any[]) {
        this.placeListDataSubject.next(data)
    }

    async loadPlaceListData() {
        const dada = await getApiWithAuth('/base/location/getAll')
        this.setPlaceListData(dada)
    }

    async loadAllGlobalData() {
        await this.loadPlaceListData()
    }
}