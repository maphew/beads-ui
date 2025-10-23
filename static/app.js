// Theme management
(function() {
    const html = document.documentElement;

    // Get system preference
    function getSystemTheme() {
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }

    // Get stored theme or default to system preference
    function getStoredTheme() {
        return localStorage.getItem('theme') || 'auto';
    }

    // Apply theme
    function applyTheme(theme) {
        if (theme === 'auto') {
            theme = getSystemTheme();
        }
        html.setAttribute('data-theme', theme);
        updateToggleButton(theme);
    }

    // Update toggle button icon/text
    function updateToggleButton(currentTheme) {
        const themeToggle = document.getElementById('theme-toggle');
        if (!themeToggle) return;

        const isDark = currentTheme === 'dark';
        themeToggle.innerHTML = isDark ? '☀️' : '🌙';
        themeToggle.title = isDark ? 'Switch to light mode' : 'Switch to dark mode';
    }

    // Listen for system theme changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
        const storedTheme = getStoredTheme();
        if (storedTheme === 'auto') {
            applyTheme('auto');
        }
    });

    // Initialize theme when DOM is ready
    function initTheme() {
        applyTheme(getStoredTheme());

        // Add toggle event listener
        const themeToggle = document.getElementById('theme-toggle');
        if (themeToggle) {
            themeToggle.addEventListener('click', function() {
                const currentTheme = getStoredTheme();
                let newTheme;

                if (currentTheme === 'auto') {
                    newTheme = getSystemTheme() === 'dark' ? 'light' : 'dark';
                } else {
                    newTheme = currentTheme === 'dark' ? 'light' : 'dark';
                }

                localStorage.setItem('theme', newTheme);
                applyTheme(newTheme);
            });
        }
    }

    // Run when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initTheme);
    } else {
        initTheme();
    }
})();

// Live reload WebSocket connection
(function() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = protocol + '//' + window.location.host + '/ws';
    const ws = new WebSocket(wsUrl);

    ws.onmessage = function(event) {
        if (event.data === 'reload') {
            window.location.reload();
        }
    };

    ws.onclose = function() {
        // Optionally reconnect after a delay
        setTimeout(function() {
            window.location.reload();
        }, 1000);
    };
})();
