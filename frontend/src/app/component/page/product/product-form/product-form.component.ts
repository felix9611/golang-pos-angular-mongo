import { Component, OnInit } from '@angular/core'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzFormModule } from 'ng-zorro-antd/form'
import moment from 'moment'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { ProductFormDto } from './interface'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { ActivatedRoute, Router } from '@angular/router'
import { debounceTime, Observable, Observer, Subject, Subscription, timer } from 'rxjs'
import { MatButtonModule } from '@angular/material/button'
import { MatIconModule } from '@angular/material/icon'
import { NzUploadChangeParam, NzUploadFile, NzUploadModule } from 'ng-zorro-antd/upload'
import { NzIconModule } from 'ng-zorro-antd/icon'
import { uploadImgToBase64 } from '../../../../../tool/imageUpload'
import { FileViewComponent } from '../../../components/file-view-dialog/file-view-dialog.component'
import { findMenuItem } from '../../../tool-function'
import { UserStoreService } from '../../../../../state/user.service'
import { ListOption } from '../../../components/filter-list/filter-list.component'
import { FilterListApiSmall } from '../../../components/filter-list-small/filter-list.component'
// import { UploadComponentComponent } from '../../../components/upload-component/upload-component.component'

@Component({
    // selector: 'app-footer',
    standalone: true,
    imports: [
        CommonModule, 
        NzFormModule,
        NzButtonModule, 
        FormsModule, 
        NzModalModule, 
        NzTableModule, 
        NzInputModule, 
        NzPaginationModule,
        NzSelectModule,
        NzDatePickerModule,
        NzCheckboxModule,
        NzInputNumberModule,
        MatIconModule, 
        MatButtonModule, 
        NzUploadModule,
        NzButtonModule, 
        NzIconModule,
        NzModalModule,
        FileViewComponent,
        FilterListApiSmall
    ],
    templateUrl: './product-form.component.html',
    styleUrl: './product-form.component.css',
})
export class ProductFormComponent implements OnInit {
    private rightSubscription: Subscription
    constructor(
        private route: ActivatedRoute, 
        private routeTo: Router,
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {
        this.changeEvent$.pipe(debounceTime(300)).subscribe(event => {
            this.preAction(event.file.originFileObj);
        })
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Product List', 'product-list')
            this.userRightInside = {
                read: answer?.read ?? false,
                write: answer.write ?? false,
                update: answer.update ?? false,
                delete: answer.delete ?? false,
                upload: answer.upload ?? false
                 // keep default value
            }
        })
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }

    userRightInside: any = {
        read: false,
        write: false,
        update: false
    }

    private changeEvent$ = new Subject<NzUploadChangeParam>()

    editForm: ProductFormDto = {
        id: '',
        productCode: '',
        productName: '',
        itemCode: '',
        brandCode: '',
        brandName: '',
        typeId: '',
        deptId: '',
        vendorId: '',
        unit: '',
        costPrice: 0,
        retailPrice: 0,
        description:  '',
        remark:  '',
        taxType: '',
        uploadProductFiles: [],
        productListFiles: []
    }

    taxInformation: boolean = false

    theId: any = ''

    fileUpdloadList: any = []

    ngOnInit() {
        this.route.queryParams.subscribe((x: any) => {
            if (x.id) {
                this.theId = x.id
                this.getOne()
            }
        })

        this.loadCodeTypeList()
    }

    beforeUpload = (file: NzUploadFile, _fileList: NzUploadFile[]): Observable<boolean> =>
        new Observable((observer: Observer<boolean>) => {
        const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
        if (!isJpgOrPng) {
            this.message.error('You can only upload JPG or PNGfile!')
            observer.complete()
       //     return
        }
        const isLt2M = file.size! / 1024 / 1024 < 3
        if (!isLt2M) {
                this.message.error('Image must smaller than 3MB!')
             //   observer.complete()
  //      return
        }
    
        observer.next(isJpgOrPng && isLt2M)
        observer.complete()
    })

    finalFileList: any[] = []
    handleChange(event: any) {
        switch(event.type) {
            case 'start':
                this.changeEvent$.next(event)
            break

            case 'removed':
                this.finalFileList = this.finalFileList.filter(f => f.fileName !== event.file.originFileObj.name)
            break
        }
        console.log(event.type)
        console.log(this.finalFileList)
        
    }

    async preAction(fileObj: any) {
        const response: any = await uploadImgToBase64(fileObj)
        this.finalFileList.push({
            fileName: fileObj.name,
            fileType: fileObj.type,
            base64: response.data
        })
    }

    async getOne() {
        this.editForm = await getApiWithAuth(`/product/product-list/one/${ this.theId}`)
        console.log(this.editForm, 'k')
    }

    async submitForm() {
        const url = this.editForm.id ? '/product/product-list/update' : '/product/product-list/create'
        this.editForm.uploadProductFiles = this.finalFileList

        const res = await postApiWithAuth(url, this.editForm)

        if (!res.error) {
            this.message.info('Data save successfully!')
            timer(2500).subscribe(() => {
                this.routeTo.navigate(['/product-list'])
            })
        } else {
            this.message.error(res.error)
        }
    }

    resetForm() {
        this.editForm = {
            id: '',
            productCode: '',
            productName: '',
            itemCode: '',
            brandCode: '',
            brandName: '',
            typeId: '',
            deptId: '',
            vendorId: '',
            unit: '',
            costPrice: 0,
            retailPrice: 0,
            description:  '',
            remark:  '',
            taxType: '',
            uploadProductFiles: [],
            productListFiles: []
        }
    }

    backToList() {
        this.routeTo.navigate([`/product-list`])
    }

    fileDialogVisible: boolean = false
    openFileDialog() {
        this.fileDialogVisible = true
    }

    closeFileDialog() {
        this.fileDialogVisible = false
    }

    codeTypeLists: any[] = []
    async loadCodeTypeList() {
        const res = await getApiWithAuth('/base/code-type/get-type/Product-TaxType')
        this.codeTypeLists = res.datas
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
//  Vendor data
    vendorTotal: number = 0
    vendors: any = []
    vendorOptions: ListOption[] = [];
    vendorSearchForm: any = {
        name: '',
        page: 1,
        limit: 10
    }
    loadVendorLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.vendorSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.vendors = res ?? [];
            this.vendorTotal = res.length ?? 0;
        
            this.vendorOptions = this.vendors.map((a: any) => ({
                value: a.id,
                label: `${a.vendorCode} - ${a.vendorName}`,
            }))
    
            return this.vendorOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }
    // Vendor data
    taxListTotal: number = 0
    taxLists: any = []
    taxListOptions: ListOption[] = [];
    taxListSearchForm: any = {
        nameCode: '',
        tax: '',
        page: 1,
        limit: 10
    }
    loadTaxListLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.vendorSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.taxLists = res ?? []
            this.taxListTotal = res.length ?? 0
        
            this.taxListOptions = this.taxLists.map((x: any) => ({
                value: x.taxType,
                label: `${x.countryName} ${x.taxCode} (${Number(x.taxRate) * 100}%)`
            }))
    
            return this.taxListOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }

    /* new filter field */
}