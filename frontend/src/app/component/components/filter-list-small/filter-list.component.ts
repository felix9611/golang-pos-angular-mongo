import {
  Component,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  Output,
  inject,
} from '@angular/core';
import { OnChanges, OnDestroy, OnInit, SimpleChanges, signal } from '@angular/core';
import { mapListResponse } from './filter-list-api-data';
import { FilterListPopover } from './filter-list-popover';

export interface ListOption {
  value: string;
  label: string;
}

@Component({
  selector: 'app-filter-list-small',
  imports: [FilterListPopover],
  standalone: true,
  template: `
    <div
      class="filter-list [&_button]:box-border [&_button]:[font:inherit] [&_button]:cursor-pointer [&_button:focus-visible]:[outline:2px_solid_#b09bd7] [&_button:focus-visible]:outline-offset-2 [&_button:disabled]:opacity-50 [&_button:disabled]:cursor-default"
    >
      <button
        type="button"
        #filterTrigger
        class="filter-trigger flex w-full items-center justify-between gap-[10px] [border:1px_solid_#e5dfec] rounded-[6px] px-[12px] py-[0.5rem] bg-white text-left text-[13px] text-[#756780] [&_>_span:first-child]:overflow-hidden [&_>_span:first-child]:text-ellipsis [&_>_span:first-child]:whitespace-nowrap [&[aria-expanded='true']]:border-[#b09bd7] [&[aria-expanded='true']]:[box-shadow:0_0_0_3px_#cccccc] max-[760px]:min-h-[40px] max-[760px]:text-[16px]"
        [attr.aria-label]="label"
        [attr.aria-expanded]="expanded"
        aria-haspopup="dialog"
        (click)="toggle()"
      >
        <span>{{ summary() }}</span
        ><span>⌄</span>
      </button>
      @if (selectedIds && selectedIds.length > 1 && multipleDisplay === true) {
        <div
          class="selected-options mt-[8px] flex flex-wrap gap-[6px] [&_button]:[border:1px_solid_#e8dff2] [&_button]:bg-[#f6f1fb] [&_button]:px-[9px] [&_button]:py-[6px] [&_button]:rounded-[5px] [&_button]:text-[12px] [&_button]:text-[#9575b1] max-[760px]:[&_button]:min-h-[44px] max-[760px]:[&_button]:text-[13px]"
          aria-live="polite"
          aria-label="Selected items"
        >
          @for (id of selectedIds; track id) {
            <button type="button" (click)="remove(id)" [attr.aria-label]="'Remove ' + name(id)">
              {{ name(id) }} ×
            </button>
          }
        </div>
      }
      @if (expanded && !disabled) {
        <app-filter-list-popover
          [label]="label"
          [anchor]="filterTrigger"
          (closed)="expanded = false"
        >
          <div class="filter-menu bg-white p-[10px]">
            <input
              class="box-border [font:inherit] block w-full m-[0_0_8px] [border:1px_solid_#e5dfec] rounded-[6px] p-[10px] bg-[#f5f5f5] text-[13px] text-[#53485f] focus-visible:[outline:2px_solid_#c2c2c2] focus-visible:outline-offset-2 max-[760px]:min-h-[46px] max-[760px]:text-[16px]"
              type="search"
              [name]="searchKeys.join(',')"
              [value]="query"
              (input)="onInput($event)"
              (keydown.enter)="submitSearch($event)"
              [attr.aria-label]="'Search ' + label"
              [placeholder]="'Search ' + label + '...'"
            />
            <div
              class="api-state text-[12px] text-[#887397] [&_button]:border-0 [&_button]:bg-[#f6f1fb] [&_button]:rounded-[5px] [&_button]:p-[8px] [&_button]:text-[#8664bc] max-[760px]:[&_button]:min-h-[44px] max-[760px]:[&_button]:text-[13px]"
              aria-live="polite"
            >
              @if (loading()) {
                <p>Loading...</p>
              }
              @if (error()) {
                <p role="alert">{{ error() }}</p>
                <button type="button" (click)="loadPage(failedPage)">Reload</button>
              }
            </div>
            <div
              class="filter-options max-h-[220px] overflow-y-auto overscroll-contain [&_>_button]:flex [&_>_button]:items-center [&_>_button]:justify-between [&_>_button]:w-full [&_>_button]:gap-[10px] [&_>_button]:border-0 [&_>_button]:rounded-[5px] [&_>_button]:p-[10px] [&_>_button]:bg-white [&_>_button]:text-left [&_>_button]:text-[#8f819e] [&_>_button:hover]:bg-[#cccccc] [&_>_button:hover]:text-[#858585] [&_>_button.is-selected]:bg-[#f7f7f7] [&_>_button.is-selected]:text-[#8664bc] [&_>_button_>_span:first-child]:min-w-0 [&_>_button_>_span:first-child]:[overflow-wrap:anywhere] [&_b]:block [&_b]:text-[13px] [&_b]:font-medium [&_small]:block [&_small]:text-[11px] [&_small]:mt-[4px] max-[760px]:[&_>_button]:min-h-[56px] max-[760px]:[&_b]:text-[14px] max-[760px]:[&_small]:text-[12px]"
              [attr.aria-label]="label + 'Options'"
            >
              @for (option of filtered(); track option.value) {
                <button
                  type="button"
                  [class.is-selected]="selectedIds.includes(option.value)"
                  [attr.aria-pressed]="selectedIds.includes(option.value)"
                  (click)="choose(option.value)"
                >
                  <span
                    ><b>{{ option.label }}</b
                    ><small>{{ option.value }}</small></span
                  ><span>{{ selectedIds.includes(option.value) ? '✓' : '＋' }}</span>
                </button>
              } @empty {
                @if (!loading() && !error()) {
                  <p class="filter-empty text-[12px] text-[#887397]">Not found any natched result</p>
                }
              }
            </div>
            @if (hasMore() && !error()) {
              <button
                type="button"
                class="load-more w-full min-h-[44px] mt-[8px] [border:1px_solid_#e5dfec] rounded-[6px] bg-[#f7f7f7] text-[#8664bc]"
                [disabled]="loading()"
                (click)="loadPage(page + 1)"
              >
                Loading More
              </button>
            }
            <div
              class="filter-menu-footer [border-top:1px_solid_#f0ebf5] mt-[6px] pt-[8px] flex items-center justify-between gap-[8px] [&_small]:text-[12px] [&_small]:text-[#887397] [&_button]:border-0 [&_button]:bg-[#f7f7f7] [&_button]:rounded-[5px] [&_button]:p-[8px] [&_button]:text-[#8664bc] max-[760px]:[&_button]:min-h-[44px] max-[760px]:[&_button]:text-[13px]"
            >
              <small>{{ filtered().length }} options</small
              ><button type="button" (click)="clear()">Clear All</button>
            </div>
          </div>
        </app-filter-list-popover>
      }
    </div>
  `,
})
export class FilterListApiSmall implements OnChanges, OnInit, OnDestroy {
  @Input() apiUrl = '/sys/user/list';
  @Input({ required: true }) searchDataKey: string | string[] = '';
  @Input() requestBody: Record<string, unknown> = { page: 1, limit: 10 };
  @Input() postApi?: (url: string, body: Record<string, unknown>) => Promise<unknown>;
  @Input() mapResponse: (response: unknown) => ListOption[] = mapListResponse;
  @Input() value: string[] = [];
  @Input() singleValue: string | null = null;
  @Input() label = '';
  @Input() multiple = false;
  @Input() multipleDisplay = false;
  @Input() disabled = false;
  @Output() valueChange = new EventEmitter<string[]>();
  @Output() singleValueChange = new EventEmitter<string>();

