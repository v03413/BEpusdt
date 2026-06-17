(function () {
    'use strict';

    var config = {};
    var methods = [];
    var selectedNetworkId = '';
    var selectedCurrency = '';
    var selectedMethod = null;
    var paymentDetail = null;
    var tradeId = '';
    var paymentAmountValue = '';
    var totalSeconds = 0;
    var countdownTimer = null;
    var statusCheckTimer = null;
    var currentLang = 'zh';
    var currentStatusKey = 'status.waitingPayment';

    var translations = {
        zh: {
            'page.title': 'BEpusdt - 收银台',
            'brand.title': 'BEpusdt 收银台',
            'brand.subtitle': '商户支付会话',
            'lang.group': '语言',
            'theme.group': '主题',
            'theme.toLight': '切换浅色主题',
            'theme.toDark': '切换深色主题',
            'checkout.aria': '支付流程',
            'checkout.eyebrow': '安全支付',
            'status.waitingPayment': '等待付款',
            'status.waitingConfirm': '等待链上确认',
            'status.timeout': '订单已超时',
            'status.success': '支付成功',
            'steps.aria': '支付步骤',
            'steps.choose': '选择支付方式',
            'steps.transfer': '扫码或复制转账',
            'order.title': '订单信息',
            'order.merchantOrder': '商户订单',
            'order.product': '商品',
            'order.tradeId': '交易编号',
            'order.amount': '订单金额',
            'method.title': '选择网络与币种',
            'method.network': '网络',
            'method.currency': '币种',
            'method.selectNetwork': '请选择网络',
            'method.selectCurrency': '请选择币种',
            'method.selectMethod': '请选择支付方式',
            'method.availableCurrencies': '可用币种：{{currencies}}',
            'method.warningTitle': '请仔细核对网络和币种',
            'method.warningBody': '转账前请确认钱包中选择的网络和币种与页面完全一致。转错网络或币种可能导致资产永久丢失。',
            'summary.needPay': '需支付',
            'summary.riskNote': '下一步将锁定本次支付方式，并显示专属收款地址和二维码。转账金额必须与页面展示金额一致。',
            'timer.remaining': '剩余时间',
            'actions.next': '下一步',
            'actions.generating': '生成付款信息...',
            'actions.back': '返回上一步',
            'actions.copyAddress': '复制地址',
            'actions.copied': '已复制',
            'transfer.payAmount': '支付金额',
            'transfer.copyAmount': '复制支付金额',
            'transfer.info': '转账信息',
            'transfer.title': '扫码或复制地址完成付款',
            'transfer.address': '收款地址',
            'instruction.networkDefault': '使用当前页面选择的网络',
            'instruction.networkUse': '使用当前页面选择的 {{network}} 网络',
            'instruction.amount': '精准转入金额',
            'instruction.autoConfirm': '转账后保持页面打开，系统会自动确认订单',
            'message.createFailed': '生成付款信息失败',
            'message.networkError': '网络错误，请稍后重试',
            'message.noNetworks': '未加载到可用支付网络，请检查钱包配置。',
            'overlay.timeoutTitle': '支付已超时',
            'overlay.timeoutBody': '请返回商户重新发起订单，超时订单请勿继续转账。',
            'overlay.returnMerchant': '返回商户',
            'overlay.successTitle': '支付成功',
            'overlay.successBody': '订单已确认到账，即将返回商户页面。'
        },
        'zh-Hant': {
            'page.title': 'BEpusdt - 收銀台',
            'brand.title': 'BEpusdt 收銀台',
            'brand.subtitle': '商戶支付會話',
            'lang.group': '語言',
            'theme.group': '主題',
            'theme.toLight': '切換淺色主題',
            'theme.toDark': '切換深色主題',
            'checkout.aria': '支付流程',
            'checkout.eyebrow': '安全支付',
            'status.waitingPayment': '等待付款',
            'status.waitingConfirm': '等待鏈上確認',
            'status.timeout': '訂單已逾時',
            'status.success': '支付成功',
            'steps.aria': '支付步驟',
            'steps.choose': '選擇支付方式',
            'steps.transfer': '掃碼或複製轉帳資訊',
            'order.title': '訂單資訊',
            'order.merchantOrder': '商戶訂單',
            'order.product': '商品',
            'order.tradeId': '交易編號',
            'order.amount': '訂單金額',
            'method.title': '選擇網路與幣種',
            'method.network': '網路',
            'method.currency': '幣種',
            'method.selectNetwork': '請選擇網路',
            'method.selectCurrency': '請選擇幣種',
            'method.selectMethod': '請選擇支付方式',
            'method.availableCurrencies': '可用幣種：{{currencies}}',
            'method.warningTitle': '請仔細核對網路和幣種',
            'method.warningBody': '轉帳前請確認錢包中選擇的網路和幣種與頁面完全一致。轉錯網路或幣種可能導致資產永久遺失。',
            'summary.needPay': '需支付',
            'summary.riskNote': '下一步將鎖定本次支付方式，並顯示專屬收款地址和二維碼。轉帳金額必須與頁面顯示金額一致。',
            'timer.remaining': '剩餘時間',
            'actions.next': '下一步',
            'actions.generating': '正在產生付款資訊...',
            'actions.back': '返回上一步',
            'actions.copyAddress': '複製地址',
            'actions.copied': '已複製',
            'transfer.payAmount': '支付金額',
            'transfer.copyAmount': '複製支付金額',
            'transfer.info': '轉帳資訊',
            'transfer.title': '掃描二維碼或複製地址完成付款',
            'transfer.address': '收款地址',
            'instruction.networkDefault': '使用目前頁面選擇的網路',
            'instruction.networkUse': '使用目前頁面選擇的 {{network}} 網路',
            'instruction.amount': '精準轉入金額',
            'instruction.autoConfirm': '轉帳後請保持頁面開啟，系統會自動確認訂單',
            'message.createFailed': '產生付款資訊失敗',
            'message.networkError': '網路錯誤，請稍後重試',
            'message.noNetworks': '未載入可用支付網路，請檢查錢包設定。',
            'overlay.timeoutTitle': '支付已逾時',
            'overlay.timeoutBody': '請返回商戶重新發起訂單，逾時訂單請勿繼續轉帳。',
            'overlay.returnMerchant': '返回商戶',
            'overlay.successTitle': '支付成功',
            'overlay.successBody': '訂單已確認到帳，即將返回商戶頁面。'
        },
        en: {
            'page.title': 'BEpusdt - Checkout',
            'brand.title': 'BEpusdt Checkout',
            'brand.subtitle': 'Merchant payment session',
            'lang.group': 'Language',
            'theme.group': 'Theme',
            'theme.toLight': 'Switch to light theme',
            'theme.toDark': 'Switch to dark theme',
            'checkout.aria': 'Payment flow',
            'checkout.eyebrow': 'Secure payment',
            'status.waitingPayment': 'Waiting for payment',
            'status.waitingConfirm': 'Waiting for on-chain confirmation',
            'status.timeout': 'Order expired',
            'status.success': 'Payment successful',
            'steps.aria': 'Payment steps',
            'steps.choose': 'Choose payment method',
            'steps.transfer': 'Scan or copy transfer details',
            'order.title': 'Order details',
            'order.merchantOrder': 'Merchant order',
            'order.product': 'Product',
            'order.tradeId': 'Transaction ID',
            'order.amount': 'Order amount',
            'method.title': 'Choose network and currency',
            'method.network': 'Network',
            'method.currency': 'Currency',
            'method.selectNetwork': 'Select network',
            'method.selectCurrency': 'Select currency',
            'method.selectMethod': 'Select payment method',
            'method.availableCurrencies': 'Available tokens: {{currencies}}',
            'method.warningTitle': 'Check network and token carefully',
            'method.warningBody': 'Before sending, make sure the network and token selected in your wallet exactly match this page. Transfers sent to the wrong network or token may be permanently lost.',
            'summary.needPay': 'Amount due',
            'summary.riskNote': 'The next step locks this payment method and shows the dedicated address and QR code. Transfer exactly the displayed amount.',
            'timer.remaining': 'Time left',
            'actions.next': 'Next',
            'actions.generating': 'Generating payment details...',
            'actions.back': 'Back',
            'actions.copyAddress': 'Copy address',
            'actions.copied': 'Copied',
            'transfer.payAmount': 'Payment amount',
            'transfer.copyAmount': 'Copy payment amount',
            'transfer.info': 'Transfer details',
            'transfer.title': 'Scan or copy the address to pay',
            'transfer.address': 'Receiving address',
            'instruction.networkDefault': 'Use the selected network',
            'instruction.networkUse': 'Use the selected {{network}} network',
            'instruction.amount': 'Transfer the exact amount',
            'instruction.autoConfirm': 'Keep this page open after transfer; the system will confirm automatically',
            'message.createFailed': 'Failed to generate payment details',
            'message.networkError': 'Network error, please try again later',
            'message.noNetworks': 'No payment network loaded. Please check wallet configuration.',
            'overlay.timeoutTitle': 'Payment expired',
            'overlay.timeoutBody': 'Please return to the merchant and create a new order. Do not transfer to an expired order.',
            'overlay.returnMerchant': 'Return to merchant',
            'overlay.successTitle': 'Payment successful',
            'overlay.successBody': 'The order has been confirmed and will return to the merchant page shortly.'
        },
        ru: {
            'page.title': 'BEpusdt - Касса',
            'brand.title': 'BEpusdt Касса',
            'brand.subtitle': 'Платежная сессия магазина',
            'lang.group': 'Язык',
            'theme.group': 'Тема',
            'theme.toLight': 'Переключить на светлую тему',
            'theme.toDark': 'Переключить на темную тему',
            'checkout.aria': 'Процесс оплаты',
            'checkout.eyebrow': 'Безопасная оплата',
            'status.waitingPayment': 'Ожидается оплата',
            'status.waitingConfirm': 'Ожидается подтверждение в сети',
            'status.timeout': 'Заказ просрочен',
            'status.success': 'Оплата успешна',
            'steps.aria': 'Шаги оплаты',
            'steps.choose': 'Выберите способ оплаты',
            'steps.transfer': 'Сканируйте или скопируйте данные',
            'order.title': 'Данные заказа',
            'order.merchantOrder': 'Заказ магазина',
            'order.product': 'Товар',
            'order.tradeId': 'Номер транзакции',
            'order.amount': 'Сумма заказа',
            'method.title': 'Выберите сеть и валюту',
            'method.network': 'Сеть',
            'method.currency': 'Валюта',
            'method.selectNetwork': 'Выберите сеть',
            'method.selectCurrency': 'Выберите валюту',
            'method.selectMethod': 'Выберите способ оплаты',
            'method.availableCurrencies': 'Доступные валюты: {{currencies}}',
            'method.warningTitle': 'Внимательно проверьте сеть и валюту',
            'method.warningBody': 'Перед переводом убедитесь, что сеть и валюта в вашем кошельке полностью совпадают с этой страницей. Перевод в неверной сети или валюте может привести к безвозвратной потере средств.',
            'summary.needPay': 'К оплате',
            'summary.riskNote': 'На следующем шаге способ оплаты будет зафиксирован, появятся адрес и QR-код. Переведите точно указанную сумму.',
            'timer.remaining': 'Осталось',
            'actions.next': 'Далее',
            'actions.generating': 'Формируем данные для оплаты...',
            'actions.back': 'Назад',
            'actions.copyAddress': 'Скопировать адрес',
            'actions.copied': 'Скопировано',
            'transfer.payAmount': 'Сумма оплаты',
            'transfer.copyAmount': 'Скопировать сумму оплаты',
            'transfer.info': 'Данные перевода',
            'transfer.title': 'Сканируйте QR-код или скопируйте адрес',
            'transfer.address': 'Адрес получателя',
            'instruction.networkDefault': 'Используйте выбранную сеть',
            'instruction.networkUse': 'Используйте выбранную сеть {{network}}',
            'instruction.amount': 'Переведите точную сумму',
            'instruction.autoConfirm': 'После перевода оставьте страницу открытой, система подтвердит оплату автоматически',
            'message.createFailed': 'Не удалось сформировать данные для оплаты',
            'message.networkError': 'Ошибка сети, попробуйте позже',
            'message.noNetworks': 'Не загружены платежные сети. Проверьте настройки кошельков.',
            'overlay.timeoutTitle': 'Время оплаты истекло',
            'overlay.timeoutBody': 'Вернитесь в магазин и создайте новый заказ. Не переводите средства по просроченному заказу.',
            'overlay.returnMerchant': 'Вернуться в магазин',
            'overlay.successTitle': 'Оплата успешна',
            'overlay.successBody': 'Заказ подтвержден. Сейчас вы будете возвращены на страницу магазина.'
        },
        vi: {
            'page.title': 'BEpusdt - Thanh toán',
            'brand.title': 'BEpusdt Thanh toán',
            'brand.subtitle': 'Phiên thanh toán của cửa hàng',
            'lang.group': 'Ngôn ngữ',
            'theme.group': 'Giao diện',
            'theme.toLight': 'Chuyển sang giao diện sáng',
            'theme.toDark': 'Chuyển sang giao diện tối',
            'checkout.aria': 'Quy trình thanh toán',
            'checkout.eyebrow': 'Thanh toán an toàn',
            'status.waitingPayment': 'Đang chờ thanh toán',
            'status.waitingConfirm': 'Đang chờ xác nhận trên blockchain',
            'status.timeout': 'Đơn hàng đã hết hạn',
            'status.success': 'Thanh toán thành công',
            'steps.aria': 'Các bước thanh toán',
            'steps.choose': 'Chọn cách thanh toán',
            'steps.transfer': 'Quét hoặc sao chép thông tin',
            'order.title': 'Thông tin đơn hàng',
            'order.merchantOrder': 'Đơn hàng cửa hàng',
            'order.product': 'Sản phẩm',
            'order.tradeId': 'Mã giao dịch',
            'order.amount': 'Số tiền đơn hàng',
            'method.title': 'Chọn mạng và đồng tiền',
            'method.network': 'Mạng',
            'method.currency': 'Đồng tiền',
            'method.selectNetwork': 'Chọn mạng',
            'method.selectCurrency': 'Chọn đồng tiền',
            'method.selectMethod': 'Chọn cách thanh toán',
            'method.availableCurrencies': 'Đồng tiền khả dụng: {{currencies}}',
            'method.warningTitle': 'Kiểm tra kỹ mạng và đồng tiền',
            'method.warningBody': 'Trước khi chuyển, hãy đảm bảo mạng và đồng tiền trong ví khớp hoàn toàn với trang này. Chuyển sai mạng hoặc sai đồng tiền có thể khiến tài sản mất vĩnh viễn.',
            'summary.needPay': 'Cần thanh toán',
            'summary.riskNote': 'Bước tiếp theo sẽ khóa cách thanh toán này và hiển thị địa chỉ cùng mã QR riêng. Vui lòng chuyển đúng số tiền hiển thị.',
            'timer.remaining': 'Còn lại',
            'actions.next': 'Tiếp tục',
            'actions.generating': 'Đang tạo thông tin thanh toán...',
            'actions.back': 'Quay lại',
            'actions.copyAddress': 'Sao chép địa chỉ',
            'actions.copied': 'Đã sao chép',
            'transfer.payAmount': 'Số tiền thanh toán',
            'transfer.copyAmount': 'Sao chép số tiền',
            'transfer.info': 'Thông tin chuyển khoản',
            'transfer.title': 'Quét mã QR hoặc sao chép địa chỉ để thanh toán',
            'transfer.address': 'Địa chỉ nhận',
            'instruction.networkDefault': 'Dùng mạng đã chọn',
            'instruction.networkUse': 'Dùng mạng {{network}}',
            'instruction.amount': 'Chuyển đúng số tiền',
            'instruction.autoConfirm': 'Sau khi chuyển tiền, hãy giữ trang này mở để hệ thống tự xác nhận',
            'message.createFailed': 'Không tạo được thông tin thanh toán',
            'message.networkError': 'Lỗi mạng, vui lòng thử lại sau',
            'message.noNetworks': 'Chưa tải được mạng thanh toán. Vui lòng kiểm tra cấu hình ví.',
            'overlay.timeoutTitle': 'Thanh toán đã hết hạn',
            'overlay.timeoutBody': 'Vui lòng quay lại cửa hàng và tạo đơn mới. Không chuyển tiền cho đơn đã hết hạn.',
            'overlay.returnMerchant': 'Quay lại cửa hàng',
            'overlay.successTitle': 'Thanh toán thành công',
            'overlay.successBody': 'Đơn hàng đã được xác nhận. Bạn sẽ được đưa về trang cửa hàng ngay sau đây.'
        },
        tr: {
            'page.title': 'BEpusdt - Ödeme',
            'brand.title': 'BEpusdt Ödeme',
            'brand.subtitle': 'Mağaza ödeme oturumu',
            'lang.group': 'Dil',
            'theme.group': 'Tema',
            'theme.toLight': 'Açık temaya geç',
            'theme.toDark': 'Koyu temaya geç',
            'checkout.aria': 'Ödeme akışı',
            'checkout.eyebrow': 'Güvenli ödeme',
            'status.waitingPayment': 'Ödeme bekleniyor',
            'status.waitingConfirm': 'Zincir onayı bekleniyor',
            'status.timeout': 'Siparişin süresi doldu',
            'status.success': 'Ödeme başarılı',
            'steps.aria': 'Ödeme adımları',
            'steps.choose': 'Ödeme yöntemini seç',
            'steps.transfer': 'QR tara veya bilgileri kopyala',
            'order.title': 'Sipariş bilgileri',
            'order.merchantOrder': 'Mağaza siparişi',
            'order.product': 'Ürün',
            'order.tradeId': 'İşlem numarası',
            'order.amount': 'Sipariş tutarı',
            'method.title': 'Ağ ve para birimi seç',
            'method.network': 'Ağ',
            'method.currency': 'Para birimi',
            'method.selectNetwork': 'Ağ seç',
            'method.selectCurrency': 'Para birimi seç',
            'method.selectMethod': 'Ödeme yöntemi seç',
            'method.availableCurrencies': 'Kullanılabilir tokenler: {{currencies}}',
            'method.warningTitle': 'Ağı ve para birimini dikkatlice kontrol edin',
            'method.warningBody': 'Göndermeden önce cüzdanınızda seçilen ağ ve para biriminin bu sayfayla tamamen aynı olduğundan emin olun. Yanlış ağa veya para birimine yapılan transferler kalıcı kayba yol açabilir.',
            'summary.needPay': 'Ödenecek tutar',
            'summary.riskNote': 'Sonraki adımda bu ödeme yöntemi sabitlenir, size özel adres ve QR kodu gösterilir. Lütfen ekrandaki tutarın aynısını gönderin.',
            'timer.remaining': 'Kalan süre',
            'actions.next': 'Devam',
            'actions.generating': 'Ödeme bilgileri hazırlanıyor...',
            'actions.back': 'Geri',
            'actions.copyAddress': 'Adresi kopyala',
            'actions.copied': 'Kopyalandı',
            'transfer.payAmount': 'Ödeme tutarı',
            'transfer.copyAmount': 'Ödeme tutarını kopyala',
            'transfer.info': 'Transfer bilgileri',
            'transfer.title': 'Ödemek için QR kodu tara veya adresi kopyala',
            'transfer.address': 'Alıcı adresi',
            'instruction.networkDefault': 'Seçilen ağı kullan',
            'instruction.networkUse': '{{network}} ağını kullan',
            'instruction.amount': 'Tutarı tam olarak gönder',
            'instruction.autoConfirm': 'Transferden sonra sayfayı açık bırakın; sistem ödemeyi otomatik onaylar',
            'message.createFailed': 'Ödeme bilgileri oluşturulamadı',
            'message.networkError': 'Ağ hatası, lütfen tekrar deneyin',
            'message.noNetworks': 'Ödeme ağı yüklenemedi. Lütfen cüzdan ayarlarını kontrol edin.',
            'overlay.timeoutTitle': 'Ödeme süresi doldu',
            'overlay.timeoutBody': 'Lütfen mağazaya dönüp yeni bir sipariş oluşturun. Süresi dolan siparişe transfer yapmayın.',
            'overlay.returnMerchant': 'Mağazaya dön',
            'overlay.successTitle': 'Ödeme başarılı',
            'overlay.successBody': 'Sipariş onaylandı. Kısa süre içinde mağaza sayfasına yönlendirileceksiniz.'
        },
        ja: {
            'page.title': 'BEpusdt - お支払い',
            'brand.title': 'BEpusdt お支払い',
            'brand.subtitle': '加盟店のお支払いセッション',
            'lang.group': '言語',
            'theme.group': 'テーマ',
            'theme.toLight': 'ライトテーマに切り替え',
            'theme.toDark': 'ダークテーマに切り替え',
            'checkout.aria': 'お支払い手続き',
            'checkout.eyebrow': '安全なお支払い',
            'status.waitingPayment': 'お支払い待ち',
            'status.waitingConfirm': 'オンチェーン承認待ち',
            'status.timeout': '注文の有効期限が切れました',
            'status.success': 'お支払い完了',
            'steps.aria': 'お支払い手順',
            'steps.choose': '支払い方法を選択',
            'steps.transfer': 'QRコードまたは送金情報',
            'order.title': '注文情報',
            'order.merchantOrder': '加盟店注文番号',
            'order.product': '商品',
            'order.tradeId': '取引番号',
            'order.amount': '注文金額',
            'method.title': 'ネットワークと通貨を選択',
            'method.network': 'ネットワーク',
            'method.currency': '通貨',
            'method.selectNetwork': 'ネットワークを選択',
            'method.selectCurrency': '通貨を選択',
            'method.selectMethod': '支払い方法を選択',
            'method.availableCurrencies': '利用可能な通貨：{{currencies}}',
            'method.warningTitle': 'ネットワークと通貨を必ず確認してください',
            'method.warningBody': '送金前に、ウォレットで選択したネットワークと通貨がこのページと完全に一致していることを確認してください。誤ったネットワークまたは通貨への送金は、資産を永久に失う可能性があります。',
            'summary.needPay': 'お支払い金額',
            'summary.riskNote': '次へ進むと今回の支払い方法が確定し、専用のアドレスとQRコードが表示されます。表示された金額ちょうどを送金してください。',
            'timer.remaining': '残り時間',
            'actions.next': '次へ',
            'actions.generating': '支払い情報を作成中...',
            'actions.back': '戻る',
            'actions.copyAddress': 'アドレスをコピー',
            'actions.copied': 'コピーしました',
            'transfer.payAmount': '支払い金額',
            'transfer.copyAmount': '支払い金額をコピー',
            'transfer.info': '送金情報',
            'transfer.title': 'QRコードを読み取るか、アドレスをコピーしてお支払いください',
            'transfer.address': '受取アドレス',
            'instruction.networkDefault': '選択したネットワークを使用',
            'instruction.networkUse': '{{network}} ネットワークを使用',
            'instruction.amount': '表示どおりの金額を送金',
            'instruction.autoConfirm': '送金後はこのページを開いたままお待ちください。システムが自動で確認します',
            'message.createFailed': '支払い情報を作成できませんでした',
            'message.networkError': '通信エラーです。時間をおいてもう一度お試しください',
            'message.noNetworks': '利用可能な支払いネットワークを読み込めません。ウォレット設定をご確認ください。',
            'overlay.timeoutTitle': 'お支払い期限が切れました',
            'overlay.timeoutBody': '加盟店ページに戻り、新しい注文を作成してください。期限切れの注文には送金しないでください。',
            'overlay.returnMerchant': '加盟店へ戻る',
            'overlay.successTitle': 'お支払いが完了しました',
            'overlay.successBody': '注文の入金確認が完了しました。まもなく加盟店ページへ戻ります。'
        },
        ko: {
            'page.title': 'BEpusdt - 결제',
            'brand.title': 'BEpusdt 결제',
            'brand.subtitle': '가맹점 결제 세션',
            'lang.group': '언어',
            'theme.group': '테마',
            'theme.toLight': '라이트 모드로 전환',
            'theme.toDark': '다크 모드로 전환',
            'checkout.aria': '결제 절차',
            'checkout.eyebrow': '안전 결제',
            'status.waitingPayment': '결제 대기 중',
            'status.waitingConfirm': '온체인 확인 대기 중',
            'status.timeout': '주문 시간이 만료되었습니다',
            'status.success': '결제가 완료되었습니다',
            'steps.aria': '결제 단계',
            'steps.choose': '결제 수단 선택',
            'steps.transfer': 'QR 스캔 또는 정보 복사',
            'order.title': '주문 정보',
            'order.merchantOrder': '가맹점 주문',
            'order.product': '상품',
            'order.tradeId': '거래 번호',
            'order.amount': '주문 금액',
            'method.title': '네트워크와 코인 선택',
            'method.network': '네트워크',
            'method.currency': '코인',
            'method.selectNetwork': '네트워크 선택',
            'method.selectCurrency': '코인 선택',
            'method.selectMethod': '결제 수단 선택',
            'method.availableCurrencies': '사용 가능한 코인: {{currencies}}',
            'method.warningTitle': '네트워크와 코인을 반드시 확인하세요',
            'method.warningBody': '송금 전에 지갑에서 선택한 네트워크와 코인이 이 페이지와 완전히 일치하는지 확인하세요. 잘못된 네트워크나 코인으로 송금하면 자산을 영구적으로 잃을 수 있습니다.',
            'summary.needPay': '결제할 금액',
            'summary.riskNote': '다음 단계에서 이번 결제 수단이 확정되며 전용 주소와 QR 코드가 표시됩니다. 표시된 금액과 정확히 같은 금액을 보내주세요.',
            'timer.remaining': '남은 시간',
            'actions.next': '다음',
            'actions.generating': '결제 정보를 생성 중...',
            'actions.back': '이전',
            'actions.copyAddress': '주소 복사',
            'actions.copied': '복사됨',
            'transfer.payAmount': '결제 금액',
            'transfer.copyAmount': '결제 금액 복사',
            'transfer.info': '송금 정보',
            'transfer.title': 'QR 코드를 스캔하거나 주소를 복사해 결제하세요',
            'transfer.address': '받는 주소',
            'instruction.networkDefault': '선택한 네트워크 사용',
            'instruction.networkUse': '{{network}} 네트워크 사용',
            'instruction.amount': '정확한 금액 송금',
            'instruction.autoConfirm': '송금 후 이 페이지를 열어 두면 시스템이 자동으로 확인합니다',
            'message.createFailed': '결제 정보를 생성하지 못했습니다',
            'message.networkError': '네트워크 오류입니다. 잠시 후 다시 시도해 주세요',
            'message.noNetworks': '결제 네트워크를 불러오지 못했습니다. 지갑 설정을 확인해 주세요.',
            'overlay.timeoutTitle': '결제 시간이 만료되었습니다',
            'overlay.timeoutBody': '가맹점으로 돌아가 새 주문을 만들어 주세요. 만료된 주문에는 송금하지 마세요.',
            'overlay.returnMerchant': '가맹점으로 돌아가기',
            'overlay.successTitle': '결제가 완료되었습니다',
            'overlay.successBody': '주문 입금이 확인되었습니다. 곧 가맹점 페이지로 이동합니다.'
        }
    };

    var tokenIcons = {
        USDT: '/payment/assets/img/tokens/usdt.svg',
        USDC: '/payment/assets/img/tokens/usdc.svg',
        TRX: '/payment/assets/img/tokens/trx.svg',
        BNB: '/payment/assets/img/tokens/bnb.svg',
        ETH: '/payment/assets/img/tokens/eth.svg',
        GRAM: '/payment/assets/img/tokens/gram.svg'
    };

    var networkIcons = {
        tron: '/payment/assets/img/networks/tron.svg',
        bsc: '/payment/assets/img/networks/bsc.svg',
        ethereum: '/payment/assets/img/networks/ethereum.svg',
        solana: '/payment/assets/img/networks/solana.svg',
        polygon: '/payment/assets/img/networks/polygon.svg',
        arbitrum: '/payment/assets/img/networks/arbitrum.svg',
        base: '/payment/assets/img/networks/base.svg',
        aptos: '/payment/assets/img/networks/aptos.svg',
        plasma: '/payment/assets/img/networks/plasma.svg',
        xlayer: '/payment/assets/img/networks/xlayer.svg',
        ton: '/payment/assets/img/networks/ton.svg'
    };

    var languageShortLabels = {
        zh: '中',
        'zh-Hant': '繁',
        en: 'EN',
        ru: 'RU',
        vi: 'VI',
        tr: 'TR',
        ja: 'JA',
        ko: 'KO'
    };

    function $(selector) {
        return document.querySelector(selector);
    }

    function $all(selector) {
        return Array.prototype.slice.call(document.querySelectorAll(selector));
    }

    function normalizeLanguage(lang) {
        var value = String(lang || '').toLowerCase();
        if (value === 'tw' || value === 'hk' || value === 'mo') return 'zh-Hant';
        if (value.indexOf('zh-tw') === 0 || value.indexOf('zh-hk') === 0 || value.indexOf('zh-mo') === 0 || value.indexOf('zh-hant') === 0) return 'zh-Hant';
        if (value.indexOf('ru') === 0) return 'ru';
        if (value.indexOf('vi') === 0) return 'vi';
        if (value.indexOf('tr') === 0) return 'tr';
        if (value.indexOf('ja') === 0) return 'ja';
        if (value.indexOf('ko') === 0) return 'ko';
        if (value.indexOf('en') === 0) return 'en';
        if (value.indexOf('zh') === 0 || value.indexOf('cn') === 0) return 'zh';
        return 'zh';
    }

    function normalizeDefaultLanguage(lang) {
        var value = String(lang || '').trim();
        if (!value || value.toLowerCase() === 'auto') return '';
        return normalizeLanguage(value);
    }

    function detectLanguage() {
        var params = new URLSearchParams(window.location.search);
        var queryLang = params.get('lang');
        if (queryLang) return normalizeLanguage(queryLang);

        var storedLang = localStorage.getItem('bepusdt-cashier-lang');
        if (storedLang) return normalizeLanguage(storedLang);

        return normalizeDefaultLanguage(config.default_language) || normalizeLanguage(navigator.language);
    }

    function interpolate(text, values) {
        return String(text).replace(/\{\{(\w+)\}\}/g, function (_, key) {
            return values && values[key] != null ? values[key] : '';
        });
    }

    function t(key, values) {
        var source = translations[currentLang] || translations.zh;
        return interpolate(source[key] || translations.zh[key] || key, values);
    }

    function setLanguageMenuOpen(open) {
        var switcher = $('.language-switch');
        var toggle = $('#languageToggle');
        if (!switcher) return;

        switcher.classList.toggle('open', open);
        if (toggle) toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
        if (open) closeDropdowns();
    }

    function closeLanguageMenu() {
        setLanguageMenuOpen(false);
    }

    function updateLanguageControl() {
        var currentLabel = $('#currentLanguageLabel');
        if (currentLabel) currentLabel.textContent = languageShortLabels[currentLang] || currentLang.toUpperCase();

        $all('[data-lang-choice]').forEach(function (button) {
            var active = button.dataset.langChoice === currentLang;
            button.classList.toggle('active', active);
            button.setAttribute('aria-selected', active ? 'true' : 'false');
        });
    }

    function translateStatic() {
        document.documentElement.lang = currentLang === 'zh' ? 'zh-CN' : currentLang;
        document.title = t('page.title');

        $all('[data-i18n]').forEach(function (el) {
            el.textContent = t(el.dataset.i18n);
        });
        $all('[data-i18n-aria]').forEach(function (el) {
            el.setAttribute('aria-label', t(el.dataset.i18nAria));
        });
        $all('[data-i18n-title]').forEach(function (el) {
            el.setAttribute('title', t(el.dataset.i18nTitle));
        });
        updateLanguageControl();
    }

    function refreshLocalizedState() {
        translateStatic();
        applyTheme(document.documentElement.dataset.theme || localStorage.getItem('bepusdt-cashier-theme') || 'light');
        setStatusKey(currentStatusKey);

        if (selectedNetworkId) renderNetworkOptions(selectedNetworkId);
        if (selectedCurrency) renderCurrencyOptions(selectedCurrency);
        updateSelectedSummary();
        if (paymentDetail) renderPaymentDetail();
    }

    function setLanguage(lang) {
        currentLang = normalizeLanguage(lang);
        localStorage.setItem('bepusdt-cashier-lang', currentLang);
        refreshLocalizedState();
    }

    function initLanguage() {
        currentLang = detectLanguage();
        translateStatic();

        var toggle = $('#languageToggle');
        if (toggle && !toggle.dataset.bound) {
            toggle.dataset.bound = 'true';
            toggle.addEventListener('click', function (event) {
                event.stopPropagation();
                setLanguageMenuOpen(!toggle.closest('.language-switch').classList.contains('open'));
            });
        }

        $all('[data-lang-choice]').forEach(function (button) {
            if (button.dataset.bound) return;
            button.dataset.bound = 'true';
            button.addEventListener('click', function () {
                setLanguage(button.dataset.langChoice);
                closeLanguageMenu();
            });
        });
    }

    function moneyLabel(value, currency) {
        if (!value || !currency) return '--';
        return value + ' ' + currency;
    }

    function displayNetwork(network) {
        return network ? String(network).toUpperCase() : '';
    }

    function networkName(method) {
        if (!method) return '--';
        return displayNetwork(method.network) || method.token_net_name || '--';
    }

    function networkProtocolLabel(method) {
        if (!method) return '--';
        var tokenNetName = method.token_net_name || '';
        var baseNetwork = networkName(method);
        if (!tokenNetName) return baseNetwork;
        if (String(tokenNetName).toLowerCase() === String(method.network || '').toLowerCase()) return baseNetwork;
        if (String(tokenNetName).toLowerCase() === String(baseNetwork).toLowerCase()) return baseNetwork;
        return baseNetwork + ' · ' + tokenNetName;
    }

    function tokenProtocolLabel(method) {
        if (!method) return '';
        return networkProtocolLabel(method).replace(networkName(method) + ' · ', '');
    }

    function networkCurrencyList(networkId) {
        return currenciesForNetwork(networkId).join(' / ');
    }

    function networkCurrencyMeta(networkId) {
        return t('method.availableCurrencies', {currencies: networkCurrencyList(networkId)});
    }

    function currencyRank(currency) {
        var order = ['USDT', 'USDC', 'TRX', 'BNB', 'ETH'];
        var rank = order.indexOf(String(currency).toUpperCase());
        return rank === -1 ? order.length : rank;
    }

    function networkRank(network, tokenNetName) {
        var value = String(network || tokenNetName || '').toUpperCase();
        var order = ['TRON', 'TRC20', 'BSC', 'BEP20', 'ETHEREUM', 'ERC20', 'SOLANA'];
        var rank = order.indexOf(value);
        return rank === -1 ? order.length : rank;
    }

    function methodSort(a, b) {
        var diff = currencyRank(a.currency) - currencyRank(b.currency);
        if (diff !== 0) return diff;
        diff = networkRank(a.network, a.token_net_name) - networkRank(b.network, b.token_net_name);
        if (diff !== 0) return diff;
        return networkProtocolLabel(a).localeCompare(networkProtocolLabel(b));
    }

    function escapeHtml(value) {
        return String(value == null ? '' : value)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;');
    }

    function fallbackText(value) {
        var text = String(value || '').replace(/[^a-z0-9]/gi, '').slice(0, 3).toUpperCase();
        return text || '?';
    }

    function iconMarkup(src, label, className) {
        if (src) {
            return '<img class="' + className + '" src="' + src + '" alt="" loading="lazy">';
        }
        return '<span class="' + className + ' icon-fallback" aria-hidden="true">' + escapeHtml(fallbackText(label)) + '</span>';
    }

    function tokenIconPath(currency) {
        return tokenIcons[String(currency || '').toUpperCase()] || '';
    }

    function networkKey(method) {
        if (!method) return '';
        var value = String(method.network || method.token_net_name || '').toLowerCase().replace(/[\s_-]+/g, '');
        if (value.indexOf('trc20') !== -1 || value.indexOf('tron') !== -1) return 'tron';
        if (value.indexOf('bep20') !== -1 || value.indexOf('bsc') !== -1 || value.indexOf('binance') !== -1) return 'bsc';
        if (value.indexOf('erc20') !== -1 || value.indexOf('ethereum') !== -1) return 'ethereum';
        if (value.indexOf('solana') !== -1) return 'solana';
        if (value.indexOf('polygon') !== -1) return 'polygon';
        if (value.indexOf('arbitrum') !== -1) return 'arbitrum';
        if (value.indexOf('base') !== -1) return 'base';
        if (value.indexOf('aptos') !== -1) return 'aptos';
        if (value.indexOf('plasma') !== -1) return 'plasma';
        if (value.indexOf('xlayer') !== -1) return 'xlayer';
        return value;
    }

    function networkIconPath(method) {
        return networkIcons[networkKey(method)] || '';
    }

    function setImageSource(image, src) {
        if (!image) return;
        if (!src) {
            image.hidden = true;
            image.removeAttribute('src');
            return;
        }
        image.src = src;
        image.hidden = false;
    }

    function setStep(step) {
        $all('[data-step]').forEach(function (pane) {
            pane.classList.toggle('active', pane.dataset.step === step);
        });

        $all('[data-step-target]').forEach(function (button) {
            var buttonStep = button.dataset.stepTarget || '1';
            button.classList.toggle('active', buttonStep === step);
            button.classList.toggle('done', Number(buttonStep) < Number(step));
        });

        var stage = $('.checkout-card');
        if (stage) {
            stage.scrollIntoView({block: 'start'});
        }
    }

    function applyTheme(theme) {
        document.documentElement.dataset.theme = theme;
        var button = $('#themeToggle');
        if (!button) return;
        var nextTheme = theme === 'dark' ? 'light' : 'dark';
        button.textContent = theme === 'dark' ? '☀' : '🌙';
        button.dataset.nextTheme = nextTheme;
        button.setAttribute('aria-label', theme === 'dark' ? t('theme.toLight') : t('theme.toDark'));
    }

    function initTheme() {
        applyTheme(localStorage.getItem('bepusdt-cashier-theme') || 'light');
        var button = $('#themeToggle');
        if (button) {
            button.addEventListener('click', function () {
                var theme = button.dataset.nextTheme || 'dark';
                localStorage.setItem('bepusdt-cashier-theme', theme);
                applyTheme(theme);
            });
        }
    }

    function networkIdentity(method) {
        if (!method) return '';
        return String(method.network || method.token_net_name || '').toLowerCase();
    }

    function sameNetwork(method, networkId) {
        return networkIdentity(method) === networkId;
    }

    function uniqueNetworkMethods() {
        var seen = {};
        var networks = [];
        methods.forEach(function (method) {
            var id = networkIdentity(method);
            if (!seen[id]) {
                seen[id] = true;
                networks.push(method);
            }
        });
        return networks.sort(function (a, b) {
            var diff = networkRank(a.network, a.token_net_name) - networkRank(b.network, b.token_net_name);
            if (diff !== 0) return diff;
            return networkName(a).localeCompare(networkName(b));
        });
    }

    function methodsForNetwork(networkId) {
        return methods.filter(function (method) {
            return sameNetwork(method, networkId);
        }).sort(methodSort);
    }

    function currenciesForNetwork(networkId) {
        var currencies = [];
        methodsForNetwork(networkId).forEach(function (method) {
            if (currencies.indexOf(method.currency) === -1) {
                currencies.push(method.currency);
            }
        });
        return currencies.sort(function (a, b) {
            var diff = currencyRank(a) - currencyRank(b);
            if (diff !== 0) return diff;
            return String(a).localeCompare(String(b));
        });
    }

    function setDropdownOpen(type, open) {
        $all('.custom-select').forEach(function (dropdown) {
            var isTarget = dropdown.dataset.dropdown === type;
            var shouldOpen = open && isTarget;
            dropdown.classList.toggle('open', shouldOpen);
            var button = dropdown.querySelector('.custom-select-button');
            if (button) button.setAttribute('aria-expanded', shouldOpen ? 'true' : 'false');
        });
    }

    function closeDropdowns() {
        setDropdownOpen('', false);
    }

    function optionButton(content, selected) {
        var button = document.createElement('button');
        button.className = 'custom-option';
        button.type = 'button';
        button.setAttribute('role', 'option');
        button.setAttribute('aria-selected', selected ? 'true' : 'false');
        button.innerHTML = content;
        return button;
    }

    function updateSelectButton(button, icon, title, meta) {
        if (!button) return;
        button.innerHTML =
            '<span class="custom-select-value">' +
            icon +
            '<span class="custom-select-copy">' +
            '<strong>' + escapeHtml(title) + '</strong>' +
            (meta ? '<small>' + escapeHtml(meta) + '</small>' : '') +
            '</span>' +
            '</span>' +
            '<span class="custom-select-caret"></span>';
    }

    function renderPaymentSelectors() {
        methods = methods.slice().sort(methodSort);
        bindDropdownButtons();

        if (methods.length > 0) {
            var preferred = selectedMethod || methods.find(function (method) {
                return method.is_popular;
            }) || methods[0];
            selectNetwork(networkIdentity(preferred), preferred);
        } else {
            renderNetworkOptions('');
            renderCurrencyOptions('');
            updateSelectedSummary();
        }
    }

    function normalizeSelectedPayment(detail) {
        if (!detail || typeof detail !== 'object') return null;
        if (!detail.token && detail.address) {
            detail.token = detail.address;
        }

        return detail;
    }

    function selectedPaymentMethod(detail) {
        if (!detail) return null;
        var currency = String(detail.currency || '').toUpperCase();
        var network = String(detail.network || '').toLowerCase();
        var tokenNetName = String(detail.token_net_name || '').toLowerCase();

        return methods.find(function (method) {
            return String(method.currency || '').toUpperCase() === currency &&
                (String(method.network || '').toLowerCase() === network ||
                    String(method.token_net_name || '').toLowerCase() === tokenNetName);
        }) || {
            actual_amount: detail.actual_amount || detail.amount || '',
            currency: detail.currency || '',
            network: detail.network || '',
            token_net_name: detail.token_net_name || detail.network || ''
        };
    }

    function bindDropdownButtons() {
        var currencyButton = $('#currencySelectButton');
        var networkButton = $('#networkSelectButton');

        if (currencyButton && !currencyButton.dataset.bound) {
            currencyButton.dataset.bound = 'true';
            currencyButton.addEventListener('click', function (event) {
                event.stopPropagation();
                closeLanguageMenu();
                setDropdownOpen('currency', !currencyButton.closest('.custom-select').classList.contains('open'));
            });
        }

        if (networkButton && !networkButton.dataset.bound) {
            networkButton.dataset.bound = 'true';
            networkButton.addEventListener('click', function (event) {
                event.stopPropagation();
                closeLanguageMenu();
                setDropdownOpen('network', !networkButton.closest('.custom-select').classList.contains('open'));
            });
        }
    }

    function renderNetworkOptions(activeNetworkId) {
        var options = $('#networkOptions');
        if (!options) return;

        var networks = uniqueNetworkMethods();
        options.innerHTML = '';
        networks.forEach(function (method) {
            var id = networkIdentity(method);
            var selected = id === activeNetworkId;
            var option = optionButton(
                iconMarkup(networkIconPath(method), method.network || method.token_net_name, 'network-icon') +
                '<span class="custom-option-copy"><strong>' + escapeHtml(networkName(method)) + '</strong><small>' + escapeHtml(networkCurrencyMeta(id)) + '</small></span>',
                selected
            );
            option.addEventListener('click', function (event) {
                event.stopPropagation();
                selectNetwork(id);
                closeDropdowns();
            });
            options.appendChild(option);
        });

        updateNetworkOptionsState(activeNetworkId);
        updateNetworkButton(networks.find(function (method) {
            return networkIdentity(method) === activeNetworkId;
        }));
    }

    function selectNetwork(networkId, preferredMethod) {
        selectedNetworkId = networkId;
        renderNetworkOptions(networkId);

        var candidates = methodsForNetwork(networkId);
        var nextMethod = preferredMethod && sameNetwork(preferredMethod, networkId) ? preferredMethod : candidates[0];
        selectedCurrency = nextMethod ? nextMethod.currency : '';
        renderCurrencyOptions(selectedCurrency);
        selectMethod(nextMethod || null);
    }

    function renderCurrencyOptions(activeCurrency) {
        var options = $('#currencyOptions');
        if (!options) return;

        var currencies = currenciesForNetwork(selectedNetworkId);
        options.innerHTML = '';
        currencies.forEach(function (currency) {
            var method = methodsForNetwork(selectedNetworkId).find(function (item) {
                return item.currency === currency;
            });
            var selected = currency === activeCurrency;
            var option = optionButton(
                iconMarkup(tokenIconPath(currency), currency, 'token-icon') +
                '<span class="custom-option-copy"><strong>' + escapeHtml(currency) + '</strong><small>' + escapeHtml(tokenProtocolLabel(method)) + '</small></span>',
                selected
            );
            option.addEventListener('click', function (event) {
                event.stopPropagation();
                selectCurrency(currency);
                closeDropdowns();
            });
            options.appendChild(option);
        });

        updateCurrencyOptionsState(activeCurrency);
        updateCurrencyButton(activeCurrency);
    }

    function selectCurrency(currency) {
        selectedCurrency = currency;
        renderCurrencyOptions(currency);
        selectMethod(methodsForNetwork(selectedNetworkId).find(function (method) {
            return method.currency === currency;
        }) || null);
    }

    function updateCurrencyButton(currency) {
        var button = $('#currencySelectButton');
        if (currency) {
            var method = methodsForNetwork(selectedNetworkId).find(function (item) {
                return item.currency === currency;
            });
            updateSelectButton(button, iconMarkup(tokenIconPath(currency), currency, 'token-icon'), currency, tokenProtocolLabel(method));
        } else {
            updateSelectButton(button, '', t('method.selectCurrency'), '');
        }
    }

    function updateCurrencyOptionsState(currency) {
        $all('#currencyOptions .custom-option').forEach(function (option, index) {
            option.setAttribute('aria-selected', currenciesForNetwork(selectedNetworkId)[index] === currency ? 'true' : 'false');
        });
    }

    function selectMethod(method) {
        selectedMethod = method;
        if (method) selectedCurrency = method.currency;
        updateSelectedSummary();
    }

    function updateNetworkButton(method) {
        var button = $('#networkSelectButton');
        if (!method) {
            updateSelectButton(button, '', t('method.selectNetwork'), '');
            return;
        }
        updateSelectButton(
            button,
            iconMarkup(networkIconPath(method), method.network || method.token_net_name, 'network-icon'),
            networkName(method),
            networkCurrencyMeta(networkIdentity(method))
        );
    }

    function updateNetworkOptionsState(networkId) {
        var networks = uniqueNetworkMethods();
        $all('#networkOptions .custom-option').forEach(function (option, index) {
            option.setAttribute('aria-selected', networkIdentity(networks[index]) === networkId ? 'true' : 'false');
        });
    }

    function updateSelectedSummary() {
        var amount = $('#selectedAmount');
        var network = $('#selectedNetwork');
        var button = $('#continueButton');

        if (!selectedMethod) {
            if (amount) amount.textContent = '--';
            if (network) network.textContent = t('method.selectMethod');
            if (button) button.disabled = true;
            return;
        }

        if (amount) amount.textContent = moneyLabel(selectedMethod.actual_amount, selectedMethod.currency);
        if (network) network.textContent = networkProtocolLabel(selectedMethod);
        if (button) button.disabled = false;
    }

    function lockPaymentMethod() {
        if (!selectedMethod) return;

        var button = $('#continueButton');
        if (button) {
            button.disabled = true;
            button.textContent = t('actions.generating');
        }

        fetch('/api/v1/pay/update-order', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
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
                if (res.status_code !== 200 || !res.data) {
                    throw new Error(res.message || t('message.createFailed'));
                }

                paymentDetail = res.data;
                if (!paymentDetail.token && paymentDetail.payment_url) {
                    window.location.href = paymentDetail.payment_url;
                    return;
                }
                if (!paymentDetail.token) {
                    throw new Error(res.message || t('message.createFailed'));
                }

                tradeId = paymentDetail.trade_id || tradeId;
                renderPaymentDetail();
                setStep('2');
                restartStatusPolling();
            })
            .catch(function (error) {
                showMessage(error.message || t('message.networkError'));
            })
            .finally(function () {
                if (button) {
                    button.disabled = false;
                    button.textContent = t('actions.next');
                }
            });
    }

    function renderPaymentDetail() {
        if (!paymentDetail) return;

        var rawAmount = paymentDetail.actual_amount || paymentDetail.amount || '';
        var amountText = moneyLabel(rawAmount, selectedMethod.currency);
        var networkText = networkProtocolLabel(selectedMethod);
        var address = paymentDetail.token || '';
        paymentAmountValue = String(rawAmount);

        var payAmount = $('#payAmount');
        var payNetwork = $('#payNetwork');
        var payTokenIcon = $('#payTokenIcon');
        var payNetworkIcon = $('#payNetworkIcon');
        var addressEl = $('#walletAddress');
        var networkInstruction = $('#networkInstruction');
        var amountInstruction = $('#amountInstruction');

        if (payAmount) payAmount.textContent = amountText;
        if (payNetwork) payNetwork.textContent = networkText;
        setImageSource(payTokenIcon, tokenIconPath(selectedMethod.currency));
        setImageSource(payNetworkIcon, networkIconPath(selectedMethod));
        if (addressEl) addressEl.textContent = address;
        if (networkInstruction) networkInstruction.textContent = t('instruction.networkUse', {network: networkText});
        if (amountInstruction) amountInstruction.textContent = t('instruction.amount');

        renderQrCode(address);
    }

    function renderQrCode(address) {
        var qrContainer = $('#qrcode');
        if (!qrContainer) return;

        qrContainer.innerHTML = '';
        if (!address) {
            qrContainer.textContent = '--';
            return;
        }

        if (window.jQuery && window.jQuery.fn && window.jQuery.fn.qrcode) {
            window.jQuery(qrContainer).qrcode({
                text: address,
                width: 280,
                height: 280
            });
            return;
        }

        qrContainer.textContent = address;
    }

    function copyText(text, button) {
        if (!text) return;
        var original = button ? button.textContent : '';

        function done() {
            if (!button) return;
            if (button.classList.contains('icon-copy-button')) {
                button.classList.add('copied');
                setTimeout(function () {
                    button.classList.remove('copied');
                }, 1600);
                return;
            }
            button.textContent = t('actions.copied');
            setTimeout(function () {
                button.textContent = original;
            }, 1600);
        }

        if (navigator.clipboard && navigator.clipboard.writeText) {
            navigator.clipboard.writeText(text).then(done).catch(function () {
                fallbackCopy(text);
                done();
            });
        } else {
            fallbackCopy(text);
            done();
        }
    }

    function fallbackCopy(text) {
        var textarea = document.createElement('textarea');
        textarea.value = text;
        textarea.style.position = 'fixed';
        textarea.style.opacity = '0';
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand('copy');
        document.body.removeChild(textarea);
    }

    function startCountdown() {
        if (countdownTimer) clearInterval(countdownTimer);

        totalSeconds = parseInt(config.expire, 10);
        if (isNaN(totalSeconds)) totalSeconds = 0;

        function render() {
            var minutes = Math.floor(Math.max(0, totalSeconds) / 60).toString().padStart(2, '0');
            var seconds = (Math.max(0, totalSeconds) % 60).toString().padStart(2, '0');

            ['#minutes', '#minutesMirror'].forEach(function (selector) {
                var el = $(selector);
                if (el) el.textContent = minutes;
            });
            ['#seconds', '#secondsMirror'].forEach(function (selector) {
                var el = $(selector);
                if (el) el.textContent = seconds;
            });
        }

        render();
        countdownTimer = setInterval(function () {
            totalSeconds -= 1;
            render();
            if (totalSeconds <= 0) {
                clearInterval(countdownTimer);
                if (statusCheckTimer) clearInterval(statusCheckTimer);
                showTimeout();
            }
        }, 1000);
    }

    function restartStatusPolling() {
        if (statusCheckTimer) clearInterval(statusCheckTimer);
        statusCheckTimer = setInterval(checkStatus, 5000);
        checkStatus();
    }

    function checkStatus() {
        if (!tradeId) return;

        fetch('/pay/check-status/' + tradeId)
            .then(function (response) {
                return response.json();
            })
            .then(function (data) {
                if (data.status === 2) {
                    showSuccess(data);
                } else if (data.status === 3) {
                    showTimeout();
                } else if (data.status === 5) {
                    setStatusKey('status.waitingConfirm');
                }
            })
            .catch(function (error) {
                console.error('Check status error:', error);
            });
    }

    function setStatusKey(key) {
        currentStatusKey = key;
        var statusText = $('#statusText');
        if (statusText) statusText.textContent = t(key);
    }

    function showMessage(message) {
        var toast = document.createElement('div');
        toast.className = 'cashier-toast';
        toast.textContent = message;
        document.body.appendChild(toast);
        setTimeout(function () {
            toast.remove();
        }, 2600);
    }

    function showOverlay(type, title, body, actionText, action) {
        var old = $('.cashier-overlay');
        if (old) old.remove();

        var overlay = document.createElement('div');
        overlay.className = 'cashier-overlay';

        var panel = document.createElement('div');
        panel.className = 'overlay-panel ' + type;

        var mark = document.createElement('div');
        mark.className = 'overlay-mark';
        mark.textContent = type === 'success' ? '✓' : '!';

        var heading = document.createElement('h3');
        heading.textContent = title;

        var copy = document.createElement('p');
        copy.textContent = body;

        var button = document.createElement('button');
        button.className = 'button primary';
        button.type = 'button';
        button.textContent = actionText;
        button.addEventListener('click', action);

        panel.appendChild(mark);
        panel.appendChild(heading);
        panel.appendChild(copy);
        panel.appendChild(button);
        overlay.appendChild(panel);
        document.body.appendChild(overlay);
    }

    function showTimeout() {
        setStatusKey('status.timeout');
        showOverlay('timeout', t('overlay.timeoutTitle'), t('overlay.timeoutBody'), t('overlay.returnMerchant'), function () {
            window.location.href = config.return_url || '/';
        });
    }

    function showSuccess(data) {
        setStatusKey('status.success');
        if (countdownTimer) clearInterval(countdownTimer);
        if (statusCheckTimer) clearInterval(statusCheckTimer);

        showOverlay('success', t('overlay.successTitle'), t('overlay.successBody'), t('overlay.returnMerchant'), function () {
            window.location.href = data.return_url || config.return_url || '/';
        });
    }

    function bindEvents() {
        $all('[data-step-target]').forEach(function (button) {
            button.addEventListener('click', function () {
                var target = button.dataset.stepTarget || '1';
                if (target === '2' && !paymentDetail) return;
                setStep(target);
            });
        });

        var continueButton = $('#continueButton');
        if (continueButton) {
            continueButton.addEventListener('click', lockPaymentMethod);
        }

        var copyAddressButton = $('#copyAddressButton');
        if (copyAddressButton) {
            copyAddressButton.addEventListener('click', function () {
                copyText(paymentDetail ? paymentDetail.token : '', copyAddressButton);
            });
        }

        var copyAmountButton = $('#copyAmountButton');
        if (copyAmountButton) {
            copyAmountButton.addEventListener('click', function () {
                copyText(paymentAmountValue, copyAmountButton);
            });
        }

        document.addEventListener('click', function (event) {
            if (!event.target.closest('.custom-select')) {
                closeDropdowns();
            }
            if (!event.target.closest('.language-switch')) {
                closeLanguageMenu();
            }
        });

        document.addEventListener('keydown', function (event) {
            if (event.key === 'Escape') {
                closeDropdowns();
                closeLanguageMenu();
            }
        });
    }

    function init(initialConfig) {
        config = initialConfig || {};
        tradeId = config.trade_id || '';
        methods = Array.isArray(config.network) ? config.network : [];
        paymentDetail = normalizeSelectedPayment(config.selected_payment);
        if (paymentDetail) {
            selectedMethod = selectedPaymentMethod(paymentDetail);
            selectedCurrency = selectedMethod ? selectedMethod.currency : '';
            selectedNetworkId = selectedMethod ? networkIdentity(selectedMethod) : '';
        }

        initLanguage();
        initTheme();
        bindEvents();
        renderPaymentSelectors();
        startCountdown();
        restartStatusPolling();

        if (paymentDetail && selectedMethod) {
            renderPaymentDetail();
            setStep('2');
        } else if (methods.length === 0) {
            showMessage(t('message.noNetworks'));
        }
    }

    window.Payment = {
        init: init
    };
})();
