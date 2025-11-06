
async function readJsonFile(fileName) {
  try {
    const data = await fetch(fileName);
    
    const jsonObject = JSON.parse(data);
    
    console.log('Data read from file:', jsonObject);
    
    return jsonObject;

  } catch (error) {
    console.error('Error while parsing json file:', error);
    throw error; 
  }
}

//translations = readJsonFile('translation/translate.json');
translations = {
  "it": {
    "welcome_title": "Billiard Manager",
    "nav_home": "Home",
    "nav_about": "Chi siamo",
    "nav_services": "Servizi",
    "nav_contacts": "Contatti"
  },
  "en": {
    "welcome_title": "Billiard Manager",
    "nav_home": "Home",
    "nav_about": "About",
    "nav_services": "Services",
    "nav_contacts": "Contact"
  }
}

function translatePage() {
    const userLang = navigator.language || navigator.userLanguage;
    const langCode = userLang.split('-')[0];

    const languageData = translations[langCode] || translations['en'];

    document.querySelectorAll('[data-translate]').forEach(element => {
        const key = element.getAttribute('data-translate');
        if (languageData[key]) {
            element.textContent = languageData[key];
        }
    });

    document.documentElement.lang = langCode;
}

document.addEventListener('DOMContentLoaded', translatePage);
