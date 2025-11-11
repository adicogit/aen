function generateBlocks() {
    /* Number of blocks depends on number of gaming station using API*/
    const numBlocksInput = 12;
    const container = document.getElementById('container');
    const n = parseInt(numBlocksInput.value);
    
    container.innerHTML = '';

    const aspectRatio = window.innerWidth / window.innerHeight;
    const cols = Math.ceil(Math.sqrt(numBlocksInput * aspectRatio));
    const rows = Math.ceil(Math.sqrt(numBlocksInput / aspectRatio));
    
    let size = Math.floor(Math.min(window.innerWidth / cols, window.innerHeight / rows)) - 10;
    
    if (size < 20) size = 20;
    const iconFiles = ['billiard_logo_1.jpeg','billiard_logo_2.jpeg','billiard_logo_3.jpeg','billiard_logo_4.jpeg','billiard_logo_5.jpeg','billiard_logo_6.jpeg','billiard_logo_7.jpeg',
                    'billiard_logo_8.jpeg','billiard_logo_9.jpeg','billiard_logo_10.jpeg','card_logo_1.jpeg','card_logo_2.jpeg','card_logo_3.jpeg','card_logo_4.jpeg','card_logo_5.jpeg',
                    'card_logo_6.jpeg','card_logo_7.jpeg','card_logo_8.jpeg','card_logo_9.jpeg','card_logo_2.jpeg','ps_logo_1.jpeg','ps_logo_2.jpeg','ps_logo_3.jpeg','ps_logo_4.jpeg',
                    'ps_logo_5.jpeg','ps_logo_6.jpeg','ps_logo_7.jpeg']
    for (let i = 1; i <= numBlocksInput; i++) {
        let imageUrl = 'images/icons/' + iconFiles[Math.floor(Math.random() * iconFiles.length)]
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