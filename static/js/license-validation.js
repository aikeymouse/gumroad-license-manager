// License validation functionality
document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('validateLicenseForm');
    const resultDiv = document.getElementById('validationResult');
    
    if (!form || !resultDiv) {
        return; // Elements not found, probably not on the licenses page
    }
    
    form.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const licenseKey = document.getElementById('licenseKey').value.trim();
        if (!licenseKey) {
            return;
        }
        
        // Show loading state
        resultDiv.innerHTML = '<div class="loading-spinner">Validating...</div>';
        resultDiv.style.display = 'block';
        resultDiv.className = 'validation-result';
        
        // Get product ID from the page data
        const productId = window.pageData ? window.pageData.productID : '';
        
        // Make validation request
        fetch('/validate-license', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                product_id: productId,
                license_key: licenseKey
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                resultDiv.className = 'validation-result success';
                resultDiv.innerHTML = `
                    <h4>✓ Valid License</h4>
                    <div class="license-details">
                        <p><strong>Uses:</strong> ${data.uses}</p>
                        <p><strong>Purchaser:</strong> ${data.purchase.email}</p>
                        <p><strong>Product:</strong> ${data.purchase.product_name}</p>
                        <p><strong>Sale Date:</strong> ${data.purchase.sale_timestamp}</p>
                        <p><strong>Price:</strong> $${(data.purchase.price / 100).toFixed(2)}</p>
                        ${data.purchase.refunded ? '<p class="status-warning"><strong>Status:</strong> Refunded</p>' : ''}
                        ${data.purchase.disputed ? '<p class="status-warning"><strong>Status:</strong> Disputed</p>' : ''}
                        ${data.purchase.chargebacked ? '<p class="status-error"><strong>Status:</strong> Chargebacked</p>' : ''}
                    </div>
                `;
            } else {
                resultDiv.className = 'validation-result error';
                resultDiv.innerHTML = `<h4>✗ Invalid License</h4><p>${data.message || 'License key is not valid'}</p>`;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            resultDiv.className = 'validation-result error';
            resultDiv.innerHTML = '<h4>✗ Validation Error</h4><p>Failed to validate license key. Please try again.</p>';
        });
    });
});
