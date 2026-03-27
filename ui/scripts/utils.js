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

// Made with Bob
