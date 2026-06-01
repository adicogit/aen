// --- CASH REGISTER LOGIC ---
const cashBtn = document.getElementById('Cash');
const cashDrawer = document.getElementById('cashDrawer');
const closeCashDrawerBtn = document.getElementById('closeCashDrawer');
const cashChecksList = document.getElementById('cashChecksList');
const cashTotalIncoming = document.getElementById('cashTotalIncoming');
const cashTotalExpected = document.getElementById('cashTotalExpected');
const currentAccountingDayContainer = document.getElementById('currentAccountingDayContainer');
const currentAccountingDayDisplay = document.getElementById('currentAccountingDayDisplay');
const changeAccountingDayBtn = document.getElementById('changeAccountingDayBtn');
const accountingDayPicker = document.getElementById('accountingDayPicker');
const accountingDayDate = document.getElementById('accountingDayDate');
const confirmAccountingDay = document.getElementById('confirmAccountingDay');
const cancelAccountingDay = document.getElementById('cancelAccountingDay');

const cashDrawerBody = document.querySelector('.cash-drawer-body');

let previousAccountingDay = '';

function formatDateForDisplay(dateStr) {
    if (!dateStr || dateStr.length !== 8) return dateStr;
    const year = dateStr.substring(0, 4);
    const month = dateStr.substring(4, 6);
    const day = dateStr.substring(6, 8);
    return `${day}/${month}/${year}`;
}

function dateStrToInputValue(dateStr) {
    if (!dateStr || dateStr.length !== 8) return '';
    const year = dateStr.substring(0, 4);
    const month = dateStr.substring(4, 6);
    const day = dateStr.substring(6, 8);
    return `${year}-${month}-${day}`;
}

function setMaxDate() {
    if (accountingDayDate) {
        const today = new Date();
        const yyyy = today.getFullYear();
        const mm = String(today.getMonth() + 1).padStart(2, '0');
        const dd = String(today.getDate()).padStart(2, '0');
        accountingDayDate.max = `${yyyy}-${mm}-${dd}`;
    }
}

async function checkAccountingDay() {
    if (!currentAccountingDayContainer || !currentAccountingDayDisplay || !accountingDayPicker || !cashDrawerBody) return;

    try {
        const response = await fetch('/api/v1/accountingday');
        const data = await response.json();
        const day = data.currentAccountingDay;

        if (day) {
            previousAccountingDay = day;
            const label = getTranslation('cash_accounting_day', 'Giorno contabile:');
            currentAccountingDayDisplay.textContent = `${label} ${formatDateForDisplay(day)}`;
            currentAccountingDayContainer.style.display = 'flex';
            accountingDayPicker.style.display = 'none';
            cashDrawerBody.style.display = 'flex';
            loadCashChecks();
        } else {
            previousAccountingDay = '';
            currentAccountingDayContainer.style.display = 'none';
            accountingDayPicker.style.display = 'flex';
            cancelAccountingDay.style.display = 'none';
            cashDrawerBody.style.display = 'none';
            setMaxDate();
            accountingDayDate.value = '';
        }
    } catch (error) {
        console.error('Error checking accounting day:', error);
        currentAccountingDayDisplay.textContent = getTranslation('network_error', 'Errore di rete');
        currentAccountingDayContainer.style.display = 'flex';
        accountingDayPicker.style.display = 'none';
        cashDrawerBody.style.display = 'flex';
        loadCashChecks();
    }
}

function showAccountingDayPicker() {
    cancelAccountingDay.style.display = 'inline-block';
    setMaxDate();
    if (previousAccountingDay) {
        accountingDayDate.value = dateStrToInputValue(previousAccountingDay);
    }
    cashDrawerBody.style.display = 'none';
    currentAccountingDayContainer.style.display = 'none';
    accountingDayPicker.style.display = 'flex';
}

