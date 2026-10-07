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
import { DashboardReqDto, DashboardReqFilterDto } from './interface'
import { CanvasChartComponent } from '../../../components/chart/chart.component'
import { transformData, transformDataNoDate, transformDataPointsOnly, transformDate } from './function' 
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { UserStoreService } from '../../../../../state/user.service'
import { findMenuItem } from '../../../tool-function'
import { Subscription } from 'rxjs'

@Component({
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
        CanvasChartComponent,
        NzSelectModule,
        NzDatePickerModule,
    ],
    templateUrl: './stock-dashboard.component.html',
    // styleUrl: './invoice-dashboard.component.css',
})
export class StockDashboardComponent implements OnInit {
    private rightSubscription: Subscription
    constructor(
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {
            this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
                const answer = findMenuItem(data, 'Dashboard', 'dashboard')
                this.userRightInside = {
                    read: answer?.read ?? false
                     // keep default value
                }
            })
                    
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }


    userRightInside = {
        read: false
    }

    ngOnInit(): void {
        this.runSearch()
   //     this.getByTotalCost()
  //      this.getByTotalCount()
        this.loadTypeList()
        this.loadDeptList()
        this.loadLocationList()
        
    }

    typeLists: any[] = []
    async loadTypeList() {
        this.typeLists = await getApiWithAuth('/product/product-type/getAll')
    }

    deptLists: any[] = []
    async loadDeptList() {
        this.deptLists = await getApiWithAuth('/sys/department/getAll')
    }

    placeLists: any[] = []
    async loadLocationList() {
        this.placeLists = await getApiWithAuth('/base/location/getAll')
    }

    async runSearch() {
        await this.getByDeptAndQtys()
        await this.getByDeptAndPrice()
        await this.getByTypeAndPrice()
        await this.getByTypeAndQtys()
        await this.getByLocationAndQtys() 
        await this.getByLocationAndPrice()

        await this.getInvintoryRecordQtys()
        await this.getInvintoryRecordPrices()
    }

    globalFilter: DashboardReqFilterDto = {}

    byTypeAndQtysLoading: boolean = false
    byTypeAndQtys: any = {}
    async getByTypeAndQtys() {
        this.byTypeAndQtysLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'type'
        }

        const res = await this.runQueryData(dataQuery)

        const qtys = transformDataNoDate(res, 'column', false, 'typeName', 'totalCost')
        const qtysDataSet = qtys[0].dataPoints

        this.byTypeAndQtys = {
            data: [
                {
                    name: 'Total Qtys',
                    dataPoints: qtysDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Qtys'
            },
            axisX: {
                title: 'Product Types'
            },
            toolTip: {
                shared: true
            }
        }
        this.byTypeAndQtysLoading = true
    }

    byTypeAndPriceLoading: boolean = false
    byTypeAndPrice: any = {}
    async getByTypeAndPrice() {
        this.byTypeAndPriceLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'type'
        }

        const res = await this.runQueryData(dataQuery)

        const cost = transformDataNoDate(res, 'column', false, 'typeName', 'totalCost')
        const costDataSet = cost[0].dataPoints

        const retail = transformDataNoDate(res, 'column', false, 'typeName', 'totalPrice')
        const retailDataSet = retail[0].dataPoints

        this.byTypeAndPrice = {
            data: [
                {
                    name: 'Total Cost Price',
                    dataPoints: costDataSet
                },
                {
                    name: 'Total Retail Price',
                    dataPoints: retailDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Product Types',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byTypeAndPriceLoading = true
    }

    byDeptAndQtysLoading: boolean = false
    byDeptAndQtys: any = {}
    async getByDeptAndQtys() {
        this.byDeptAndQtysLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'dept'
        }

        const res = await this.runQueryData(dataQuery)
        const qtys = transformDataNoDate(res, 'column', false, 'deptName', 'totalQty')
        const qtysDataSet = qtys[0].dataPoints

        this.byDeptAndQtys = {
            data: [
                {
                    name: 'Total Qtys',
                    dataPoints: qtysDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Qtys'
            },
            axisX: {
                title: 'Departments',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byDeptAndQtysLoading = true
    }

    byDeptAndPriceLoading: boolean = false
    byDeptAndPrice: any = {}
    async getByDeptAndPrice() {
        this.byDeptAndPriceLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'dept'
        }

        const res = await this.runQueryData(dataQuery)
        const cost = transformDataNoDate(res, 'column', false, 'deptName', 'totalCost')
        const costDataSet = cost[0].dataPoints

        const retail = transformDataNoDate(res, 'column', false, 'deptName', 'totalPrice')
        const retailDataSet = retail[0].dataPoints

        this.byDeptAndPrice = {
            data: [
                {
                    name: 'Total Cost Price',
                    dataPoints: costDataSet
                },
                {
                    name: 'Total Retail Price',
                    dataPoints: retailDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Departments',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byDeptAndPriceLoading = true
    }

    byLocationAndQtysLoading: boolean = false
    byLocationAndQtys: any = {}
    async getByLocationAndQtys() {
        this.byLocationAndQtysLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'location'
        }

        const res = await this.runQueryData(dataQuery)
        const qtys = transformDataNoDate(res, 'column', false, 'placeName', 'totalQty')
        const qtysDataSet = qtys[0].dataPoints

        this.byLocationAndQtys = {
            data: [
                {
                    name: 'Total Qtys',
                    dataPoints: qtysDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Qtys'
            },
            axisX: {
                title: 'Locations',
            },
            toolTip: {
                shared: true
            }
        }
        this.byLocationAndQtysLoading = true
    }

    byLocationAndPriceLoading: boolean = false
    byLocationAndPrice: any = {}
    async getByLocationAndPrice() {
        this.byLocationAndPriceLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'location'
        }

        const res = await this.runQueryData(dataQuery)
        const cost = transformDataNoDate(res, 'column', false, 'placeName', 'totalCost')
        const costDataSet = cost[0].dataPoints

        const retail = transformDataNoDate(res, 'column', false, 'placeName', 'totalPrice')
        const retailDataSet = retail[0].dataPoints

        this.byLocationAndPrice = {
            data: [
                {
                    name: 'Total Cost Price',
                    dataPoints: costDataSet
                },
                {
                    name: 'Total Retail Price',
                    dataPoints: retailDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Locations',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byLocationAndPriceLoading = true
    }

    invintoryRecordQtysLoading: boolean = false
    invintoryRecordQtys: any = {}
    async getInvintoryRecordQtys() {
        this.invintoryRecordQtysLoading = false
        const resStockIn = await this.runInvintoryRecordQueryData({ dataTypeValue: 'stockIn' })
        const resStocMove = await this.runInvintoryRecordQueryData({ dataTypeValue: 'stockMove' })
        const resStocOut = await this.runInvintoryRecordQueryData({ dataTypeValue: 'stockOut' })

        const stockIns = transformDate(resStockIn, 'qtys', ['year', 'month'])
        const stockOuts = transformDate(resStocOut, 'qtys', ['year', 'month'])
        const stockMoves = transformDate(resStocMove, 'qtys', ['year', 'month'])

        this.invintoryRecordQtys = {
            data: [
                {
                    name: 'Stock In',
                    type: 'spline',
                    xValueType: 'dateTime',
                    dataPoints: stockIns
                },
                {
                    name: 'Stock Out',
                    type: 'spline',
                    xValueType: 'dateTime',
                    dataPoints: stockOuts
                },
                {
                    name: 'Stock Move',
                    type: 'spline',
                    xValueType: 'dateTime',
                    dataPoints: stockMoves
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Date - Month',
                valueFormatString: "YYYY - MMM",
                xValueType: "dateTime"
            },
            toolTip: {
                shared: true
            }
        }

        this.invintoryRecordQtysLoading = true
    }

    invintoryRecordPricesLoading: boolean = false
    invintoryRecordPrices: any = {}
    async getInvintoryRecordPrices() {
        this.invintoryRecordPricesLoading = false
        const resStockIn = await this.runInvintoryRecordQueryData({ dataTypeValue: 'stockIn' })
        const resStocMove = await this.runInvintoryRecordQueryData({ dataTypeValue: 'stockMove' })
        const resStocOut = await this.runInvintoryRecordQueryData({ dataTypeValue: 'stockOut' })

        const stockIns = transformDate(resStockIn, 'costs', ['year', 'month'])
        const stockOuts = transformDate(resStocOut, 'costs', ['year', 'month'])
        const stockMoves = transformDate(resStocMove, 'costs', ['year', 'month'])

        this.invintoryRecordPrices = {
            data: [
                {
                    name: 'Stock In',
                    type: 'spline',
                    xValueType: 'dateTime',
                    dataPoints: stockIns
                },
                {
                    name: 'Stock Out',
                    type: 'spline',
                    xValueType: 'dateTime',
                    dataPoints: stockOuts
                },
                {
                    name: 'Stock Move',
                    type: 'spline',
                    xValueType: 'dateTime',
                    dataPoints: stockMoves
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Date - Month',
                valueFormatString: "YYYY - MMM",
                xValueType: "dateTime"
            },
            toolTip: {
                shared: true
            }
        }

        this.invintoryRecordPricesLoading = true
    }

    async runQueryData(dataQuery: DashboardReqDto) {
        const finalQuery = {
            ...dataQuery,
            filter: this.globalFilter
        }

        return await postApiWithAuth('/product/product-list/query/data-group-by', finalQuery)
    }

    async runInvintoryRecordQueryData(dataQuery: DashboardReqDto) {
        const finalQuery = {
            dataType: dataQuery.dataTypeValue,
            filter: {
                ...this.globalFilter.placeIds && this.globalFilter.placeIds.length > 0 ? { placeIds: this.globalFilter.placeIds } : {}
            }
        }

        return await postApiWithAuth('/inventory-record/query', finalQuery)
    }

}