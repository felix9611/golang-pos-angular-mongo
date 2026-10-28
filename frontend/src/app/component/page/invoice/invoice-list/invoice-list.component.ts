import { Component } from '@angular/core'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { deleteApiWithAuth, getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzFormModule } from 'ng-zorro-antd/form'
import moment from 'moment'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { Router } from '@angular/router'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { UserStoreService } from '../../../../../state/user.service'
import { findMenuItem } from '../../../tool-function'
import { Subscription } from 'rxjs'
import { NzMessageService } from 'ng-zorro-antd/message'
import { InvoiceListRequestDto } from './interface'
import { DownloadExcelDataComponent } from '../../../components/download-excel-component/download-excel-data-component.component'

@Component({
    // selector: 'app-footer',
    standalone: true,
    imports: [ 
        NzSelectModule, 
        CommonModule, 
        NzFormModule, 
        NzButtonModule, 
        FormsModule, 
        NzModalModule, 
        NzTableModule, 
        NzInputModule, 
        NzPaginationModule,
        NzDatePickerModule,
        DownloadExcelDataComponent
    ],
    templateUrl: './invoice-list.component.html',
    styleUrl: './invoice-list.component.css',
})
export class InvoiceListComponent {
    private rightSubscription: Subscription
    private userRightSubscription: Subscription
    constructor(
        private routeTo: Router,
        private userStoreService: UserStoreService,
        private message: NzMessageService,
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Invoice List', 'invoice-list')
            console.log(answer)
            this.userRightInside = {
                read: answer?.read ?? false,
                write: answer.write ?? false,
                update: answer.update ?? false,
                delete: answer.delete ?? false,
                upload: answer.upload ?? false
                 // keep default value
            }
            this.excelFileSetting.code = answer?.excelFunctionCode ?? ''
            this.preLoadExcelSetting()
        })

        this.userRightSubscription = this.userStoreService.user$.subscribe((data: any) => {
            this.superAdmin = data.superAdmin
            if (data.superAdmin !== true) {
                this.locationList = data.locations
                this.searchForm.locationId = data.placeIds
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

    superAdmin: boolean = false

    excelFileSetting: any = {
        code: ''
    }

    dbFieldList: string[] = []
    excelFieldList: string[] = []

    searchForm: InvoiceListRequestDto = {
        page: 1,
        limit: 10,
        number: '',
        dateRange: [],
        locationId: []
    }

    userRightInside: any = {
        read: false,
        write: false,
        update: false,
        delete: false
    }


    dataLists: any[] = []
    totals: number = 0
    editFormDialog: boolean = false
    removeDialog: boolean = false
    handleId: string = ''

    ngOnInit() {
        this.loadInvoiceLists()
        this.loadLocationList()
    }

    locationList: any[] = []

    async loadLocationList() {
        if (this.superAdmin === true) {
            const data = await getApiWithAuth('/base/location/all')
            this.locationList = data
        }
    }

    async loadInvoiceLists() {
        if (this.superAdmin === true) {
            const res = await postApiWithAuth('/invoice/list', this.searchForm)
            this.dataLists = res.lists
            this.totals = res.total
        } else {
            if (this.searchForm.locationId.length === 0) {
                this.message.error('Please select least one location')
            } else {
                const res = await postApiWithAuth('/invoice/list', this.searchForm)
                this.dataLists = res.lists
                this.totals = res.total
            }
        }
    }

    handleRomeve(id: string) {
        this.handleId = id
        this.removeDialog = true
    }

    closeRemoveDialog() {
        this.removeDialog = false
    }


    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    openDetail(id: string) {
        this.routeTo.navigate([`/invoice-detail`], { queryParams: { id } })
    }

    async goToVoid() {
        const res = await await deleteApiWithAuth(`/invoice/void/${this.handleId}`)
        if (res.error) {
            this.message.error(res.error)
        } else if (!res.error) {
            this.message.success('Void successful!')
            this.closeRemoveDialog()   
            this.loadInvoiceLists()
        }
    }

    async preLoadExcelSetting() {
        const res = await getApiWithAuth(`/sys/excel-field-match/code/${this.excelFileSetting.code}`)
        this.dbFieldList = res.fieldLists.map((item: any) => item.dbFieldName)
        this.excelFieldList = res.fieldLists.map((item: any) => item.excelFieldName)
    }
}