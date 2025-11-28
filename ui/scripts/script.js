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
}

async function generateGamingBlocks() {
    /* Number of blocks depends on number of gaming station using API*/
    const gameStationIDList = await getGameStationIDList();
    const numBlocksInput = parseInt(gameStationIDList.length);
    const container = document.getElementById('containerGamingStation');
    
    container.innerHTML = '';

    const aspectRatio = window.innerWidth / window.innerHeight;
    const cols = Math.ceil(Math.sqrt(numBlocksInput * aspectRatio));
    const rows = Math.ceil(Math.sqrt(numBlocksInput / aspectRatio));
    
    let size = Math.floor(Math.min(window.innerWidth / cols, window.innerHeight / rows)) - 10;
    
    if (size < 20) size = 20;
    
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

// Seleziona l'icona del menu e il menu stesso tramite i loro ID
const menuIcon = document.getElementById('menu-icon');
const menu = document.getElementById('menu');

// Aggiungi un "ascoltatore di eventi" (event listener) per il click sull'icona
menuIcon.addEventListener('click', function() {
    // Alterna (toggle) la classe 'open' sul menu.
    // Se c'è, la toglie; se non c'è, la mette.
    menu.classList.toggle('open');
});

document.getElementById('Config').addEventListener('click', function() {
    this.classList.toggle('flipped');
    console.info('flip page')
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