function hideAccountingDayPicker() {
    accountingDayPicker.style.display = 'none';
    if (previousAccountingDay) {
        currentAccountingDayContainer.style.display = 'flex';
        cashDrawerBody.style.display = 'flex';
    } else {
        currentAccountingDayContainer.style.display = 'none';
        cashDrawerBody.style.display = 'none';
    }
}

async function hasOpenChecks() {
    try {
        const response = await fetch('/api/v1/checks');
        if (!response.ok) return false;
        const data = await response.json();
        return data.open && data.open.length > 0;
    } catch {
        return false;
    }
}

async function handleSetAccountingDay() {
    if (!accountingDayDate) return;

    const dateValue = accountingDayDate.value;
    if (!dateValue) return;

    const dateStr = dateValue.replace(/-/g, '');

    if (previousAccountingDay) {
        const open = await hasOpenChecks();
        if (open) {
            const msg = getTranslation('cash_close_checks_first', 'Chiudi prima i conti aperti per cambiare giorno contabile');
            if (typeof showErrorModal === 'function') {
                await showErrorModal(msg);
            } else {
                alert(msg);
            }
            return;
        }
    }

    try {
        const response = await fetch('/api/v1/accountingday', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ date: dateStr })
        });

        if (!response.ok) {
            const errData = await response.json();
            throw new Error(errData.error || `HTTP error! status: ${response.status}`);
        }

        await checkAccountingDay();
    } catch (error) {
        console.error('Error setting accounting day:', error);
        const errorMsg = getTranslation('api_error', 'Errore API');
        if (typeof showErrorModal === 'function') {
            await showErrorModal(`${errorMsg}: ${error.message}`);
        } else {
            alert(`${errorMsg}: ${error.message}`);
        }
    }
}

if (cashBtn && cashDrawer) {
    cashBtn.addEventListener('click', () => {
        cashDrawer.classList.toggle('open');
        if (cashDrawer.classList.contains('open')) {
            checkAccountingDay();
        }
    });
}

if (closeCashDrawerBtn && cashDrawer) {
    closeCashDrawerBtn.addEventListener('click', () => {
        cashDrawer.classList.remove('open');
    });
}

if (changeAccountingDayBtn) {
    changeAccountingDayBtn.addEventListener('click', showAccountingDayPicker);
}

if (confirmAccountingDay) {
    confirmAccountingDay.addEventListener('click', handleSetAccountingDay);
}

if (cancelAccountingDay) {
    cancelAccountingDay.addEventListener('click', hideAccountingDayPicker);
}

if (accountingDayDate) {
    accountingDayDate.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            handleSetAccountingDay();
        }
    });
    accountingDayDate.addEventListener('change', () => {
        handleSetAccountingDay();
    });
}

// Close drawer if clicking outside
document.addEventListener('click', (e) => {
    if (cashDrawer && cashDrawer.classList.contains('open')) {
        if (!cashDrawer.contains(e.target) && e.target !== cashBtn && !cashBtn.contains(e.target)) {
            if (!e.target.closest('.modal') && !e.target.closest('#confirmModal')) {
                cashDrawer.classList.remove('open');
            }
        }
    }
});

async function loadCashChecks() {
    if (!cashChecksList) return;

    cashChecksList.innerHTML = '<div class="no-checks-message">Caricamento in corso...</div>';

    try {
        const response = await fetch('/api/v1/checks');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();

        if (cashTotalIncoming) cashTotalIncoming.textContent = formatPrice(data.currentIncoming);
        if (cashTotalExpected) cashTotalExpected.textContent = formatPrice(data.currentExpected);

        cashChecksList.innerHTML = '';

        const openChecks = data.open || [];
        const closedChecks = data.closed || [];

        if (openChecks.length === 0 && closedChecks.length === 0) {
            const noChecksText = getTranslation('cash_no_checks', 'Nessun conto per oggi');
            cashChecksList.innerHTML = `<div class="no-checks-message">${noChecksText}</div>`;
            return;
        }

        openChecks.forEach(check => {
            renderCheckCard(check, 'open');
        });

        closedChecks.forEach(cc => {
            renderCheckCard(cc.check, 'paid', cc.payed);
        });

    } catch (error) {
        console.error('Error loading checks:', error);
        const errorMsg = getTranslation('network_error', 'Errore di rete');
        cashChecksList.innerHTML = `<div class="no-checks-message" style="color: #ef4444;">${errorMsg}: ${error.message}</div>`;
    }
}

