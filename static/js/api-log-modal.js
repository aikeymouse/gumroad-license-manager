// API Log Modal JavaScript
let apiCallsData = [];

// Fetch API calls data from the server
async function loadApiCallsData() {
    try {
        const response = await fetch('/api/api-calls');
        if (response.ok) {
            apiCallsData = await response.json();
        } else {
            console.error('Failed to load API calls data');
        }
    } catch (error) {
        console.error('Error loading API calls data:', error);
    }
}

function showModal(index) {
    const call = apiCallsData[index];
    if (!call) {
        console.error('API call not found at index:', index);
        return;
    }
    
    const modal = document.getElementById('apiModal');
    
    // Populate modal fields
    document.getElementById('modal-method').textContent = call.Method || call.method || '';
    document.getElementById('modal-url').textContent = call.URL || call.url || '';
    
    // Format timestamp
    const timestamp = call.Timestamp || call.timestamp || '';
    if (timestamp) {
        // If it's a Go time format, convert it to readable format
        const date = new Date(timestamp);
        document.getElementById('modal-timestamp').textContent = isNaN(date.getTime()) ? timestamp : date.toLocaleString();
    } else {
        document.getElementById('modal-timestamp').textContent = '';
    }
    
    // Format duration consistently as rounded milliseconds
    const duration = call.Duration || call.duration || 0;
    let durationMs = '';
    if (duration) {
        // Convert nanoseconds to milliseconds and round
        const ms = Math.round(duration / 1000000);
        durationMs = `${ms}ms`;
    }
    document.getElementById('modal-duration').textContent = durationMs;
    
    document.getElementById('modal-status').textContent = call.Status || call.status || '';
    
    // Format request body
    const requestBody = call.RequestBody || call.requestBody || '';
    if (requestBody) {
        try {
            const parsed = JSON.parse(requestBody);
            document.getElementById('modal-request-body').textContent = JSON.stringify(parsed, null, 2);
        } catch (e) {
            document.getElementById('modal-request-body').textContent = requestBody;
        }
    } else {
        document.getElementById('modal-request-body').textContent = 'No request body';
    }
    
    // Format response body
    const responseBody = call.ResponseBody || call.responseBody || '';
    if (responseBody) {
        try {
            const parsed = JSON.parse(responseBody);
            document.getElementById('modal-response').textContent = JSON.stringify(parsed, null, 2);
        } catch (e) {
            document.getElementById('modal-response').textContent = responseBody;
        }
    } else {
        document.getElementById('modal-response').textContent = 'No response body';
    }
    
    // Show/hide error section
    const error = call.Error || call.error || '';
    if (error) {
        document.getElementById('modal-error-section').style.display = 'block';
        document.getElementById('modal-error').textContent = error;
    } else {
        document.getElementById('modal-error-section').style.display = 'none';
    }
    
    modal.style.display = 'block';
}

// Modal functionality
document.addEventListener('DOMContentLoaded', function() {
    const modal = document.getElementById('apiModal');
    const closeBtn = document.querySelector('.close');
    
    if (closeBtn) {
        closeBtn.onclick = function() {
            modal.style.display = 'none';
        }
    }
    
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = 'none';
        }
    }
    
    // Load API calls data when the page loads
    loadApiCallsData();
});
