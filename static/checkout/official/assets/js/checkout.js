(function () {
    'use strict';

    var i18nReady = false;
    var lang = 'zh';
    var SVG_C = '<svg width="W" height="W" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>';
    var SVG_K = '<svg width="W" height="W" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>';
    var WEB3 = '/checkout/official/assets/web3icons';
    var iconCache = {};
    var methods = [];
    var selCur = '', selMethod = null;
    var cfg = {}, tradeId = '';
    var cdTimer = null, stTimer = null;

    function detectLang() {
        try {
            return ((navigator.language || navigator.userLanguage || 'en').toLowerCase().indexOf('zh') === 0) ? 'zh' : 'en';
        } catch (e) { return 'en'; }
    }

    function initI18n() {
        lang = detectLang();
        return new Promise(function (resolve) {
            if (typeof i18next === 'undefined') return resolve();
            i18next.init({ lng: lang, debug: false, resources: {} }, function (err) {
                if (err) return resolve();
                fetch('/checkout/official/assets/locales/' + lang + '.json')
                    .then(function (r) { return r.json(); })
                    .then(function (d) {
                        i18next.addResourceBundle(lang, 'translation', d);
                        i18nReady = true;
                        applyI18n();
                        resolve();
                    })
                    .catch(function () { resolve(); });
            });
        });
    }

    function t(key, def) {
        if (i18nReady && typeof i18next !== 'undefined') {
            var v = i18next.t(key);
            if (v && v !== key) return v;
        }
        return def != null ? def : key;
    }

    function applyI18n() {
        if (!i18nReady || typeof i18next === 'undefined') return;
        document.querySelectorAll('[data-i18n]').forEach(function (el) {
            var k = el.getAttribute('data-i18n');
            if (!k) return;
            if (k.charAt(0) === '[') {
                var m = k.match(/\[(.+?)\](.+)/);
                if (!m) return;
                var v = i18next.t(m[2]);
                if (m[1] === 'html') el.innerHTML = v;
                else el.setAttribute(m[1], v);
                return;
            }
            var tr = i18next.t(k);
            if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') el.placeholder = tr;
            else el.innerHTML = tr;
        });
        try {
            document.title = i18next.t('pageTitle');
            document.documentElement.lang = lang === 'zh' ? 'zh-CN' : 'en';
        } catch (e) {}
    }

    function switchLang(l) {
        if (l !== 'zh' && l !== 'en') { console.warn('Use "zh" or "en"'); return; }
        if (typeof i18next === 'undefined') return;
        lang = l;
        fetch('/checkout/official/assets/locales/' + l + '.json')
            .then(function (r) { return r.json(); })
            .then(function (d) {
                i18next.addResourceBundle(l, 'translation', d, true, true);
                return i18next.changeLanguage(l);
            })
            .then(function () { applyI18n(); console.log('✓ Language: ' + l); })
            .catch(function (e) { console.error(e); });
    }

    function showToast(msg, type) {
        var el = document.getElementById('cmusToast');
        if (!el) return;
        el.textContent = msg;
        el.className = 'toast ' + (type || '');
        el.classList.add('show');
        clearTimeout(el._t);
        el._t = setTimeout(function () { el.classList.remove('show'); }, 3000);
    }

    function copyText(text, msg, iconEl, sm) {
        if (!text) return;
        var sz = sm ? 14 : 16;
        var ok = function () {
            showToast(msg || t('toastCopied'), 'success');
            if (iconEl) {
                iconEl.innerHTML = SVG_K.replace(/W/g, sz);
                iconEl.classList.add('copied');
                clearTimeout(iconEl._r);
                iconEl._r = setTimeout(function () {
                    iconEl.innerHTML = SVG_C.replace(/W/g, sz);
                    iconEl.classList.remove('copied');
                }, 1500);
            }
        };
        var fb = function () {
            var ta = document.createElement('textarea');
            ta.value = text;
            ta.style.cssText = 'position:fixed;opacity:0;top:0;left:0;';
            document.body.appendChild(ta);
            ta.select();
            try { document.execCommand('copy'); ok(); }
            catch (e) { showToast(t('toastCopyFailed', '复制失败'), 'error'); }
            document.body.removeChild(ta);
        };
        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).then(ok).catch(fb);
        } else fb();
    }

    function buildDropdown(id, items, onSelect) {
        var w = document.getElementById(id);
        if (!w) return;
        var trigger = w.querySelector('.select-trigger');
        var dd = w.querySelector('.select-dropdown');
        if (!trigger || !dd) return;

        dd.innerHTML =
            '<div class="dropdown-search-wrap">' +
            '<svg class="dropdown-search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>' +
            '<input class="dropdown-search" type="text" placeholder="' + t('searchPlaceholder', '搜索') + '" autocomplete="off">' +
            '</div><div class="dropdown-list"></div>';

        var search = dd.querySelector('.dropdown-search');
        var list = dd.querySelector('.dropdown-list');

        function render(filter) {
            list.innerHTML = '';
            var arr = filter ? items.filter(function (it) { return it.label.toLowerCase().indexOf(filter.toLowerCase()) !== -1; }) : items;
            arr.forEach(function (it) {
                var el = document.createElement('div');
                el.className = 'dropdown-item';
                el.dataset.value = it.value;
                el.innerHTML =
                    '<img src="' + it.iconSrc + '" alt="' + it.label + '" class="item-img" onerror="this.style.visibility=\'hidden\'">' +
                    '<span class="item-label">' + it.label + '</span>' +
                    (it.badge ? '<span class="item-badge">' + it.badge + '</span>' : '') +
                    '<svg class="item-check-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>';
                el.addEventListener('click', function (e) {
                    e.stopPropagation();
                    list.querySelectorAll('.dropdown-item').forEach(function (d) { d.classList.remove('selected'); });
                    el.classList.add('selected');
                    var txt = trigger.querySelector('.select-text');
                    var oldImg = trigger.querySelector('.select-img');
                    if (oldImg) oldImg.remove();
                    var img = document.createElement('img');
                    img.className = 'select-img';
                    img.src = it.iconSrc;
                    img.alt = it.label;
                    img.addEventListener('error', function () { img.style.visibility = 'hidden'; });
                    trigger.insertBefore(img, trigger.firstChild);
                    if (txt) { txt.textContent = it.label; txt.classList.remove('placeholder'); }
                    dd.classList.remove('open');
                    trigger.classList.remove('open');
                    if (onSelect) onSelect(it.value, it);
                });
                list.appendChild(el);
            });
        }

        render('');
        search.addEventListener('input', function () { render(this.value); });
        search.addEventListener('click', function (e) { e.stopPropagation(); });

        if (!trigger.dataset.bound) {
            trigger.dataset.bound = '1';
            trigger.addEventListener('click', function (e) {
                e.stopPropagation();
                var open = dd.classList.contains('open');
                document.querySelectorAll('.select-dropdown').forEach(function (d) { d.classList.remove('open'); });
                document.querySelectorAll('.select-trigger').forEach(function (tt) { tt.classList.remove('open'); });
                if (!open) {
                    dd.classList.add('open');
                    trigger.classList.add('open');
                    var s = dd.querySelector('.dropdown-search');
                    if (s) { s.value = ''; render(''); s.focus(); }
                }
            });
        }
    }

    function resetTrigger(id, ph) {
        var w = document.getElementById(id);
        if (!w) return;
        var trigger = w.querySelector('.select-trigger');
        if (!trigger) return;
        var img = trigger.querySelector('.select-img');
        if (img) img.remove();
        var txt = trigger.querySelector('.select-text');
        if (txt) { txt.textContent = ph; txt.classList.add('placeholder'); }
        var list = w.querySelector('.dropdown-list');
        if (list) list.querySelectorAll('.dropdown-item').forEach(function (d) { d.classList.remove('selected'); });
    }

    function tokenIcon(c) { return WEB3 + '/token/' + (c || '').toUpperCase() + '.svg'; }
    function netIcon(n) { return WEB3 + '/network/' + (n || '').toLowerCase() + '.svg'; }

    function preloadIcons(srcs, cb) {
        var pending = srcs.length;
        if (!pending) return cb();
        srcs.forEach(function (src) {
            if (iconCache[src]) { if (--pending === 0) cb(); return; }
            fetch(src)
                .then(function (r) { return r.ok ? r.blob() : null; })
                .then(function (blob) {
                    if (blob) {
                        var reader = new FileReader();
                        reader.onload = function () { iconCache[src] = reader.result; if (--pending === 0) cb(); };
                        reader.readAsDataURL(blob);
                    } else { iconCache[src] = src; if (--pending === 0) cb(); }
                })
                .catch(function () { iconCache[src] = src; if (--pending === 0) cb(); });
        });
    }
    function cached(src) { return iconCache[src] || src; }

    function initSelectionUI() {
        if (!methods.length) return;
        var currencies = Array.from(new Set(methods.map(function (m) { return m.currency; })));
        var srcs = currencies.map(tokenIcon).concat(methods.map(function (m) { return netIcon(m.network); }));
        preloadIcons(srcs, function () {
            buildDropdown('currencySelect', currencies.map(function (c) {
                return { value: c, label: c, iconSrc: cached(tokenIcon(c)), badge: '' };
            }), function (val) {
                selCur = val;
                selMethod = null;
                resetTrigger('networkSelect', t('selectNetwork', '选择网络'));
                updateNetworkDropdown();
                updateAmount();
                updatePayBtn();
            });
            buildDropdown('networkSelect', [], function () {});
            updateAmount();
            updatePayBtn();
        });
    }

    function updateNetworkDropdown() {
        var arr = methods.filter(function (m) { return m.currency === selCur; });
        buildDropdown('networkSelect', arr.map(function (m) {
            return {
                value: m.token_net_name,
                label: m.token_net_name.toUpperCase(),
                iconSrc: cached(netIcon(m.network)),
                badge: m.is_popular ? t('hotBadge', '热门') : '',
                fullData: m
            };
        }), function (_, item) {
            selMethod = item.fullData;
            updateAmount();
            updatePayBtn();
        });
    }

    function updateAmount() {
        var aEl = document.getElementById('payAmountCrypto');
        var nEl = document.getElementById('payNetworkTag');
        var lineEl = document.getElementById('networkTagSep');
        var rowEl = aEl ? aEl.closest('.amount-crypto-row') : null;
        if (!aEl) return;
        if (selMethod) {
            aEl.textContent = selMethod.actual_amount + ' ' + selMethod.currency;
            if (nEl) nEl.textContent = t('networkPrefix', '区块网络 · ') + selMethod.token_net_name;
            if (lineEl) lineEl.style.display = 'flex';
            if (rowEl) rowEl.style.display = '';
        } else {
            aEl.textContent = '--';
            if (nEl) nEl.textContent = '';
            if (lineEl) lineEl.style.display = 'none';
            if (rowEl) rowEl.style.display = 'none';
        }
    }

    function updatePayBtn() {
        var b = document.getElementById('payBtn');
        if (b) b.disabled = !selMethod;
    }

    function startStatusCheck() {
        if (stTimer) clearInterval(stTimer);
        stTimer = setInterval(checkStatus, 5000);
        checkStatus();
    }

    function checkStatus() {
        if (!tradeId) return;
        fetch('/api/v1/pay/info', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ trade_id: tradeId })
        })
            .then(function (r) { return r.json(); })
            .then(function (res) {
                if (res.status_code !== 200) return;
                var d = res.data;
                if (d.status === 5) showConfirming();
                else if (d.status === 2) {
                    stopTimers();
                    hideConfirming();
                    showSuccess(d);
                } else if (d.status === 4) {
                    showCanceled(d);
                }
            })
            .catch(function (e) { console.error(e); });
    }

    function stopTimers() {
        if (stTimer) {
            clearInterval(stTimer);
            stTimer = null;
        }
        if (cdTimer) {
            clearInterval(cdTimer);
            cdTimer = null;
        }
    }

    function hideCheckoutStages() {
        var selection = document.getElementById('selectionStage');
        var payment = document.getElementById('paymentStage');
        if (selection) selection.style.display = 'none';
        if (payment) payment.style.display = 'none';
    }

    function startCountdown(timerEl, expiredAt, createdAt, onExpire) {
        if (cdTimer) clearInterval(cdTimer);
        var total = expiredAt - createdAt;
        if (total <= 0) total = 1;
        var wrap = timerEl ? timerEl.closest('.timer-badge-wrap') : null;
        var bar = wrap ? wrap.querySelector('.timer-ring-bar') : null;

        function tick() {
            var rem = Math.max(0, expiredAt - Math.floor(Date.now() / 1000));
            var hh = Math.floor(rem / 3600);
            var mm = Math.floor((rem % 3600) / 60);
            var ss = rem % 60;
            if (timerEl) timerEl.textContent = ('0' + hh).slice(-2) + ':' + ('0' + mm).slice(-2) + ':' + ('0' + ss).slice(-2);
            var pct = rem / total;
            if (bar) {
                bar.style.strokeDashoffset = (1 - pct) * 100;
                bar.classList.toggle('danger', pct <= 0.10);
                bar.classList.toggle('warn', pct > 0.10 && pct <= 0.30);
            }
            if (timerEl) {
                timerEl.classList.toggle('danger', pct <= 0.10);
                timerEl.classList.toggle('warn', pct > 0.10 && pct <= 0.30);
            }
            if (rem <= 0) {
                clearInterval(cdTimer);
                if (onExpire) onExpire();
            }
        }
        tick();
        cdTimer = setInterval(tick, 1000);
    }

    function showConfirming() {
        var m = document.getElementById('confirmingModal');
        if (m && m.style.display === 'none') m.style.display = 'flex';
    }
    function hideConfirming() {
        var m = document.getElementById('confirmingModal');
        if (m) m.style.display = 'none';
    }

    function showSuccess(data) {
        var modal = document.getElementById('successModal');
        if (!modal) return;
        modal.style.display = 'flex';
        if (data.trade_url) {
            var sec = document.getElementById('txHashSection');
            if (sec) sec.style.display = 'block';
            var link = document.getElementById('txHashLink');
            if (link) link.href = data.trade_url;
            var disp = document.getElementById('txHashDisplayText');
            if (disp) {
                var u = data.trade_url;
                disp.textContent = u.length > 40 ? u.slice(0, 20) + '...' + u.slice(-16) : u;
            }
        }
        var ret = data.redirect_url || (cfg && cfg.return_url) || '';
        var btn = document.getElementById('returnBtn');
        if (btn) {
            if (ret) btn.href = ret;
            else btn.style.display = 'none';
        }
        showToast(t('paymentSuccessToast', '支付成功'), 'success');
    }

    function showCanceled(data) {
        stopTimers();
        hideConfirming();
        hideCheckoutStages();
        var modal = document.getElementById('canceledModal');
        if (!modal) {
            var ret = (data && data.redirect_url) || (cfg && cfg.return_url) || '/';
            modal = document.createElement('div');
            modal.id = 'canceledModal';
            modal.className = 'modal-overlay';
            modal.innerHTML =
                '<div class="modal-card"><div class="modal-body">' +
                '<div style="margin-bottom:16px;display:flex;justify-content:center;">' +
                    '<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">' +
                        '<circle cx="12" cy="12" r="10"/>' +
                        '<path d="M15 9l-6 6"/>' +
                        '<path d="M9 9l6 6"/>' +
                    '</svg>' +
                '</div>' +
                '<div class="modal-title" style="color:#475569;">' + t('canceledTitle', '订单已取消') + '</div>' +
                '<p class="modal-subtitle">' + t('canceledMessage', '该订单已取消，不能继续付款。<br>如需支付，请重新发起订单。') + '</p>' +
                '<a href="' + ret + '" class="return-btn">' + t('returnBtn', '返回商户平台') + '</a>' +
                '</div></div>';
            document.body.appendChild(modal);
        }
        modal.style.display = 'flex';
        showToast(t('orderCanceledToast', '订单已取消'), 'error');
    }

    function showTimeout() {
        if (document.getElementById('timeoutModal')) return;
        var ret = cfg.return_url || '/';
        var ov = document.createElement('div');
        ov.id = 'timeoutModal';
        ov.className = 'modal-overlay';
        ov.innerHTML =
            '<div class="modal-card"><div class="modal-body">' +
            '<div style="margin-bottom:16px;display:flex;justify-content:center;">' +
                '<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">' +
                    '<circle cx="12" cy="12" r="10"/>' +
                    '<polyline points="12 6 12 12 16 14"/>' +
                '</svg>' +
            '</div>' +
            '<div class="modal-title" style="color:#ef4444;">' + t('timeoutTitle', '支付已超时') + '</div>' +
            '<p class="modal-subtitle">' + t('timeoutMessage', '很抱歉，支付时间已超时。<br>请重新发起支付。') + '</p>' +
            '<a href="' + ret + '" class="return-btn">' + t('returnBtn', '返回商户平台') + '</a>' +
            '</div></div>';
        document.body.appendChild(ov);
    }

    function createTransaction() {
        if (cfg && cfg.status && cfg.status !== 1) {
            showToast(t('toastOrderNotPayable', '当前订单状态不允许继续付款'), 'error');
            return;
        }
        if (!selMethod) return;
        fetch('/api/v1/pay/update-order', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ trade_id: tradeId, currency: selMethod.currency, network: selMethod.network })
        })
            .then(function (r) { return r.json(); })
            .then(function (res) {
                if (res.status_code === 200 && res.data && res.data.payment_url) {
                    window.location.href = res.data.payment_url;
                } else showToast(res.message || t('toastCreateFailed', '创建交易失败'), 'error');
            })
            .catch(function () { showToast(t('toastNetworkError', '网络错误'), 'error'); });
    }

    window.Payment = {
        init: function (config) {
            cfg = config || {};
            tradeId = cfg.trade_id;
            document.addEventListener('click', function () {
                document.querySelectorAll('.select-dropdown').forEach(function (d) { d.classList.remove('open'); });
                document.querySelectorAll('.select-trigger').forEach(function (x) { x.classList.remove('open'); });
            });
            startCountdown(document.getElementById('timerDisplay'), parseInt(cfg.expired_at) || 0, parseInt(cfg.created_at) || 0, showTimeout);
            startStatusCheck();
            var payBtn = document.getElementById('payBtn');
            if (payBtn && !payBtn.dataset.bound) {
                payBtn.dataset.bound = '1';
                payBtn.addEventListener('click', createTransaction);
            }
            var caBtn = document.getElementById('copyAmountBtn');
            var caIcon = document.getElementById('copyAmountIcon');
            if (caBtn) caBtn.addEventListener('click', function () {
                if (selMethod) copyText(selMethod.actual_amount, t('toastAmountCopied', '金额已复制'), caIcon, false);
            });
            fetch('/api/v1/pay/methods', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ trade_id: tradeId })
            })
                .then(function (r) { return r.json(); })
                .then(function (res) {
                    if (res.status_code === 200 && res.data && Array.isArray(res.data.methods) && res.data.methods.length) {
                        methods = res.data.methods;
                        initSelectionUI();
                    } else {
                        methods = [];
                        showToast(res.message || t('toastLoadFailed', '无法加载付款网络'), 'error');
                    }
                })
                .catch(function () { showToast(t('toastNetworkError', '网络错误'), 'error'); });
        },
        initQrPage: function (config) {
            cfg = config || {};
            tradeId = cfg.trade_id;
            startCountdown(document.getElementById('timerDisplayQ'), parseInt(cfg.expired_at) || 0, parseInt(cfg.created_at) || 0, showTimeout);
            startStatusCheck();
            var caBtn = document.getElementById('copyAmountQBtn');
            var caIcon = document.getElementById('copyAmountQIcon');
            if (caBtn) caBtn.addEventListener('click', function () {
                if (window._qrAmount) copyText(window._qrAmount, t('toastAmountCopied', '金额已复制'), caIcon, false);
            });
            var adBtn = document.getElementById('copyAddressBtn');
            var adIcon = document.getElementById('copyAddressIcon');
            if (adBtn) adBtn.addEventListener('click', function () {
                var el = document.getElementById('walletAddress');
                if (el && el.textContent !== '--') copyText(el.textContent, t('toastAddressCopied', '地址已复制'), adIcon, true);
            });
        },
        initI18n: initI18n,
        applyI18n: applyI18n,
        t: t,
        showCanceled: showCanceled,
        switchLang: switchLang
    };
})();