  loading = signal(false);
  error = signal('');
  hasMore = signal(false);
  options = signal<ListOption[]>([]);
  expanded = false;
  query = '';
  page = 0;
  failedPage = 1;
  private initialized = false;
  private destroyed = false;
  private version = 0;
  private searchTimer?: ReturnType<typeof setTimeout>;
  private timers = new Set<ReturnType<typeof setTimeout>>();
  private names = new Map<string, string>();
  private element = inject(ElementRef);

  ngOnInit() {
    this.initialized = true;
    this.readQuery();
    void this.reload();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (
      this.initialized &&
      ['apiUrl', 'requestBody', 'postApi', 'mapResponse', 'searchDataKey'].some(
        (key) => changes[key],
      )
    ) {
      this.readQuery();
      this.names.clear();
      void this.reload();
    }
  }

  get searchKeys(): string[] {
    const keys = Array.isArray(this.searchDataKey) ? this.searchDataKey : [this.searchDataKey];
    return [...new Set(keys.map((key) => key.trim()).filter(Boolean))];
  }

  private readQuery() {
    const values = this.searchKeys.map((key) => this.requestBody[key]);
    this.query =
      values.find((value): value is string => typeof value === 'string' && value.length > 0) ?? '';
  }

  toggle() {
    this.expanded = !this.expanded;
    if (this.expanded && this.page === 0 && !this.loading()) void this.reload();
  }

  onInput(event: Event) {
    this.search((event.target as HTMLInputElement).value);
  }

  search(query: string) {
    this.query = query;
    clearTimeout(this.searchTimer);
    ++this.version;
    this.options.set([]);
    this.page = 0;
    this.hasMore.set(false);
    this.error.set('');
    this.loading.set(true);
    this.searchTimer = setTimeout(() => void this.loadPage(1), 300);
  }

