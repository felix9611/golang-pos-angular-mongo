import { CommonModule, isPlatformBrowser, isPlatformServer } from '@angular/common'
import { Component, ComponentFactoryResolver, EventEmitter, HostListener, Inject, Input, OnInit, Output, PLATFORM_ID, ViewChild, ViewContainerRef } from '@angular/core'
import { CanvasJSAngularChartsModule } from '@canvasjs/angular-charts'

@Component({
    selector: 'app-chart',
    standalone: true,
    imports: [CommonModule, CanvasJSAngularChartsModule],
    templateUrl: './chart.component.html',
    styleUrls: ['./chart.component.scss'],
})
export class CanvasChartComponent implements OnInit {

    @Input() chartData: any = {}
    @Output() chartDataChange = new EventEmitter<any>()

    constructor() {
    }
    ngOnInit(): void {
        this.chartOptions = this.chartData
    }

    chartOptions = {
        legend: {
            horizontalAlign: 'top'
        },
        axisY: {
            title: 'Amount'
        },
        data: [{
            type: "stackedColumn",
            name: "Facebook",
            showInLegend: true,
            dataPoints: [
              { label: "Qtr 1", y: 19729 },
              { label: "Qtr 2", y: 22127 },
              { label: "Qtr 3", y: 12654 },
              { label: "Qtr 4", y: 22914 }
            ]
          }, {
              type: "stackedColumn",
              name: "Twitter",
              showInLegend: true,
              dataPoints: [
                { label: "Qtr 1", y: 4288 },
                { label: "Qtr 2", y: 6390 },
                { label: "Qtr 3", y: 3510 },
                { label: "Qtr 4", y: 3876 }
              ]
          }, {
              type: "stackedColumn",
               name: "Instagram",
              showInLegend: true,
              dataPoints: [
                { label: "Qtr 1", y: 5338 },
                { label: "Qtr 2", y: 8670 },
                { label: "Qtr 3", y: 4779 },
                { label: "Qtr 4", y: 9415 }
              ]
          }]  
    }

    updateChartData(newData: any) {
        this.chartDataChange.emit(newData);
    }

    
}
