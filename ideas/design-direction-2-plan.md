### Design Direction 2: Kanban Workflow - Implementation Plan

**Tagline**: "Agile board lanes for status-driven issue management and quick drag-insights."

**Rationale**: Enhances status visibility and priority sorting from ideas (v0 likes filters, Fossil's color-coding); heuristic evaluation: natural and intuitive (lanes match workflow stages), flexibility (easy reordering hints), consistency (Pico grid); addresses density by grouping related issues, reducing cognitive load vs. flat tables.

**Visual System**:
- **Color Palette**: Primary (#007bff, #ffffff light; #0056b3, #121212 dark) - WCAG AA. Lanes: Open (#d1ecf1, #0c5460), In-Progress (#fff3cd, #856404), Closed (#d4edda, #155724). CSS vars: `--lane-open: hsl(180, 25%, 90%)` etc.
- **Tone**: Dynamic, workflow-oriented; elevation for cards in lanes.
- **Elevation Strategy**: Box-shadow on cards; hover lift.
- **Density**: Medium (lane columns, card spacing).

**Typography**: Pico defaults; bold for lane headers.

**Spacing/Sizing Scales**: Pico base; custom `--lane-gap: 1rem` for columns.

**Iconography**: Inline SVG for status icons (e.g., circle for open); no libs.

**Information Architecture & Navigation**:
- Global: Header with search, "Kanban" view toggle.
- Secondary: Lane tabs; detail nav for deps.
- Search/Filter: Per-lane or global.

**Responsive Behavior**: Mobile: Single column stack; tablet: 2 lanes; desktop: 3+ lanes; dark mode via prefers-color-scheme.

**Accessibility**: Pico focus; 44px targets; keyboard (arrow keys for lanes); reduced motion; validation patterns.

**Microinteractions**: Transitions on card move hints (simulate drag); fade for updates.

**Component Mapping**: Buttons (.btn), links, forms, inputs, tables → kanban (.grid > .card), tags (.label), breadcrumbs, tabs (.tab), pagination, alerts, modals, empty (lane message), loading (skeleton), error (alert).

**Page Layouts**:

1. **Dashboard/Index**: Lane columns with stats header.
   ```html
   <!DOCTYPE html>
   <html lang="en" data-theme="auto">
   <head><meta charset="UTF-8"><title>Kanban - Beady</title></head>
   <body>
   <main class="container">
     <header><h1>Issue Kanban</h1><nav><a href="/list">List View</a></nav></header>
     <div class="grid"><!-- Lanes: Open, In-Progress, Closed -->
       <section class="card"><h3>Open</h3><div><!-- Cards --></div></section>
     </div>
   </main>
   </body>
   </html>
   ```

2. **Ready Work**: Filtered lanes.
   - Kanban with exclude filter form.

3. **Issue Detail**: Lanes for deps/blockers.
   ```html
   <!DOCTYPE html>
   <html lang="en" data-theme="auto">
   <head><meta charset="UTF-8"><title>bd-123 - Beady</title></head>
   <body>
   <main class="container">
     <nav aria-label="breadcrumb"><ol><li><a href="/">Home</a></li><li>bd-123</li></ol></nav>
     <article class="card"><h2>Implement feature</h2><!-- Details --></article>
     <section><h3>Dependencies</h3><div class="grid"><!-- Mini lanes --></div></section>
   </main>
   </body>
   </html>
   ```

4. **Settings**: Board config form.
   - Form for lane visibility.

**Implementation Plan**: Convert tables to kanban grid; add lane CSS vars; enable drag hints via JS; rollout: Lanes first, then features; risks: JS dependency (fallback to list); QA: Lane navigation, contrast.

**Success Metrics**: Faster status changes (20% time reduction); error rate down via visual cues; action discoverability up (user flows); satisfaction from workflow alignment.