  submitSearch(event: Event) {
    event.preventDefault();
    event.stopPropagation();
    this.query = (event.target as HTMLInputElement).value;
    void this.reload();
  }

  async reload() {
    clearTimeout(this.searchTimer);
    this.options.set([]);
    this.page = 0;
    this.hasMore.set(false);
    await this.loadPage(1);
  }

  async loadPage(page: number) {
    clearTimeout(this.searchTimer);
    const version = ++this.version;
    this.failedPage = page;
    this.error.set('');
    if (typeof this.postApi !== 'function' || !this.apiUrl.trim()) {
      this.loading.set(false);
      this.error.set('Not configured API：Please bind [postApi]="loadUserLists" 及 apiUrl。');
      return;
    }
    const requestedLimit = Number(this.requestBody['limit']);
    const limit = Number.isInteger(requestedLimit) && requestedLimit > 0 ? requestedLimit : 10;
    if (!this.searchKeys.length) {
      this.loading.set(false);
      this.error.set('Please set searchDataKey in searchDataKey。');
      return;
    }
    const filters = Object.fromEntries(this.searchKeys.map((key) => [key, this.query.trim()]));
    const body = { ...this.requestBody, ...filters, page, limit };
    this.loading.set(true);
    let timer: ReturnType<typeof setTimeout> | undefined;
    try {
      // Invoke the parent's async callback directly. No HttpClient registration required.
      const request = this.postApi(this.apiUrl, body);
      const result = await Promise.race([
        request,
        new Promise<never>((_, reject) => {
          timer = setTimeout(() => reject(new Error('API timeout, please retry again.')), 15000);
          this.timers.add(timer);
        }),
      ]);
      if (version !== this.version || this.destroyed) return;
      if (result === undefined) throw new Error('Father Component must setyp return array options or API response');
      let rows: ListOption[];
      try {
        rows = this.mapResponse(result);
        if (
          !Array.isArray(rows) ||
          rows.some((row) => !row || typeof row.value !== 'string' || typeof row.label !== 'string')
        ) {
          throw new Error('Invalid options');
        }
      } catch {
        throw new Error('Format error：Please set return { value, label }[]，or set mapResponse。');
      }
      rows.forEach((row) => this.names.set(row.value, row.label));
      this.options.set(
        Array.from(
          new Map(
            [...(page === 1 ? [] : this.options()), ...rows].map((row) => [row.value, row]),
          ).values(),
        ),
      );
      this.page = page;
      this.hasMore.set(rows.length >= limit);
    } catch (error) {
      if (version === this.version && !this.destroyed) {
        this.error.set(error instanceof Error ? error.message : 'Options Loading failure');
      }
    } finally {
      if (timer) {
        clearTimeout(timer);
        this.timers.delete(timer);
      }
      if (version === this.version && !this.destroyed) this.loading.set(false);
    }
  }

  filtered() {
    return this.options();
  }
  name(id: string) {
    this.options();
    return this.names.get(id) || id;
  }
  get selectedIds(): string[] {
    if (this.multiple) return this.value;
    return this.singleValue !== null
      ? this.singleValue
        ? [this.singleValue]
        : []
      : this.value.slice(0, 1);
  }
  summary() {
    return this.selectedIds.length
      ? this.selectedIds.map((id) => this.name(id)).join('、')
      : 'Please select ' + this.label;
  }
  choose(id: string) {
    if (this.multiple) {
      this.valueChange.emit(
        this.value.includes(id) ? this.value.filter((value) => value !== id) : [...this.value, id],
      );
    } else {
      this.singleValueChange.emit(id);
      this.valueChange.emit([id]);
      this.expanded = false;
    }
  }
  clear() {
    if (!this.multiple) this.singleValueChange.emit('');
    this.valueChange.emit([]);
    this.expanded = false;
  }
  remove(id: string) {
    if (!this.multiple) {
      this.clear();
      return;
    }
    this.valueChange.emit(this.value.filter((value) => value !== id));
  }

  @HostListener('document:click', ['$event']) outside(event: MouseEvent) {
    if (!this.element.nativeElement.contains(event.target)) this.expanded = false;
  }

  ngOnDestroy() {
    this.destroyed = true;
    ++this.version;
    clearTimeout(this.searchTimer);
    this.timers.forEach((timer) => clearTimeout(timer));
    this.timers.clear();
  }

  selectedItem(id: string) {
    if (this.selectedIds.includes(id)) {
      return '✓'
    } else {
      if (this.multiple === true) {
        return '＋'
      } else {
        return ''
      }
    }
  }
}