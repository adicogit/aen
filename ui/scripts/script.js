async function getGameStationIDList(){
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
    }
    return list;
}

async function getGameStation(id){
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
    }
    return gameStation;
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
    
    for (let i = 0; i < numBlocksInput; i++) {
        let id = gameStationIDList[i]
        let gameStation = await getGameStation(id)
        let imageUrl = gameStation.iconPath
        const block = document.createElement('div');
        block.classList.add('block_station');
        block.textContent = "postazione " + i;
        block.style.width = `${size}px`;
        block.style.height = `${size}px`;
        // Block icon must be retrieved using API
        block.style.backgroundImage = `url('${imageUrl}')`; 

        const block_top = document.createElement('div');
        block_top.classList.add('block_station_top');
        block_top.textContent = "prova";

        const block_bottom = document.createElement('div');
        block_bottom.classList.add('block_station_bottom');
        block_bottom.textContent = "prova sotto";

        // Add two part to main block
        block.appendChild(block_top);
        block.appendChild(block_bottom);

        container.appendChild(block);
    }
}

async function generateConfigBlocks() {
    const gameStationIDList = await getGameStationIDList();
    const numBlocksInput = parseInt(gameStationIDList.length);
    const container = document.getElementById('configBlocksContainer');
    
    container.innerHTML = '';

    // Calculate grid to fit all blocks as squares (including the + block)
    const totalBlocks = numBlocksInput + 1;
    const cols = Math.ceil(Math.sqrt(totalBlocks));
    const rows = Math.ceil(totalBlocks / cols);
    
    // Calculate square size based on available space with padding
    const availableWidth = window.innerWidth - 20;
    const availableHeight = window.innerHeight - 20;
    
    let size = Math.floor(Math.min(availableWidth / cols, availableHeight / rows)) - 10;
    
    if (size < 100) size = 100; // Minimum size for usability
    
    for (let i = 0; i < numBlocksInput; i++) {
        let id = gameStationIDList[i];
        let gameStation = await getGameStation(id);
        let imageUrl = gameStation.iconPath;
        
        const block = document.createElement('div');
        block.classList.add('block_station');
        block.style.width = `${size}px`;
        block.style.height = `${size}px`;
        block.style.backgroundImage = `url('${imageUrl}')`;

        const block_top = document.createElement('div');
        block_top.classList.add('block_station_top');
        block_top.textContent = `Config Station ${i + 1}`;

        const block_bottom = document.createElement('div');
        block_bottom.classList.add('block_station_bottom');
        block_bottom.textContent = 'Settings';

        block.appendChild(block_top);
        block.appendChild(block_bottom);
        container.appendChild(block);
    }
    
    // Additional block with "+" sign
    const addBlock = document.createElement('div');
    addBlock.classList.add('block_station');
    addBlock.style.width = `${size}px`;
    addBlock.style.height = `${size}px`;
    addBlock.textContent = '+';
    addBlock.style.display = 'flex';
    addBlock.style.alignItems = 'center';
    addBlock.style.justifyContent = 'center';
    addBlock.style.fontSize = `${size * 0.5}px`;
    addBlock.setAttribute('data-translate-title', 'add_station_button');
    addBlock.setAttribute('title', 'Add new station');
    container.appendChild(addBlock);
    
    // Translate dynamically created content
    if (typeof translateDynamicContent === 'function') {
        translateDynamicContent();
    }
}

async function generateBlocks() {
    generateGamingBlocks()
    generateConfigBlocks()
}

document.addEventListener('DOMContentLoaded', () => {
    generateBlocks();
});

window.addEventListener('resize', () => {
    generateBlocks();
});

document.getElementById('Config').addEventListener('click', function() {
    const container = document.getElementById('container');
    container.classList.toggle('flipped');
    console.info('flip page');
});

// Set backgound image on page load
async function setbackgroundImage() {
    // Specify URL to be used to load UI config
    const API_URL = '/api/v1/uiconfig'; 

    try {
        // GET needed info
        const response = await fetch(API_URL);

        if (!response.ok) {
            throw new Error(`Error HTTP! Status: ${response.status}`);
        }

        // parse received JSON
        const data = await response.json();

        // Get background path from received json
        const imageUrl = data.background_image; 

        if (imageUrl) {
            // Set background image using CSS
            document.body.style.backgroundImage = `url('${imageUrl}')`;

            console.log(`Background set using: ${imageUrl}`);
        } else {
            console.error("Retrieved JSON does not contain valid 'background_image' field.");
        }

    } catch (error) {
        console.error("ERROR in loading background image:", error);
        // Use default image
        document.body.style.backgroundColor = `url('/images/background/background.webp')`; 
    }
}

// invoke backgroun image set on page loading
setbackgroundImage();
