// Shared utility functions

/**
 * Format price from cents to euros
 * @param {number} priceInCents - Price in cents
 * @returns {string} Formatted price string (e.g., "€12.50")
 */
function formatPrice(priceInCents) {
    const euros = (priceInCents || 0) / 100;
    return `€${euros.toFixed(2)}`;
}

// Function to show custom error modal
function showErrorModal(message) {
    return new Promise((resolve) => {
        const modal = document.getElementById('errorModal');
        const messageElement = document.getElementById('errorMessage');
        const okBtn = document.getElementById('errorOk');
        
        // Set the message
        messageElement.textContent = message;
        
        // Show the modal
        modal.classList.add('show');
        
        // Handle OK button
        const handleOk = () => {
            modal.classList.remove('show');
            cleanup();
            resolve();
        };
        
        // Cleanup function to remove event listeners
        const cleanup = () => {
            okBtn.removeEventListener('click', handleOk);
            modal.removeEventListener('click', handleModalClick);
        };
        
        // Close modal when clicking outside
        const handleModalClick = (e) => {
            if (e.target === modal) {
                handleOk();
            }
        };
        
        // Add event listeners
        okBtn.addEventListener('click', handleOk);
        modal.addEventListener('click', handleModalClick);
    });
}

async function requireAccountingDay() {
    try {
        const response = await fetch('/api/v1/accountingday');
        const data = await response.json();
        if (data.currentAccountingDay) {
            return true;
        }
        const msg = getTranslation('cash_no_accounting_day', 'Nessun giorno contabile impostato. Imposta il giorno contabile dalla cassa.');
        await showErrorModal(msg);
        return false;
    } catch (error) {
        console.error('Error checking accounting day:', error);
        const msg = getTranslation('network_error', 'Errore di rete');
        await showErrorModal(msg);
        return false;
    }
}
