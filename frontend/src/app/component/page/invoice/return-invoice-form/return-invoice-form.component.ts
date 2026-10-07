import { CommonModule } from '@angular/common'
import { Component, OnInit } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { MatIconModule } from '@angular/material/icon'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzTableModule } from 'ng-zorro-antd/table'
import { UpdateReturnInvoiceDto } from './interface'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { ActivatedRoute, Router } from '@angular/router'
import { NzMessageService } from 'ng-zorro-antd/message'
import { UserStoreService } from '../../../../../state/user.service'
import { debounceTime, Subscription, timer } from 'rxjs'
import { findMenuItem } from '../../../tool-function'

@Component({
    imports: [
            CommonModule,
            FormsModule,
            NzFormModule,
            NzSelectModule,
            NzInputModule,
            NzInputNumberModule,
            NzTableModule,
            MatIconModule,
            NzDatePickerModule,
            NzCheckboxModule, 
        ],
        templateUrl: './return-invoice-form.component.html'
})
export class ReturnInvoiceFormComponent implements OnInit {
    private rightSubscription: Subscription
    constructor(
            private route: ActivatedRoute, 
            private routeTo: Router,
            private message: NzMessageService,
            private userStoreService: UserStoreService
    ) {
         this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Create Return', 'create-return-invoice')
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

    ngOnInit(): void {
        this.loadLocationList()
    }

    returnInvoiceForm: UpdateReturnInvoiceDto = {
        returnInvoiceNo: '',
        returnReason: '',
        returnLocationId: '',
        returnDatail: '',
        returnMethod: '',
        processMethod: '',
        refund: false,
        refundAmount: 0,
        refundMethod: '',
        returnToVendor: false,
        returnItems: [],

        _id: ''
    }

    itemSubmitForm: any = {
        productCode: '',
        productId: '',
        qty: 0,
        price: 0,
        taxType: '',
        taxCode: '',
        taxRate: 0,
        taxAmount: 0,
        product: {}
    }

    oldQtys: number = 0

    invoiceDetail: any = {}

    goBackList() {
        this.routeTo.navigate(['invoice-list'])
    }

    async returnInvoiceNoChanged(event: any) {
        if (event)  {
            const data = await getApiWithAuth(`/invoice/number/${event}`)
            this.invoiceDetail = data
        }
    }

    locationList: any[] = []
    async loadLocationList() {
        const data = await getApiWithAuth('/base/location/getAll')
        this.locationList = data
    }

    searchInvoiceItem() {
        const data = this.invoiceDetail.invoiceItems.find((x: any) => x.product.productCode === this.itemSubmitForm.productCode)

        this.itemSubmitForm.productId = data.product._id
        this.itemSubmitForm.product = data.product   
        this.itemSubmitForm.price = data.price
        this.itemSubmitForm.taxRate = data.taxRate      
        this.itemSubmitForm.taxAmount = data.taxAmount
        this.itemSubmitForm.taxType = data.taxType
        this.itemSubmitForm.taxCode = data.taxCode
        this.itemSubmitForm.qty = data.qty  
        this.oldQtys = data.qty
        
    }

    qtyChanged(event: any) {
        const oldAvg = this.itemSubmitForm.price / this.oldQtys
        this.itemSubmitForm.price = event * oldAvg
        this.itemSubmitForm.taxAmount = (event * oldAvg * (1 + this.itemSubmitForm.taxRate ? this.itemSubmitForm.taxRate : 0)).toFixed(2)
        this.itemSubmitForm.qty = event
        this.oldQtys = event
    }

    addToList() {
        const { productId, productCode, qty, price, taxRate, taxCode, taxType, taxAmount, product } = this.itemSubmitForm
        const addToListItem = {
            productId,
            productCode,
            qty,
            price,
            taxRate,
            taxCode, 
            taxType,
            taxAmount,
            product,
            returnInvoiceId: this.invoiceDetail._id // Assuming returnInvoiceId is the _id of the Invoice DETAIL
        }

        this.returnInvoiceForm.refundAmount = this.returnInvoiceForm.refundAmount +Number(this.itemSubmitForm.price)
        
        this.returnInvoiceForm.returnItems.push(addToListItem)
        this.itemSubmitForm = {}
    }

    removeItem(index: number) {
        this.returnInvoiceForm.refundAmount -= Number(this.returnInvoiceForm.returnItems[index].price)
        this.returnInvoiceForm.returnItems.splice(index, 1)
    }


    async submitForm() {
        const data = await postApiWithAuth('/return-invoice/create', this.returnInvoiceForm)
        if (data._id) {
            this.returnInvoiceForm = {
                returnInvoiceNo: '',
                returnReason: '',
                returnLocationId: '',
                returnDatail: '',
                returnMethod: '',
                processMethod: '',
                refund: false,
                refundAmount: 0,
                refundMethod: '',
                returnToVendor: false,
                returnItems: [],
        
                _id: ''
            }

            this.message.info('Data save successfully!')

            timer(2500).subscribe(() => {
                this.routeTo.navigate([`/return-invoice-detail`], { queryParams: { id: data._id } })
             //   this.routeTo.navigate([`/return-invoice-list`])
            })
        }
    }

}
