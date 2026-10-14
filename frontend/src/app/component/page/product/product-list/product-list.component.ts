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
import { QRcodeComponent } from '../../../components/qr-code/qr-code.component'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { UserStoreService } from '../../../../../state/user.service'
import { findMenuItem } from '../../../tool-function'
import { Subscription } from 'rxjs'
import { UploadDialogComponent } from '../../../components/upload-dialog-component/upload-dialog-component.component'
import { DownloadExcelTemplateComponent } from '../../../components/download-template-component/download-template-component.component'
import { NzMessageService } from 'ng-zorro-antd/message'
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
        QRcodeComponent,
        NzDatePickerModule,
        DownloadExcelTemplateComponent,
        UploadDialogComponent,
        DownloadExcelDataComponent
    ],
    templateUrl: './product-list.component.html',
    styleUrl: './product-list.component.css',
})
export class ProductListComponent {
    private rightSubscription: Subscription
    constructor(
        private routeTo: Router,
        private userStoreService: UserStoreService,
        private message: NzMessageService,
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Product List', 'product-list')
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
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }

    searchForm: any = {
        page: 1,
        limit: 10
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
        this.loadProductLists()
        this.loadTypeList()
        this.loadDeptList()
    }

    typeLists: any[] = []
    async loadTypeList() {
        this.typeLists = await getApiWithAuth('/product/product-type/all')
    }

    deptLists: any[] = []
    async loadDeptList() {
        const data = await getApiWithAuth('/sys/department/all')
        this.deptLists = data.data
    }

    async loadProductLists() {
        const res = await postApiWithAuth('/product/product-list/list', this.searchForm)
        this.dataLists = res.lists
        this.totals = res.total
    }

    handleRomeve(id: string) {
        this.handleId = id
        this.removeDialog = true
    }

    closeRemoveDialog() {
        this.removeDialog = false
    }


    dateFormat(data: string) {
        return data ? moment(data).format('DD-MM-YYYY HH:MM') : null
    }

    openEdit(id: string) {
        this.routeTo.navigate([`/product-detail`], { queryParams: { id } })
    }

    goToCreate() {
        this.routeTo.navigate(['/product-create'])
    }

    async goToVoid() {
        const res = await await deleteApiWithAuth(`/product/product-list/void/${this.handleId}`)
        if (res.error) {
            this.message.error(res.error)
        } else if (!res.error) {
            this.message.success('Delete successful!')
            this.closeRemoveDialog()
            this.loadProductLists()
        }
    }

    qrCodeDialog: boolean = false
    qrCodeString: string = ''
    qrCodeRandomHtml: string = ''

    openQrCodeDialog(data: any) {
        this.qrCodeDialog = true

        const qrCodeString = `${data.productCode}|${data.productName}|${data.producttype.typeCode??''}|${data.producttype.typeName??''}|${data.department.deptCode??''}|${data.department.deptName??''}`
        this.qrCodeString = qrCodeString

        this.qrCodeRandomHtml = `
            <div class="px-3 text-left grid grid-cols-1">
                <div><span class="font-bold">Product Code:</span>${data.productCode}</div>
                <div><span class="font-bold">Product Name:</span> ${data.productName}</div>
                <div><span class="font-bold">Type Code:</span> ${data.producttype.typeCode??''}</div>
                <div><span class="font-bold">Type Name:</span> ${data.producttype.typeName??''}</div>
                <div><span class="font-bold">Department Code:</span> ${data.department.deptCode??''}</div>
                <div><span class="font-bold">Department Code:</span> ${data.department.deptName??''}</div>
            </div>
        `
    }

    openQrCodeDialogClose(event: any) {
        this.qrCodeDialog = false
    }

    repairRecordDialog: boolean = false
    handleData: any = {}

    openRepairRecordDialog(data: any) {
        this.repairRecordDialog = true
        this.handleId = data.id
    }

    excelFileSetting: any = {
        code: ''
    }

    dbFieldList: string[] = []
    excelFieldList: string[] = []
    async preLoadExcelSetting() {
        const res = await getApiWithAuth(`/sys/excel-field-match/code/${this.excelFileSetting.code}`)
        this.dbFieldList = res.fieldLists.map((item: any) => item.dbFieldName)
        this.excelFieldList = res.fieldLists.map((item: any) => item.excelFieldName)
    }
}