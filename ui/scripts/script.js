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
    const iconFiles = ['billiard_logo_1.jpg','billiard_logo_2.jpg','billiard_logo_3.jpg','billiard_logo_4.jpg','billiard_logo_5.jpg','billiard_logo_6.jpg','billiard_logo_7.jpg',
                    'billiard_logo_8.jpg','billiard_logo_9.jpg','billiard_logo_10.jpg','card_logo_1.jpg','card_logo_2.jpg','card_logo_3.jpg','card_logo_4.jpg','card_logo_5.jpg',
                    'card_logo_6.jpg','card_logo_7.jpg','card_logo_8.jpg','card_logo_9.jpg','card_logo_2.jpg','ps_logo_1.jpg','ps_logo_2.jpg','ps_logo_3.jpg','ps_logo_4.jpg',
                    'ps_logo_5.jpg','ps_logo_6.jpg','ps_logo_7.jpg']
    for (let i = 1; i <= numBlocksInput; i++) {
        let imageUrl = 'images/icons/' + iconFiles[Math.floor(Math.random() * iconFiles.length)]
        const block = document.createElement('div');
        block.classList.add('block');
        block.textContent = i;
        block.style.width = `${size}px`;
        block.style.height = `${size}px`;
        // Block icon must be retrieved using API
        block.style.backgroundImage = `url('${imageUrl}')`; 
        container.appendChild(block);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    generateBlocks();
});

window.addEventListener('resize', () => {
    generateBlocks();
});
