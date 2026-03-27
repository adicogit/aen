
let translations = null;
window.translations = null; // Make translations globally accessible

async function loadTranslations() {
  try {
    const response = await fetch('../translation/translate.json');
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const jsonObject = await response.json();
    
    console.log('Translations loaded from file:', jsonObject);
    
    return jsonObject;

  } catch (error) {
    console.error('Error while loading translations:', error);
    throw error;
  }
}

async function translatePage() {
    if (!translations) {
        translations = await loadTranslations();
        window.translations = translations; // Update global reference
    }
    
    const userLang = navigator.language || navigator.userLanguage;
    const langCode = userLang.split('-')[0];

    const languageData = translations[langCode] || translations['en'];

    document.querySelectorAll('[data-translate]').forEach(element => {
        const key = element.getAttribute('data-translate');
        if (languageData[key]) {
            element.textContent = languageData[key];
        }
    });

    // Handle alt attribute translations
    document.querySelectorAll('[data-translate-alt]').forEach(element => {
        const key = element.getAttribute('data-translate-alt');
        if (languageData[key]) {
            element.alt = languageData[key];
        }
    });

    // Handle title attribute translations (tooltips)
    document.querySelectorAll('[data-translate-title]').forEach(element => {
        const key = element.getAttribute('data-translate-title');
        if (languageData[key]) {
            element.title = languageData[key];
        }
    });

    document.documentElement.lang = langCode;
}

// Function to translate dynamically created elements
async function translateDynamicContent() {
    await translatePage();
}

// Helper function to get translated text with fallback (synchronous)
function getTranslation(key, fallback = '') {
    try {
        // Check if translations exist and are not a Promise
        if (!window.translations || typeof window.translations.then === 'function') {
            console.warn('Translations not loaded yet, using fallback');
            return fallback || key;
        }
        
        const userLang = navigator.language || navigator.userLanguage;
        const langCode = userLang.split('-')[0];
        const languageData = window.translations[langCode] || window.translations['en'];
        
        if (!languageData) {
            console.warn('Language data not found, using fallback');
            return fallback || key;
        }
        
        const translation = languageData[key];
        return translation || fallback || key;
    } catch (error) {
        console.error('Translation error:', error);
        return fallback || key;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    translatePage().catch(error => {
        console.error('Failed to translate page:', error);
    });
});
