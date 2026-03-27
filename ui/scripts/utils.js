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

// Made with Bob
