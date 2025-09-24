// Gumroad License Manager JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // Initialize the application
    initializeApp();
});

function initializeApp() {
    // Add loading states to buttons
    addLoadingStates();
    
    // Add search functionality
    addSearchFunctionality();
    
    // Add keyboard shortcuts
    setupKeyboardShortcuts();
    
    // Add copy functionality for license keys
    addCopyFunctionality();
}

function addLoadingStates() {
    const buttons = document.querySelectorAll('.btn, .view-licenses');
    
    buttons.forEach(button => {
        button.addEventListener('click', function(e) {
            // Don't add loading state for external links
            if (this.href && this.href.startsWith('http') && !this.href.includes(window.location.hostname)) {
                return;
            }
            
            // Add loading state
            const originalText = this.textContent;
            this.innerHTML = '<span class="loading"></span> Loading...';
            this.style.pointerEvents = 'none';
            
            // Remove loading state after navigation or timeout
            setTimeout(() => {
                this.textContent = originalText;
                this.style.pointerEvents = 'auto';
            }, 2000);
        });
    });
}

function addSearchFunctionality() {
    // Add search box to licenses page
    if (window.location.pathname.startsWith('/licenses/')) {
        addLicenseSearch();
    }
    
    // Add search box to API log page
    if (window.location.pathname === '/api-log') {
        addAPILogSearch();
    }
}

function addLicenseSearch() {
    const table = document.querySelector('table');
    if (!table) return;
    
    const searchDiv = document.createElement('div');
    searchDiv.style.marginBottom = '20px';
    searchDiv.innerHTML = `
        <input type="text" id="licenseSearch" placeholder="Search licenses..." 
               style="width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 16px;">
    `;
    
    table.parentNode.insertBefore(searchDiv, table);
    
    const searchInput = document.getElementById('licenseSearch');
    const rows = table.querySelectorAll('tbody tr');
    
    searchInput.addEventListener('input', function() {
        const searchTerm = this.value.toLowerCase();
        
        rows.forEach(row => {
            const rowText = row.textContent.toLowerCase();
            if (rowText.includes(searchTerm)) {
                row.style.display = '';
            } else {
                row.style.display = 'none';
            }
        });
    });
}

function addAPILogSearch() {
    const table = document.querySelector('table');
    if (!table) return;
    
    const searchDiv = document.createElement('div');
    searchDiv.style.marginBottom = '20px';
    searchDiv.innerHTML = `
        <div style="display: flex; gap: 10px; flex-wrap: wrap;">
            <input type="text" id="apiLogSearch" placeholder="Search API calls..." 
                   style="flex: 1; min-width: 200px; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 16px;">
            <select id="statusFilter" style="padding: 10px; border: 1px solid #ddd; border-radius: 4px;">
                <option value="">All Status</option>
                <option value="200">Success (200)</option>
                <option value="error">Errors</option>
            </select>
        </div>
    `;
    
    table.parentNode.insertBefore(searchDiv, table);
    
    const searchInput = document.getElementById('apiLogSearch');
    const statusFilter = document.getElementById('statusFilter');
    const rows = table.querySelectorAll('tbody tr');
    
    function filterRows() {
        const searchTerm = searchInput.value.toLowerCase();
        const statusValue = statusFilter.value;
        
        rows.forEach(row => {
            const rowText = row.textContent.toLowerCase();
            const statusCell = row.cells[3]; // Status column
            const statusText = statusCell ? statusCell.textContent : '';
            
            let showRow = true;
            
            // Text search filter
            if (searchTerm && !rowText.includes(searchTerm)) {
                showRow = false;
            }
            
            // Status filter
            if (statusValue) {
                if (statusValue === 'error' && statusText === '200') {
                    showRow = false;
                } else if (statusValue === '200' && statusText !== '200') {
                    showRow = false;
                }
            }
            
            row.style.display = showRow ? '' : 'none';
        });
    }
    
    searchInput.addEventListener('input', filterRows);
    statusFilter.addEventListener('change', filterRows);
}

function setupKeyboardShortcuts() {
    document.addEventListener('keydown', function(e) {
        // Escape to go back
        if (e.key === 'Escape') {
            const backLink = document.querySelector('.back-link');
            if (backLink) {
                window.location.href = backLink.href;
            }
        }
        
        // Ctrl/Cmd + F to focus search
        if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
            const searchInput = document.querySelector('#productSearch, #licenseSearch, #apiLogSearch');
            if (searchInput) {
                e.preventDefault();
                searchInput.focus();
            }
        }
    });
}

function addCopyFunctionality() {
    const licenseKeys = document.querySelectorAll('.license-key');
    
    licenseKeys.forEach(licenseKey => {
        licenseKey.style.cursor = 'pointer';
        licenseKey.title = 'Click to copy';
        
        licenseKey.addEventListener('click', function() {
            const text = this.textContent;
            
            if (navigator.clipboard) {
                navigator.clipboard.writeText(text).then(() => {
                    showCopySuccess(this);
                });
            } else {
                // Fallback for older browsers
                const textArea = document.createElement('textarea');
                textArea.value = text;
                document.body.appendChild(textArea);
                textArea.select();
                document.execCommand('copy');
                document.body.removeChild(textArea);
                showCopySuccess(this);
            }
        });
    });
}

function showCopySuccess(element) {
    const originalBg = element.style.backgroundColor;
    element.style.backgroundColor = '#d4edda';
    element.style.borderColor = '#c3e6cb';
    
    setTimeout(() => {
        element.style.backgroundColor = originalBg;
        element.style.borderColor = '#ddd';
    }, 1000);
}

// Utility function to format dates
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
}

// Add some visual enhancements
function addVisualEnhancements() {
    // Add hover effects to table rows
    const tableRows = document.querySelectorAll('tbody tr');
    tableRows.forEach(row => {
        row.addEventListener('mouseenter', function() {
            this.style.backgroundColor = '#f8f9fa';
        });
        
        row.addEventListener('mouseleave', function() {
            this.style.backgroundColor = '';
        });
    });
}
