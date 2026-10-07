import { Component, OnInit, ViewChild, AfterViewInit } from '@angular/core'
import { getApiWithAuth, postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzFormModule } from 'ng-zorro-antd/form'
import { FormsModule } from '@angular/forms'
import { CommonModule } from '@angular/common'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzTableModule } from 'ng-zorro-antd/table'
import { CreateInvoiceForm, InvoicePaymentDto, MmeberForm, ProductFormDto } from './interface'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { MatButtonModule } from '@angular/material/button'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzMessageService } from 'ng-zorro-antd/message'
import { timer } from 'rxjs'
import { Router } from '@angular/router'

@Component({
    templateUrl: './sale-point.component.html',
    imports: [
        NzFormModule,
        FormsModule,
        CommonModule,
        NzSelectModule,
        NzInputModule,
        NzInputNumberModule,
        NzTableModule,
        NzButtonModule,
        MatButtonModule, 
        NzModalModule
    ]
})
export class SalePointComponent implements OnInit {
    constructor(
        private routeTo: Router,
        private message: NzMessageService,
    ) {}

    ngOnInit(): void {
        this.loadLocationList()
        this.orgTotal = 0
    }

    locationList: any[] = []
    async loadLocationList() {
        const data = await getApiWithAuth('/base/location/getAll')
        this.locationList = data
        this.inputedProductList = []
        this.inputedPaymentList = []
    }

    invoiceForm: CreateInvoiceForm = {
        memberId: '',
        locationId: '',
        locationCode: '',
        totalAmount: 0,
        discount: 0,
        discountType: '',
        taxTotal: 0,
        remark: '',
        invoiceItems: [],
        invoicePayments: []
    }

    enterProductForm: any = {
        productCode: '',
        productId: '',
        productName: '',
        locationId: '',
        qty: 0,
        retailPrice: 0,
        discount: 0,
        discountType: '',
        allTotalPrice: 0,
        taxRate: 0,
        taxRateDisplay: 0,
        taxAmount: 0,
        taxType: '',
        taxCode: ''
    }

    productForm: ProductFormDto = {
        _id: '',
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
        description: '',
        remark: '',
        taxType: ''
    }

    taxInfoFilter = {
        taxType: '',
        countryCode: ''
    }

    memberSearch: string = ''

    discountList = [
        { type: '%' },
        { type: '$' },
    ]
    
    payMethodCategory: any = [
        { key: 'Cash' },
        { key: 'IC Card' },
        { key: 'Credit Card' },
        { key: 'E-payment' },
    ]

    inputedProductList: any[] = []
    inputedPaymentList: any[] = []

    async productCodeChanged(event: any) {
        if (event) {
            
            const data = await getApiWithAuth(`/product/product-list/code/${event}`)
            this.enterProductForm.productId = data._id
            this.enterProductForm.retailPrice = data.retailPrice
            this.productForm = data

            this.taxInfoFilter.taxType = data.taxType
        }
    }

    locationIdChanged(event: any) {
        if (event)  {
            const data = this.locationList.find((item: any) => item._id === event)
            this.taxInfoFilter.countryCode = data.country
            this.findTaxInfoForSalePoint()
        }
    }

    qtyChanged(event: any) {
        this.enterProductForm.qty = event
        this.enterProductForm.allTotalPrice = event * this.enterProductForm.retailPrice
        this.enterProductForm.taxAmount =  (event * (1 + this.enterProductForm.taxRate) * this.enterProductForm.retailPrice).toFixed(2)
    }

    typeChanged(event: any) {
        if (event === '%') {
            this.enterProductForm.allTotalPrice = this.enterProductForm.retailPrice * (1 - this.enterProductForm.discount /100) * this.enterProductForm.qty
            this.enterProductForm.taxAmount =  ((1 + this.enterProductForm.taxRate) * this.enterProductForm.retailPrice * (1 - this.enterProductForm.discount /100) * this.enterProductForm.qty).toFixed(2)
        } else if (event === '$') {
            this.enterProductForm.allTotalPrice = (this.enterProductForm.retailPrice - this.enterProductForm.discount) * this.enterProductForm.qty
            this.enterProductForm.taxAmount = ((this.enterProductForm.retailPrice  - this.enterProductForm.discount) * (1 + this.enterProductForm.taxRate) * this.enterProductForm.qty).toFixed(2)
        }
    }

