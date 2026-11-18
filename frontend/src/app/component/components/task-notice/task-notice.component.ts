import { Component, ElementRef, HostListener, ViewChild } from '@angular/core'
import { getApiWithAuth } from '../../../../tool/httpRequest-auth';
import moment from 'moment';

@Component({
    selector: 'app-tasks-notice',
    standalone: true,
    styleUrls: [],
    template: `
         <button
            class="relative flex w-[2.7rem] h-[2.7rem] cursor-pointer items-center justify-center rounded-full border-0 bg-[#f0f2f5] text-[#353945] transition-colors hover:bg-[#e4e6eb] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-blue-500"
            mat-icon-button
            [attr.aria-expanded]="expanded"
            aria-haspopup="dialog"
            (click)="expanded = !expanded"
        >
            <svg
                aria-hidden="true"
                width="22"
                height="22"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9Z" />
                <path d="M10 21h4" />
            </svg>
            @if (tasksNumber > 0) {
                <span
                    aria-hidden="true"
                    class="absolute -right-[3px] -top-[3px] flex h-[21px] min-w-[21px] items-center justify-center rounded-full border-2 border-solid border-white bg-[#e41e3f] px-[4px] text-[11px] font-bold leading-none text-white"
                >{{ tasksNumber > 9 ? '99+' : tasksNumber }}</span >
            }
        </button>
        @if (expanded) {
            <section
                role="dialog"
                class="absolute top-[4rem] z-[41] max-[780px]:right-[6rem] max-[1024px]:right-[16rem] right-[16rem] max-[490px]:w-[250px] max-[640px]:w-[320px] w-[380px] overflow-hidden rounded-[12px] border border-solid border-[#e4e6eb] bg-white text-left text-[#1c1e21] shadow-[0_8px_32px_#00000026]"
            >
                <header class="flex items-center justify-between px-[18px] pt-[17px]">
                    <h2 class="m-0 text-[1.3rem] font-extrabold leading-tight">Notify</h2>
                </header>
                <div class="px-[18px] pb-[12px]">
                    <span
                        class="inline-flex rounded-full bg-[#e7f3ff] px-[14px] py-[2px] text-[14px] font-semibold text-[#1877f2]"
                        >All</span
                    >
                    <p class="pt-[0.2rem] mb-0 mt-[1rem] text-[0.85rem] text-[#65676b]" role="status">
                        Total of {{ tasksNumber > 9 ? '10+' : tasksNumber }} assigned tasks need to be followed up
                    </p>
                </div>
                <div
                    class="max-h-[min(540px,calc(100dvh-225px))] overflow-y-auto overscroll-contain px-[8px] pb-[10px]"
                >
                    @for (job of jobs; track job._id) {
                        <div class="mb-[2px] flex w-full cursor-pointer flex-col gap-[8px] rounded-[8px] border-0 bg-white px-[10px] py-[13px] text-left hover:bg-[#f0f2f5] focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-blue-500">
                            <span class="flex w-full items-center justify-between gap-[12px]">
                                <time
                                class="text-[12px] leading-normal text-[#65676b]"
                                [attr.datetime]="validDate(job.createdAt)"
                                    >{{
                                        validDate(job.createdAt)
                                        ? validDate(job.createdAt)
                                        : 'No date'
                                    }}
                                </time>
                                <span
                                    class="shrink-0 rounded-full px-[9px] py-[4px] text-[11px] font-semibold"
                                    [class]="statusNameClass(job.status)">
                                    {{ statusName(job.status) }}
                                </span>
                            </span>
                            <span class="text-[1rem] font-semibold leading-[1.6] [overflow-wrap:anywhere]">
                                {{ job.title }}
                            </span>
                        </div>
                    }
                </div>
            </section>
        }
    `
})
export class TasksNoticeComponent {
    @ViewChild('trigger') trigger?: ElementRef<HTMLButtonElement>

    expanded: boolean = false
    tasksNumber: number = 0
    jobs: any[] = []
    private reloadTimer?: ReturnType<typeof setInterval>

    close() {
        this.expanded = false;
        this.trigger?.nativeElement.focus();
    }

    @HostListener('document:keydown.escape') onEscape() {
        if (this.expanded) this.close();
    }

    ngOnInit() {
        this.reloadTimer = setInterval(() => {
        void this.loadSyncJobLists()
        }, 10_000);
    }

    ngOnDestroy() {
        clearInterval(this.reloadTimer)
    }

    async loadSyncJobLists(): Promise<void> {
        const data = await getApiWithAuth('/tasks/sync')
        this.jobs = data
        this.tasksNumber = data.length
    }

    validDate(data?: string) {
        return data ? moment(data).format('DD-MM-YYYY HH:mm') : null
    }

    statusName(status: number) {
        return ({ 
            1: 'To do', 
            2: 'In progress', 
            3: 'Pending Review', 
            4: 'In review',
            5: 'Done' 
        } as Record<number, string>)[status] || '';
    }

    statusNameClass(status: number) {
        return ({ 
            1: 'bg-[#e3e3e3] text-[#8b91a1]', 
            2: 'bg-[#d0c9ff] text-[#6c5ce7]', 
            3: 'bg-[#ffe4bd] text-[#e5a345]', 
            4: 'bg-[#ffcdb0] text-[#f57124]' 
        } as Record<number, string>)[status] || '';
    }
}