// Constants
const STATUS_UPDATE_INTERVAL_MS = 5 * 60 * 1000; // 5 minutes
const MIN_BLOCK_SIZE = 100; // Minimum size in pixels for usability
const CONTAINER_PADDING = 20; // Total padding for container
const BLOCK_MARGIN = 10; // Margin between blocks
const STATUS_CACHE_TTL_MS = 10 * 1000; // Cache status for 10 seconds

// Status cache to avoid redundant API calls
const statusCache = new Map();

// Function to show custom confirmation modal
function showConfirmModal(message) {
    return new Promise((resolve) => {
        const modal = document.getElementById('confirmModal');
        const messageElement = document.getElementById('confirmMessage');
        const confirmBtn = document.getElementById('confirmOk');
        const cancelBtn = document.getElementById('confirmCancel');

        // Set the message
        messageElement.textContent = message;

        // Show the modal
        modal.classList.add('show');

        // Handle confirm
        const handleConfirm = () => {
            modal.classList.remove('show');
            cleanup();
            resolve(true);
        };

        // Handle cancel
        const handleCancel = () => {
            modal.classList.remove('show');
            cleanup();
            resolve(false);
        };

        // Cleanup function to remove event listeners
        const cleanup = () => {
            confirmBtn.removeEventListener('click', handleConfirm);
            cancelBtn.removeEventListener('click', handleCancel);
            modal.removeEventListener('click', handleModalClick);
        };

        // Close modal when clicking outside
        const handleModalClick = (e) => {
            if (e.target === modal) {
                handleCancel();
            }
        };

        // Add event listeners
        confirmBtn.addEventListener('click', handleConfirm);
        cancelBtn.addEventListener('click', handleCancel);
        modal.addEventListener('click', handleModalClick);
    });
}

// Shared helper function to create status icons
function createIcon(type, color, stationId, action = null, additionalClass = '') {
    const wrapper = document.createElement('div');
    wrapper.classList.add('status-icon');
    if (additionalClass) {
        wrapper.classList.add(additionalClass);
    }
    wrapper.style.color = color;

    if (action) {
        wrapper.style.cursor = 'pointer';
        wrapper.addEventListener('click', async () => {
            if (action === 'stop') {
                const confirmMessage = getTranslation('confirm_stop_game', 'Sei sicuro di voler chiudere la partita?');
                const confirmed = await showConfirmModal(confirmMessage);
                if (confirmed) {
                    sendGameStationAction(stationId, action);
                }
            } else {
                sendGameStationAction(stationId, action);
            }
        });
    }

    let svgContent = '';
    if (type === 'play') {
        svgContent = `<svg viewBox="0 0 24 24" fill="currentColor" width="100%" height="100%"><path d="M8 5v14l11-7z"/></svg>`;
    } else if (type === 'pause') {
        svgContent = `<svg viewBox="0 0 24 24" fill="currentColor" width="100%" height="100%"><path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/></svg>`;
    } else if (type === 'stop') {
        svgContent = `<svg viewBox="0 0 24 24" fill="currentColor" width="100%" height="100%"><path d="M6 6h12v12H6z"/></svg>`;
    } else if (type === 'restart') {
        svgContent = `<svg viewBox="0 0 24 24" fill="currentColor" width="100%" height="100%"><path d="M12 5V1L7 6l5 5V7c3.31 0 6 2.69 6 6s-2.69 6-6 6-6-2.69-6-6H4c0 4.42 3.58 8 8 8s8-3.58 8-8-3.58-8-8-8z"/></svg>`;
    }
    wrapper.innerHTML = svgContent;
    return wrapper;
}

// Shared function to render status icons based on status value
function renderStatusIcons(status, stationId, container) {
    if (status === 1) {
        container.appendChild(createIcon('play', '#28a745', stationId, 'start'));
    }

    if (status === 0) {
        container.appendChild(createIcon('pause', '#ffc107', stationId, 'suspend'));
    }

    if (status === 2) {
        const suspendedContainer = document.createElement('div');
        suspendedContainer.classList.add('suspended-icons-container');
        suspendedContainer.appendChild(createIcon('stop', '#dc3545', stationId, 'stop', 'stop-icon'));
        suspendedContainer.appendChild(createIcon('restart', '#17a2b8', stationId, 'start', 'restart-icon'));
        container.appendChild(suspendedContainer);
    }
}

