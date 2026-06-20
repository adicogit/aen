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


// Function to render consumption list
function renderConsumptionList(items, searchTerm = '') {
    const consumptionList = document.getElementById('consumptionList');
    consumptionList.innerHTML = '';
    
    // Filter items based on search term
    const filteredItems = items.filter(item => 
        item.name.toLowerCase().includes(searchTerm.toLowerCase())
    );
    
    if (filteredItems.length === 0) {
        const message = getTranslation('no_consumption_found', 'Nessuna consumazione trovata');
        consumptionList.innerHTML = `<div class="empty-message">${message}</div>`;
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

// Function to update save button state
function updateSaveButtonState() {
    const saveBtn = document.getElementById('saveConsumption');
    if (saveBtn) {
        if (selectedConsumptions.length === 0) {
            saveBtn.disabled = true;
            saveBtn.classList.add('disabled');
        } else {
            saveBtn.disabled = false;
            saveBtn.classList.remove('disabled');
        }
    }
}

// Function to render selected consumptions
function renderSelectedConsumptions() {
    const selectedContainer = document.getElementById('selectedConsumptions');
    selectedContainer.innerHTML = '';
    
    if (selectedConsumptions.length === 0) {
        const message = getTranslation('no_item_selected', 'Nessun articolo selezionato');
        selectedContainer.innerHTML = `<div class="empty-message">${message}</div>`;
        updateSaveButtonState();
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
    
    updateSaveButtonState();
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
    // Button should be disabled if no items selected, but check anyway
    if (selectedConsumptions.length === 0) {
        return;
    }
    
    if (!currentStationId) {
        console.error('No station ID set');
        return;
    }
    
    const hasDay = await requireAccountingDay();
    if (!hasDay) return;

    try {
        // Send consumption data to server
        const response = await fetch(`/api/v1/gamestations/${currentStationId}/consumption`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                items: selectedConsumptions.map(item => item.id)
            })
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('Consumptions saved:', result);
        
        // Save station ID before closing modal (which sets it to null)
        const stationId = currentStationId;
        
        // Close modal
        closeConsumptionModal();
        
        // Update the station's display
        if (typeof updateSingleBlockStatus === 'function') {
            await updateSingleBlockStatus(stationId);
        }
        
    } catch (error) {
        console.error('Error saving consumptions:', error);
        const message = getTranslation('error_saving_consumptions', 'Errore nel salvare le consumazioni');
        await showErrorModal(message);
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
