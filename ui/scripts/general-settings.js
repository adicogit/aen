// general-settings.js

document.addEventListener('DOMContentLoaded', () => {
    initGeneralSettings();
});

// Cache for general settings configuration
let currentGeneralSettings = {
    background_image: '',
    theme: 'light',
    billiard_room_name: ''
};

async function initGeneralSettings() {
    // Load config from API
    await fetchGeneralSettings();
    
    // Setup event listeners for form
    setupGeneralSettingsListeners();
}

async function fetchGeneralSettings() {
    try {
        const response = await fetch('/api/v1/uiconfig');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        
        currentGeneralSettings = {
            background_image: data.background_image || '/images/background/background.webp',
            theme: data.theme || 'light',
            billiard_room_name: data.billiard_room_name || ''
        };
        
        // Populate inputs in HTML
        const nameInput = document.getElementById('billiardRoomName');
        if (nameInput) {
            nameInput.value = currentGeneralSettings.billiard_room_name;
        }
        
        const urlInput = document.getElementById('backgroundImageUrl');
        if (urlInput) {
            urlInput.value = currentGeneralSettings.background_image;
        }
        
        // Preview background
        updateBackgroundPreview(currentGeneralSettings.background_image);
        
        // Update radio select
        const radios = document.getElementsByName('uiTheme');
        radios.forEach(radio => {
            if (radio.value === currentGeneralSettings.theme) {
                radio.checked = true;
            }
        });
        
        // Select cards styling
        updateThemeCardSelection(currentGeneralSettings.theme);

        // Apply theme to document
        applyTheme(currentGeneralSettings.theme);
        
    } catch (error) {
        console.error("Error loading general settings:", error);
    }
}

function applyTheme(theme) {
    if (theme === 'dark') {
        document.body.classList.remove('theme-light');
        document.body.classList.add('theme-dark');
    } else {
        document.body.classList.remove('theme-dark');
        document.body.classList.add('theme-light');
    }
}

function updateThemeCardSelection(selectedTheme) {
    const lightCard = document.getElementById('theme-light-card');
    const darkCard = document.getElementById('theme-dark-card');
    
    if (lightCard && darkCard) {
        if (selectedTheme === 'dark') {
            darkCard.classList.add('selected');
            lightCard.classList.remove('selected');
        } else {
            lightCard.classList.add('selected');
            darkCard.classList.remove('selected');
        }
    }
}

function updateBackgroundPreview(url) {
    const previewImg = document.getElementById('backgroundPreviewImg');
    if (previewImg) {
        if (url) {
            previewImg.src = url;
            previewImg.style.display = 'block';
        } else {
            previewImg.src = '';
            previewImg.style.display = 'none';
        }
    }
}

function setupGeneralSettingsListeners() {
    // Radio buttons event logic
    const radios = document.getElementsByName('uiTheme');
    radios.forEach(radio => {
        radio.addEventListener('change', (e) => {
            const selectedTheme = e.target.value;
            updateThemeCardSelection(selectedTheme);
            // Apply it immediately as a preview
            applyTheme(selectedTheme);
        });
    });
    
    // Also allow clicking anywhere on the label/card to trigger check
    const lightCard = document.getElementById('theme-light-card');
    const darkCard = document.getElementById('theme-dark-card');
    
    if (lightCard) {
        lightCard.addEventListener('click', (e) => {
            // Only trigger if not clicking the radio input directly (since it propagates and handles it)
            if (e.target.tagName !== 'INPUT') {
                const radio = lightCard.querySelector('input[type="radio"]');
                if (radio && !radio.checked) {
                    radio.checked = true;
                    // Trigger change manually
                    const event = new Event('change');
                    radio.dispatchEvent(event);
                }
            }
        });
    }
    
    if (darkCard) {
        darkCard.addEventListener('click', (e) => {
            if (e.target.tagName !== 'INPUT') {
                const radio = darkCard.querySelector('input[type="radio"]');
                if (radio && !radio.checked) {
                    radio.checked = true;
                    const event = new Event('change');
                    radio.dispatchEvent(event);
                }
            }
        });
    }

    // Handle file upload for background
    const fileInput = document.getElementById('backgroundImageFile');
    const urlInput = document.getElementById('backgroundImageUrl');
    
    if (fileInput && urlInput) {
        fileInput.addEventListener('change', (e) => {
            const file = e.target.files[0];
            if (file) {
                const fakePath = `/images/background/${file.name}`;
                urlInput.value = fakePath;
                updateBackgroundPreview(fakePath);
            }
        });
    }
    
    // Save button
    const saveBtn = document.getElementById('saveGeneralSettings');
    if (saveBtn) {
        saveBtn.addEventListener('click', async () => {
            const roomName = document.getElementById('billiardRoomName').value;
            const bgUrl = document.getElementById('backgroundImageUrl').value;
            const selectedRadio = document.querySelector('input[name="uiTheme"]:checked');
            const selectedTheme = selectedRadio ? selectedRadio.value : 'light';
            
            const payload = {
                background_image: bgUrl,
                theme: selectedTheme,
                billiard_room_name: roomName
            };
            
            try {
                const response = await fetch('/api/v1/uiconfig', {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload)
                });
                
                if (!response.ok) {
                    throw new Error(`Failed to save settings: ${response.statusText}`);
                }
                
                // Set body background image immediately if it's correct
                if (bgUrl) {
                    // Validate basic image path structure (mirroring script.js safety check)
                    if (/^[a-zA-Z0-9\/._-]+$/.test(bgUrl)) {
                        document.body.style.backgroundImage = `url('${bgUrl}')`;
                    }
                }
                
                // Update cached settings
                currentGeneralSettings = payload;
                
                // Go back to main panel (flip container back)
                const container = document.getElementById('container');
                if (container) {
                    container.classList.remove('flipped');
                }
                
            } catch (error) {
                console.error("Error saving general settings:", error);
                if (typeof showErrorModal === 'function') {
                    showErrorModal("Errore nel salvataggio delle impostazioni.");
                } else {
                    alert("Errore nel salvataggio delle impostazioni.");
                }
            }
        });
    }
    
    // Cancel button
    const cancelBtn = document.getElementById('cancelGeneralSettings');
    if (cancelBtn) {
        cancelBtn.addEventListener('click', () => {
            // Restore from cached settings
            document.getElementById('billiardRoomName').value = currentGeneralSettings.billiard_room_name;
            document.getElementById('backgroundImageUrl').value = currentGeneralSettings.background_image;
            updateBackgroundPreview(currentGeneralSettings.background_image);
            
            const radios = document.getElementsByName('uiTheme');
            radios.forEach(radio => {
                if (radio.value === currentGeneralSettings.theme) {
                    radio.checked = true;
                }
            });
            updateThemeCardSelection(currentGeneralSettings.theme);
            applyTheme(currentGeneralSettings.theme);
            
            // Go back (flip container)
            const container = document.getElementById('container');
            if (container) {
                container.classList.remove('flipped');
            }
        });
    }
}