    addToList() {
        const { productId, productCode, qty, allTotalPrice, discount, discountType, taxAmount, taxRate, taxCode, taxType } = this.enterProductForm
        const addToList = {
            productId,
            productCode,
            qty,
            price: allTotalPrice,
            discount,
            discountType,
            taxAmount, 
            taxRate,
            taxCode, 
            taxType
        }

        this.inputedProductList.push(addToList)
        this.orgTotal += Number(taxAmount)

        this.enterProductForm = {
            productCode: '',
            productId: '',
            productName: '',
            locationId: '',
            qty: 0,
            retailPrice: 0,
            disscount: 0,
            disscountType: '',
            allTotalPrice: 0,
            taxAmount: 0, 
            taxType: '',
            taxCode: '',
            taxRate: 0, 
            taxRateDisplay: 0
        }
    }

    orgTotal: number = 0
      

    // RIGHT

    removeItem(index: number) {
        this.orgTotal -= Number(this.inputedProductList[index].taxAmount)
        this.inputedProductList.splice(index, 1)
    }

    memberLists: any[] = []
    async findMmebers() {
        this.memberLists = await postApiWithAuth('/member/member-list/list-member', { name: this.memberSearch })
    }

    memberForm: MmeberForm  = {
        _id: '',
        memberCode: '',
        name: '',
        phone: '',
        email: ''
    }


    fillMemberInfo(data: any) {
        this.memberForm = data
        this.invoiceForm.memberId = data._id
    }

    totalTypeChnage(event: any) {
        if (event === '%') {
            this.invoiceForm.totalAmount = parseFloat((this.orgTotal * (1 - this.invoiceForm.discount /100)).toFixed(2))
        } else if (event === '$') {
            this.invoiceForm.totalAmount = parseFloat((this.orgTotal - this.enterProductForm.disscount).toFixed(2))
        }
    }


    paymentForm: InvoicePaymentDto = {
        method: '',
        amount: 0,
        balance: 0,
        findRedemption: 0
    }

    addPaymentList() {
        const { method, amount } = this.paymentForm
        const addToList = {
            method,
            amount,
            paymentTime: new Date().toISOString()
        }

        this.inputedPaymentList.push(addToList)
        
        this.paymentForm.findRedemption = amount > this.paymentForm.balance ? amount - this.paymentForm.balance : 0
        this.paymentForm.balance = this.invoiceForm.totalAmount - amount
        

    }

    lastPaymentDialogOpen: boolean = false

    openPaymentDialg() {
        this.invoiceForm.totalAmount = this.invoiceForm.totalAmount ? this.invoiceForm.totalAmount : this.orgTotal
        this.paymentForm.balance = this.invoiceForm.totalAmount ? this.invoiceForm.totalAmount : this.orgTotal
        this.lastPaymentDialogOpen = true
    }

    closePaymentDialg() {
        this.lastPaymentDialogOpen = false
    }

    removePayItem(index: number, amount: any) {
        this.paymentForm.balance = this.paymentForm.balance + amount.amount
        this.paymentForm.findRedemption = this.paymentForm.findRedemption - amount.amount
        this.inputedPaymentList.splice(index, 1)
    }

    async submitPOSale() {

        const finalDataSubmit = {
            ...this.invoiceForm,
            discount: this.invoiceForm.discount ? this.invoiceForm.discount / 100 : 0,
            invoiceItems: this.inputedProductList,
            invoicePayments: this.inputedPaymentList
        }

        const res = await postApiWithAuth('/invoice/create', finalDataSubmit)
        if (res) {
            this.message.success('Invoice created successfully!')
            timer(2500).subscribe(() => {
                this.routeTo.navigate(['/invoice-detail'], { queryParams: { id: res._id}})
            })
        } else {
            this.message.error(res.msg)
        }
    }

    

    async findTaxInfoForSalePoint() {
        const res = await postApiWithAuth('/base/tax-information/sales-point', this.taxInfoFilter)
        this.enterProductForm.taxRate = res.taxRate
        this.enterProductForm.taxRateDisplay = res.taxRate * 100
        this.enterProductForm.taxCode = res.taxCode
        this.enterProductForm.taxType = res.taxType
    }
    
}
