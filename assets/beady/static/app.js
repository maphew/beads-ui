// Basic JavaScript for Beads UI

console.log('Beads UI loaded');

// Username management
function initUsername(serverUsername) {
    // Check if username is already stored
    let username = localStorage.getItem('beady-username');

    // If not stored and server provided one, use it
    if (!username && serverUsername) {
        username = serverUsername;
        try {
            localStorage.setItem('beady-username', username);
            console.log('Username initialized from server:', username);
        } catch (e) {
            console.warn('Could not save username to localStorage:', e);
        }
    }

    return username || 'web-user';
}

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
    const statusCheckboxes = document.querySelectorAll('input[name="status"]');
    const priorityCheckboxes = document.querySelectorAll('input[name="priority"]');

    // Add debounced listener to search input (500ms delay)
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            if (filterTimeout) clearTimeout(filterTimeout);
            filterTimeout = setTimeout(applyFilters, 500);
        });
    }

    // Add immediate listeners to checkboxes (no delay)
    statusCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', applyFilters);
    });
    priorityCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', applyFilters);
    });

    // Restore filter values from URL
    const urlParams = new URLSearchParams(window.location.search);
    if (searchInput && urlParams.has('search')) {
        searchInput.value = urlParams.get('search');
    }

    // Restore checked status checkboxes
    const statusValues = urlParams.getAll('status');
    statusCheckboxes.forEach(checkbox => {
        checkbox.checked = statusValues.includes(checkbox.value);
    });

    // Restore checked priority checkboxes
    const priorityValues = urlParams.getAll('priority');
    priorityCheckboxes.forEach(checkbox => {
        checkbox.checked = priorityValues.includes(checkbox.value);
    });
}

// Server connection monitoring
let connectionCheckInterval = null;
let serverOnline = true;

function checkServerConnection() {
    fetch('/api/stats', {
        method: 'GET',
        cache: 'no-cache'
    })
    .then(response => {
        if (response.ok) {
            updateConnectionStatus(true);
        } else {
            updateConnectionStatus(false);
        }
    })
    .catch(error => {
        updateConnectionStatus(false);
    });
}

function updateConnectionStatus(online) {
    const shutdownBtn = document.getElementById('shutdown-btn');
    if (!shutdownBtn) return;

    if (online !== serverOnline) {
        serverOnline = online;
        if (online) {
            shutdownBtn.textContent = 'Shutdown';
            shutdownBtn.disabled = false;
            shutdownBtn.classList.remove('contrast');
            shutdownBtn.classList.add('secondary', 'outline');
        } else {
            shutdownBtn.textContent = 'Server Offline';
            shutdownBtn.disabled = true;
            shutdownBtn.classList.remove('secondary', 'outline');
            shutdownBtn.classList.add('contrast');
        }
    }
}

function startConnectionMonitoring() {
    // Check connection every 3 seconds
    if (!connectionCheckInterval) {
        connectionCheckInterval = setInterval(checkServerConnection, 3000);
    }
}

function stopConnectionMonitoring() {
    if (connectionCheckInterval) {
        clearInterval(connectionCheckInterval);
        connectionCheckInterval = null;
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
                // Start monitoring connection to detect when server is offline
                startConnectionMonitoring();
            })
            .catch(error => {
                console.error('Shutdown error:', error);
                updateConnectionStatus(false);
            });
        }
    });

    // Start connection monitoring on page load
    startConnectionMonitoring();
}

// View selector functionality
document.addEventListener('DOMContentLoaded', function() {
    // Initialize username (use server-provided username if available)
    initUsername(window.beadyServerUsername);

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