(function () {
    var orderData = null, gTradeId = '';
    window._qrAmount = '';

    function _t(k, d) { return (window.Payment && window.Payment.t) ? window.Payment.t(k, d) : d; }

    function showSelection() {
        document.getElementById('selectionStage').style.display = 'block';
        document.getElementById('paymentStage').style.display = 'none';
        var d = orderData;
        document.getElementById('orderMoneyS').textContent = d.money || '--';
        document.getElementById('orderFiatS').textContent = d.fiat || '';
        document.getElementById('orderIdS').textContent = d.order_id || '--';
        bindHelp('helpBtnS', d.support_url);
        Payment.init({
            expired_at: d.expired_at,
            created_at: d.created_at,
            trade_id: d.trade_id,
            return_url: d.redirect_url,
            status: d.status
        });
    }

    function updateReselectButton() {
        var btn = document.getElementById('reselectPaymentBtn');
        if (!btn) return;
        btn.style.display = orderData && orderData.reselect ? 'inline-flex' : 'none';
        if (!btn.dataset.bound) {
            btn.dataset.bound = '1';
            btn.addEventListener('click', function () {
                showSelection();
            });
        }
    }

    function showPayment() {
        document.getElementById('selectionStage').style.display = 'none';
        document.getElementById('paymentStage').style.display = 'block';
        var d = orderData;
        var currency = (d.network && d.network.crypto) ? d.network.crypto : (d.selected_payment ? d.selected_payment.currency : 'USDT');
        var netName = (d.network && d.network.name) ? d.network.name : (d.selected_payment ? d.selected_payment.token_net_name : '');
        var amount = d.actual_amount || '--';
        window._qrAmount = amount;
        document.getElementById('orderMoneyQ').textContent = d.money || '--';
        document.getElementById('orderFiatQ').textContent = d.fiat || '';
        document.getElementById('payAmountQ').textContent = amount + ' ' + currency;
        document.getElementById('payNetworkQ').textContent = _t('networkPrefix', '区块网络 · ') + netName;
        document.getElementById('orderIdQ').textContent = d.order_id || '--';
        document.getElementById('walletAddress').textContent = d.token || d.address || '--';
        var addr = document.getElementById('addressLabelQ');
        if (addr) addr.textContent = _t('receivingAddress', '收款地址');
        $('#qrcode').empty().qrcode({ text: d.token || d.address || '', width: 200, height: 200 });
        updateReselectButton();
        bindHelp('helpBtnQ', d.support_url);
        Payment.initQrPage({
            expired_at: d.expired_at,
            created_at: d.created_at,
            trade_id: d.trade_id,
            return_url: d.redirect_url,
            status: d.status
        });
    }

    function bindHelp(id, url) {
        var el = document.getElementById(id);
        if (!el) return;
        if (url) {
            el.href = url;
            el.removeAttribute('aria-disabled');
            el.style.pointerEvents = '';
            el.style.opacity = '';
        } else {
            el.href = '#';
            el.setAttribute('aria-disabled', 'true');
            el.style.pointerEvents = 'none';
            el.style.opacity = '0.35';
        }
    }

    function loadOrder() {
        fetch('/api/v1/pay/info', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ trade_id: gTradeId })
        })
            .then(function (r) { return r.json(); })
            .then(function (res) {
                if (res.status_code !== 200) {
                    var t = document.getElementById('cmusToast');
                    if (t) {
                        t.textContent = res.message || _t('toastLoadOrderFailed', '加载订单失败');
                        t.className = 'toast error show';
                    }
                    return;
                }
                orderData = res.data;
                if (orderData.status === 4) {
                    if (window.Payment && Payment.showCanceled) Payment.showCanceled(orderData);
                    return;
                }
                if (!orderData.trade_type || !orderData.token) showSelection();
                else showPayment();
            })
            .catch(function (e) {
                var t = document.getElementById('cmusToast');
                if (t) {
                    t.textContent = _t('toastNetworkError', '网络错误') + ': ' + e.message;
                    t.className = 'toast error show';
                }
            });
    }

    document.addEventListener('DOMContentLoaded', function () {
        var parts = window.location.pathname.split('/');
        gTradeId = parts[parts.length - 1];
        if (!gTradeId) return;
        var ready = (window.Payment && Payment.initI18n) ? Payment.initI18n() : Promise.resolve();
        ready.then(loadOrder);
    });
})();