async function getGameStationIDList() {
    // Specify URL to be used to load UI config
    const API_URL = '/api/v1/gamestations';
    let list;

    try {
        // GET needed info
        const response = await fetch(API_URL);

        if (!response.ok) {
            throw new Error(`Error HTTP! Status: ${response.status}`);
        }

        // parse received JSON
        const gameStations = await response.json();

        // Get number of game stations
        list = gameStations.id;
    } catch (error) {
        console.error("ERROR in loading list og game station IDs:", error);
        return [];
    }
    return list;
}

async function getGameStation(id) {
    // Specify URL to be used to load UI config
    const API_URL = '/api/v1/gamestations/' + id;
    let gameStation;

    try {
        // GET needed info
        const response = await fetch(API_URL);

        if (!response.ok) {
            throw new Error(`Error HTTP! Status: ${response.status}`);
        }

        // parse received JSON
        gameStation = await response.json();

    } catch (error) {
        console.error("ERROR in loading game station " + id + " :", error);
        return null;
    }
    return gameStation;
}

async function getGameStationStatus(id, useCache = true) {
    // Check cache first if enabled
    if (useCache) {
        const cached = statusCache.get(id);
        if (cached && (Date.now() - cached.timestamp) < STATUS_CACHE_TTL_MS) {
            return cached.data;
        }
    }

    // Fetch only the status of a specific game station
    // API returns: { id: string, status: int, cost: int }
    const API_URL = '/api/v1/gamestations/' + id + '/status';

    try {
        // GET status info
        const response = await fetch(API_URL);

        if (!response.ok) {
            throw new Error(`Error HTTP! Status: ${response.status}`);
        }

        // parse received JSON
        const data = await response.json();

        // Cache the result
        statusCache.set(id, {
            data: data,
            timestamp: Date.now()
        });

        // Return the entire data object with status and cost
        return data;

    } catch (error) {
        console.error("ERROR in loading game station " + id + " status:", error);
        return null;
    }
}

// Function to send action to game station
async function sendGameStationAction(stationId, action) {
    const API_URL = '/api/v1/gamestations/' + stationId + '/action';

    try {
        const response = await fetch(API_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ action: action })
        });

        if (!response.ok) {
            const errorMsg = getTranslation('api_error', 'API Error');
            await showErrorModal(`${errorMsg}: ${response.status}`);
            return false;
        }

        const data = await response.json();

        // Check if result is success
        if (data.result === 'success') {
            console.log(`Action "${action}" sent successfully to game station ${stationId}`);

            // Update only this specific block's status
            await updateSingleBlockStatus(stationId);

            return true;
        } else {
            // Show error popup with translated message
            const errorMsg = getTranslation('action_failed', 'Action failed');
            await showErrorModal(`${errorMsg}: ${data.result || 'Unknown error'}`);
            return false;
        }
    } catch (error) {
        console.error(`ERROR sending action "${action}" to game station ${stationId}:`, error);
        const errorMsg = getTranslation('network_error', 'Network error');
        await showErrorModal(`${errorMsg}: ${error.message}`);
        return false;
    }
}

