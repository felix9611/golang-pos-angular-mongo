import { CommonModule } from '@angular/common'
import { Component, OnInit } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzTableModule } from 'ng-zorro-antd/table'
import { Subscription } from 'rxjs'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { ListProductLocationtRequestDto } from '../interface'
import { DownloadExcelDataComponent } from '../../../components/download-excel-component/download-excel-data-component.component'
import { Router } from '@angular/router'
import { UserStoreService } from '../../../../../state/user.service'
import { NzMessageService } from 'ng-zorro-antd/message'
import { findMenuItem } from '../../../tool-function'

@Component({
    selector: 'product-location',
    standalone: true,
    templateUrl: './product-location.component.html',
    imports: [
        NzSelectModule, 
        CommonModule, 
        NzFormModule, 
        NzButtonModule, 
        FormsModule, 
        NzTableModule, 
        NzInputModule, 
        NzPaginationModule,
        DownloadExcelDataComponent
    ]
})
export class ProductLocationListComponent {
    private rightSubscription: Subscription
    private userRightSubscription: Subscription
    constructor(
        private userStoreService: UserStoreService,
        private message: NzMessageService,
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Product Location', 'product-location-lists')
            this.excelFileSetting.code = answer?.excelFunctionCode ?? ''
            this.preLoadExcelSetting()
        })

        this.userRightSubscription = this.userStoreService.user$.subscribe((data: any) => {
            this.superAdmin = data.superAdmin
            if (data.superAdmin !== true) {
                this.locationList = data.locations
                this.searchForm.locationIds = data.placeIds

                this.loadProductLocationLists()
            }
        })
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
        if (this.userStoreService.user$) {
            this.userRightSubscription.unsubscribe()
        }
    }

    ngOnInit() {
        this.loadLocationList()
        if (this.superAdmin === true) {
            this.loadProductLocationLists()
        }
    }

    superAdmin: boolean = false

    excelFileSetting: any = {
        code: ''
    }
    
    searchForm: ListProductLocationtRequestDto = {
        page: 1,
        limit: 10,
        locationIds: [],
        assetCode: ''
    }

    totals: number = 0

    dataLists: any[] = []

    async loadProductLocationLists() {
        if (this.superAdmin === true) {
            const res = await postApiWithAuth('/inventory/product-location-list', this.searchForm)
            this.dataLists = res.lists
            this.totals = res.total
        } else {
            if (this.searchForm.locationIds.length === 0) {
                this.message.error('Please select least one location')
            } else {
                const res = await postApiWithAuth('/inventory/product-location-list', this.searchForm)
                this.dataLists = res.lists
                this.totals = res.total
            }
        }
    }

    locationList: any[] = []
    async loadLocationList() {
        if (this.superAdmin === true) {
            const data = await getApiWithAuth('/base/location/all')
            this.locationList = data.data
        }
    }

    dbFieldList: string[] = []
    excelFieldList: string[] = []
    async preLoadExcelSetting() {
        const res = await getApiWithAuth(`/sys/excel-field-match/code/${this.excelFileSetting.code}`)
        this.dbFieldList = res.fieldLists.map((item: any) => item.dbFieldName)
        this.excelFieldList = res.fieldLists.map((item: any) => item.excelFieldName)
    }
}