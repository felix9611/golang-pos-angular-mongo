import { CommonModule } from '@angular/common'
import { Component, OnInit } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzTableModule } from 'ng-zorro-antd/table'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { ListInvRecordDto } from '../interface'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import moment from 'moment'
import { DownloadExcelDataComponent } from '../../../components/download-excel-component/download-excel-data-component.component'
import { findMenuItem } from '../../../tool-function'
import { Subscription } from 'rxjs'
import { UserStoreService } from '../../../../../state/user.service'

@Component({
    standalone: true,
    templateUrl: './inventory-record.component.html',
    imports: [
        NzSelectModule, 
        CommonModule, 
        NzFormModule, 
        NzButtonModule, 
        FormsModule, 
        NzTableModule, 
        NzInputModule, 
        NzPaginationModule,
        NzDatePickerModule,
        DownloadExcelDataComponent
    ]
})
export class InventoryRecordComponent implements OnInit {
    private rightSubscription: Subscription
    constructor(
        private userStoreService: UserStoreService,
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Inventory Record', 'inventory-record')
            this.excelFileSetting.code = answer?.excelFunctionCode ?? ''
            this.preLoadExcelSetting()
        })
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }

    ngOnInit() {
        this.loadInventoryRecordLists()
    }

    excelFileSetting: any = {
        code: ''
    }

    searchForm: ListInvRecordDto = {
        page: 1,
        limit: 10
    }

    totals: number = 0

    dataLists: any[] = []

    async loadInventoryRecordLists() {
        const res = await postApiWithAuth('/inventory-record/list', this.searchForm)
        this.dataLists = res.lists
        this.totals = res.total
    }

    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    dbFieldList: string[] = []
    excelFieldList: string[] = []
    async preLoadExcelSetting() {
        const res = await getApiWithAuth(`/sys/excel-field-match/code/${this.excelFileSetting.code}`)
        this.dbFieldList = res.fieldLists.map((item: any) => item.dbFieldName)
        this.excelFieldList = res.fieldLists.map((item: any) => item.excelFieldName)
    }
}