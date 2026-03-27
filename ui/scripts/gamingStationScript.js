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

async function getGameStationStatus(id) {
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
        
        // Return the entire data object with status and cost
        return data;

    } catch (error) {
        console.error("ERROR in loading game station " + id + " status:", error);
        return null;
    }
}

// Function to get translated text
async function getTranslation(key) {
    if (!translations) {
        translations = await loadTranslations();
    }
    const userLang = navigator.language || navigator.userLanguage;
    const langCode = userLang.split('-')[0];
    const languageData = translations[langCode] || translations['en'];
    return languageData[key] || key;
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
            const errorMsg = await getTranslation('api_error') || 'API Error';
            alert(`${errorMsg}: ${response.status}`);
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
            const errorMsg = await getTranslation('action_failed') || 'Action failed';
            alert(`${errorMsg}: ${data.result || 'Unknown error'}`);
            return false;
        }
    } catch (error) {
        console.error(`ERROR sending action "${action}" to game station ${stationId}:`, error);
        const errorMsg = await getTranslation('network_error') || 'Network error';
        alert(`${errorMsg}: ${error.message}`);
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

    // Fetch the status for this specific station
    const statusData = await getGameStationStatus(stationId);
    
    if (statusData === null) {
        return;
    }

    const status = statusData.status;
    const cost = statusData.cost;
    
    // Update cost in block_station_price
    const block_price = block.querySelector('.block_station_price');
    if (block_price) {
        // If status is not initialized or Stopped (1), show zero cost, otherwise show actual cost
        block_price.textContent = (status === undefined || status === null || status === 1) ? formatCost(0) : formatCost(cost);
    }
    
    // Update status in block_station_status
    const block_status = block.querySelector('.block_station_status');
    if (block_status) {
        // Clear existing status icons (but preserve consumption icon if it exists)
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
        
        // Helper to create icon
        const createIcon = (type, color, action = null, additionalClass = '') => {
            const wrapper = document.createElement('div');
            wrapper.classList.add('status-icon');
            if (additionalClass) {
                wrapper.classList.add(additionalClass);
            }
            wrapper.style.color = color;
            
            if (action) {
                wrapper.style.cursor = 'pointer';
                wrapper.addEventListener('click', () => {
                    sendGameStationAction(stationId, action);
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
        };

        // Apply same status logic
        if (status === 1) {
            block_status.appendChild(createIcon('play', '#28a745', 'start'));
        }

        if (status === 0) {
            block_status.appendChild(createIcon('pause', '#ffc107', 'suspend'));
        }

        if (status === 2) {
            const suspendedContainer = document.createElement('div');
            suspendedContainer.classList.add('suspended-icons-container');
            suspendedContainer.appendChild(createIcon('stop', '#dc3545', 'stop', 'stop-icon'));
            suspendedContainer.appendChild(createIcon('restart', '#17a2b8', 'start', 'restart-icon'));
            block_status.appendChild(suspendedContainer);
        }
    }
}

// Helper function to translate consumption tooltip
async function translateConsumptionTooltip(element) {
    const tooltipText = await getTranslation('add_consumption');
    element.title = tooltipText;
}

// Helper function to format cost from cents to euros
function formatCost(costInCents) {
    const euros = (costInCents || 0) / 100;
    return `€${euros.toFixed(2)}`;
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
    const availableWidth = window.innerWidth - 20; // 20px total padding
    const availableHeight = window.innerHeight - 20;

    let size = Math.floor(Math.min(availableWidth / cols, availableHeight / rows)) - 10;

    if (size < 100) size = 100; // Minimum size for usability

    // Fetch all game stations in parallel for better performance
    const gameStations = await Promise.all(gameStationIDList.map(id => getGameStation(id)));

    for (let i = 0; i < numBlocksInput; i++) {
        let gameStation = gameStations[i]
        let imageUrl = gameStation.iconPath
        // Validate and sanitize imageUrl to prevent CSS injection
        const urlPattern = /^(https?:\/\/|\/|\.\/)[^\s'"()]*\.(jpg|jpeg|png|gif|webp|svg)$/i;
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
        block_price.textContent = (status === undefined || status === null || status === 1) ? formatCost(0) : formatCost(gameStation.cost);

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
        consumptionIcon.addEventListener('click', (e) => {
            e.stopPropagation(); // Prevent triggering other click events
            if (typeof openConsumptionModal === 'function') {
                openConsumptionModal(gameStation.id);
            }
        });
        block_status.appendChild(consumptionIcon);

        // Helper to create icon with optional action
        const createIcon = (type, color, action = null, additionalClass = '') => {
            const wrapper = document.createElement('div');
            wrapper.classList.add('status-icon');
            if (additionalClass) {
                wrapper.classList.add(additionalClass);
            }
            wrapper.style.color = color;
            
            if (action) {
                wrapper.style.cursor = 'pointer';
                wrapper.addEventListener('click', () => {
                    sendGameStationAction(gameStation.id, action);
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
        };

        // play in verde se lo stato di gameStation non è Stopped (1)
        if (status === 1) {
            block_status.appendChild(createIcon('play', '#28a745', 'start')); // Green with start action
        }

        // il simbolo di pausa in giallo se lo stato è Started (0)
        if (status === 0) {
            block_status.appendChild(createIcon('pause', '#ffc107', 'suspend')); // Yellow with suspend action
        }

        // simbolo di stop e restart verticalmente se lo stato è Suspended (2)
        if (status === 2) {
            const suspendedContainer = document.createElement('div');
            suspendedContainer.classList.add('suspended-icons-container');
            suspendedContainer.appendChild(createIcon('stop', '#dc3545', 'stop', 'stop-icon')); // Red with stop action - 60% height
            suspendedContainer.appendChild(createIcon('restart', '#17a2b8', 'start', 'restart-icon')); // Cyan with start action - 40% height
            block_status.appendChild(suspendedContainer);
        }
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
    const statusData = await Promise.all(gameStationIDList.map(id => getGameStationStatus(id)));

    // Helper to create icon (same as in generateGamingBlocks)
    const createIcon = (type, color, stationId, action = null, additionalClass = '') => {
        const wrapper = document.createElement('div');
        wrapper.classList.add('status-icon');
        if (additionalClass) {
            wrapper.classList.add(additionalClass);
        }
        wrapper.style.color = color;
        
        if (action) {
            wrapper.style.cursor = 'pointer';
            wrapper.addEventListener('click', () => {
                sendGameStationAction(stationId, action);
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
    };

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
                block_price.textContent = (status === undefined || status === null || status === 1) ? formatCost(0) : formatCost(cost);
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

                // Apply same status logic as in generateGamingBlocks
                if (status === 1) {
                    block_status.appendChild(createIcon('play', '#28a745', stationId, 'start')); // Green with start action
                }

                if (status === 0) {
                    block_status.appendChild(createIcon('pause', '#ffc107', stationId, 'suspend')); // Yellow with suspend action
                }

                if (status === 2) {
                    const suspendedContainer = document.createElement('div');
                    suspendedContainer.classList.add('suspended-icons-container');
                    suspendedContainer.appendChild(createIcon('stop', '#dc3545', stationId, 'stop', 'stop-icon')); // Red with stop action - 60% height
                    suspendedContainer.appendChild(createIcon('restart', '#17a2b8', stationId, 'start', 'restart-icon')); // Cyan with start action - 40% height
                    block_status.appendChild(suspendedContainer);
                }
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
    }, 300000); // 5 minutes = 300,000 milliseconds
    
    console.log('Automatic status updates enabled (every 5 minutes)');
}
