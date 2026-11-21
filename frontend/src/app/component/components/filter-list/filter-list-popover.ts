import {
  AfterViewInit,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnDestroy,
  Output,
  ViewChild,
} from '@angular/core';

@Component({
  selector: 'app-filter-list-popover',
  standalone: true,
  template: `
    <div
      #panel
      popover="auto"
      role="dialog"
      [attr.aria-label]="label || 'Select'"
      (toggle)="onToggle($event)"
      (keydown.escape)="$event.stopPropagation()"
    >
      <ng-content />
    </div>
  `,
  styles: `
    [popover] {
      box-sizing: border-box;
      position: fixed;
      inset: auto;
      margin: 0;
      padding: 0;
      border: 1px solid #e5dfec;
      border-radius: 14px;
      background: #fff;
      color: #756780;
      box-shadow: 0 4px 18px #21143226;
      overflow: auto;
      overscroll-behavior: contain;
      font: inherit;
    }
    [popover]::backdrop {
      background: transparent;
      pointer-events: none;
    }
  `,
})
export class FilterListPopover implements AfterViewInit, OnDestroy {
  @Input() label = '';
  @Input({ required: true }) anchor!: HTMLElement;
  @Output() closed = new EventEmitter<void>();
  @ViewChild('panel', { static: true }) panel!: ElementRef<HTMLElement>;
  private observer?: ResizeObserver;

  private position = () => {
    const panel = this.panel.nativeElement;
    const anchor = this.anchor.getBoundingClientRect();
    const gap = 6;
    const margin = 12;
    const viewportHeight = window.innerHeight;
    const below = Math.max(0, viewportHeight - anchor.bottom - gap - margin);
    const above = Math.max(0, anchor.top - gap - margin);
    const openAbove = below < Math.min(panel.scrollHeight, 300) && above > below;
    panel.style.width = `${Math.min(Math.max(anchor.width, 280), window.innerWidth - margin * 2)}px`;
    panel.style.maxHeight = `${openAbove ? above : below}px`;
    panel.style.left = `${Math.max(margin, Math.min(anchor.left, window.innerWidth - panel.offsetWidth - margin))}px`;
    panel.style.top = `${openAbove ? Math.max(margin, anchor.top - gap - panel.offsetHeight) : anchor.bottom + gap}px`;
  };

  ngAfterViewInit() {
    const panel = this.panel.nativeElement;
    panel.showPopover();
    this.position();
    panel.querySelector<HTMLInputElement>('input')?.focus({ preventScroll: true });
    window.addEventListener('resize', this.position);
    window.addEventListener('scroll', this.position, true);
    this.observer = new ResizeObserver(this.position);
    this.observer.observe(panel);
    this.observer.observe(this.anchor);
  }

  onToggle(event: Event) {
    if ((event as ToggleEvent).newState === 'closed') {
      this.closed.emit();
    }
  }

  ngOnDestroy() {
    const panel = this.panel.nativeElement;
    if (panel.contains(document.activeElement)) this.anchor.focus({ preventScroll: true });
    this.observer?.disconnect();
    window.removeEventListener('resize', this.position);
    window.removeEventListener('scroll', this.position, true);
  }
}