// Function to update a single block's status
async function updateSingleBlockStatus(stationId) {
    const container = document.getElementById('containerGamingStation');
    const block = container.querySelector(`[data-station-id="${stationId}"]`);

    if (!block) {
        console.error(`Block with station ID ${stationId} not found`);
        return;
    }

    // Fetch the status for this specific station (bypass cache to get fresh data)
    const statusData = await getGameStationStatus(stationId, false);

    if (statusData === null) {
        return;
    }

    const status = statusData.status;
    const cost = statusData.cost;

    // Update cost in block_station_price
    const block_price = block.querySelector('.block_station_price');
    if (block_price) {
        // If status is not initialized or Stopped (1), show zero cost, otherwise show actual cost
        block_price.textContent = (status === undefined || status === null || status === 1) ? formatPrice(0) : formatPrice(cost);
    }

    // Update status in block_station_status
    const block_status = block.querySelector('.block_station_status');
    if (block_status) {
        // Clear existing status icons (but preserve consumption icon if it exists)
        const consumptionIcon = block_status.querySelector('.consumption-icon');
        block_status.innerHTML = '';

        // Re-add consumption icon if it existed
        if (consumptionIcon) {
            // Update disabled state based on new status
            if (status === 1) {
                consumptionIcon.classList.add('disabled');
            } else {
                consumptionIcon.classList.remove('disabled');
            }

            // Re-attach the click event listener
            consumptionIcon.addEventListener('click', (e) => {
                e.stopPropagation();

                // Don't allow click if disabled
                if (consumptionIcon.classList.contains('disabled')) {
                    return;
                }

                if (typeof openConsumptionModal === 'function') {
                    openConsumptionModal(stationId);
                } else {
                    console.error("openConsumptionModal not available");
                }
            });
            block_status.appendChild(consumptionIcon);
        }

        // Apply status icons using shared function
        renderStatusIcons(status, stationId, block_status);
    }
}

// Helper function to translate consumption tooltip
function translateConsumptionTooltip(element) {
    const tooltipText = getTranslation('add_consumption', 'Add consumption');
    element.title = tooltipText;
}


