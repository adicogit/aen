// Consumption Modal Management
let currentStationId = null;
let allConsumptions = [];
let selectedConsumptions = [];

// Function to fetch all warehouse items
async function fetchWarehouseItems() {
    try {
        // First, get the list of item IDs
        const response = await fetch('/api/v1/warehouseitems');
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const data = await response.json();
        const itemIds = data.id || [];
        
        // Fetch details for each item
        const itemPromises = itemIds.map(async (itemId) => {
            const itemResponse = await fetch(`/api/v1/warehouseitems/${itemId}`);
            if (!itemResponse.ok) {
                console.error(`Failed to fetch item ${itemId}`);
                return null;
            }
            return await itemResponse.json();
        });
        
        const items = await Promise.all(itemPromises);
        return items.filter(item => item !== null);
    } catch (error) {
        console.error('Error fetching warehouse items:', error);
        return [];
    }
}

// Function to format price from cents to euros
function formatPrice(priceInCents) {
    const euros = (priceInCents || 0) / 100;
    return `€${euros.toFixed(2)}`;
}

// Function to render consumption list
function renderConsumptionList(items, searchTerm = '') {
    const consumptionList = document.getElementById('consumptionList');
    consumptionList.innerHTML = '';
    
    // Filter items based on search term
    const filteredItems = items.filter(item => 
        item.name.toLowerCase().includes(searchTerm.toLowerCase())
    );
    
    if (filteredItems.length === 0) {
        consumptionList.innerHTML = '<div class="empty-message">Nessuna consumazione trovata</div>';
        return;
    }
    
    filteredItems.forEach(item => {
        const itemElement = document.createElement('div');
        itemElement.classList.add('consumption-item');
        itemElement.dataset.itemId = item.id;
        
        itemElement.innerHTML = `
            <span class="consumption-item-name">${item.name}</span>
            <span class="consumption-item-price">${formatPrice(item.publicPrice)}</span>
        `;
        
        itemElement.addEventListener('click', () => addConsumption(item));
        consumptionList.appendChild(itemElement);
    });
}

// Function to add a consumption to selected list
function addConsumption(item) {
    // Add to selected consumptions array
    selectedConsumptions.push({
        id: item.id,
        name: item.name,
        price: item.publicPrice
    });
    
    renderSelectedConsumptions();
}

// Function to remove a consumption from selected list
function removeConsumption(index) {
    selectedConsumptions.splice(index, 1);
    renderSelectedConsumptions();
}

// Function to render selected consumptions
function renderSelectedConsumptions() {
    const selectedContainer = document.getElementById('selectedConsumptions');
    selectedContainer.innerHTML = '';
    
    if (selectedConsumptions.length === 0) {
        selectedContainer.innerHTML = '<div class="empty-message">Nessun articolo selezionato</div>';
        return;
    }
    
    selectedConsumptions.forEach((item, index) => {
        const itemElement = document.createElement('div');
        itemElement.classList.add('selected-item');
        
        itemElement.innerHTML = `
            <div class="selected-item-info">
                <div class="selected-item-name">${item.name}</div>
                <div class="selected-item-price">${formatPrice(item.price)}</div>
            </div>
            <button class="remove-btn" data-index="${index}">
                <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
                </svg>
            </button>
        `;
        
        const removeBtn = itemElement.querySelector('.remove-btn');
        removeBtn.addEventListener('click', () => removeConsumption(index));
        
        selectedContainer.appendChild(itemElement);
    });
}

// Function to open consumption modal
async function openConsumptionModal(stationId) {
    currentStationId = stationId;
    selectedConsumptions = [];
    
    const modal = document.getElementById('consumptionModal');
    modal.classList.add('show');
    
    // Fetch and display warehouse items
    allConsumptions = await fetchWarehouseItems();
    renderConsumptionList(allConsumptions);
    renderSelectedConsumptions();
}

// Function to close consumption modal
function closeConsumptionModal() {
    const modal = document.getElementById('consumptionModal');
    modal.classList.remove('show');
    currentStationId = null;
    selectedConsumptions = [];
}

// Function to save consumptions
async function saveConsumptions() {
    if (selectedConsumptions.length === 0) {
        alert('Seleziona almeno una consumazione');
        return;
    }
    
    if (!currentStationId) {
        console.error('No station ID set');
        return;
    }
    
    try {
        // Send consumption data to server
        const response = await fetch(`/api/v1/gamestations/${currentStationId}/consumption`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                items: selectedConsumptions.map(item => ({
                    id: item.id,
                    quantity: 1 // For now, always 1 item per click
                }))
            })
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('Consumptions saved:', result);
        
        // Close modal and update station status
        closeConsumptionModal();
        
        // Update the station's display
        if (typeof updateSingleBlockStatus === 'function') {
            await updateSingleBlockStatus(currentStationId);
        }
        
    } catch (error) {
        console.error('Error saving consumptions:', error);
        alert('Errore nel salvare le consumazioni');
    }
}

// Initialize modal event listeners when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    // Search field event listener
    const searchField = document.getElementById('consumptionSearch');
    if (searchField) {
        searchField.addEventListener('input', (e) => {
            renderConsumptionList(allConsumptions, e.target.value);
        });
    }
    
    // Cancel button
    const cancelBtn = document.getElementById('cancelConsumption');
    if (cancelBtn) {
        cancelBtn.addEventListener('click', closeConsumptionModal);
    }
    
    // Save button
    const saveBtn = document.getElementById('saveConsumption');
    if (saveBtn) {
        saveBtn.addEventListener('click', saveConsumptions);
    }
    
    // Close modal when clicking outside
    const modal = document.getElementById('consumptionModal');
    if (modal) {
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                closeConsumptionModal();
            }
        });
    }
});

// Made with Bob
