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
            background_image_folder: data.background_image_folder || '/images/background',
            background_image: data.background_image || 'background.webp',
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
        updateBackgroundPreview(currentGeneralSettings.background_image_folder, currentGeneralSettings.background_image);
        
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

function updateBackgroundPreview(path, fileName) {
    const previewImg = document.getElementById('backgroundPreviewImg');
    const url = path && fileName ? `${path}/${fileName}` : '';
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

    // Background chooser modal logic
    const urlInput = document.getElementById('backgroundImageUrl');
    const chooseBtn = document.getElementById('backgroundChooseBtn');
    const bgModal = document.getElementById('backgroundModal');
    const bgList = document.getElementById('backgroundList');
    const closeBgModal = document.getElementById('closeBackgroundModal');

    async function loadBackgroundImages() {
        try {
            const resp = await fetch('/api/v1/backgrounds');
            if (!resp.ok) throw new Error('Failed to load backgrounds');
            const files = await resp.json();
            // Clear
            bgList.innerHTML = '';
            files.forEach(file => {
                const thumbWrap = document.createElement('div');
                thumbWrap.style.width = '120px';
                thumbWrap.style.cursor = 'pointer';
                thumbWrap.style.textAlign = 'center';

                const img = document.createElement('img');
                img.src = `${currentGeneralSettings.background_image_folder}/${file}`;
                img.alt = file;
                img.style.maxWidth = '100%';
                img.style.height = '80px';
                img.style.objectFit = 'cover';
                img.style.border = '1px solid #ddd';
                img.style.padding = '4px';

                const label = document.createElement('div');
                label.textContent = file;
                label.style.fontSize = '12px';
                label.style.marginTop = '4px';

                thumbWrap.appendChild(img);
                thumbWrap.appendChild(label);

                thumbWrap.addEventListener('click', () => {
                    // set only filename in the input
                    if (urlInput) urlInput.value = file;
                    // update preview using folder + filename
                    updateBackgroundPreview(currentGeneralSettings.background_image_folder, file);
                    // close modal
                    if (bgModal) bgModal.classList.remove('show');
                });

                bgList.appendChild(thumbWrap);
            });
        } catch (err) {
            console.error('Unable to load background images', err);
        }
    }

    if (chooseBtn && bgModal && bgList) {
        chooseBtn.addEventListener('click', async () => {
            // ensure we have the latest folder value
            await fetchGeneralSettings();
            await loadBackgroundImages();
            bgModal.classList.add('show');
        });
    }

    if (closeBgModal && bgModal) {
        closeBgModal.addEventListener('click', () => {
            bgModal.classList.remove('show');
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
                    // build full path if user provided only filename
                    const bgPath = bgUrl.includes('/') ? bgUrl : `${currentGeneralSettings.background_image_folder}/${bgUrl}`;
                    if (/^[a-zA-Z0-9\/._-]+$/.test(bgPath)) {
                        document.body.style.backgroundImage = `url('${bgPath}')`;
                    }
                }

                // Update cached settings (preserve folder)
                currentGeneralSettings = Object.assign({}, currentGeneralSettings, {
                    background_image: bgUrl,
                    theme: selectedTheme,
                    billiard_room_name: roomName
                });
                
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
            updateBackgroundPreview(currentGeneralSettings.background_image_folder, currentGeneralSettings.background_image);
            
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
