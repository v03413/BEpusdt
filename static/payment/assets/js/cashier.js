(function () {
    'use strict';

    let i18nInitialized = false;
    let paymentMethods = [];
    let selectedCurrency = '';
    let selectedNetwork = '';
    let selectedMethod = null;
    let paymentConfig = {};

    // DOM Elements (initialized in init)
    let currencySelectEl = null;
    let networkSelectEl = null;
    let payButton = null;
    let amountValue = null;
    let networkBadge = null;

    // countdown
    let countdownTimer = null;
    let statusCheckTimer = null;
    let totalSeconds = 0;

    let tradeId = '';

    function initI18n() {
        if (i18nInitialized) {
            return Promise.resolve();
        }

        let currentLang = localStorage.getItem('payment_language') || navigator.language.split('-')[0] || 'en';
        if (currentLang !== 'zh' && currentLang !== 'en') {
            currentLang = 'en';
        }

        return new Promise(function (resolve, reject) {
            if (typeof i18next === 'undefined') {
                console.error('i18next is not loaded');
                resolve(); // Fail gracefully
                return;
            }

            i18next.init({
                lng: currentLang,
                debug: false,
                resources: {}
            }, function (err, t) {
                if (err) {
                    console.error('i18next initialization failed:', err);
                    resolve(); // Fail gracefully
                    return;
                }

                fetch('/payment/assets/locales/' + currentLang + '.json')
                    .then(function (r) {
                        return r.json();
                    })
                    .then(function (translations) {
                        i18next.addResourceBundle(currentLang, 'translation', translations);
                        i18nInitialized = true;
                        updateContent();

                        // Update language switcher if it exists
                        const languageSwitcher = document.getElementById('languageSwitcher');
                        if (languageSwitcher) {
                            languageSwitcher.value = currentLang;
                        }
                        resolve();
                    })
                    .catch(function (error) {
                        console.error('Failed to load language file:', error);
                        resolve(); // Fail gracefully
                    });
            });
        });
    }

    function changeLanguage(lang) {
        if (!i18next) return;

        const hasLanguage = i18next.hasResourceBundle(lang, 'translation');

        if (!hasLanguage) {
            fetch('/payment/assets/locales/' + lang + '.json')
                .then(function (r) {
                    return r.json();
                })
                .then(function (translations) {
                    i18next.addResourceBundle(lang, 'translation', translations);
                    return i18next.changeLanguage(lang);
                })
                .then(function () {
                    localStorage.setItem('payment_language', lang);
                    updateContent();
                })
                .catch(function (error) {
                    console.error('Failed to load language file:', error);
                });
        } else {
            i18next.changeLanguage(lang, function (err, t) {
                if (err) {
                    console.error('Language change failed:', err);
                    return;
                }
                localStorage.setItem('payment_language', lang);
                updateContent();
            });
        }
    }

    function getPaymentConfig() {
        return paymentConfig;
    }

    function replacePlaceholders(text) {
        const config = getPaymentConfig();

        let token = (config.token || '--').toUpperCase();
        let network = selectedMethod ? (selectedMethod.network || '--').toUpperCase() : '--';
        let networkName = selectedMethod ? (selectedMethod.token_net_name || '--').toUpperCase() : '--';

        if (selectedMethod) {
            token = selectedMethod.currency.toUpperCase();
            network = selectedMethod.network.toUpperCase();
            networkName = selectedMethod.token_net_name.toUpperCase();
        } else if (selectedCurrency) {
            token = selectedCurrency.toUpperCase();
        }

        return text
            .replace(/\{\{token\}\}/g, token)
            .replace(/\{\{network\}\}/g, network)
            .replace(/\{\{networkName\}\}/g, networkName);
    }

    function updateContent() {
        if (!i18next) return;

        const elements = document.querySelectorAll('[data-i18n]');
        elements.forEach(function (element) {
            // Special handling for dropdowns to preserve structure
            if (element.classList.contains('coin-select')) {
                const sp = element.querySelector('span');
                if (sp) {
                    const key = element.getAttribute('data-i18n');
                    let translation = i18next.t(key);
                    sp.textContent = translation;
                }
                return;
            }

            // Protect dropdowns from internal structure destruction if not caught above
            if (element.querySelector('.options-list')) {
                return;
            }

            const key = element.getAttribute('data-i18n');

            if (key.startsWith('[')) {
                const matches = key.match(/\[(.+?)\](.+)/);
                if (matches) {
                    const attr = matches[1];
                    const transKey = matches[2];
                    let translation = i18next.t(transKey);
                    translation = replacePlaceholders(translation);

                    if (attr === 'html') {
                        element.innerHTML = translation;
                    } else {
                        element.setAttribute(attr, translation);
                    }
                }
            } else {
                let translation = i18next.t(key);
                translation = replacePlaceholders(translation);

                if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
                    element.placeholder = translation;
                } else {
                    element.innerHTML = translation;
                }
            }
        });
    }

    function t(key, defaultValue) {
        if (i18nInitialized && i18next) {
            let translation = i18next.t(key);
            return replacePlaceholders(translation);
        }
        return defaultValue || key;
    }

    function setupDropdown(el, onSelectCallback) {
        if (!el) return;

        el.addEventListener('click', function (e) {
            e.stopPropagation();
            const otherSelects = document.querySelectorAll('.coin-select');
            otherSelects.forEach(function (other) {
                if (other !== el) other.classList.remove('active');
            });
            el.classList.toggle('active');
        });

        // The logic for populating the list is in renderOptions,
        // but the event listener for global click is needed once.
    }

    function initSelection() {
        if (!paymentMethods || paymentMethods.length === 0) return;

        // Extract unique currencies
        const currencies = [...new Set(paymentMethods.map(function (m) {
            return m.currency;
        }))];

        // Render Currency Options
        renderOptions(currencySelectEl, currencies.map(function (c) {
            return {value: c, label: c, badge: ''};
        }), function (val) {
            selectedCurrency = val;
            updateNetworkOptions();
            // We need to update UI for crypto-icon when currency changes
            updateUI();
        });

        // Default Select First Currency removed to allow "Please select" state
        // Initial UI update to ensure consistent state
        updateUI();
    }

    function updateNetworkOptions() {
        const networks = paymentMethods
            .filter(function (m) {
                return m.currency === selectedCurrency;
            })
            .map(function (m) {
                return {
                    value: m.token_net_name,
                    label: m.token_net_name.toUpperCase(),
                    fullData: m
                };
            });

        renderOptions(networkSelectEl, networks, function (val, item) {
            selectedNetwork = val;
            selectedMethod = item.fullData;
            updateUI();
            updateContent();
        });

        // Reset selection but do NOT auto-select
        selectedNetwork = '';
        selectedMethod = null;

        // Reset dropdown UI
        if (networkSelectEl) {
            const sp = networkSelectEl.querySelector('span');
            if (sp) {
                sp.setAttribute('data-i18n', 'payment.selectNetwork');
                // Let updateContent() handle the text update
            }
        }

        updateUI();
        updateContent();
    }

    function renderOptions(parent, items, onSelect) {
        if (!parent) return;

        let list = parent.querySelector('.options-list');
        if (!list) {
            list = document.createElement('div');
            list.className = 'options-list';
            parent.appendChild(list);
        }

        list.innerHTML = '';

        items.forEach(function (item) {
            const itemEl = document.createElement('div');
            itemEl.className = 'option-item';
            itemEl.setAttribute('data-value', item.value);

            let html = '<span class="option-text">' + item.label + '</span>';
            if (item.badge) {
                html += '<span class="option-badge">' + item.badge + '</span>';
            }
            itemEl.innerHTML = html;

            itemEl.addEventListener('click', function (e) {
                e.stopPropagation();
                parent.classList.remove('active');
                selectOption(parent, item.value);
                if (onSelect) onSelect(item.value, item);
            });

            list.appendChild(itemEl);
        });
    }

    function selectOption(parent, value) {
        if (!parent) return;

        const items = parent.querySelectorAll('.option-item');
        items.forEach(function (item) {
            item.classList.remove('selected');
        });

        let selected = null;
        items.forEach(function (item) {
            if (item.getAttribute('data-value') === value) selected = item;
        });

        if (selected) {
            selected.classList.add('selected');
            const text = selected.querySelector('.option-text').textContent;
            const sp = parent.querySelector('span');
            if (sp) {
                sp.textContent = text || value;
                // Avoid i18n overwriting the selected value
                sp.removeAttribute('data-i18n');
            }
        }
    }

    function updateUI() {
        // Update crypto-icon, default to '--' if no selection
        const cryptoIcon = document.querySelector('.crypto-icon');
        if (cryptoIcon) {
            cryptoIcon.textContent = selectedCurrency ? selectedCurrency.toUpperCase() : '--';
        }

        const subtitleSpan = document.querySelector('.payment-subtitle span[data-i18n]');
        const subtitleP = document.querySelector('.payment-subtitle');
        let networkBadge = document.querySelector('.network-badge');

        if (!selectedMethod) {
            if (payButton) payButton.disabled = true;
            if (amountValue) amountValue.textContent = '---';

            // Revert to subtitle2 (Please select...)
            if (subtitleSpan) {
                subtitleSpan.setAttribute('data-i18n', 'payment.subtitle2');
            }
            // Remove badge if exists
            if (networkBadge) {
                networkBadge.remove();
                networkBadge = null; // Clear reference
            }

            return;
        }

        // Switch to subtitle (Please use...)
        if (subtitleSpan) {
            subtitleSpan.setAttribute('data-i18n', 'payment.subtitleCashier');
        }

        // Inject network badge if missing
        if (!networkBadge && subtitleP) {
            networkBadge = document.createElement('span');
            networkBadge.className = 'network-badge';
            networkBadge.setAttribute('data-i18n', 'payment.networkBadge');
            subtitleP.appendChild(networkBadge);
        }

        // Update Amount
        if (amountValue) {
            amountValue.textContent = selectedMethod.actual_amount + ' ' + selectedMethod.currency;
        }

        // Update Network Badge
        if (networkBadge) {
            networkBadge.textContent = selectedMethod.token_net_name;
        }

        // Enable Pay Button
        if (payButton) {
            payButton.disabled = false;
        }
    }

    function createTransaction() {
        if (!selectedMethod) return;
        if (payButton && payButton.disabled) return;

        fetch('/api/v1/pay/update-order', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                trade_id: tradeId,
                currency: selectedMethod.currency,
                network: selectedMethod.network
            })
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (res) {
                if (res.status_code === 200 && res.data.payment_url) {
                    window.location.href = res.data.payment_url;
                } else {
                    showMessage(res.message || t('payment.createTransactionFailed', 'Failed to create transaction'));
                }
            })
            .catch(function (err) {
                console.error(err);
                showMessage(t('payment.networkError', 'Network error'));
            });
    }

    function showMessage(msg) {
        if (window.layer && window.layer.msg) {
            window.layer.msg(msg);
        } else {
            alert(msg);
        }
    }

    function copyToClipboard(text, element, isButton) {
        if (!text) return;

        navigator.clipboard.writeText(text).then(function () {
            showCopySuccess(element, isButton);
        }).catch(function (err) {
            console.error('复制失败: ', err);
            fallbackCopy(text, element, isButton);
        });
    }

    function showCopySuccess(element, isButton) {
        const t = window.i18next ? window.i18next.t.bind(window.i18next) : function (key) {
            const translations = {
                'payment.copied': '已复制!',
                'payment.copySuccess': '✓ 已复制!'
            };
            return translations[key] || key;
        };

        if (isButton) {
            const originalText = element.textContent;
            element.textContent = t('payment.copied');
            element.style.background = '#48bb78';
            setTimeout(function () {
                element.textContent = originalText;
                element.style.background = '#667eea';
            }, 2000);
        } else {
            const originalText = element.textContent;
            const originalColor = element.style.color;
            element.textContent = t('payment.copySuccess');
            element.style.color = '#48bb78';
            setTimeout(function () {
                element.textContent = originalText;
                element.style.color = originalColor;
            }, 2000);
        }
    }

    function fallbackCopy(text, element, isButton) {
        const t = window.i18next ? window.i18next.t.bind(window.i18next) : function (key) {
            const translations = {
                'payment.copyFailed': '复制失败，请手动复制'
            };
            return translations[key] || key;
        };

        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.opacity = '0';
        document.body.appendChild(textArea);
        textArea.select();
        try {
            document.execCommand('copy');
            showCopySuccess(element, isButton);
        } catch (err) {
            alert(t('payment.copyFailed'));
        }
        document.body.removeChild(textArea);
    }

    function copyAmount() {
        const amountEl = document.getElementById('payAmount');
        let text = '';
        if (selectedMethod && selectedMethod.actual_amount) {
            text = selectedMethod.actual_amount;
        } else if (amountEl) {
            // Fallback to text content if no method selected (though usually it's ---)
            const content = amountEl.textContent.trim();
            if (content !== '---') {
                text = content.split(' ')[0];
            }
        }

        if (text) {
            copyToClipboard(text, amountEl, false);
        }
    }

    function startCountdown() {
        if (countdownTimer) clearInterval(countdownTimer);

        const config = getPaymentConfig();
        // Ensure expire is a number
        let expireTime = parseInt(config.expire);
        if (isNaN(expireTime)) {
            expireTime = 0;
        }
        totalSeconds = expireTime;

        const minutesEl = document.getElementById('minutes');
        const secondsEl = document.getElementById('seconds');
        const countdownBanner = document.querySelector('.countdown-banner');

        function updateDisplay() {
            if (!minutesEl || !secondsEl) return;

            const minutes = Math.floor(Math.max(0, totalSeconds) / 60);
            const seconds = Math.max(0, totalSeconds) % 60;

            if (totalSeconds <= 0) {
                minutesEl.textContent = '00';
                secondsEl.textContent = '00';
            } else {
                minutesEl.textContent = minutes.toString().padStart(2, '0');
                secondsEl.textContent = seconds.toString().padStart(2, '0');
            }

            if (countdownBanner) {
                if (totalSeconds <= 300) {
                    countdownBanner.style.background = 'linear-gradient(135deg, #ff4757, #ff3838)';
                }
                if (totalSeconds <= 60) {
                    countdownBanner.style.animation = 'urgentBlink 1s infinite';
                }
            }
        }

        updateDisplay(); // Initial display

        countdownTimer = setInterval(function () {
            totalSeconds--;
            updateDisplay();

            if (totalSeconds <= 0) {
                clearInterval(countdownTimer);
                if (statusCheckTimer) clearInterval(statusCheckTimer);
                showTimeoutMessage();
            }
        }, 1000);
    }

    function showTimeoutMessage() {
        // Prevent multiple overlays
        if (document.querySelector('.status-overlay')) return;

        const overlay = document.createElement('div');
        overlay.className = 'status-overlay';
        overlay.style.cssText = 'position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.8); display: flex; align-items: center; justify-content: center; z-index: 9999;';

        const modal = document.createElement('div');
        modal.style.cssText = 'background: white; padding: 30px; border-radius: 12px; text-align: center; max-width: 400px; width: 90%; box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);';

        const config = getPaymentConfig();

        modal.innerHTML = '<div style="font-size: 48px; margin-bottom: 20px;">⏰</div>' +
            '<h3 style="color: #e53e3e; margin-bottom: 15px;">' + t('payment.paymentTimeout', 'Payment Timeout') + '</h3>' +
            '<p style="color: #666; margin-bottom: 25px; line-height: 1.5;">' +
            t('payment.timeoutMessage', 'Sorry, the payment time has expired.<br>Please restart payment.') +
            '</p>' +
            '<button onclick="location.href=\'' + (config.return_url || '/') + '\'" style="background: #667eea; color: white; border: none; padding: 12px 24px; border-radius: 6px; cursor: pointer; font-size: 14px; margin-right: 10px;">' +
            t('payment.returnToMerchant', 'Return to Merchant') +
            '</button>';

        overlay.appendChild(modal);
        document.body.appendChild(overlay);
    }

    function showWaitingConfirmation(data) {
        if (document.getElementById('waiting-overlay')) return;

        const overlay = document.createElement('div');
        overlay.id = 'waiting-overlay';
        overlay.style.cssText = 'position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.8); display: flex; align-items: center; justify-content: center; z-index: 9999;';

        const modal = document.createElement('div');
        modal.style.cssText = 'background: white; padding: 30px; border-radius: 12px; text-align: center; max-width: 400px; width: 90%; box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);';

        modal.innerHTML = '<div style="font-size: 48px; margin-bottom: 20px; animation: spin 2s linear infinite;">⏳</div>' +
            '<h3 style="color: #667eea; margin-bottom: 15px;">' + t('payment.waitingConfirmation', 'Waiting for Confirmation') + '</h3>' +
            '<p style="color: #666; margin-bottom: 20px; line-height: 1.5;">' +
            '<span style="color: #059669; font-weight: 600; background: linear-gradient(135deg, #ecfdf5, #f0fdf4); padding: 8px 12px; border-radius: 6px; border-left: 3px solid #10b981;">' + t('payment.confirmationMessage', 'Payment detected, confirming...') + '</span><br>' +
            '<small style="color: #999;">' + t('payment.confirmationNote', 'Please wait patiently') + '</small>' +
            '</p>' +
            '<div style="display: flex; justify-content: center; margin-bottom: 20px;">' +
            '<div style="width: 40px; height: 40px; border: 3px solid #f3f3f3; border-top: 3px solid #667eea; border-radius: 50%; animation: spin 1s linear infinite;"></div>' +
            '</div>' +
            '<p style="color: #999; font-size: 12px;">' + t('payment.estimatedTime', 'Estimated time: 1-3 mins') + '</p>';

        const style = document.createElement('style');
        style.textContent = '@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }';
        document.head.appendChild(style);

        overlay.appendChild(modal);
        document.body.appendChild(overlay);

        // Ensure we keep checking
        if (!statusCheckTimer) {
            statusCheckTimer = setInterval(checkPaymentStatus, 5000);
        }
    }

    function showSuccessMessage(data) {
        const waitingOverlay = document.getElementById('waiting-overlay');
        if (waitingOverlay) waitingOverlay.remove();

        // Prevent multiple overlays
        if (document.querySelector('.success-overlay')) return;

        const overlay = document.createElement('div');
        overlay.className = 'success-overlay';
        overlay.style.cssText = 'position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.8); display: flex; align-items: center; justify-content: center; z-index: 9999;';

        const modal = document.createElement('div');
        modal.style.cssText = 'background: white; padding: 30px; border-radius: 12px; text-align: center; max-width: 400px; width: 90%; box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);';

        modal.innerHTML = '<div style="font-size: 48px; margin-bottom: 20px;">✅</div>' +
            '<h3 style="color: #48bb78; margin-bottom: 15px;">' + t('payment.paymentSuccess', 'Payment Success!') + '</h3>' +
            '<p style="color: #666; margin-bottom: 20px; line-height: 1.5;">' +
            t('payment.transactionHash', 'Transaction Hash:') + '<br>' +
            '<code style="background: #e6f0fa; color: #2563eb; padding: 4px 8px; border-radius: 4px; font-size: 12px; word-break: break-all;">' +
            (data.trade_hash || t('payment.processing', 'Processing...')) +
            '</code>' +
            '</p>' +
            '<button onclick="location.href=\'' + (data.return_url || '/') + '\'" style="background: #48bb78; color: white; border: none; padding: 12px 24px; border-radius: 6px; cursor: pointer; font-size: 14px;">' +
            t('payment.returnToMerchant', 'Return to Merchant') +
            '</button>';

        overlay.appendChild(modal);
        document.body.appendChild(overlay);

        if (countdownTimer) clearInterval(countdownTimer);
        if (statusCheckTimer) clearInterval(statusCheckTimer);
    }

    function checkPaymentStatus() {
        if (!tradeId) return;

        fetch('/pay/check-status/' + tradeId)
            .then(function (r) {
                return r.json();
            })
            .then(function (data) {
                // status: 1=pending, 2=success, 3=expired, 5=waiting_confirmation
                if (data.status === 1) {
                    // Continue polling (timer handles generic polling, waiting confirms uses specific)
                    // If not in waiting confirmation mode, we rely on the main timer
                } else if (data.status === 2) {
                    showSuccessMessage(data);
                } else if (data.status === 3) {
                    showTimeoutMessage();
                } else if (data.status === 5) {
                    showWaitingConfirmation(data);
                }
            })
            .catch(function (err) {
                console.error('Check status error:', err);
            });
    }

    function init(config) {
        paymentConfig = config || {};
        tradeId = paymentConfig.trade_id;

        // Initialize DOM Elements
        currencySelectEl = document.getElementById('currency-select');
        networkSelectEl = document.getElementById('network-select');
        payButton = document.querySelector('.pay-button');
        amountValue = document.querySelector('.amount-value');
        networkBadge = document.querySelector('.network-badge');

        // Close dropdowns when clicking outside
        document.addEventListener('click', function () {
            const selects = document.querySelectorAll('.coin-select');
            selects.forEach(function (el) {
                el.classList.remove('active');
            });
        });

        setupDropdown(currencySelectEl);
        setupDropdown(networkSelectEl);

        if (payButton) {
            payButton.addEventListener('click', createTransaction);
        }

        initI18n().then(function () {
            if (paymentConfig.network && Array.isArray(paymentConfig.network) && paymentConfig.network.length > 0) {
                paymentMethods = paymentConfig.network;
                initSelection();
            } else {
                showMessage(t('payment.loadPaymentNetworkFailed', 'Failed to load payment network, please check if the wallet has been added'));
            }

            startCountdown();
            // Poll status every 5s
            statusCheckTimer = setInterval(checkPaymentStatus, 5000);
            // Initial check
            checkPaymentStatus();
        });
    }

    // Expose necessary functions
    window.copyAmount = copyAmount;
    window.changeLanguage = changeLanguage;
    window.Payment = {
        init: init
    };
})();
