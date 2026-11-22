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
import { ListOption } from '../../../components/filter-list/filter-list.component'
import { FilterListApiSmall } from '../../../components/filter-list-small/filter-list.component'

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
        DownloadExcelDataComponent,
        FilterListApiSmall
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
            this.manager = data.manager
            if (data.superAdmin !== true) {
                this.searchForm.locationIds = data.placeIds
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
        this.loadProductLocationLists()
    }

    manager: boolean = false
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
        const res = await postApiWithAuth('/inventory/product-location-list', this.searchForm)
        this.dataLists = res.lists
        this.totals = res.total
    }

    dbFieldList: string[] = []
    excelFieldList: string[] = []
    async preLoadExcelSetting() {
        const res = await getApiWithAuth(`/sys/excel-field-match/code/${this.excelFileSetting.code}`)
        this.dbFieldList = res.fieldLists.map((item: any) => item.dbFieldName)
        this.excelFieldList = res.fieldLists.map((item: any) => item.excelFieldName)
    }
    // Location data
    placeTotal: number = 0
    places: any = []
    placeOptions: ListOption[] = [];
    placeSearchForm: any = {
        placeCode: '',
        placeName: '',
        page: 1,
        limit: 10
    }
    loadPlaceLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.placeSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.places = res ?? [];
            this.placeTotal = res.length ?? 0;
        
            this.placeOptions = this.places.map((a: any) => ({
                value: a.id,
                label: `${a.placeCode}`,
            }))
    
            return this.placeOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }
}