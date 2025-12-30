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
        // block.textContent = gameStation.name; // Removed: will be added as a separate element
        block.style.width = `${size}px`;
        block.style.height = `${size}px`;
        // Block icon must be retrieved using API
        if (sanitizedUrl) {
            block.style.backgroundImage = `url('${sanitizedUrl}')`;
        }

        const block_top = document.createElement('div');
        block_top.classList.add('block_station_top');
        block_top.textContent = "prova";

        const block_bottom = document.createElement('div');
        block_bottom.classList.add('block_station_bottom');
        block_bottom.textContent = "prova sotto";

        // Add dedicated name element
        const block_name = document.createElement('div');
        block_name.classList.add('station-name');
        block_name.textContent = gameStation.name;

        // Add two part to main block
        block.appendChild(block_top);
        block.appendChild(block_bottom);
        block.appendChild(block_name); // Append name last to be on top

        container.appendChild(block);
    }
}

async function generateBlocks() {
    generateGamingBlocks()
}
