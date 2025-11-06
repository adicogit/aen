function generateBlocks() {
    /* Number of blocks depends on number of gaming station */
    const numBlocksInput = 12;
    const container = document.getElementById('container');
    const n = parseInt(numBlocksInput.value);
    
    container.innerHTML = '';

    const aspectRatio = window.innerWidth / window.innerHeight;
    const cols = Math.ceil(Math.sqrt(numBlocksInput * aspectRatio));
    const rows = Math.ceil(Math.sqrt(numBlocksInput / aspectRatio));
    
    let size = Math.floor(Math.min(window.innerWidth / cols, window.innerHeight / rows)) - 10;
    
    if (size < 20) size = 20;
    /*if (size > 200) size = 200;*/

    for (let i = 1; i <= numBlocksInput; i++) {
        const block = document.createElement('div');
        block.classList.add('block');
        block.textContent = i;
        block.style.width = `${size}px`;
        block.style.height = `${size}px`; 
        container.appendChild(block);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    generateBlocks();
});

window.addEventListener('resize', () => {
    generateBlocks();
});
