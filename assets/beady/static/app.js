// Basic JavaScript for Beads UI

console.log('Beads UI loaded');

// Theme functionality
function applyTheme(theme) {
    const html = document.documentElement;
    if (theme === 'auto') {
        html.removeAttribute('data-theme');
    } else {
        html.setAttribute('data-theme', theme);
    }
}

function initTheme() {
    const themeSelect = document.getElementById('theme-select');
    const savedTheme = (typeof localStorage !== 'undefined') ? localStorage.getItem('beady.theme') : 'auto';
    
    // Set initial theme
    applyTheme(savedTheme);
    if (themeSelect) {
        themeSelect.value = savedTheme;
    }
    
    // Listen for theme changes
    if (themeSelect) {
        themeSelect.addEventListener('change', function() {
            const theme = this.value;
            applyTheme(theme);
            try {
                localStorage.setItem('beady.theme', theme);
            } catch (e) {
                // ignore storage errors (e.g., privacy modes)
            }
        });
    }
}

// View selector functionality
document.addEventListener('DOMContentLoaded', function() {
    // Initialize theme
    initTheme();
    
    const viewSelect = document.getElementById('view-select');
    const views = {
        grid: document.getElementById('grid-view'),
        kanban: document.getElementById('kanban-view'),
        timeline: document.getElementById('timeline-view')
    };

    function switchView(view) {
        Object.keys(views).forEach(key => {
            if (views[key]) {
                views[key].style.display = key === view ? 'block' : 'none';
            }
        });
    }

    // Restore selection from localStorage if available, otherwise use the current select value or 'timeline'
    const saved = (typeof localStorage !== 'undefined') ? localStorage.getItem('beady.view') : null;
    const initial = saved && views[saved] ? saved : (viewSelect ? viewSelect.value : 'timeline');
    switchView(initial);
    if (viewSelect && saved && views[saved]) {
        viewSelect.value = saved;
    }

    if (viewSelect) {
        viewSelect.addEventListener('change', function() {
            const v = this.value;
            switchView(v);
            try {
                localStorage.setItem('beady.view', v);
            } catch (e) {
                // ignore storage errors (e.g., privacy modes)
            }
        });
    }
});
