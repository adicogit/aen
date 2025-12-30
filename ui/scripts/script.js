
document.addEventListener('DOMContentLoaded', () => {
    generateBlocks();
});

window.addEventListener('resize', () => {
    generateBlocks();
});

document.getElementById('Config').addEventListener('click', function () {
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
            // Validate imageUrl format before using
            if (!/^[a-zA-Z0-9\/._-]+$/.test(imageUrl)) {
                console.error("Invalid background image URL format. URL rejected for security reasons.");
                return;
            }

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

// Logic for Config Slider
const dots = document.querySelectorAll('.dot');
const panels = document.querySelectorAll('.config-subpanel');

dots.forEach(dot => {
    dot.addEventListener('click', () => {
        const index = dot.getAttribute('data-index');

        // Remove active class from all dots and panels
        dots.forEach(d => d.classList.remove('active'));
        panels.forEach(p => p.classList.remove('active'));

        // Add active class to clicked dot and corresponding panel
        dot.classList.add('active');
        if (panels[index]) {
            panels[index].classList.add('active');
        }
    });
});

