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
import { Subscription, timer } from 'rxjs'
import { Router } from '@angular/router'
import { UserStoreService } from '../../../../state/user.service'

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
    private userRightSubscription: Subscription
    constructor(
        private routeTo: Router,
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {
        this.userRightSubscription = this.userStoreService.user$.subscribe((data: any) => {
            this.superAdmin = data.superAdmin
            if (data.superAdmin !== true) {
                this.locationList = data.locations
                if (this.locationList.length === 1) {
                    this.invoiceForm.locationId = this.locationList[0].id
                }
            }
        })
    }

    ngOnDestroy() {
        if (this.userStoreService.user$) {
            this.userRightSubscription.unsubscribe()
        }
    }

    superAdmin: boolean = false

    ngOnInit(): void {
        this.loadLocationList()
    }

    locationList: any[] = []
    async loadLocationList() {
        if (this.superAdmin === true) {
            const data = await getApiWithAuth('/base/location/all')
            this.locationList = data.data
        }
        
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
            this.enterProductForm.productId = data.id
            this.enterProductForm.retailPrice = data.retailPrice
            this.enterProductForm.disscount = 0 
            this.productForm = data

            if (data.id) {
                this.qtyChanged(1)
            }
            
            this.taxInfoFilter.taxType = data.taxType
            if (this.locationList.length === 1) {
                this.taxInfoFilter.countryCode = this.invoiceForm.locationCode
                this.findTaxInfoForSalePoint()
            }
        }
    }

    locationIdChanged(event: any) {
        if (event)  {
            const data = this.locationList.find((item: any) => item.id === event)
            this.taxInfoFilter.countryCode = data.country
            this.findTaxInfoForSalePoint()
        }
    }

    qtyChanged(event: any) {
        this.enterProductForm.qty = event
        this.enterProductForm.allTotalPrice = (event * this.enterProductForm.retailPrice).toFixed(2)
        this.enterProductForm.taxAmount =  (event * (1 + this.enterProductForm.taxRate) * this.enterProductForm.retailPrice).toFixed(2)
    }

    typeChanged(event: any) {
        if (event === '%') {
            this.enterProductForm.allTotalPrice = (this.enterProductForm.retailPrice * (1 - this.enterProductForm.discount /100) * this.enterProductForm.qty).toFixed(2)
            this.enterProductForm.taxAmount =  (this.enterProductForm.retailPrice * (1 - this.enterProductForm.discount /100) * (1 + this.enterProductForm.taxRate) * this.enterProductForm.qty).toFixed(2)
        } else if (event === '$') {
            this.enterProductForm.allTotalPrice = ((this.enterProductForm.retailPrice - this.enterProductForm.discount) * this.enterProductForm.qty).toFixed(2)
            this.enterProductForm.taxAmount = ((this.enterProductForm.retailPrice  - this.enterProductForm.discount) * (1 + this.enterProductForm.taxRate) * this.enterProductForm.qty).toFixed(2)
        }
    }

    addToList() {
        const { productId, productCode, qty, allTotalPrice, discount, discountType, taxAmount, taxRate, taxCode, taxType } = this.enterProductForm
        const addToList = {
            productId,
            productCode,
            qty,
            price: Number(allTotalPrice),
            discount: Number(discount),
            discountType,
            taxAmount: Number(taxAmount), 
            taxRate: Number(taxRate),
            taxCode, 
            taxType
        }

        this.inputedProductList.push(addToList)
        this.orgTotalAfterTax += Number(taxAmount)
        this.orgTotal += Number(allTotalPrice)

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
      orgTotalAfterTax: number = 0

    // RIGHT

    removeItem(index: number) {
        this.orgTotal -= Number(this.inputedProductList[index].price)
        this.orgTotalAfterTax -= Number(this.inputedProductList[index].taxAmount)
        this.inputedProductList.splice(index, 1)
        
        if (this.inputedProductList.length === 0) {
            this.orgTotal = 0
            this.orgTotalAfterTax = 0
        }
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
        this.invoiceForm.memberId = data.id
    }

    totalTypeChnage(event: any) {
        if (event === '%') {
            this.invoiceForm.totalAmount = parseFloat((this.orgTotal).toFixed(2))
            this.invoiceForm.taxTotal = parseFloat((this.orgTotalAfterTax * (1 - this.invoiceForm.discount /100)).toFixed(2))
        } else {
            this.invoiceForm.totalAmount = parseFloat((this.orgTotal).toFixed(2))
            this.invoiceForm.taxTotal = parseFloat((Number(this.orgTotalAfterTax) - Number(this.invoiceForm.discount)).toFixed(2))
        }
    }


    paymentForm: InvoicePaymentDto = {
        method: '',
        amount: 0,
        balance: '',
        findRedemption: ''
    }

    addPaymentList() {
        const { method, amount } = this.paymentForm
        const addToList = {
            method,
            amount,
            paymentTime: new Date().toISOString()
        }

        this.inputedPaymentList.push(addToList)
        
        this.paymentForm.findRedemption = (Number(amount) > this.invoiceForm.taxTotal ? Number(amount) - Number(this.invoiceForm.taxTotal) : 0).toFixed(2)
        this.paymentForm.balance = (this.invoiceForm.taxTotal - amount).toFixed(2)
    }

    lastPaymentDialogOpen: boolean = false

    openPaymentDialg() {
        this.invoiceForm.totalAmount = this.invoiceForm.totalAmount ? this.invoiceForm.totalAmount : this.orgTotalAfterTax
        this.paymentForm.balance = (this.invoiceForm.taxTotal ? this.invoiceForm.taxTotal : this.orgTotalAfterTax).toFixed(2)
        this.lastPaymentDialogOpen = true
    }

    closePaymentDialg() {
        this.lastPaymentDialogOpen = false
    }

    removePayItem(index: number, amount: any) {
        this.paymentForm.balance = this.paymentForm.balance + amount.amount
        this.paymentForm.findRedemption = (Number(this.paymentForm.findRedemption) - Number(amount.amount)).toFixed(2)
        this.inputedPaymentList.splice(index, 1)
    }

    async submitPOSale() {

        const finalDataSubmit = {
            ...this.invoiceForm,
            discount: this.invoiceForm.discount ? (this.invoiceForm.discountType === '%' ? this.invoiceForm.discount / 100 : this.invoiceForm.discount) : 0,
            invoiceItems: this.inputedProductList,
            invoicePayments: this.inputedPaymentList
        }

        const res = await postApiWithAuth('/invoice/create', finalDataSubmit)
        if (res.id) {
            this.message.success('Invoice created successfully!')
            timer(2500).subscribe(() => {
                this.routeTo.navigate(['/invoice-detail'], { queryParams: { id: res.id}})
            })
        } else if (res.error) {
            this.message.error(res.error)
        }
    }

    async findTaxInfoForSalePoint() {
        const res = await postApiWithAuth('/base/tax-information/sales-point', this.taxInfoFilter)
        if (res.error) {
            this.message.error(res.error)
        } else {
            this.enterProductForm.taxRate = res.taxRate
            this.enterProductForm.taxRateDisplay = res.taxRate * 100
            this.enterProductForm.taxCode = res.taxCode
            this.enterProductForm.taxType = res.taxType
        }
    }
}
