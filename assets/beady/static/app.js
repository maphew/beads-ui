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
    const savedTheme = (typeof localStorage !== 'undefined') ? (localStorage.getItem('beady.theme') || 'auto') : 'auto';
    
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

// Immediate filter functionality
let filterTimeout = null;

function applyFilters() {
    const form = document.getElementById('filter-form');
    if (!form) return;

    const formData = new FormData(form);
    const params = new URLSearchParams(formData);

    // Update URL without reload
    const newUrl = `${window.location.pathname}?${params.toString()}`;
    window.history.pushState({}, '', newUrl);

    // Fetch filtered content and reload page
    // Simple implementation: just reload the page with new params
    window.location.href = newUrl;
}

function initFilters() {
    const searchInput = document.getElementById('search-input');
    const statusSelect = document.getElementById('status-select');
    const prioritySelect = document.getElementById('priority-select');

    // Add debounced listener to search input (500ms delay)
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            if (filterTimeout) clearTimeout(filterTimeout);
            filterTimeout = setTimeout(applyFilters, 500);
        });
    }

    // Add immediate listeners to selects (no delay)
    if (statusSelect) {
        statusSelect.addEventListener('change', applyFilters);
    }
    if (prioritySelect) {
        prioritySelect.addEventListener('change', applyFilters);
    }

    // Restore filter values from URL
    const urlParams = new URLSearchParams(window.location.search);
    if (searchInput && urlParams.has('search')) {
        searchInput.value = urlParams.get('search');
    }
    if (statusSelect && urlParams.has('status')) {
        statusSelect.value = urlParams.get('status');
    }
    if (prioritySelect && urlParams.has('priority')) {
        prioritySelect.value = urlParams.get('priority');
    }
}

// Shutdown functionality
function initShutdown() {
    const shutdownBtn = document.getElementById('shutdown-btn');
    if (!shutdownBtn) return;

    shutdownBtn.addEventListener('click', function() {
        if (confirm('Are you sure you want to shutdown the server?')) {
            fetch('/api/shutdown', {
                method: 'POST'
            })
            .then(response => response.json())
            .then(data => {
                console.log('Shutdown initiated:', data);
                // Show a message to the user
                shutdownBtn.textContent = 'Shutting down...';
                shutdownBtn.disabled = true;
            })
            .catch(error => {
                console.error('Shutdown error:', error);
                alert('Error shutting down server');
            });
        }
    });
}

// View selector functionality
document.addEventListener('DOMContentLoaded', function() {
    // Initialize theme
    initTheme();

    // Initialize filters
    initFilters();

    // Initialize shutdown button
    initShutdown();

    const viewRadios = document.querySelectorAll('input[name="view"]');
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

    // Restore selection from localStorage if available, otherwise use 'timeline'
    const saved = (typeof localStorage !== 'undefined') ? localStorage.getItem('beady.view') : null;
    const initial = saved && views[saved] ? saved : 'timeline';
    switchView(initial);

    // Set the correct radio button as checked
    viewRadios.forEach(radio => {
        if (radio.value === initial) {
            radio.checked = true;
        }

        // Add change listener to each radio button
        radio.addEventListener('change', function() {
            if (this.checked) {
                const v = this.value;
                switchView(v);
                try {
                    localStorage.setItem('beady.view', v);
                } catch (e) {
                    // ignore storage errors (e.g., privacy modes)
                }
            }
        });
    });
});
