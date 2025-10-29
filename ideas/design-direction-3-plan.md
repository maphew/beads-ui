### Design Direction 3: Timeline Insights - Implementation Plan

**Tagline**: "Event-driven timelines and dependency graphs for issue history and insights."

**Rationale**: Inspired by Fossil's timeline/graph for commit history and ticket events; addresses missing event views and graphs from current UI; heuristic: information scent (timelines aid context), recognition (graphs visualize deps), flexibility (sortable events); improves density via compact timelines over tables.

**Visual System**:
- **Color Palette**: Primary (#28a745, #ffffff light; #1e7e34, #1a1a1a dark) - WCAG AA. Events: Created (#cce5ff, #004085), Updated (#f8f9fa, #6c757d). CSS vars: `--event-created: hsl(210, 100%, 95%)` etc.
- **Tone**: Analytical, data-focused; subtle elevation for graph areas.
- **Elevation Strategy**: Border/shadow on timeline items.
- **Density**: Medium-high (timeline lists, graph embeds).

**Typography**: Pico defaults; code font for events.

**Spacing/Sizing Scales**: Pico; custom `--timeline-margin: 0.75rem`.

**Iconography**: SVG for event types (clock, arrow); system-safe.

**Information Architecture & Navigation**:
- Global: Header with "Timeline" tab, search.
- Secondary: Graph view toggle; event filters.
- Search: Timeline query.

**Responsive Behavior**: Mobile: Vertical timeline; tablet/desktop: Horizontal graph; dark mode.

**Accessibility**: Focus on events; keyboard scroll; reduced motion; form validation.

**Microinteractions**: Animate timeline scroll; fade for graph loads.

**Component Mapping**: Buttons, links, forms, inputs, tables → timelines (ul.timeline), graphs (canvas or SVG), tags, breadcrumbs, tabs, pagination, alerts, modals, empty (no events message), loading, error.

**Page Layouts**:

1. **Dashboard/Index**: Timeline with stats.
   ```html
   <!DOCTYPE html>
   <html lang="en" data-theme="auto">
   <head><meta charset="UTF-8"><title>Timeline - Beady</title></head>
   <body>
   <main class="container">
     <header><h1>Issue Timeline</h1></header>
     <ul class="timeline"><!-- Event items --></ul>
   </main>
   </body>
   </html>
   ```

2. **Ready Work**: Timeline filtered.
   - Add filter to timeline.

3. **Issue Detail**: Graph and event list.
   ```html
   <!DOCTYPE html>
   <html lang="en" data-theme="auto">
   <head><meta charset="UTF-8"><title>bd-123 - Beady</title></head>
   <body>
   <main class="container">
     <nav aria-label="breadcrumb"><ol><li><a href="/">Home</a></li><li>bd-123</li></ol></nav>
     <article class="card"><h2>Bug fix</h2><!-- Details --></article>
     <section><h3>Dependency Graph</h3><svg><!-- Inline graph --></svg></section>
     <ul class="timeline"><!-- Events --></ul>
   </main>
   </body>
   </html>
   ```

4. **Settings**: Graph config.
   - Form for display options.

**Implementation Plan**: Add timeline component; integrate graphs; CSS vars for events; rollout: Timeline → graphs; risks: Graph rendering (fallback); QA: Event accessibility, graph clarity.

**Success Metrics**: Improved dependency understanding (fewer blockers); task time for reviews down 25%; error rate in edits; discoverability via timelines; satisfaction from insights.