function renderCheckCard(check, status, payedAmount = null) {
    const card = document.createElement('div');
    card.classList.add('check-card');

    const statusText = status === 'open' ? getTranslation('cash_status_open', 'Aperto') : getTranslation('cash_status_paid', 'Pagato');
    const badgeClass = status === 'open' ? 'open' : 'paid';

    let itemsHtml = '';
    if (check.ItemList && check.ItemList.length > 0) {
        const itemsTitle = getTranslation('cash_items', 'Consumazioni');
        itemsHtml = `
            <div class="check-items-section">
                <div class="check-items-title">${itemsTitle}</div>
        `;
        check.ItemList.forEach(item => {
            itemsHtml += `
                <div class="check-item-row">
                    <span>${item.Name}</span>
                    <span>${formatPrice(item.PublicPrice)}</span>
                </div>
            `;
        });
        itemsHtml += `</div>`;
    }

    const tableLabel = getTranslation('cash_table_name', 'Tavolo');
    const durationLabel = getTranslation('cash_duration', 'Durata');

    card.innerHTML = `
        <div class="check-card-header">
            <div>
                <div class="check-station-name">${check.GameStationName}</div>
                <div class="check-id-sub">ID: ${check.ID.substring(0, 8)}...</div>
            </div>
            <span class="check-status-badge ${badgeClass}">${statusText}</span>
        </div>
        <div class="check-card-details">
            <div class="check-detail-row">
                <span>${tableLabel}:</span>
                <span>${check.GameStationName}</span>
            </div>
            <div class="check-detail-row">
                <span>${durationLabel}:</span>
                <span>${check.Duration} min</span>
            </div>
        </div>
        ${itemsHtml}
        <div class="check-card-footer">
            <div class="check-price-total">
                ${formatPrice(payedAmount !== null ? payedAmount : check.Price)}
            </div>
            ${status === 'open' ? `
                <button class="btn-pay-check" data-id="${check.ID}" data-price="${check.Price}">
                    ${getTranslation('cash_pay_button', 'Paga')}
                </button>
            ` : ''}
        </div>
    `;

    if (status === 'open') {
        const payBtn = card.querySelector('.btn-pay-check');
        if (payBtn) {
            payBtn.addEventListener('click', async (e) => {
                e.stopPropagation();
                const checkId = payBtn.getAttribute('data-id');
                const price = parseInt(payBtn.getAttribute('data-price'));
                await handlePayCheck(checkId, price, check.GameStationName);
            });
        }
    }

    cashChecksList.appendChild(card);
}

async function handlePayCheck(checkId, price, tableName) {
    const confirmMsg = `${getTranslation('cash_pay_confirm', 'Conferma pagamento per ')} ${tableName} (${formatPrice(price)})?`;

    let confirmed = false;
    if (typeof showConfirmModal === 'function') {
        confirmed = await showConfirmModal(confirmMsg);
    } else {
        confirmed = window.confirm(confirmMsg);
    }

    if (!confirmed) return;

    try {
        const response = await fetch(`/api/v1/checks/${checkId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ payment: price })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        await loadCashChecks();

        if (typeof generateBlocks === 'function') {
            generateBlocks();
        }
    } catch (error) {
        console.error('Error paying check:', error);
        const errorMsg = getTranslation('api_error', 'Errore API');
        if (typeof showErrorModal === 'function') {
            await showErrorModal(`${errorMsg}: ${error.message}`);
        } else {
            alert(`${errorMsg}: ${error.message}`);
        }
    }
}
