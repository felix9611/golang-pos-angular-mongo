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
import { FilterListApiSmall, ListOption } from '../../../components/filter-list-small/filter-list.component'

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
        DownloadExcelDataComponent,
        FilterListApiSmall
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
        limit: 10,
        typeIds: [],
        deptIds: []
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
//  Department data
    departmentTotal: number = 0
    departments: any = []
    departmentOptions: ListOption[] = [];
    departmentSearchForm: any = {
        name: '',
        page: 1,
        limit: 10
    }
    loadDeptLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.departmentSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.departments = res ?? [];
            this.departmentTotal = res.length ?? 0;
        
            this.departmentOptions = this.departments.map((a: any) => ({
                value: a.id,
                label: `${a.deptCode} - ${a.deptName}`,
            }))
    
            return this.departmentOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }
//  Product Type data
    productTypeTotal: number = 0
    productTypes: any = []
    productTypeOptions: ListOption[] = [];
    productTypeSearchForm: any = {
        name: '',
        page: 1,
        limit: 10
    }
    loadTypeLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.productTypeSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.productTypes = res ?? [];
            this.productTypeTotal = res.length ?? 0;
        
            this.productTypeOptions = this.productTypes.map((a: any) => ({
                value: a.id,
                label: `${a.typeCode} - ${a.typeName}`,
            }))
    
            return this.productTypeOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }
}