// Basic JavaScript for Beads UI

console.log('Beads UI loaded');

// View selector functionality
document.addEventListener('DOMContentLoaded', function() {
    const viewSelect = document.getElementById('view-select');
    const views = {
        grid: document.getElementById('grid-view'),
        kanban: document.getElementById('kanban-view'),
        timeline: document.getElementById('timeline-view')
    };

    function switchView(view) {
        Object.keys(views).forEach(key => {
            views[key].style.display = key === view ? 'block' : 'none';
        });
    }

    if (viewSelect) {
        viewSelect.addEventListener('change', function() {
            switchView(this.value);
        });
    }
});
