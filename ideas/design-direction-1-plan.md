### Design Direction 1: Dense Grid View - Implementation Plan

**Tagline**: "Compact and scannable issue grids with smart filters for rapid task triage."

**Rationale**: Addresses low density and filtering gaps from v0.1 and pico versions; ties to Fossil's sortable columns and color-coded status. Heuristic focus: visibility (status colors), efficiency (dense cards over tables), consistency (Pico alignment).

**Visual System**:
- **Color Palette**: Primary (#0066cc, #ffffff for light; #4d9de0, #1a1a1a for dark) - WCAG AA compliant. Status: Open (#d4edda, #155724), In-Progress (#fff3cd, #856404), Closed (#f8d7da, #721c24). Use CSS vars: `--status-open: hsl(120, 25%, 90%)` etc.
- **Tone**: Professional, neutral; elevation via shadows for cards.
- **Elevation Strategy**: Subtle box-shadow on cards/hover.
- **Density**: High (compact cards, reduced spacing).

**Typography**: Pico defaults (sans-serif, headings h1-h6); semantic use for hierarchy.

**Spacing/Sizing Scales**: Pico's block-spacing (1rem base); custom var `--dense-spacing: 0.5rem` for tighter grids.

**Iconography**: Inline SVG for actions (e.g., search icon); system-safe fonts.

**Information Architecture & Navigation**:
- Global nav: Header with logo, search bar, filters (status/priority), "Ready" tab.
- Secondary: Breadcrumbs on detail pages; tabs for deps/blockers.
- Search: Global input with filters.

**Responsive Behavior**: Mobile-first; stack cards vertically; tablet/desktop grid (3-4 cols); dark mode via `prefers-color-scheme`.

**Accessibility**: Focus outlines (Pico default); 44px hit targets; keyboard nav (Tab/Enter); reduced motion (prefers-reduced-motion); form validation (Pico roles).

**Microinteractions**: CSS transitions (0.2s ease) on hover/focus; fade for loading.

**Component Mapping**: Buttons (Pico .btn), links (a), forms (form), inputs (input), tables → cards (.card), tags (.label), breadcrumbs (nav > ul), pagination (if needed), alerts (article), modals (dialog), empty states (custom .empty), loading states (.skeleton via custom var).

**Page Layouts** (with HTML prototypes; add Pico via CDN `<link href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">`):

1. **Dashboard/Index**: Stats grid, filter bar, card grid of issues.
   ```html
   <!DOCTYPE html>
   <html lang="en" data-theme="auto">
   <head><meta charset="UTF-8"><title>Dashboard - Beady</title></head>
   <body>
   <main class="container">
     <header><h1>Beady Dashboard</h1><nav><a href="/ready">Ready Work</a></nav></header>
     <div class="grid">
       <article class="card"><h3>Total: 42</h3></article><!-- Repeat for stats -->
     </div>
     <form><input type="search" placeholder="Search issues"><select name="status"><option>Open</option></select><button>Filter</button></form>
     <div class="grid"><!-- Cards for issues --></div>
   </main>
   </body>
   </html>
   ```

2. **Ready Work**: Filtered card list.
   - Similar to above, with exclude input.

3. **Issue Detail**: Card layout with sections.
   ```html
   <!DOCTYPE html>
   <html lang="en" data-theme="auto">
   <head><meta charset="UTF-8"><title>bd-123 - Beady</title></head>
   <body>
   <main class="container">
     <nav aria-label="breadcrumb"><ol><li><a href="/">Home</a></li><li>bd-123</li></ol></nav>
     <article class="card"><h2>Fix login bug</h2><p>Status: Open</p><!-- Details --></article>
     <section><h3>Dependencies</h3><ul><!-- Links --></ul></section>
   </main>
   </body>
   </html>
   ```

4. **Settings**: Form for filters/preferences.
   - Form with inputs for defaults.

**Implementation Plan**: 
- Refactor tables to cards first; add CSS vars for colors; enable dark mode; implement search/sort via JS.
- Rollout: Core pages → add features; mitigate risks with fallback styles; QA: Contrast check, mobile test, form validation.

**Success Metrics**: Reduce task time for triage by 30% (dense grids); lower error rate in filters via validation; improve discoverability (user tests); satisfaction via surveys.