async function generateGamingBlocks() {
    /* Number of blocks depends on number of gaming station using API*/
    const gameStationIDList = await getGameStationIDList();
    const numBlocksInput = parseInt(gameStationIDList.length);
    const container = document.getElementById('containerGamingStation');

    container.innerHTML = '';

    // Calculate grid to fit all blocks as squares
    const cols = Math.ceil(Math.sqrt(numBlocksInput));
    const rows = Math.ceil(numBlocksInput / cols);

    // Calculate square size based on available space with padding
    const availableWidth = window.innerWidth - CONTAINER_PADDING;
    const availableHeight = window.innerHeight - CONTAINER_PADDING;

    let size = Math.floor(Math.min(availableWidth / cols, availableHeight / rows)) - BLOCK_MARGIN;

    if (size < MIN_BLOCK_SIZE) size = MIN_BLOCK_SIZE; // Minimum size for usability

    // Fetch all game stations in parallel for better performance
    const gameStations = await Promise.all(gameStationIDList.map(id => getGameStation(id)));

    for (let i = 0; i < numBlocksInput; i++) {
        let gameStation = gameStations[i]
        let imageUrl = gameStation.iconPath
        // Validate and sanitize imageUrl to prevent CSS injection
        const urlPattern = /^(https?:\/\/|\/|\.\/|[a-zA-Z0-9])[^\s'"()]*\.(jpg|jpeg|png|gif|webp|svg)$/i;
        const sanitizedUrl = (imageUrl && urlPattern.test(imageUrl)) ? imageUrl.replace(/['"()]/g, '') : '';

        const block = document.createElement('div');
        block.classList.add('block_station');
        block.dataset.stationId = gameStation.id; // Store station ID for later reference
        // block.textContent = gameStation.name; // Removed: will be added as a separate element
        block.style.width = `${size}px`;
        block.style.height = `${size}px`;
        // Block icon must be retrieved using API
        if (sanitizedUrl) {
            block.style.backgroundImage = `url('${sanitizedUrl}')`;
        }

        // Status Logic - declare status first
        const status = gameStation.status; // 0: Started, 1: Stopped, 2: Suspended

        const block_price = document.createElement('div');
        block_price.classList.add('block_station_price');
        // If status is not initialized or Stopped (1), show zero cost, otherwise show actual cost
        block_price.textContent = (status === undefined || status === null || status === 1) ? formatPrice(0) : formatPrice(gameStation.cost);

        const block_status = document.createElement('div');
        block_status.classList.add('block_station_status');

        // Add consumption icon
        const consumptionIcon = document.createElement('div');
        consumptionIcon.classList.add('consumption-icon');
        consumptionIcon.innerHTML = `<svg viewBox="0 0 24 24" fill="currentColor" width="100%" height="100%">
            <path d="M11 9H9V2H7v7H5V2H3v7c0 2.12 1.66 3.84 3.75 3.97V22h2.5v-9.03C11.34 12.84 13 11.12 13 9V2h-2v7zm5-3v8h2.5v8H21V2c-2.76 0-5 2.24-5 4z"/>
        </svg>`;
        consumptionIcon.title = 'add_consumption'; // Will be translated
        consumptionIcon.dataset.translate = 'add_consumption';

        // Enable/disable based on status (enabled only when game is active: status 0 or 2)
        if (status === 1) {
            consumptionIcon.classList.add('disabled');
        }

        consumptionIcon.addEventListener('click', (e) => {
            e.stopPropagation(); // Prevent triggering other click events

            // Don't allow click if disabled
            if (consumptionIcon.classList.contains('disabled')) {
                return;
            }

            if (typeof openConsumptionModal === 'function') {
                openConsumptionModal(gameStation.id);
            } else {
                console.error("openConsumptionModal not available");
            }
        });
        block_status.appendChild(consumptionIcon);

        // Apply status icons using shared function
        renderStatusIcons(status, gameStation.id, block_status);
        const block_name = document.createElement('div');
        block_name.classList.add('station-name');
        block_name.textContent = gameStation.name;

        // Add parts to main block
        block.appendChild(block_price);
        block.appendChild(block_status);
        block.appendChild(block_name); // Append name last to be on top

        // Translate consumption icon tooltip after adding to DOM
        translateConsumptionTooltip(consumptionIcon);

        container.appendChild(block);
    }
}

// Function to update only the status of existing blocks
async function updateBlockStationStatus() {
    const gameStationIDList = await getGameStationIDList();
    const container = document.getElementById('containerGamingStation');
    const blocks = container.querySelectorAll('.block_station');

    // Fetch status and cost for all game stations in parallel using the status API
    // Bypass cache for periodic updates to ensure fresh data
    const statusData = await Promise.all(gameStationIDList.map(id => getGameStationStatus(id, false)));

    // Update each block's status and cost
    blocks.forEach((block, index) => {
        if (index < statusData.length && statusData[index] !== null) {
            const data = statusData[index];
            const status = data.status;
            const cost = data.cost;
            const stationId = gameStationIDList[index];

            // Update cost in block_station_price
            const block_price = block.querySelector('.block_station_price');
            if (block_price) {
                // If status is not initialized or Stopped (1), show zero cost, otherwise show actual cost
                block_price.textContent = (status === undefined || status === null || status === 1) ? formatPrice(0) : formatPrice(cost);
            }

            // Update status in block_station_status
            const block_status = block.querySelector('.block_station_status');
            if (block_status) {
                // Clear existing status icons but preserve consumption icon
                const consumptionIcon = block_status.querySelector('.consumption-icon');
                block_status.innerHTML = '';

                // Re-add consumption icon if it existed
                if (consumptionIcon) {
                    // Re-attach the click event listener
                    consumptionIcon.addEventListener('click', (e) => {
                        e.stopPropagation();
                        if (typeof openConsumptionModal === 'function') {
                            openConsumptionModal(stationId);
                        }
                    });
                    block_status.appendChild(consumptionIcon);
                }

                // Apply status icons using shared function
                renderStatusIcons(status, stationId, block_status);
            }
        }
    });

    console.log('Block station status updated at:', new Date().toLocaleTimeString());
}

async function generateBlocks() {
    await generateGamingBlocks();

    // Set up automatic status updates every 5 minutes (300,000 milliseconds)
    // Clear any existing interval to prevent duplicates
    if (window.statusUpdateInterval) {
        clearInterval(window.statusUpdateInterval);
    }

    window.statusUpdateInterval = setInterval(() => {
        updateBlockStationStatus();
    }, STATUS_UPDATE_INTERVAL_MS);

    console.log('Automatic status updates enabled (every 5 minutes)');
}
